#!/bin/bash

# the following grep will filter out every line containing an `env` annotation
# and prints all entries found that are invalid. this is required to see where CI might fail

RED=$(echo -e "\033[0;31m")
GREEN=$(echo -e "\033[0;32m")
NORM=$(echo -e "\033[0m")

# create a here doc function to be printed in case of introductionVersion annotation errors
# note that tabs are used intentionally. they are removed by cat but are required to make the code readable. 
print_introduction_version_examples() {
	cat <<-EOL
		  ${GREEN}Valid examples:${NORM}

		  introductionVersion:"pre5.0"
		  introductionVersion:"5.0" deprecationVersion:"7.0.0"
		  introductionVersion:"4.9.3-rc5"
		  introductionVersion:"5.10.100.15"
		  introductionVersion:"0.0"
		  introductionVersion:"release" deprecationVersion:"7.0.0"
		  introductionVersion:"releaseX" deprecationVersion:"7.0.0"
		  introductionVersion:"Addams" deprecationVersion:"7.0.0"

		  ${RED}Invalid examples:${NORM}

		  introductionVersion:""
		  introductionVersion:"  "
		  introductionVersion:"-"
		  introductionVersion:"--"
		  introductionVersion:"51.0cheesecake" deprecationVersion:"7.0.0"
		  introductionVersion:"50..0.1-cheesecake"
		  introductionVersion:"15"
		  introductionVersion:"54blueberry"
		  introductionVersion:"5-lasagna"
		  introductionVersion:"5.lasagna-rc1"
		  introductionVersion:"4.9.3rc-5" deprecationVersion:"7.0.0"
		  introductionVersion:"5B-rc1" deprecationVersion:"7.0.0"

		  See the envvar life cycle in the dev docs for more details:
		  https://owncloud.dev/services/general-info/envvars/envvar-naming-scopes/
	EOL
}

ERROR=0

# note that our semver is not fully compliant because we do not have a patch version for new envvars
# if you want to test, use the examples above and take the blocks separated by `|` (or)
SEMVER_REGEX='(?:introductionVersion:\")(?:(?:[^\dA-Za-z]*?)\"|(?:\d+[A-z]+).*?\"|(?:\d*?\")|(?:\d*?[-\.]+\D.*?\")|(?:\d*?[\.]{2,}.*?\")|(?:\d+\.\d+[A-z]+\")|(?:(?:\d+\.)+\d+[A-z]+.*?\"))'

# filter out the the following paths
EXCLUDE_PATHS='_test.go|vendor/'

# query the code
QI=$(git grep -n "env:" -- '*.go' | grep -v -E "${EXCLUDE_PATHS}" | grep --color=always -P "${SEMVER_REGEX}")

# count the results found
RESULTS_INTRO=$(echo "${QI}"|wc -l)

# add a new line after each hit, eol identified via "´ which only appears at the end of each envvar string definition
QUERY_INTRO=$(echo "${QI}" | sed "s#\"\`#\"\`\n#g")

echo "Checking introductionVersion annotations"

if [ "${QUERY_INTRO}" != "" ] && [ "${RESULTS_INTRO}" -gt 0 ]; then
  echo
  echo "==============================================================================================="
  echo ${RED}"The following item(s) contain invalid or missing introductionVersion annotation(s):"${NORM}
  echo "==============================================================================================="
  echo
  echo "$QUERY_INTRO"
  echo
  echo ${GREEN}"${RESULTS_INTRO} items found"${NORM}
  echo
  print_introduction_version_examples
  echo
  ERROR=1
else
  echo "All introductionVersion annotations are valid"
  echo
fi

# the following grep will filter out every line containing an `env` annotation and
# will print each line that has an invalid `desc` annotation.

# query the code
QUERY_DESC=$(git grep -n "env:" -- '*.go' | grep -v -E "${EXCLUDE_PATHS}" | grep -v --color=always  -P "desc:\".{10,}\"")

# count the results found
RESULTS_DESC=$(echo "${QUERY_DESC}"|wc -l)

echo "Checking description annotations"

if [ "${QUERY_DESC}" != "" ] && [ "${RESULTS_DESC}" -gt 0 ]; then
  echo
  echo "==============================================================================================="
  echo ${RED}"The following item(s) contain invalid or missing description annotation:"${NORM}
  echo "==============================================================================================="
  echo "$QUERY_DESC"
  echo
  echo ${GREEN}"${RESULTS_DESC} items found"${NORM}
  ERROR=1
else
  echo "All description annotations are valid"
fi
exit ${ERROR}
