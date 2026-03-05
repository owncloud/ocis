package githubextractor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	apiBase       = "https://api.github.com"
	requestDelay  = 1 * time.Second
	retryDelay    = 2 * time.Second
	maxRetries    = 3
)

type Extractor struct {
	client *http.Client
	token  string
}

func NewExtractor(client *http.Client, token string) *Extractor {
	return &Extractor{
		client: client,
		token:  token,
	}
}

type PRInfo struct {
	Number int `json:"number"`
	Head   struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

type PRCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

type CommitListItem struct {
	SHA     string `json:"sha"`
	HTMLURL string `json:"html_url"`
	Commit  struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

type CombinedStatus struct {
	State    string          `json:"state"`
	Statuses []StatusContext `json:"statuses"`
}

type StatusContext struct {
	Context     string `json:"context"`
	State       string `json:"state"`
	TargetURL   string `json:"target_url"`
	Description string `json:"description"`
}

func (e *Extractor) GetPRInfo(repo string, prNumber int) (*PRInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/pulls/%d", apiBase, repo, prNumber)
	var prInfo PRInfo
	return &prInfo, e.get(url, &prInfo)
}

func (e *Extractor) GetPRCommits(repo string, prNumber int) ([]PRCommit, error) {
	url := fmt.Sprintf("%s/repos/%s/pulls/%d/commits", apiBase, repo, prNumber)
	var commits []PRCommit
	return commits, e.get(url, &commits)
}

func (e *Extractor) GetCommits(repo, branch string, maxCommits int, since string) ([]CommitListItem, error) {
	var commits []CommitListItem
	page := 1

	for len(commits) < maxCommits {
		url := fmt.Sprintf("%s/repos/%s/commits?sha=%s&per_page=100&page=%d", apiBase, repo, branch, page)
		if since != "" {
			url += "&since=" + since
		}

		resp, err := e.request(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
		}

		var pageCommits []CommitListItem
		if err := json.NewDecoder(resp.Body).Decode(&pageCommits); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if len(pageCommits) == 0 {
			break
		}

		commits = append(commits, pageCommits...)
		if len(commits) >= maxCommits {
			commits = commits[:maxCommits]
			break
		}

		page++
		time.Sleep(requestDelay)
	}

	return commits, nil
}

func (e *Extractor) GetCommitStatus(repo, sha string) (*CombinedStatus, error) {
	url := fmt.Sprintf("%s/repos/%s/commits/%s/status", apiBase, repo, sha)
	var status CombinedStatus
	return &status, e.get(url, &status)
}

// GetPRsForCommit returns PRs associated with a commit (e.g. merged PR that introduced the commit).
// See https://docs.github.com/en/rest/commits/commits#list-pull-requests-associated-with-a-commit
func (e *Extractor) GetPRsForCommit(repo, commitSHA string) ([]PRInfo, error) {
	url := fmt.Sprintf("%s/repos/%s/commits/%s/pulls", apiBase, repo, commitSHA)
	var prs []PRInfo
	return prs, e.get(url, &prs)
}

func isRetryableNetworkError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "broken pipe") ||
		strings.Contains(s, "connection reset") ||
		strings.Contains(s, "connection refused") ||
		strings.Contains(s, "EOF")
}

// get makes an HTTP GET request, decodes JSON response, and handles rate limiting
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
			if isRetryableNetworkError(err) && attempt < maxRetries-1 {
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, lastErr
}

// Pipeline represents a single Drone pipeline run
type Pipeline struct {
	ID       int     `json:"id"`
	Status   string  `json:"status"`
	Started  int64   `json:"started"`
	Finished int64   `json:"finished"`
	Duration float64 `json:"duration_minutes"`
	Trigger  string  `json:"trigger"`
}

// OriginalCommit represents an individual commit from the PR branch before squashing
type OriginalCommit struct {
	SHA       string     `json:"sha"`
	Title     string     `json:"title"`
	Pipelines []Pipeline `json:"pipelines,omitempty"`
}

// ExtractedCommit represents a commit with PR information
type ExtractedCommit struct {
	PR              int              `json:"pr"`
	SHA             string           `json:"sha"`
	Date            string           `json:"date"`
	Title           string           `json:"title"`
	HTMLURL         string           `json:"html_url"`
	Pipelines       []Pipeline       `json:"pipelines,omitempty"`
	OriginalCommits []OriginalCommit `json:"original_commits,omitempty"`
}

// ExtractedData represents the output structure
type ExtractedData struct {
	Commits []ExtractedCommit `json:"commits"`
}

// extractSubject extracts the first line of the commit message
func extractSubject(message string) string {
	if idx := strings.Index(message, "\n"); idx >= 0 {
		return strings.TrimSpace(message[:idx])
	}
	return strings.TrimSpace(message)
}

// ExtractCommitData fetches commits and extracts them into structured format
func (e *Extractor) ExtractCommitData(repo, branch string, maxCommits int, since string) (*ExtractedData, error) {
	commits, err := e.GetCommits(repo, branch, maxCommits, since)
	if err != nil {
		return nil, fmt.Errorf("fetching commits: %w", err)
	}

	extractedCommits := make([]ExtractedCommit, 0, len(commits))
	originalCommitSHAs := make(map[string]bool) // Track SHAs that appear as original commits

	for _, commit := range commits {
		title := extractSubject(commit.Commit.Message)

		// Fetch PR information from GitHub API
		prNumber := 0
		var originalCommits []OriginalCommit
		prs, err := e.GetPRsForCommit(repo, commit.SHA)
		if err == nil && len(prs) > 0 {
			// Use the first PR number if multiple PRs are associated
			prNumber = prs[0].Number

			// Fetch original commits from the PR before squashing
			prCommits, err := e.GetPRCommits(repo, prNumber)
			if err == nil && len(prCommits) > 0 {
				originalCommits = make([]OriginalCommit, 0, len(prCommits))
				for _, prCommit := range prCommits {
					// Skip if this is the same commit (don't self-reference)
					if prCommit.SHA == commit.SHA {
						continue
					}
					originalCommits = append(originalCommits, OriginalCommit{
						SHA:   prCommit.SHA,
						Title: extractSubject(prCommit.Commit.Message),
					})
					// Track this SHA as an original commit
					originalCommitSHAs[prCommit.SHA] = true
				}
			}
		}

		extractedCommits = append(extractedCommits, ExtractedCommit{
			PR:              prNumber,
			SHA:             commit.SHA,
			Date:            commit.Commit.Author.Date,
			Title:           title,
			HTMLURL:         commit.HTMLURL,
			OriginalCommits: originalCommits,
		})
	}

	// Filter out commits that are already listed as original commits in merge commits
	// These are PR commits that appear in master history but should only show under the merge commit
	filteredCommits := make([]ExtractedCommit, 0, len(extractedCommits))
	for _, commit := range extractedCommits {
		// Keep the commit ONLY if it's NOT in the originalCommitSHAs set
		// If a commit appears in any other commit's original_commits, it should not be a top-level entry
		if !originalCommitSHAs[commit.SHA] {
			filteredCommits = append(filteredCommits, commit)
		}
	}

	return &ExtractedData{
		Commits: filteredCommits,
	}, nil
}
