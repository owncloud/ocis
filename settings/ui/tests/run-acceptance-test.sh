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
	export TEST_INFRA_DIRECTORY=$testFolder/tests
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

cp -r "$WEB_PATH"/tests/acceptance/stepDefinitions "$testFolder"
cp "$WEB_PATH"/tests/acceptance/setup.js "$testFolder"

export SERVER_HOST=${SERVER_HOST:-https://localhost:9200}
export BACKEND_HOST=${BACKEND_HOST:-https://localhost:9200}
export TEST_TAGS=${TEST_TAGS:-"not @skip"}

yarn run acceptance-tests "$1"

status=$?
exit $status
