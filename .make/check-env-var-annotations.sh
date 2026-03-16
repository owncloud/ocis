#!/bin/bash

# The following grep will filter out every line containing an `env` annotation.
# It will ignore every line that already has a valid `introductionVersion` annotation.

RED=$(echo -e "\033[0;31m")
GREEN=$(echo -e "\033[0;32m")
NORM=$(echo -e "\033[0m")

# create a here doc function to be printed in case of introductionVersion annotation errors
# note that tabs are used intentionally. they are removed by cat but are required to make the code readable. 
print_introduction_version_examples() {
	cat <<-EOL
		  ${GREEN}Valid examples:${NORM}

		  introductionVersion:"pre5.0"
		  introductionVersion:"5.0"
		  introductionVersion:"4.9.3-rc5"
		  introductionVersion:"5.0.1-cheesecake"
		  introductionVersion:"5.10.100.15"
		  introductionVersion:"0.0"
		  introductionVersion:"releaseX"  # acceptable alphabetical version
		  introductionVersion:"Addams"    # another alphabetical example such as a release name

		  ${RED}Invalid examples:${NORM}

		  introductionVersion:"5.0cheesecake"
		  introductionVersion:"5"
		  introductionVersion:"5blueberry"
		  introductionVersion:"5-lasagna"
		  introductionVersion:"4.9.3rc5"

		  See the dev docs for more details: https://owncloud.dev/services/general-info/envvars/envvar-naming-scopes/
	EOL
}

ERROR=0

SEMVER_REGEX="([0-9]|[1-9][0-9]*)(\.([0-9]|[1-9][0-9]*)){1,2}(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?"
ALPHA_REGEX="[A-Za-z]+[A-Za-z0-9-]*"

QI=$(git grep -n "env:" -- '*.go' |grep -v -P "introductionVersion:\"($SEMVER_REGEX|(pre5\\.0)|($ALPHA_REGEX))\""|grep -v "_test.go"|grep -v "vendor/")

# add a new line after each hit, eol identified via "´ which only appears at the end of each envvar string definition
QUERY_INTRO=$(echo "${QI}" | sed "s#\"\`#\"\`\n#g")

RESULTS_INTRO=$(echo "${QUERY_INTRO}"|wc -l)

echo "Checking introductionVersion annotations"

if [ "${QUERY_INTRO}" != "" ] && [ "${RESULTS_INTRO}" -gt 0 ]; then
  echo
  echo "==============================================================================================="
  echo ${RED}"The following ${RESULTS_INTRO} items contain invalid or missing introductionVersion annotation(s):"${NORM}
  echo "==============================================================================================="
  echo
  echo "$QUERY_INTRO"
  echo
  print_introduction_version_examples
  echo
  ERROR=1
else
  echo "All introductionVersion annotations are valid"
  echo
fi

# The following grep will filter out every line containing an `env` annotation
# it will ignore every line that has allready a valid `desc` annotation

QUERY_DESC=$(git grep -n "env:" -- '*.go' |grep -v -P "desc:\".{10,}\""|grep -v "_test.go"|grep -v "vendor/")

RESULTS_DESC=$(echo "${QUERY_DESC}"|wc -l)

echo "Checking description annotations"

if [ "${QUERY_DESC}" != "" ] && [ "${RESULTS_DESC}" -gt 0 ]; then
  echo
  echo "==============================================================================================="
  echo ${RED}"The following ${RESULTS_DESC} items contain invalid or missing description annotation:"${NORM}
  echo "==============================================================================================="
  echo "$QUERY_DESC"
  ERROR=1
else
  echo "All description annotations are valid"
fi
exit ${ERROR}
