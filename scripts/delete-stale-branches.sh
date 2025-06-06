#!/bin/bash
set -euo pipefail

# Configurable parameters (set via env or script args)
DRY_RUN="${DRY_RUN:-true}"
VERBOSE="${VERBOSE:-false}"

DAYS_BEFORE_STALE="${DAYS_BEFORE_STALE:-90}"
IGNORE_BRANCHES_REGEX="^(main$|master$|stable-|release-|docs$|docs-stable-)"
MAX_BRANCHES_PER_RUN="${MAX_BRANCHES_PER_RUN:-5}"

# Global variables
branch_info_file=""

# Colors and styles
readonly B="\033[1m"  # Bold  
readonly Y="\033[33m" # Yellow
readonly G="\033[90m" # Gray
readonly R="\033[0m"  # Reset

# Script metadata
readonly SCRIPT_NAME="${0##*/}"
IFS=$'\n\t'

# Helper functions
log_error() {
    echo -e "${R}ERROR: $*${R}" >&2
}

log_header() {
    echo
    echo -e "${B}$*${R}"
}

log_info() {  
    echo -e "$*$"
}

# Helper: get ISO8601 date N days ago (cross-platform)
get_iso_date_n_days_ago() {
    local days_ago="$1"
    if ! [[ "$days_ago" =~ ^[0-9]+$ ]]; then
        log_error "Days must be a positive integer"
        return 1
    fi

    if date --version >/dev/null 2>&1; then
        # GNU date (Linux)
        date -d "$days_ago days ago" --iso-8601=seconds
    else
        # BSD date (macOS)
        date -v-"$days_ago"d +"%Y-%m-%dT%H:%M:%S%z"
    fi
}

# Function: Convert ISO 8601 date to epoch seconds (portable: Linux/macOS)
iso_to_epoch() {
    local iso="$1"
    if date --version >/dev/null 2>&1; then
        # GNU date (Linux)
        date -d "$iso" "+%s"
    else
        # BSD date (macOS)
        local fixed
        if [[ "$iso" =~ Z$ ]]; then
            # Replace Z with +0000
            fixed="${iso/Z/+0000}"
        else
            # Remove colon in timezone for BSD date compatibility
            fixed=$(echo "$iso" | sed -E 's/([0-9]{2}):([0-9]{2})$/\1\2/')
        fi
        date -j -f "%Y-%m-%dT%H:%M:%S%z" "$fixed" "+%s"
    fi
}

# Usage: get days difference from now
get_days_diff_from_now() {
    local iso_date="$1"
    local now_epoch
    now_epoch=$(date "+%s")
    local date_epoch
    date_epoch=$(iso_to_epoch "$iso_date")
    echo $(( (now_epoch - date_epoch) / 86400 ))
}

validate_requirements() {
    local missing_deps=()
    
    # Check for required commands
    for cmd in git date; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            missing_deps+=("$cmd")
        fi
    done

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        exit 1
    fi

    # Validate numeric parameters
    if ! [[ "$DAYS_BEFORE_STALE" =~ ^[0-9]+$ ]]; then
        log_error "DAYS_BEFORE_STALE must be a positive integer"
        exit 1
    fi

    if ! [[ "$MAX_BRANCHES_PER_RUN" =~ ^[0-9]+$ ]]; then
        log_error "MAX_BRANCHES_PER_RUN must be a positive integer"
        exit 1
    fi
}

# Parse repository information from git remote URL
parse_repo_info() {
    local remote_url
    remote_url=$(git remote get-url origin)
    
    if [[ $remote_url =~ github\.com[:/][^/]+/[^/]+\.git$ ]]; then
        # Extract owner and repo using parameter expansion for robustness
        # Remove protocol and .git suffix
        local url_no_proto=${remote_url#*github.com[:/]}
        local url_no_git=${url_no_proto%.git}
        OWNER=$(echo "$url_no_git" | cut -d'/' -f1)
        REPO=$(echo "$url_no_git" | cut -d'/' -f2)
        
        if [[ -z "$OWNER" || -z "$REPO" ]]; then
            log_error "OWNER or REPO is empty after parsing"
            return 1
        fi
    else
        log_error "Could not parse owner/repo from remote URL: $remote_url"
        return 1
    fi
}

main() {
    # Initialize temporary file
    branch_info_file=$(mktemp)
    trap 'rm -f "$branch_info_file"' EXIT

    # Ensure we have all git information
    git fetch --prune origin || true

    # Validate script requirements
    validate_requirements

    # Parse repository information
    parse_repo_info || exit 1

    # Print header
    log_header "Repository: $OWNER/$REPO"
    log_header "Configuration:"
    log_info "  Today: $(date)"
    log_info "  Days before stale: $DAYS_BEFORE_STALE"
    log_info "  Max branches per run: $MAX_BRANCHES_PER_RUN"
    log_info "  Dry run: $DRY_RUN"

    # Get all branches
    local branches_all branches_all_count branches_ignored branches_candidates
    branches_all=$(git branch -r | grep -v '\->' | sed 's/origin\///' | xargs -I{} echo "{}")
    branches_all_count=$(echo "$branches_all" | wc -l)
    branches_ignored=$(echo "$branches_all" | grep -E "$IGNORE_BRANCHES_REGEX" || true)

    log_header "Branch Filtering:"
    log_info "  IGNORE_BRANCHES_REGEX: ${Y}${IGNORE_BRANCHES_REGEX}${R}"
    log_info "  Ignored branches:"
    echo "$branches_ignored"

    branches_candidates=$(echo "$branches_all" | grep -Ev "$IGNORE_BRANCHES_REGEX" || true)

    # Collect branch info
    for branch in $branches_candidates; do
        local last_commit_date last_commit_author last_commit_sha age_days
        last_commit_date=$(git log -1 --format=%aI "origin/$branch")
        last_commit_author=$(git log -1 --format=%an "origin/$branch")
        last_commit_sha=$(git log -1 --format=%H "origin/$branch")
        age_days=$(get_days_diff_from_now "$last_commit_date")
        
        if [[ "$age_days" -lt "$DAYS_BEFORE_STALE" ]]; then
            continue
        fi
        printf "%d\t%s\t%s\t%s\t%s\n" "$age_days" "$last_commit_sha" "$last_commit_date" "$branch" "$last_commit_author" >> "$branch_info_file"
    done

    # Sort and get oldest branches
    if [[ "$VERBOSE" == "true" ]]; then
        log_header "Stale branches by age:"
        sort -nr "$branch_info_file" | while IFS=$'\t' read -r age_days sha date branch author; do
            printf "  ${G}%-4d days${R}  %s  ${G}%s${R}  by %s\n" "$age_days" "$branch" "$date" "$author"
        done
    fi

    local oldest_branches
    oldest_branches=$(sort -nr "$branch_info_file" | head -n "$MAX_BRANCHES_PER_RUN")
    if [[ -z "$oldest_branches" ]]; then
        log_header "No branches to delete."
        exit 0
    fi

    log_header "Deleting oldest branches (max $MAX_BRANCHES_PER_RUN):"
    while IFS=$'\t' read -r age_days sha date branch author; do
        echo -e "  $branch (${G}$age_days days old${R})  by $author"
        if [[ "$DRY_RUN" == "false" ]]; then
            git push origin --delete "${branch}" >/dev/null 2>&1
        fi
    done <<< "$oldest_branches"

    # Print summary
    local branches_remaining branches_remaining_count deleted_count
    branches_remaining=$(git branch -r | grep -v '\->' | sed 's/origin\///' | xargs -I{} echo "{}")
    branches_remaining_count=$(echo "$branches_remaining" | wc -l)
    deleted_count=$((branches_all_count - branches_remaining_count))

    log_header "Summary:"
    if [[ "$DRY_RUN" != "false" ]]; then
        log_info "  DRY_RUN: No branches deleted"
    fi
    log_info "  Branches before: $branches_all_count"
    log_info "  Branches after:  $branches_remaining_count"
    log_info "  Deleted:        $deleted_count"
    echo
}

main
