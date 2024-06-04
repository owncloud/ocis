#!/bin/bash

# The following grep will filter out every line containing an `env` annotation
# it will ignore every line that already has a valid `introductionVersion` annotation
#
# valid examples:
#
# introductionVersion:"%%NEXT%%"
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

QUERY_INTRO=$(git grep -n "env:" -- '*.go' |grep -v -P "introductionVersion:\"($SEMVER_REGEX|(pre5\.0)|(%%NEXT%%))\""|grep -v "_test.go"|grep -v "vendor/")
RESULTS_INTRO=$(echo "${QUERY_INTRO}"|wc -l)
if [ "${QUERY_INTRO}" != "" ] && [ "${RESULTS_INTRO}" -gt 0 ]; then
  echo "==============================================================================================="
  echo "The following ${RESULTS_INTRO} items contain an invalid or missing introductionVersion annotation:"
  echo "==============================================================================================="
  echo "$QUERY_INTRO"
  ERROR=1
else
  echo "All introductionVersion annotations are valid"
fi

# The following grep will filter out every line containing an `env` annotation
# it will ignore every line that has allready a valid `desc` annotation

QUERY_DESC=$(git grep -n "env:" -- '*.go' |grep -v -P "desc:\".{10,}\""|grep -v "_test.go"|grep -v "vendor/")

RESULTS_DESC=$(echo "${QUERY_DESC}"|wc -l)
if [ "${QUERY_DESC}" != "" ] && [ "${RESULTS_DESC}" -gt 0 ]; then
  echo "==============================================================================================="
  echo "The following ${RESULTS_DESC} items contain an invalid or missing description annotation:"
  echo "==============================================================================================="
  echo "$QUERY_DESC"
  ERROR=1
else
  echo "All description annotations are valid"
fi
exit ${ERROR}
