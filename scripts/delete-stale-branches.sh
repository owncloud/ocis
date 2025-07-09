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
    echo -e "$*"
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

# Get stale branches with age filtering (atomic git operation)
get_stale_branches() {
    local days_threshold="$1"
    local temp_file=$(mktemp)
    
    # Ensure temp file cleanup on function exit/interruption
    trap 'rm -f "$temp_file"' RETURN
    
    # Single atomic git command gets all branch info at once, write to file to avoid SIGPIPE
    # Reference: TabrisJS uses git for-each-ref -> file -> process for automation across repos    
    # https://tabris.com/iterate-over-branches-in-your-git-repository/
    git for-each-ref --format='%(refname:short)|%(objectname)|%(committerdate:iso8601-strict)|%(authorname)' refs/remotes/origin/ > "$temp_file"
    
    local stale_info=""
    while IFS='|' read -r refname sha date author; do
        local branch=${refname#origin/}
        [[ "$branch" == "HEAD" ]] && continue
        
        # Skip ignored branches
        [[ "$branch" =~ $IGNORE_BRANCHES_REGEX ]] && continue
        
        # Calculate age and filter immediately
        # Convert git date format "2025-05-20 09:48:29 +0200" to ISO8601 "2025-05-20T09:48:29+0200"
        local iso_date="${date/ /T}"      # Replace first space with T
        iso_date="${iso_date/ /}"         # Remove space before timezone
        local age_days=$(get_days_diff_from_now "$iso_date")
        if [[ "$age_days" -ge "$days_threshold" ]]; then
            stale_info="$stale_info$age_days"$'\t'"$sha"$'\t'"$date"$'\t'"$branch"$'\t'"$author"$'\n'
        fi
    done < "$temp_file"
    
    echo "${stale_info%$'\n'}"
}

validate_requirements() {
    local missing_deps=()
    
    # Check for required commands
    for cmd in git date curl; do
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

# Delete branch via GitHub REST API (uses GITHUB_TOKEN)
# Args: branch-name
# src: https://docs.github.com/en/rest/git/refs#delete-a-reference
# Same as gh CLI but without 40MB binary
github_api_delete_branch() {
    local branch="$1"

    if [[ -z "${GITHUB_TOKEN:-}" ]]; then
        log_error "GITHUB_TOKEN not set, cannot call GitHub API to delete $branch"
        return 1
    fi

    local api_url="https://api.github.com/repos/${OWNER}/${REPO}/git/refs/heads/${branch}"
    local curl_verbosity="-s"
    [[ "$VERBOSE" == "true" ]] && curl_verbosity="-v"

    rsp=$(curl $curl_verbosity -w '\n%{http_code}' \
        -X DELETE \
        -H "Authorization: token ${GITHUB_TOKEN}" \
        -H "Accept: application/vnd.github+json" \
        "$api_url")
    http_code=$(printf '%s\n' "$rsp" | tail -n1)
    [[ "$VERBOSE" == "true" ]] && echo "$rsp"

    if [[ "$http_code" == "204" ]]; then
        return 0
    else
        log_error "GitHub API returned HTTP $http_code while deleting $branch"
        return 1
    fi
}

main() {
    # Validate script requirements
    validate_requirements

    git config --global user.email "droneci@placeholder.com"
    git config --global user.name "Drone CI"

    # Ensure we have all git information
    git fetch --prune origin || true

    # Parse repository information
    parse_repo_info || exit 1

    # Print header
    log_header "Repository: $OWNER/$REPO"
    log_header "Configuration:"
    log_info "  Today: $(date)"
    log_info "  Days before stale: $DAYS_BEFORE_STALE"
    log_info "  Max branches per run: $MAX_BRANCHES_PER_RUN"
    log_info "  Dry run: $DRY_RUN"

    local stale_branches
    stale_branches=$(get_stale_branches "$DAYS_BEFORE_STALE")
    
    if [[ -z "$stale_branches" ]]; then
        log_header "No stale branches found."
        exit 0
    fi

    # Sort and show stale branches
    if [[ "$VERBOSE" == "true" ]]; then
        log_header "Stale branches by age:"
        echo "$stale_branches" | sort -nr | while IFS=$'\t' read -r age_days sha date branch author; do
            printf "  ${G}%-4d days${R}  %s  ${G}%s${R}  by %s\n" "$age_days" "$branch" "${date%%T*}" "$author"
        done
    fi

    local oldest_branches
    oldest_branches=$(echo "$stale_branches" | sort -nr | head -n "$MAX_BRANCHES_PER_RUN")

    log_header "Deleting oldest branches (max $MAX_BRANCHES_PER_RUN):"
    while IFS=$'\t' read -r age_days sha date branch author; do
        echo -e "  $branch (${G}$age_days days old${R})  by $author"
        if [[ "$DRY_RUN" == "false" ]]; then
            # Attempt git deletion, fall back to REST if that fails
            git push origin --delete "$branch" >/dev/null 2>&1 || \
            github_api_delete_branch "$branch" || \
            log_error "Failed to delete $branch with both git and REST"
        fi
    done <<< "$oldest_branches"

    # Print summary
    local stale_count deleted_count
    stale_count=$(printf "%s\n" "$stale_branches" | wc -l)
    deleted_count=$(printf "%s\n" "$oldest_branches" | wc -l)

    log_header "Summary:"
    if [[ "$DRY_RUN" != "false" ]]; then
        log_info "  DRY_RUN: No branches deleted"
    fi
    log_info "  Stale branches found: $stale_count"
    log_info "  Branches deleted: $deleted_count"
    echo
}

main
