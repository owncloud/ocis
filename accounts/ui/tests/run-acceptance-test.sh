#!/bin/bash

if [ -z "$WEB_PATH" ]
then
	echo "WEB_PATH env variable is not set, cannot find files for tests infrastructure"
	exit 1
fi

if [ -z "$WEB_UI_CONFIG" ]
then
	echo "WEB_UI_CONFIG env variable is not set, cannot find web config file"
	exit 1
fi

if [ -z "$1" ]
then
	echo "Features path not given, exiting test run"
	exit 1
fi

trap clean_up SIGHUP SIGINT SIGTERM

if [ -z "$TEST_INFRA_DIRECTORY" ]
then
	cleanup=true
	testFolder=$(mktemp -d -p .)
	printf "creating folder $testFolder for Test infrastructure setup\n\n"
	export TEST_INFRA_DIRECTORY=$(realpath $testFolder)
fi

clean_up() {
	if $cleanup
	then
		if [ -d "$testFolder" ]; then
			printf "\n\n\n\nDeleting folder $testFolder Test infrastructure setup..."
			rm -rf "$testFolder"
		fi
	fi
}

trap clean_up SIGHUP SIGINT SIGTERM EXIT

cp -r $(ls -d "$WEB_PATH"/tests/acceptance/* | grep -v 'node_modules') "$testFolder"

export SERVER_HOST=${SERVER_HOST:-https://localhost:9200}
export BACKEND_HOST=${BACKEND_HOST:-https://localhost:9200}
export TEST_TAGS=${TEST_TAGS:-"not @skip"}

cucumber-js --retry 1 \
						--require-module @babel/register \
						--require-module @babel/polyfill \
						--require ${TEST_INFRA_DIRECTORY}/setup.js \
						--require ui/tests/acceptance/stepDefinitions \
						--require ${TEST_INFRA_DIRECTORY}/stepDefinitions \
						--format @cucumber/pretty-formatter \
						-t ${TEST_TAGS:-not @skip and not @skipOnOC10}

status=$?
exit $status
