#!/usr/bin/env bash

SCRIPT_PATH=$(dirname "$0")
SCRIPT_PATH=$(cd "${SCRIPT_PATH}" && pwd) # absolute path
PROJECT_ROOT=$(cd "$SCRIPT_PATH/../../" && pwd)
SCRIPT_PATH_REL=${SCRIPT_PATH//"$PROJECT_ROOT/"/}

FILTER_SUITES=""
EXCLUDE_SUITES=""
FEATURE_PATHS=""
GLOB_FEATURE_PATHS=""
FEATURE_PATHS_FROM_ARG=""

SKIP_RUN_PARTS=true
RUN_PART=""
TOTAL_PARTS=""

HELP_COMMAND="
COMMAND [options] [paths]

Available options:
    --suites        - suites to run. Comma separated values (folder names)
                      e.g.: --suites smoke,shares
    --xsuites       - exclude suites from running. Comma separated values
                      e.g.: --xsuites spaces,search
    --run-part      - part to run out of total parts (groups)
                      e.g.: --run-part 2 (runs part 2 out of 4)
    --total-parts   - total number of groups to divide into
                      e.g.: --total-parts 4 (suites will be divided into 4 groups)
    --type          - type of tests to run. (Default: 'playwright' tests)
                      e.g.: --type playwright to run Playwright specs.
    --help, -h      - show cli options

Available env variables:
    TEST_SUITES     - Comma separated list of suites to run. (Will be ignored if --suites is provided)
    FEATURE_FILES   - Comma separated list of feature files to run. (Will be ignored if feature paths are provided)
    TEST_TYPE       - Type of tests to run. (Default: 'playwright' tests)
                      e.g.: TEST_TYPE='playwright' to run Playwright specs.
"

function log() {
    case $1 in
    info)
        echo -e "\e[0mINF: $2\e[0m"
        ;;
    error)
        echo -e "\e[31mERR: $2\e[0m"
        ;;
    warn)
        echo -e "\e[93mWRN: $2\e[0m"USAGE:
        ;;
    cmd)
        echo -e "\e[96mUSAGE: $2\e[0m"
        ;;
    *)
        echo -e "\e[0m$1\e[0m"
        ;;
    esac
}

while [[ $# -gt 0 ]]; do
    key="$1"
    case ${key} in
    --type)
        TEST_TYPE="$2"
        shift 2
        ;;
    --suites)
        FILTER_SUITES=$(echo "$2" | sed -E "s/,/\n/g")
        shift 2
        ;;
    --xsuites)
        EXCLUDE_SUITES=$(echo "$2" | sed -E "s/,/\n/g")
        shift 2
        ;;
    --run-part)
        SKIP_RUN_PARTS=false
        RUN_PART=$2
        shift 2
        ;;
    --total-parts)
        SKIP_RUN_PARTS=false
        TOTAL_PARTS=$2
        shift 2
        ;;
    --help | -h)
        log "$HELP_COMMAND"
        exit 0
        ;;
    *)
        if [[ $1 =~ ^-.* ]]; then
            log error "Unknown option: '$1'"
            log "$HELP_COMMAND"
            exit 1
        fi
        FEATURE_PATHS_FROM_ARG+=" $1" # maintain the white space
        shift
        ;;
    esac
done

# defaullt browser: chromium
if [[ -z $BROWSER ]]; then
    BROWSER="chromium"
fi

FEATURES_DIR="${SCRIPT_PATH}/../e2e/specs"
FEATURES_DIR=$(cd "$FEATURES_DIR" && pwd) # get absolute path
E2E_COMMAND="pnpm test:e2e:playwright --project=$BROWSER" # run command defined in package.json
ALL_SUITES=$(find "${FEATURES_DIR}"/ -type d | sort | rev | cut -d"/" -f1 | rev | grep -v '^[[:space:]]*$')

function getFeaturePaths() {
    local paths
    local real_paths=""
    local file_path
    local line_number
    local check_path
    local runner_path
    local spec_path

    # $1    - paths to suite or feature file
    paths=$(echo "$1" | xargs)
    for path in $paths; do
        file_path="${path%%:*}"
        line_number="${path#*:}"
        [[ "$file_path" == "$line_number" ]] && line_number=""

        if [[ "$TEST_TYPE" == "playwright" ]]; then
            spec_path=$(cd "$SCRIPT_PATH/../e2e" && pwd)

            if [[ -f "$file_path" || -d "$file_path" ]]; then
                runner_path="$file_path"
            elif [[ -f "$PROJECT_ROOT/$file_path" || -d "$PROJECT_ROOT/$file_path" ]]; then
                runner_path="$PROJECT_ROOT/$file_path"
            elif [[ -f "$spec_path/$file_path" || -d "$spec_path/$file_path" ]]; then
                runner_path="$spec_path/$file_path"
            elif [[ -f "$spec_path/specs/$file_path" || -d "$spec_path/specs/$file_path" ]]; then
                runner_path="$spec_path/specs/$file_path"
            else
                log error "File or folder doesn't exist: '$file_path'"
                log info "Tried: '$file_path', '$PROJECT_ROOT/$file_path', '$spec_path/$file_path', and '$spec_path/specs/$file_path'"
                exit 1
            fi

            [[ "$runner_path" == "$PROJECT_ROOT/"* ]] && runner_path="${runner_path#"$PROJECT_ROOT/"}"
            # Playwright testDir is tests/e2e/specs.
            runner_path="${runner_path#tests/e2e/specs/}"
            runner_path="${runner_path#tests/e2e/}"
            runner_path="${runner_path#specs/}"
        else
            check_path=$(echo "./$file_path" | cut -d ":" -f1)
            if [[ ! -f "$check_path" && ! -d "$check_path" ]]; then
                log error "File or folder doesn't exist: '$check_path'"
                log info "Path must be relative to '$SCRIPT_PATH_REL'"
                exit 1
            fi
            runner_path="$SCRIPT_PATH/$file_path"
        fi

        if [[ -n "$line_number" && "$line_number" =~ ^[0-9]+$ ]]; then
            real_paths+=" $runner_path:$line_number"
        else
            real_paths+=" $runner_path"
        fi
    done
    FEATURE_PATHS=$(echo "$real_paths" | xargs) # remove trailing white spaces
}

function runE2E() {
    if [[ ! -d "$PROJECT_ROOT" ]]; then
        log error "Project root doesn't exist: '$PROJECT_ROOT'"
    fi
    cd "$PROJECT_ROOT" || exit 1
    if [[ -n $GLOB_FEATURE_PATHS ]]; then
        $E2E_COMMAND "$GLOB_FEATURE_PATHS" # run without expanding glob pattern
    else
        # shellcheck disable=SC2086
        $E2E_COMMAND $FEATURE_PATHS # do not enclose paths with quote
    fi
    exit $?
}

function checkSuites() {
    # $1    - suites (separated by space or newline)
    for e_suite in $1; do
        exists=false
        for a_suite in $ALL_SUITES; do
            if [[ "$e_suite" == "$a_suite" ]]; then
                exists=true
            fi
        done
        if [[ "$exists" == false ]]; then
            log error "Suite doesn't exist: '$e_suite'"
            exit 1
        fi
    done
}

function buildSuitesPattern() {
    local delimiter=","
    local bracket_open="{"
    local bracket_close="}"

    CURRENT_SUITES_COUNT=$(echo "$1" | wc -w) # count words

    if [[ "$TEST_TYPE" == "playwright" ]]; then
        delimiter="|"
        bracket_open="("
        bracket_close=")"
    fi

    suites=$(echo "$1" | xargs | sed -E "s/( )+/$delimiter/g")
    if [[ $CURRENT_SUITES_COUNT -gt 1 ]]; then
        suites="$bracket_open${suites}$bracket_close"
    fi
    GLOB_FEATURE_PATHS="$FEATURES_DIR/$suites"
}

if [[ -n $TEST_SUITES ]] && [[ -z "$FILTER_SUITES" ]]; then
    FILTER_SUITES=$(echo "$TEST_SUITES" | sed -E "s/,/\n/g")
fi
if [[ -n $FEATURE_FILES ]] && [[ -z "$FEATURE_PATHS_FROM_ARG" ]]; then
    FEATURE_PATHS_FROM_ARG=$(echo "$FEATURE_FILES" | sed -E "s/,/ /g")
fi

# 1. [RUN E2E] run features from provided paths
if [[ -n $FEATURE_PATHS_FROM_ARG && "$SKIP_RUN_PARTS" == true ]]; then
    getFeaturePaths "$FEATURE_PATHS_FROM_ARG"
    log info "Running e2e using paths. All cli options will be discarded"
    runE2E
fi

# check if suites exist
if [[ -n $FILTER_SUITES ]]; then
    checkSuites "$FILTER_SUITES"
    ALL_SUITES=$FILTER_SUITES
fi
if [[ -n $EXCLUDE_SUITES ]]; then
    checkSuites "$EXCLUDE_SUITES"
fi

# exclude suites from running
if [[ -n $EXCLUDE_SUITES ]]; then
    for exclude_suite in $EXCLUDE_SUITES; do
        ALL_SUITES=$(echo "${ALL_SUITES/$exclude_suite/}" | sed -E "/^( )*$/d") # remove suite and trim empty lines
    done
fi

if [[ "$SKIP_RUN_PARTS" != true ]]; then
    if [[ -z $RUN_PART ]]; then
        log error "Missing '--run-part'"
        log cmd "--run-part <number>"
        exit 1
    fi
    if [[ -z $TOTAL_PARTS ]]; then
        log error "Missing '--total-parts'"
        log cmd "--total-parts <number>"
        exit 1
    fi

    ALL_SUITES_COUNT=$(echo "${ALL_SUITES}" | wc -l)
    SUITES_PER_RUN=$((ALL_SUITES_COUNT / TOTAL_PARTS))
    REMAINING_SUITES=$((ALL_SUITES_COUNT - (TOTAL_PARTS * SUITES_PER_RUN)))

    if [[ ${RUN_PART} -le ${REMAINING_SUITES} ]]; then
        SUITES_PER_RUN=$((SUITES_PER_RUN + 1))
        PREVIOUS_SUITES_COUNT=$(((RUN_PART - 1) * SUITES_PER_RUN))
    else
        PREV_MAX_SUITES=$((REMAINING_SUITES * (SUITES_PER_RUN + 1)))
        PREV_MIN_SUITES=$((((RUN_PART - 1) - REMAINING_SUITES) * SUITES_PER_RUN))
        PREVIOUS_SUITES_COUNT=$((PREV_MAX_SUITES + PREV_MIN_SUITES))
    fi

    GRAB_SUITES_UPTO=$((PREVIOUS_SUITES_COUNT + SUITES_PER_RUN))
    ALL_SUITES=$(echo "${ALL_SUITES}" | head -n "$GRAB_SUITES_UPTO" | tail -n "$SUITES_PER_RUN")
fi

buildSuitesPattern "$ALL_SUITES"
# 2. [RUN E2E] run the suites
runE2E
