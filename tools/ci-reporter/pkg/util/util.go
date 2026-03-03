package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ParseSinceFlag parses a date/time string in multiple formats:
// - YYYY-MM-DD: absolute date (e.g., "2026-02-15")
// - RFC3339: timestamp (e.g., "2026-02-15T10:30:00Z")
// - Nd: relative days (e.g., "7d", "30d")
//
// Returns the parsed time in RFC3339 format.
func ParseSinceFlag(since string) (string, error) {
	if since == "" {
		return "", nil
	}

	// Try parsing as YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", since); err == nil {
		return t.Format(time.RFC3339), nil
	}

	// Try parsing as RFC3339 (backward compatibility)
	if t, err := time.Parse(time.RFC3339, since); err == nil {
		return t.Format(time.RFC3339), nil
	}

	// Try parsing as relative days format (e.g., "7d", "30d")
	if strings.HasSuffix(since, "d") {
		daysStr := strings.TrimSuffix(since, "d")
		days, err := time.ParseDuration(daysStr + "h")
		if err == nil {
			// Convert days to hours (Nd = N*24h)
			hours := days.Hours()
			t := time.Now().Add(-time.Duration(hours*24) * time.Hour)
			return t.Format(time.RFC3339), nil
		}
		// Try parsing as integer
		var numDays int
		if _, err := fmt.Sscanf(daysStr, "%d", &numDays); err == nil {
			t := time.Now().AddDate(0, 0, -numDays)
			return t.Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("invalid format: use YYYY-MM-DD, RFC3339, or 'Nd' (e.g., '7d')")
}

// ExtractSubject extracts the first line (subject) from a commit message.
func ExtractSubject(message string) string {
	lines := strings.Split(message, "\n")
	return strings.TrimSpace(lines[0])
}

// RedactGitHubToken redacts GitHub tokens from command strings for safe logging.
// Matches --github-token=<token> or --github-token <token> patterns.
func RedactGitHubToken(command string) string {
	// Match --github-token=<token>
	pattern1 := regexp.MustCompile(`--github-token=[^\s]+`)
	command = pattern1.ReplaceAllString(command, "--github-token=***")

	// Match --github-token <token>
	pattern2 := regexp.MustCompile(`--github-token\s+[^\s]+`)
	command = pattern2.ReplaceAllString(command, "--github-token ***")

	return command
}
