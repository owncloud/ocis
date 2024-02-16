#!/bin/bash

# The following grep will filter out every line containing an `env` annotation
# it will ignore every line that already has a valid `introductionVersion` annotation
#
# valid examples:
#
# introductionVersion:"pre5.0"
# introductionVersion:"5.0"
# introductionVersion:"4.9.3-rc5"
# introductionVersion:"5.0.1-cheesecake"
# introductionVersion:"5.10.100.15"
# introductionVersion:"0.0"
#
# invalid examples:
#
# introductionVersion:"5.0cheesecake"
# introductionVersion:"5"
# introductionVersion:"5blueberry"
# introductionVersion:"5-lasagna"
# introductionVersion:"4.9.3rc5"

ERROR=0

SEMVER_REGEX="([0-9]|[1-9][0-9]*)(\.([0-9]|[1-9][0-9]*)){1,2}(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?"

QUERY_INTRO=$(git grep "env:" -- '*.go' |grep -v -P "introductionVersion:\"($SEMVER_REGEX|(pre5\.0))\""|grep -v "_test.go"|grep -v "vendor/")

RESULTS_INTRO=$(echo "${QUERY_INTRO}"|wc -l)
if [ "${RESULTS_INTRO}" -gt 0 ]; then
  echo "==============================================================================================="
  echo "The following ${RESULTS_INTRO} files contain an invalid introductionVersion annotation:"
  echo "==============================================================================================="
  echo "$QUERY_INTRO"
  ERROR=1
fi

# The following grep will filter out every line containing an `env` annotation
# it will ignore every line that has allready a valid `desc` annotation

QUERY_DESC=$(git grep "env:" -- '*.go' |grep -v -P "desc:\".{10,}\""|grep -v "_test.go"|grep -v "vendor/")

RESULTS_DESC=$(echo "${QUERY_DESC}"|wc -l)
if [ "${RESULTS_DESC}" -gt 0 ]; then
  echo "==============================================================================================="
  echo "The following ${RESULTS_DESC} files contain an invalid description annotation:"
  echo "==============================================================================================="
  echo "$QUERY_DESC"
  ERROR=1
fi
exit ${ERROR}
