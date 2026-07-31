// Usage examples:
//   go run ./cmd/list-flaky-tests --since 14d
//   go run ./cmd/list-flaky-tests --since 3m --workflow acceptance-tests.yml
//   go run ./cmd/list-flaky-tests --github-token $GITHUB_TOKEN --since 30d
//   go run ./cmd/list-flaky-tests --repo owncloud/ocis --since 7d --out report.json

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/tools/ci-reporter/pkg/ghaextractor"
	"github.com/owncloud/ocis/v2/tools/ci-reporter/pkg/util"
)

const (
	defaultRepo     = "owncloud/ocis"
	defaultWorkflow = "acceptance-tests.yml"
)

type Summary struct {
	Args                  string `json:"args"`
	Repo                  string `json:"repo"`
	Workflow              string `json:"workflow"`
	DateRange             Range  `json:"date_range"`
	TotalRuns             int    `json:"total_runs"`
	RestartedRuns         int    `json:"restarted_runs"`
	RestartRatePct        float64 `json:"restart_rate_pct"`
	TotalTimeLostMinutes  float64 `json:"total_time_lost_minutes"`
}

type Range struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type FlakyJob struct {
	JobName          string  `json:"job_name"`
	FailThenPassCount int    `json:"fail_then_pass_count"`
	FailureRatePct   float64 `json:"failure_rate_pct"`
}

type Report struct {
	Summary   Summary    `json:"summary"`
	FlakyJobs []FlakyJob `json:"flaky_jobs"`
}

func main() {
	repo := flag.String("repo", defaultRepo, "GitHub repository (owner/repo)")
	workflow := flag.String("workflow", defaultWorkflow, "Workflow filename (e.g. acceptance-tests.yml)")
	since := flag.String("since", "14d", "Date filter: YYYY-MM-DD, RFC3339, or 'Nd'/'Nm' (e.g. '7d', '3m')")
	outFile := flag.String("out", "", "Output file path (default: stdout)")
	githubToken := flag.String("github-token", "", "GitHub token (falls back to GITHUB_TOKEN env)")
	flag.Parse()

	token := strings.TrimSpace(*githubToken)
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}
	if token != "" {
		fmt.Fprintf(os.Stderr, "Using GitHub token\n")
	} else {
		fmt.Fprintf(os.Stderr, "No token provided, unauthenticated (60 req/hr limit)\n")
	}

	sinceRaw := *since
	// Support Nm month notation (not in ParseSinceFlag)
	sinceRFC, err := parseSince(sinceRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing --since: %v\n", err)
		os.Exit(1)
	}

	command := strings.Join(os.Args[1:], " ")
	command = util.RedactGitHubToken(command)

	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 2,
	}
	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}
	extractor := ghaextractor.NewExtractor(client, token)

	fmt.Fprintf(os.Stderr, "Fetching workflow runs for %s/%s since %s...\n", *repo, *workflow, sinceRaw)
	runs, err := extractor.GetWorkflowRuns(*repo, *workflow, sinceRFC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching runs: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Found %d completed runs\n", len(runs))

	report, err := buildReport(extractor, *repo, *workflow, runs, command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building report: %v\n", err)
		os.Exit(1)
	}

	var w io.Writer = os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		w = f
		fmt.Fprintf(os.Stderr, "Writing JSON to %s...\n", *outFile)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func buildReport(extractor *ghaextractor.Extractor, repo, workflow string, runs []ghaextractor.WorkflowRun, command string) (*Report, error) {
	// The API returns one record per attempt, so a pipeline retried twice yields 3 records.
	// Deduplicate by head_sha, keeping the highest-attempt record for each unique pipeline.
	type pipelineRecord struct {
		latest     ghaextractor.WorkflowRun
		maxAttempt int
	}
	bySHA := make(map[string]*pipelineRecord)
	for _, run := range runs {
		r := bySHA[run.HeadSHA]
		if r == nil {
			bySHA[run.HeadSHA] = &pipelineRecord{latest: run, maxAttempt: run.RunAttempt}
		} else if run.RunAttempt > r.maxAttempt {
			r.latest = run
			r.maxAttempt = run.RunAttempt
		}
	}

	pipelines := make([]pipelineRecord, 0, len(bySHA))
	for _, r := range bySHA {
		pipelines = append(pipelines, *r)
	}

	var minDate, maxDate string
	for _, r := range pipelines {
		d := r.latest.CreatedAt
		if minDate == "" || d < minDate {
			minDate = d
		}
		if maxDate == "" || d > maxDate {
			maxDate = d
		}
	}

	totalPipelines := len(pipelines)
	restartedRuns := 0
	var timeLostMinutes float64
	// job name -> count of pipelines where that job was retried (any re-run, regardless of outcome)
	flakyCount := make(map[string]int)
	for i, r := range pipelines {
		run := r.latest
		fmt.Fprintf(os.Stderr, "  [%d/%d] pipeline %d (max attempt %d)...\n", i+1, len(pipelines), run.ID, r.maxAttempt)

		wasRestarted := r.maxAttempt > 1
		if !wasRestarted {
			continue
		}
		restartedRuns++

		// Fetch all jobs across all attempts for this run to find which jobs were re-run
		// and measure actual wall-clock time lost on retry attempts.
		jobs, err := extractor.GetRunJobs(repo, run.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠  run %d: %v\n", run.ID, err)
			continue
		}

		// Track which job names appeared in attempt > 1 (= developer had to wait for them again)
		retriedJobNames := make(map[string]bool)
		var retryStart, retryEnd time.Time
		for _, j := range jobs {
			if j.RunAttempt <= 1 || j.Name == "all-acceptance-tests" {
				continue
			}
			retriedJobNames[j.Name] = true
			if j.StartedAt != "" && j.CompletedAt != "" {
				if s, err := time.Parse(time.RFC3339, j.StartedAt); err == nil {
					if retryStart.IsZero() || s.Before(retryStart) {
						retryStart = s
					}
				}
				if e, err := time.Parse(time.RFC3339, j.CompletedAt); err == nil {
					if e.After(retryEnd) {
						retryEnd = e
					}
				}
			}
		}

		for name := range retriedJobNames {
			flakyCount[name]++
		}
		if !retryStart.IsZero() && retryEnd.After(retryStart) {
			timeLostMinutes += retryEnd.Sub(retryStart).Minutes()
		}
	}

	restartRatePct := 0.0
	if totalPipelines > 0 {
		restartRatePct = float64(restartedRuns) / float64(totalPipelines) * 100.0
	}

	flakyJobs := make([]FlakyJob, 0, len(flakyCount))
	for name, count := range flakyCount {
		ratePct := 0.0
		if totalPipelines > 0 {
			ratePct = float64(count) / float64(totalPipelines) * 100.0
		}
		flakyJobs = append(flakyJobs, FlakyJob{
			JobName:           name,
			FailThenPassCount: count,
			FailureRatePct:    ratePct,
		})
	}
	sort.Slice(flakyJobs, func(i, j int) bool {
		return flakyJobs[i].FailThenPassCount > flakyJobs[j].FailThenPassCount
	})

	return &Report{
		Summary: Summary{
			Args:                 command,
			Repo:                 repo,
			Workflow:             workflow,
			DateRange:            Range{Start: minDate, End: maxDate},
			TotalRuns:            totalPipelines,
			RestartedRuns:        restartedRuns,
			RestartRatePct:       restartRatePct,
			TotalTimeLostMinutes: timeLostMinutes,
		},
		FlakyJobs: flakyJobs,
	}, nil
}


// parseSince extends util.ParseSinceFlag with Nm (months) and Nw (weeks) support.
func parseSince(s string) (string, error) {
	if numStr, ok := strings.CutSuffix(s, "m"); ok {
		var months int
		if _, err := fmt.Sscanf(numStr, "%d", &months); err == nil {
			return time.Now().AddDate(0, -months, 0).Format(time.RFC3339), nil
		}
	}
	if numStr, ok := strings.CutSuffix(s, "w"); ok {
		var weeks int
		if _, err := fmt.Sscanf(numStr, "%d", &weeks); err == nil {
			return time.Now().AddDate(0, 0, -weeks*7).Format(time.RFC3339), nil
		}
	}
	return util.ParseSinceFlag(s)
}
