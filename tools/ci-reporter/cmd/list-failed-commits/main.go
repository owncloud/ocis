// Usage examples:
//   go run ./cmd/list-failed-commits --max-commits 20
//   go run ./cmd/list-failed-commits --max-commits 50 --out failed.json
//   go run ./cmd/list-failed-commits --github-token $GITHUB_TOKEN --max-commits 100
//   go run ./cmd/list-failed-commits --repo owncloud/ocis --branch master --max-commits 20
//   go run ./cmd/list-failed-commits --since 7d --max-commits 50
//   go run ./cmd/list-failed-commits --since 30d --out failed.json
//   go run ./cmd/list-failed-commits --since 2026-01-01 --max-commits 100
//   go run ./cmd/list-failed-commits --since 2026-02-15 --out recent.json

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

	"github.com/owncloud/ocis/v2/tools/ci-reporter/pkg/droneextractor"
	"github.com/owncloud/ocis/v2/tools/ci-reporter/pkg/githubextractor"
	"github.com/owncloud/ocis/v2/tools/ci-reporter/pkg/util"
)

const (
	defaultRepo      = "owncloud/ocis"
	defaultBranch    = "master"
	requestDelay     = 1 * time.Second
	defaultLimit     = 20
	defaultUnlimited = 999999999
)

type RunConfig struct {
	Repo                      string
	Branch                    string
	MaxCommits                int
	SinceRaw                  string
	SinceParsed               string
	OutFile                   string
	Token                     string
	Command                   string
	FetchFailedPipelineHistory bool
}

type FailureCount struct {
	StageName string `json:"stage_name"`
	StepName  string `json:"step_name"`
	Count     int    `json:"count"`
}

type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type Raport struct {
	Args          string    `json:"args"`
	Repo          string    `json:"repo"`
	Branch        string    `json:"branch"`
	DateRange     DateRange `json:"date_range"`
	TotalCommits  int       `json:"total_commits"`
	FailedCommits int       `json:"failed_commits"`
	Percentage    float64   `json:"percentage"`
}

type RaportMetadata struct {
	Args         string
	Repo         string
	Branch       string
	DateRange    DateRange
	TotalCommits int
}

type Report struct {
	Raport    *Raport        `json:"raport"`
	Count     []FailureCount `json:"count"`
	Pipelines []FailedCommit `json:"pipelines"`
}

type FailedCommit struct {
	PR                    int                                  `json:"pr,omitempty"`
	SHA                   string                               `json:"sha"`
	Date                  string                               `json:"date"`
	Subject               string                               `json:"subject"`
	HTMLURL               string                               `json:"html_url"`
	CombinedState         string                               `json:"combined_state"`
	FailedContexts        []StatusContext                      `json:"failed_contexts"`
	PipelineInfo          *droneextractor.PipelineInfo         `json:"pipeline_info,omitempty"`
	FailedPipelineHistory []droneextractor.PipelineInfoSummary `json:"failed_pipeline_history,omitempty"`
	PrPipelineHistory     []droneextractor.PipelineInfoSummary `json:"pr_pipeline_history,omitempty"`
}

type StatusContext struct {
	Context     string `json:"context"`
	State       string `json:"state"`
	TargetURL   string `json:"target_url"`
	Description string `json:"description"`
}

func main() {
	repo := flag.String("repo", defaultRepo, "GitHub repository (owner/repo)")
	branch := flag.String("branch", defaultBranch, "Branch name")
	maxCommits := flag.Int("max-commits", 0, "Maximum commits to check (default: 20 if no --since, unlimited if --since set)")
	since := flag.String("since", "", "Date filter: YYYY-MM-DD, RFC3339 timestamp, or 'Nd' for last N days (e.g., '7d', '30d')")
	outFile := flag.String("out", "", "Output file path (default: stdout)")
	githubToken := flag.String("github-token", "", "GitHub token (optional; falls back to env GITHUB_TOKEN)")
	failedPipelineHistory := flag.Bool("failed_pipeline_history", false, "if set, fetch failed_pipeline_history per commit (extra Drone API calls)")
	flag.Parse()

	token := strings.TrimSpace(*githubToken)
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}

	if token != "" {
		fmt.Fprintf(os.Stderr, "Using GitHub token\n")
	} else {
		fmt.Fprintf(os.Stderr, "No token provided, unauthenticated, 60 req/hr limit\n")
	}

	command := strings.Join(os.Args[1:], " ")
	command = util.RedactGitHubToken(command)

	parsedSince, err := util.ParseSinceFlag(*since)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing --since flag: %v\n", err)
		os.Exit(1)
	}

	// Determine final max commits based on git-like UX:
	// - No flags: default to 20
	// - --since only: unlimited (all in range)
	// - --max-commits only: use explicit value
	// - Both: use explicit value (both filters apply)
	finalMaxCommits := *maxCommits
	if finalMaxCommits == 0 {
		if *since != "" {
			finalMaxCommits = defaultUnlimited // Unlimited when --since is used
		} else {
			finalMaxCommits = defaultLimit // Default for no flags
		}
	}

	config := RunConfig{
		Repo:                       *repo,
		Branch:                     *branch,
		MaxCommits:                 finalMaxCommits,
		SinceRaw:                   *since,
		SinceParsed:                parsedSince,
		OutFile:                    *outFile,
		Token:                      token,
		Command:                    command,
		FetchFailedPipelineHistory: *failedPipelineHistory,
	}

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cfg RunConfig) error {
	// Configure Transport to prevent stale connection reuse during long pagination
	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
		MaxIdleConnsPerHost: 2,
	}
	client := &http.Client{
		Timeout:   60 * time.Second, // Increased for long-running pagination
		Transport: transport,
	}
	githubExtractor := githubextractor.NewExtractor(client, cfg.Token)
	droneExtractor := droneextractor.NewExtractor(client)

	// Build scanning message
	var scanMsg string
	if cfg.SinceRaw != "" {
		scanMsg = fmt.Sprintf("Scanning commits since %s from %s/%s", cfg.SinceRaw, cfg.Repo, cfg.Branch)
	} else {
		scanMsg = fmt.Sprintf("Scanning up to %d commits from %s/%s", cfg.MaxCommits, cfg.Repo, cfg.Branch)
	}
	fmt.Fprintf(os.Stderr, "%s\n", scanMsg)
	commits, err := githubExtractor.GetCommits(cfg.Repo, cfg.Branch, cfg.MaxCommits, cfg.SinceParsed)
	if err != nil {
		return fmt.Errorf("fetching commits: %w", err)
	}

	var failedCommits []FailedCommit
	for _, commit := range commits {
		status, err := githubExtractor.GetCommitStatus(cfg.Repo, commit.SHA)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠  %s: %v\n", commit.SHA[:7], err)
			time.Sleep(requestDelay)
			continue
		}

		if status.State == "failure" || status.State == "error" {
			var failedCtxs []StatusContext
			var pipelineInfo *droneextractor.PipelineInfo

			for _, s := range status.Statuses {
				if s.State != "success" {
					// Skip Drone entries - they'll be in pipeline_info instead
					if droneExtractor.IsDroneURL(s.TargetURL) {
						info, err := droneExtractor.Extract(s.TargetURL)
						if err == nil {
							pipelineInfo = info
						}
						// Don't add to failed_contexts to avoid redundancy
						continue
					}
					failedCtxs = append(failedCtxs, toStatusContext(s))
				}
			}
			// TODO: Refactor redundant field assignment - prPipelineHistory always assigned, failedPipelineHistory conditional
			var failedPipelineHistory []droneextractor.PipelineInfoSummary
			var prPipelineHistory []droneextractor.PipelineInfoSummary
			var prNumber int
			var headRef string

			// Get PR for this commit (master commit -> PR ID -> PR branch)
			prs, err := githubExtractor.GetPRsForCommit(cfg.Repo, commit.SHA)
			if err == nil && len(prs) > 0 {
				prNumber = prs[0].Number
				headRef = prs[0].Head.Ref

				// Get failed builds by PR branch (black box: all logic in droneextractor)
				failedBuilds, err := droneExtractor.GetFailedBuildsByPRBranch(droneextractor.BaseURL, cfg.Repo, headRef)
				if err == nil {
					prPipelineHistory = failedBuilds
					if cfg.FetchFailedPipelineHistory {
						failedPipelineHistory = failedBuilds
					}
				}
			}

			// Initialize empty slice if flag is set but no PR/builds found
			if cfg.FetchFailedPipelineHistory && failedPipelineHistory == nil {
				failedPipelineHistory = []droneextractor.PipelineInfoSummary{}
			}

			subject := util.ExtractSubject(commit.Commit.Message)
			failedCommits = append(failedCommits, FailedCommit{
				PR:                    prNumber,
				SHA:                   commit.SHA,
				Date:                  commit.Commit.Author.Date,
				Subject:               subject,
				HTMLURL:               commit.HTMLURL,
				CombinedState:         status.State,
				FailedContexts:        failedCtxs,
				PipelineInfo:          pipelineInfo,
				FailedPipelineHistory: failedPipelineHistory,
				PrPipelineHistory:     prPipelineHistory,
			})

			// Single line output: ✗ sha subject -> PR #num (branch: name)
			if prNumber > 0 && headRef != "" {
				fmt.Fprintf(os.Stderr, "✗ %s  %s -> PR #%d (branch: %s)\n", commit.SHA[:7], subject, prNumber, headRef)
			} else {
				fmt.Fprintf(os.Stderr, "✗ %s  %s\n", commit.SHA[:7], subject)
			}
		}

		time.Sleep(requestDelay)
	}

	fmt.Fprintf(os.Stderr, "\nResult: %d failed / %d total\n", len(failedCommits), len(commits))

	var minDate, maxDate string
	for _, c := range commits {
		date := c.Commit.Author.Date
		if minDate == "" || date < minDate {
			minDate = date
		}
		if maxDate == "" || date > maxDate {
			maxDate = date
		}
	}

	metadata := RaportMetadata{
		Args:         cfg.Command,
		Repo:         cfg.Repo,
		Branch:       cfg.Branch,
		DateRange:    DateRange{StartDate: minDate, EndDate: maxDate},
		TotalCommits: len(commits),
	}

	report := generateReport(failedCommits, metadata)

	var w io.Writer = os.Stdout
	if cfg.OutFile != "" {
		f, err := os.Create(cfg.OutFile)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
		fmt.Fprintf(os.Stderr, "Writing JSON to %s...\n", cfg.OutFile)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}

func toStatusContext(ctx githubextractor.StatusContext) StatusContext {
	return StatusContext{
		Context:     ctx.Context,
		State:       ctx.State,
		TargetURL:   ctx.TargetURL,
		Description: ctx.Description,
	}
}


func generateReport(failedCommits []FailedCommit, metadata RaportMetadata) *Report {
	countMap := make(map[string]int)

	for _, commit := range failedCommits {
		if commit.PipelineInfo == nil {
			continue
		}
		for _, stage := range commit.PipelineInfo.PipelineStages {
			for _, step := range stage.Steps {
				if step.Status == "failure" {
					key := stage.StageName + "/" + step.StepName
					countMap[key]++
				}
			}
		}
	}

	var counts []FailureCount
	for key, count := range countMap {
		parts := strings.Split(key, "/")
		if len(parts) >= 2 {
			counts = append(counts, FailureCount{
				StageName: parts[0],
				StepName:  strings.Join(parts[1:], "/"),
				Count:     count,
			})
		}
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	percentage := 0.0
	if metadata.TotalCommits > 0 {
		percentage = float64(len(failedCommits)) / float64(metadata.TotalCommits) * 100.0
	}

	raport := &Raport{
		Args:          metadata.Args,
		Repo:          metadata.Repo,
		Branch:        metadata.Branch,
		DateRange:     metadata.DateRange,
		TotalCommits:  metadata.TotalCommits,
		FailedCommits: len(failedCommits),
		Percentage:    percentage,
	}

	return &Report{
		Raport:    raport,
		Count:     counts,
		Pipelines: failedCommits,
	}
}
