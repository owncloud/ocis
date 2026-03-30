#!/bin/bash

# this script can run from everywhere in the ocis repo because it uses git grep

# The following grep will filter out every line containing an `env` annotation AND
# containing a non semver `introductionVersion` annotation that is used as placeholder.
# note that invalid `introductionVersion` annotations are already covered via:
# /.make/check-env-var-annotations.sh as part of the CI

RED=$(echo -e "\033[0;31m")
GREEN=$(echo -e "\033[0;32m")
NORM=$(echo -e "\033[0m")


# build the correct regex
IV_REGEX='(?:introductionVersion:\")(?!pre5\.0)(?:[^\d].*?\")'
IV_NAME=introductionVersion

RV_REGEX='(?:removalVersion:\")(?!pre5\.0)(?:[^\d].*?\")'
RV_NAME=removalVersion

EXCLUDE_PATHS='_test.go|vendor/'

# create a here doc function to be printed when the option is selected
# note that tabs are used intentionally. they are removed by cat but are required to make the code readable. 
print_version_examples() {
	cat <<-EOL
		  ${GREEN}Valid $1 examples:${NORM}

		  $1:"releaseX"  # acceptable alphabetical version
		  $1:"Addams"    # another alphabetical example such as a release name
		  $1:"Addams.8"  # another alphabetical example such as a release name plus a version
		  $1:"%%NEXT%%"  # a dummy placeholder as release name

		  See the dev docs for more details: https://owncloud.dev/services/general-info/envvars/envvar-naming-scopes/
	EOL
}

# ask what the script should search for
echo "Select one of the following keys to search for:"
echo
echo "  1) introductionVersion"
echo "  2) removalVersion"
echo "  3) Print valid matching search pattern examples"
echo
read -s -n 1 n
case $n in
  1) USEREGEX="${IV_REGEX}"
     NAME="${IV_NAME}"
     ;;
  2) USEREGEX="${RV_REGEX}"
     NAME="${RV_NAME}"
     ;;
  3) echo
     print_version_examples "${IV_NAME}"
     echo
     print_version_examples "${RV_NAME}"
     echo
     exit
     ;;
  *) echo "Invalid option"
     exit
     ;;
esac

# query the code
QS=$(git grep -n "env:" -- '*.go' | grep -v -E "${EXCLUDE_PATHS}" | grep --color=always -P "${USEREGEX}")

# count the results found
RESULTS_COUNT=$(echo "${QS}"|wc -l)

# add a new line after each hit, eol identified via "´ which only appears at the end of each envvar string definition
QUERY_RESULT=$(echo "${QS}" | sed "s#\"\`#\"\`\n#g")

echo "Checking ${NAME} annotations"

if [ "${QUERY_RESULT}" != "" ] && [ "${RESULTS_COUNT}" -gt 0 ]; then
  echo
  echo "==============================================================================================="
  echo ${GREEN}"The following items contain placeholder ${NAME} annotation(s):"${NORM}
  echo "==============================================================================================="
  echo
  echo "$QUERY_RESULT"
  echo
  echo ${GREEN}"${RESULTS_COUNT} items found"${NORM}
  echo
else
  echo ${GREEN}"No matching placeholders for ${NAME} annotations found."${NORM}
  echo
fi
