package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	since, until, period, err := getTimeframe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	logs, err := getGitLog(since, until)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	csvRows := make([][]string, 0, 1000)

	for _, logLine := range logs {
		logParts := strings.Split(logLine, " ")
		if len(logParts) < 2 {
			continue
		}

		commit := logParts[0]
		date := logParts[1]

		diffLines, err := getGitDiff(commit)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		addedTests, changedTests, deletedTests := 0, 0, 0
		// var inScenarios bool

		for i, line := range diffLines {
			switch {
			case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
				if strings.Contains(line, "Scenario:") {
					addedTests++
				} else if strings.Contains(line, "Scenario Outline:") {
					addedTests += countAddedTestsInExamples(diffLines, i)
				}
			case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
				if strings.Contains(line, "Scenario") {
					deletedTests++
				}
			case strings.Contains(line, "@@ Feature:"):
				changedTests := 0
				for i, line := range diffLines {
					if strings.Contains(line, "@@ Feature:") {
						inScenarios, changed := checkChangedTests(diffLines, i)
						if !inScenarios {
							changedTests += changed
						}
					}
				}
			}
		}

		csvRows = append(csvRows, []string{"API Test", date, strconv.Itoa(addedTests), strconv.Itoa(changedTests), strconv.Itoa(deletedTests), commit})
	}

	// Ensure the directory exists
	reportDir := "tests/qa-activity-report/reports"
	if err := os.MkdirAll(reportDir, os.ModePerm); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Define the path for the CSV file
	filePath := fmt.Sprintf("%s/QA_Activity_Report_%s.csv", reportDir, period)
	if err := generateCSV(csvRows, filePath); err != nil {
		fmt.Println("Error writing CSV report:", err)
	} else {
		fmt.Println("CSV report generated successfully. You can find it in", filePath)
	}
}

func getTimeframe() (since string, until string, period string, err error) {
	monthStr := os.Getenv("MONTH")
	yearStr := os.Getenv("YEAR")
	daysStr := os.Getenv("DAYS")

	if monthStr != "" && yearStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			return "", "", "", fmt.Errorf("invalid month: %w", err)
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return "", "", "", fmt.Errorf("invalid year: %w", err)
		}
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, -1)
		since = startDate.Format("2006-01-02")
		until = endDate.Format("2006-01-02")
		period = fmt.Sprintf("%02d_%04d", month, year)
	} else if daysStr != "" {
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return "", "", "", fmt.Errorf("invalid number of days: %w", err)
		}
		until = time.Now().Format("2006-01-02")
		since = time.Now().AddDate(0, 0, -days).Format("2006-01-02")
		period = fmt.Sprintf("Last_%d_days", days)
	} else {
		return "", "", "", fmt.Errorf("please provide either MONTH and YEAR or DAYS")
	}

	return since, until, period, nil
}

func getGitLog(since, until string) ([]string, error) {
	cmd := exec.Command("git", "log", "--since="+since, "--until="+until, "--pretty=format:%H %ad", "--date=short", "--", "tests/acceptance/features", ":(exclude)tests/acceptance/features/bootstrap/")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return strings.Split(out.String(), "\n"), nil
}

func getGitDiff(commit string) ([]string, error) {
	cmd := exec.Command("git", "diff", commit+"~1", commit, "--", "tests/acceptance/features", ":(exclude)tests/acceptance/features/bootstrap/")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return strings.Split(out.String(), "\n"), nil
}

func countAddedTestsInExamples(diffLines []string, startIndex int) int {
	var inExamples bool
	addedTests := 0

	for j := startIndex + 1; j < len(diffLines); j++ {
		exampleLine := diffLines[j]
		if strings.HasPrefix(exampleLine, "+") && strings.Contains(exampleLine, "Examples:") {
			inExamples = true
			continue
		} else if inExamples {
			trimmedLine := strings.TrimSpace(exampleLine)
			if strings.HasPrefix(trimmedLine, "+") && strings.Contains(trimmedLine, "|") {
				// Count a string if it starts with "+" and contains "|"
				addedTests++
			} else if strings.TrimSpace(exampleLine) == "" || !strings.HasPrefix(trimmedLine, "+") || !strings.HasPrefix(trimmedLine, "|") {
				// Abort counting when a row that does not belong to the table is encountered
				break
			}
		}
	}

	// We have one line | resource | which is not a test line. So we deleted one line from addedTest
	// Examples:
	// 	| resource      |
	// 	| testfile.txt  |
	// 	| FolderToShare |
	if inExamples {
		addedTests--
	}
	return addedTests
}

func checkChangedTests(diffLines []string, startIndex int) (bool, int) {
	var inScenarios bool
	changedTests := 0

	for j := startIndex + 1; j < len(diffLines); j++ {
		scenarioLine := diffLines[j]
		if strings.HasPrefix(scenarioLine, "+") || strings.HasPrefix(scenarioLine, "-") {
			// If there are changes and the string contains the word "Scenario", set inScenarios to true
			if strings.Contains(scenarioLine, "Scenario") {
				inScenarios = true
				break
			}
			// If the line no longer starts with "-" "+", then the change block has ended
			if !strings.HasPrefix(scenarioLine, "-") && !strings.HasPrefix(scenarioLine, "+") {
				break
			}
		}
	}

	// If we didn't find "Scenario" in the changes, increase the changedTests counter
	if !inScenarios {
		changedTests++
	}
	return inScenarios, changedTests
}

func generateCSV(csvRows [][]string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Test-Type", "Date", "Tests Added", "Tests Changed", "Tests Deleted", "commit-ID"}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, row := range csvRows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
