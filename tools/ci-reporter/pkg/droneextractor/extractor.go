package droneextractor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

const (
	// BaseURL is the default Drone instance URL
	BaseURL = "https://drone.owncloud.com"
)

type Extractor struct {
	client *http.Client
}

func NewExtractor(client *http.Client) *Extractor {
	return &Extractor{client: client}
}

// IsDroneURL checks if a URL is a Drone CI URL
func (e *Extractor) IsDroneURL(url string) bool {
	return strings.Contains(url, "drone.owncloud.com")
}

func (e *Extractor) Extract(buildURL string) (*PipelineInfo, error) {
	apiURL := e.buildURLToAPIURL(buildURL)
	if apiURL == "" {
		return nil, fmt.Errorf("invalid build URL: %s", buildURL)
	}

	buildData, err := e.fetchDroneAPI(apiURL)
	if err != nil {
		return nil, fmt.Errorf("fetching Drone API: %w", err)
	}

	stages := e.parseStagesFromAPI(buildData, buildURL)

	var durationMinutes float64
	if buildData.Finished > 0 && buildData.Started > 0 {
		durationMinutes = float64(buildData.Finished-buildData.Started) / 60.0
	}

	return &PipelineInfo{
		BuildURL:        buildURL,
		Started:         buildData.Started,
		Finished:        buildData.Finished,
		DurationMinutes: durationMinutes,
		PipelineStages:  stages,
	}, nil
}

func (e *Extractor) buildURLToAPIURL(buildURL string) string {
	if !strings.Contains(buildURL, "drone.owncloud.com") {
		return ""
	}
	parts := strings.Split(buildURL, "/")
	if len(parts) < 5 {
		return ""
	}
	repo := parts[3] + "/" + parts[4]
	buildNum := parts[5]
	return fmt.Sprintf("https://drone.owncloud.com/api/repos/%s/builds/%s", repo, buildNum)
}

func (e *Extractor) fetchDroneAPI(url string) (*droneBuildResponse, error) {
	resp, err := e.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var buildData droneBuildResponse
	if err := json.NewDecoder(resp.Body).Decode(&buildData); err != nil {
		return nil, err
	}

	return &buildData, nil
}

func (e *Extractor) parseStagesFromAPI(buildData *droneBuildResponse, buildURL string) []PipelineStage {
	var stages []PipelineStage

	for _, apiStage := range buildData.Stages {
		normalizedStatus := e.normalizeStatus(apiStage.Status)
		if normalizedStatus == "success" || normalizedStatus == "skipped" {
			continue
		}

		stage := PipelineStage{
			StageNumber: apiStage.Number,
			StageName:   apiStage.Name,
			Status:      normalizedStatus,
			Steps:       []PipelineStep{},
		}

		for _, apiStep := range apiStage.Steps {
			normalizedStepStatus := e.normalizeStatus(apiStep.Status)
			if normalizedStepStatus == "skipped" {
				continue
			}
			stepURL := fmt.Sprintf("%s/%d/%d", buildURL, apiStage.Number, apiStep.Number)
			step := PipelineStep{
				StepNumber: apiStep.Number,
				StepName:   apiStep.Name,
				Status:     normalizedStepStatus,
				URL:        stepURL,
			}

			if normalizedStepStatus == "failure" {
				logs, err := e.fetchStepLogs(stepURL, buildURL, apiStage.Number, apiStep.Number)
				if err == nil {
					lines := strings.Split(logs, "\n")
					if len(lines) > 500 && (strings.Contains(logs, "--- Failed scenarios:") || strings.Contains(logs, "Scenario:")) {
						logs = e.extractBehatFailures(logs)
					}
					step.Logs = logs
				}
			}

			stage.Steps = append(stage.Steps, step)
		}

		stages = append(stages, stage)
	}

	return stages
}

func (e *Extractor) normalizeStatus(status string) string {
	switch strings.ToLower(status) {
	case "success", "passed":
		return "success"
	case "failure", "failed", "error":
		return "failure"
	case "running", "pending", "started":
		return "running"
	case "skipped":
		return "skipped"
	default:
		return status
	}
}

// droneBuildResponse matches the structure returned by Drone API:
// GET /api/repos/{owner}/{repo}/builds/{buildNumber}
// Verified against: https://drone.owncloud.com/api/repos/owncloud/ocis/builds/51629
type droneBuildResponse struct {
	Started  int64        `json:"started"`
	Finished int64        `json:"finished"`
	Stages   []droneStage `json:"stages"`
}

type droneStage struct {
	Number int         `json:"number"` // Stage number (1, 2, 3, ...)
	Name   string      `json:"name"`   // Stage name (e.g., "API-wopi", "coding-standard-php8.4")
	Status string      `json:"status"` // "success", "failure", "running", etc.
	Steps  []droneStep `json:"steps"`  // Array of steps in this stage
}

type droneStep struct {
	Number int    `json:"number"` // Step number within stage
	Name   string `json:"name"`   // Step name (e.g., "clone", "test-acceptance-api")
	Status string `json:"status"` // "success", "failure", "skipped", etc.
}

type droneLogEntry struct {
	Pos  int    `json:"pos"`
	Out  string `json:"out"`
	Time int64  `json:"time"`
}

type BuildInfo struct {
	BuildNumber int    `json:"build_number"`
	Status      string `json:"status"` // "success", "failure", "error"
	Started     int64  `json:"started"`
	Finished    int64  `json:"finished"`
	CommitSHA   string `json:"commit_sha"`
}

type droneBuildListItem struct {
	Number   int    `json:"number"`
	Status   string `json:"status"`
	Event    string `json:"event"`  // "pull_request", "push", etc.
	Source   string `json:"source"` // Source branch for PR builds
	Target   string `json:"target"` // Target branch for PR builds
	After    string `json:"after"`  // Commit SHA
	Started  int64  `json:"started"`
	Finished int64  `json:"finished"`
}

func (e *Extractor) fetchStepLogs(stepURL, buildURL string, stageNum, stepNum int) (string, error) {
	apiLogsURL := e.buildLogsAPIURL(buildURL, stageNum, stepNum)
	if apiLogsURL != "" {
		logs, err := e.fetchLogsFromAPI(apiLogsURL)
		if err == nil {
			return logs, nil
		}
		if err != nil && !strings.Contains(err.Error(), "401") && !strings.Contains(err.Error(), "403") {
			return "", err
		}
	}

	return e.fetchLogsFromHTML(stepURL)
}

func (e *Extractor) buildLogsAPIURL(buildURL string, stageNum, stepNum int) string {
	if !strings.Contains(buildURL, "drone.owncloud.com") {
		return ""
	}
	parts := strings.Split(buildURL, "/")
	if len(parts) < 5 {
		return ""
	}
	repo := parts[3] + "/" + parts[4]
	buildNum := parts[5]
	return fmt.Sprintf("https://drone.owncloud.com/api/repos/%s/builds/%s/logs/%d/%d", repo, buildNum, stageNum, stepNum)
}

func (e *Extractor) fetchLogsFromAPI(apiURL string) (string, error) {
	resp, err := e.client.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var logEntries []droneLogEntry
	if err := json.NewDecoder(resp.Body).Decode(&logEntries); err != nil {
		return "", err
	}

	var logs strings.Builder
	for _, entry := range logEntries {
		logs.WriteString(entry.Out)
	}

	return logs.String(), nil
}

func (e *Extractor) fetchLogsFromHTML(stepURL string) (string, error) {
	resp, err := e.client.Get(stepURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return e.extractLogsFromHTML(string(body))
}

func (e *Extractor) extractLogsFromHTML(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	var logs strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "pre" || n.Data == "code" {
				text := e.getTextContent(n)
				if text != "" {
					logs.WriteString(text)
					logs.WriteString("\n")
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)
	return strings.TrimSpace(logs.String()), nil
}

func (e *Extractor) getTextContent(n *html.Node) string {
	var text strings.Builder
	var collectText func(*html.Node)
	collectText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			collectText(c)
		}
	}
	collectText(n)
	return text.String()
}

func (e *Extractor) extractBehatFailures(logs string) string {
	var result strings.Builder

	// ANSI red code pattern
	redCodePattern := regexp.MustCompile(`\x1b\[31m|\033\[31m|\[31m`)
	ansiStripPattern := regexp.MustCompile(`\x1b\[[0-9;]*m|\033\[[0-9;]*m`)

	// Split by Scenario: or Scenario Outline: (case-insensitive)
	// Use regex to find all scenario starts while preserving the marker
	scenarioPattern := regexp.MustCompile(`(?i)(Scenario(?:\s+Outline)?:\s*)`)
	matches := scenarioPattern.FindAllStringIndex(logs, -1)

	if len(matches) == 0 {
		// No scenarios found, try to extract just "--- Failed scenarios:" section
		return e.extractFailedScenariosSection(logs, ansiStripPattern)
	}

	// Extract each scenario section
	for i, match := range matches {
		start := match[0]
		end := len(logs)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}

		scenarioSection := logs[start:end]

		// Check if this section contains ANSI red codes
		if redCodePattern.MatchString(scenarioSection) {
			// Strip ANSI codes and add to result
			cleaned := ansiStripPattern.ReplaceAllString(scenarioSection, "")
			if result.Len() > 0 {
				result.WriteString("\n\n")
			}
			result.WriteString(cleaned)
		}
	}

	// Extract "--- Failed scenarios:" section
	failedScenariosSection := e.extractFailedScenariosSection(logs, ansiStripPattern)
	if failedScenariosSection != "" {
		if result.Len() > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(failedScenariosSection)
	}

	extracted := strings.TrimSpace(result.String())
	if extracted == "" {
		// If extraction failed, return original logs
		return logs
	}

	return extracted
}

func (e *Extractor) extractFailedScenariosSection(logs string, ansiStripPattern *regexp.Regexp) string {
	failedScenariosMarker := "--- Failed scenarios:"
	if idx := strings.Index(logs, failedScenariosMarker); idx != -1 {
		failedScenariosSection := logs[idx:]
		cleaned := ansiStripPattern.ReplaceAllString(failedScenariosSection, "")
		return cleaned
	}
	return ""
}

// GetBuildsByCommitSHA queries Drone API for all builds matching a specific commit SHA
func (e *Extractor) GetBuildsByCommitSHA(baseURL, repo, commitSHA string) ([]BuildInfo, error) {
	var allBuilds []BuildInfo
	page := 1
	maxPages := 200 // Safety limit (25 items per page = 5000 builds max)

	for page <= maxPages {
		url := fmt.Sprintf("%s/api/repos/%s/builds?page=%d", baseURL, repo, page)
		resp, err := e.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetching build list: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		var buildList []droneBuildListItem
		if err := json.NewDecoder(resp.Body).Decode(&buildList); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decoding build list: %w", err)
		}
		resp.Body.Close()

		// Stop when we get an empty page
		if len(buildList) == 0 {
			break
		}

		// Filter builds by commit SHA
		for _, item := range buildList {
			if item.After == commitSHA {
				allBuilds = append(allBuilds, BuildInfo{
					BuildNumber: item.Number,
					Status:      item.Status,
					Started:     item.Started,
					Finished:    item.Finished,
					CommitSHA:   item.After,
				})
			}
		}

		page++
	}

	// Sort by build number (ascending = chronological)
	sort.Slice(allBuilds, func(i, j int) bool {
		return allBuilds[i].BuildNumber < allBuilds[j].BuildNumber
	})

	return allBuilds, nil
}

// GetBuildsByPRBranch queries Drone API for all builds matching PR event and source branch
func (e *Extractor) GetBuildsByPRBranch(baseURL, repo, sourceBranch string) ([]BuildInfo, error) {
	var allBuilds []BuildInfo
	page := 1
	maxPages := 200 // Safety limit (25 items per page = 5000 builds max)

	for page <= maxPages {
		url := fmt.Sprintf("%s/api/repos/%s/builds?page=%d", baseURL, repo, page)
		resp, err := e.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetching build list: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		var buildList []droneBuildListItem
		if err := json.NewDecoder(resp.Body).Decode(&buildList); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decoding build list: %w", err)
		}
		resp.Body.Close()

		// Stop when we get an empty page
		if len(buildList) == 0 {
			break
		}

		// Filter builds by event="pull_request" and source branch
		for _, item := range buildList {
			if item.Event == "pull_request" && item.Source == sourceBranch {
				allBuilds = append(allBuilds, BuildInfo{
					BuildNumber: item.Number,
					Status:      item.Status,
					Started:     item.Started,
					Finished:    item.Finished,
					CommitSHA:   item.After,
				})
			}
		}

		page++
	}

	// Sort by build number (ascending = chronological)
	sort.Slice(allBuilds, func(i, j int) bool {
		return allBuilds[i].BuildNumber < allBuilds[j].BuildNumber
	})

	return allBuilds, nil
}

// GetBuildsByPushBranch queries Drone API for all builds matching push/cron events to a specific branch (e.g., master).
// This is specifically for merged commits to main branches, not PR builds.
// It includes both "push" events (actual merges) and "cron" events (scheduled/nightly builds).
func (e *Extractor) GetBuildsByPushBranch(baseURL, repo, targetBranch string) ([]BuildInfo, error) {
	var allBuilds []BuildInfo
	page := 1
	maxPages := 200 // Safety limit (25 items per page = 5000 builds max)

	for page <= maxPages {
		url := fmt.Sprintf("%s/api/repos/%s/builds?page=%d", baseURL, repo, page)
		resp, err := e.client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("fetching build list: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		var buildList []droneBuildListItem
		if err := json.NewDecoder(resp.Body).Decode(&buildList); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decoding build list: %w", err)
		}
		resp.Body.Close()

		// Stop when we get an empty page
		if len(buildList) == 0 {
			break
		}

		// Filter builds by event="push" or event="cron" and target branch
		// "push" = actual merges to branch
		// "cron" = scheduled/nightly builds running on the branch
		for _, item := range buildList {
			if (item.Event == "push" || item.Event == "cron") && item.Target == targetBranch {
				allBuilds = append(allBuilds, BuildInfo{
					BuildNumber: item.Number,
					Status:      item.Status,
					Started:     item.Started,
					Finished:    item.Finished,
					CommitSHA:   item.After,
				})
			}
		}

		page++
	}

	// Sort by build number (ascending = chronological)
	sort.Slice(allBuilds, func(i, j int) bool {
		return allBuilds[i].BuildNumber < allBuilds[j].BuildNumber
	})

	return allBuilds, nil
}

// extractBaseURL extracts the base URL (e.g., "https://drone.owncloud.com") from a build URL
func extractBaseURL(buildURL string) string {
	if !strings.Contains(buildURL, "://") {
		return ""
	}
	parts := strings.Split(buildURL, "/")
	if len(parts) >= 3 {
		return parts[0] + "//" + parts[2]
	}
	return ""
}

// BuildInfoToPipelineSummary converts a BuildInfo to PipelineInfoSummary
func BuildInfoToPipelineSummary(build BuildInfo, baseURL, repo string) PipelineInfoSummary {
	var durationMinutes float64
	if build.Finished > 0 && build.Started > 0 {
		durationMinutes = float64(build.Finished-build.Started) / 60.0
	}

	// URL format: https://drone.owncloud.com/owncloud/ocis/{buildNumber}
	// repo format: "owncloud/ocis"
	buildURL := fmt.Sprintf("%s/%s/%d", baseURL, repo, build.BuildNumber)

	return PipelineInfoSummary{
		BuildURL:        buildURL,
		Started:         build.Started,
		Finished:        build.Finished,
		DurationMinutes: durationMinutes,
		Status:          build.Status,
	}
}

// BuildInfoToPipelineSummaryWithStages converts a BuildInfo to PipelineInfoSummary with failed stages
// by fetching full pipeline details from the build URL
func (e *Extractor) BuildInfoToPipelineSummaryWithStages(build BuildInfo, baseURL, repo string) PipelineInfoSummary {
	summary := BuildInfoToPipelineSummary(build, baseURL, repo)

	// Fetch full pipeline info to extract failed stages
	pipelineInfo, err := e.Extract(summary.BuildURL)
	if err == nil && pipelineInfo != nil {
		var failedStages []FailedStage
		for _, stage := range pipelineInfo.PipelineStages {
			if stage.Status == "failure" || stage.Status == "error" {
				failedStages = append(failedStages, FailedStage{
					StageNumber: stage.StageNumber,
					StageName:   stage.StageName,
					Status:      stage.Status,
				})
			}
		}
		summary.FailedStages = failedStages
	}

	return summary
}

// IsFailedBuild checks if a build status indicates failure
func IsFailedBuild(status string) bool {
	normalized := strings.ToLower(status)
	return normalized == "failure" || normalized == "failed" || normalized == "error"
}

// GetFailedBuildsByPRBranch returns failed builds for a PR branch with pipeline summaries
func (e *Extractor) GetFailedBuildsByPRBranch(baseURL, repo, sourceBranch string) ([]PipelineInfoSummary, error) {
	builds, err := e.GetBuildsByPRBranch(baseURL, repo, sourceBranch)
	if err != nil {
		return nil, err
	}

	var failedSummaries []PipelineInfoSummary
	for _, build := range builds {
		if IsFailedBuild(build.Status) {
			summary := e.BuildInfoToPipelineSummaryWithStages(build, baseURL, repo)
			failedSummaries = append(failedSummaries, summary)
			// Rate limiting delay for API calls in BuildInfoToPipelineSummaryWithStages
		}
	}

	return failedSummaries, nil
}
