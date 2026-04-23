package ghaextractor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiBase      = "https://api.github.com"
	requestDelay = 1 * time.Second
	retryDelay   = 2 * time.Second
	maxRetries   = 3
)

type Extractor struct {
	client *http.Client
	token  string
}

func NewExtractor(client *http.Client, token string) *Extractor {
	return &Extractor{client: client, token: token}
}

type WorkflowRun struct {
	ID         int64  `json:"id"`
	RunAttempt int    `json:"run_attempt"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	HTMLURL    string `json:"html_url"`
	HeadBranch string `json:"head_branch"`
	HeadSHA    string `json:"head_sha"`
	Event      string `json:"event"`
}

type workflowRunsPage struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type Job struct {
	ID          int64  `json:"id"`
	RunID       int64  `json:"run_id"`
	RunAttempt  int    `json:"run_attempt"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

type jobsPage struct {
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

// GetWorkflowRuns returns all runs for a workflow file created at or after `since` (RFC3339).
func (e *Extractor) GetWorkflowRuns(repo, workflow, since string) ([]WorkflowRun, error) {
	var all []WorkflowRun
	page := 1
	for {
		url := fmt.Sprintf("%s/repos/%s/actions/workflows/%s/runs?per_page=100&page=%d&status=completed", apiBase, repo, workflow, page)
		if since != "" {
			url += "&created=>=" + since
		}

		var result workflowRunsPage
		if err := e.get(url, &result); err != nil {
			return nil, err
		}
		all = append(all, result.WorkflowRuns...)
		if len(result.WorkflowRuns) < 100 {
			break
		}
		page++
	}
	return all, nil
}

// GetRunJobs returns all jobs for a given run across all attempts (filter=all).
func (e *Extractor) GetRunJobs(repo string, runID int64) ([]Job, error) {
	var all []Job
	page := 1
	for {
		url := fmt.Sprintf("%s/repos/%s/actions/runs/%d/jobs?filter=all&per_page=100&page=%d", apiBase, repo, runID, page)

		var result jobsPage
		if err := e.get(url, &result); err != nil {
			return nil, err
		}
		all = append(all, result.Jobs...)
		if len(result.Jobs) < 100 {
			break
		}
		page++
	}
	return all, nil
}

func (e *Extractor) get(url string, result any) error {
	resp, err := e.request(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	time.Sleep(requestDelay)
	return nil
}

func (e *Extractor) request(url string) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		if e.token != "" {
			req.Header.Set("Authorization", "Bearer "+e.token)
		}
		resp, err := e.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, lastErr
}
