package droneextractor

import (
	"fmt"
	"io"
	"time"

	"github.com/owncloud/ocis/v2/tools/ci-reporter/pkg/githubextractor"
)

const requestDelay = 1 * time.Second

// FetchPRBuilds fetches all Drone builds for a PR by querying:
// 1. Builds by PR branch (handles force-pushes)
// 2. Builds by individual commit SHAs (fallback for merged/deleted branches)
//
// Returns a map of commit SHA -> builds
func FetchPRBuilds(
	extractor *Extractor,
	baseURL string,
	repo string,
	prBranch string,
	commits []githubextractor.PRCommit,
	logWriter io.Writer,
) (map[string][]BuildInfo, error) {
	buildsByCommit := make(map[string][]BuildInfo)

	// Query builds by PR branch (handles force-pushes)
	if logWriter != nil {
		fmt.Fprintf(logWriter, "Fetching Drone builds for PR (by event and branch)...\n")
	}
	prBuilds, err := extractor.GetBuildsByPRBranch(baseURL, repo, prBranch)
	if err != nil {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "⚠  Failed to get builds by branch: %v\n", err)
		}
	} else {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "Found %d builds for PR branch %s\n", len(prBuilds), prBranch)
		}
		// Group builds by commit SHA
		for _, build := range prBuilds {
			buildsByCommit[build.CommitSHA] = append(buildsByCommit[build.CommitSHA], build)
		}
	}

	// Also query by current commit SHAs (for completeness)
	if logWriter != nil {
		fmt.Fprintf(logWriter, "Fetching Drone builds for current commits...\n")
	}
	for i, commit := range commits {
		if _, exists := buildsByCommit[commit.SHA]; exists {
			if logWriter != nil {
				fmt.Fprintf(logWriter, "  [%d/%d] %s: already found via PR query\n", i+1, len(commits), commit.SHA[:7])
			}
			continue
		}
		if logWriter != nil {
			fmt.Fprintf(logWriter, "  [%d/%d] Fetching builds for %s...\n", i+1, len(commits), commit.SHA[:7])
		}
		builds, err := extractor.GetBuildsByCommitSHA(baseURL, repo, commit.SHA)
		if err != nil {
			if logWriter != nil {
				fmt.Fprintf(logWriter, "    ⚠  Failed to get builds: %v\n", err)
			}
			continue
		}
		if len(builds) > 0 {
			buildsByCommit[commit.SHA] = append(buildsByCommit[commit.SHA], builds...)
			if logWriter != nil {
				fmt.Fprintf(logWriter, "    Found %d builds\n", len(builds))
			}
		}
		time.Sleep(requestDelay)
	}

	return buildsByCommit, nil
}
