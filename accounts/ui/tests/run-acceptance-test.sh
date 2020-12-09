#!/bin/bash

if [ -z "$WEB_PATH" ]
then
	echo "WEB_PATH env variable is not set, cannot find files for tests infrastructure"
	exit 1
fi

if [ -z "$OCIS_SKELETON_DIR" ]
then
	echo "OCIS_SKELETON_DIR env variable is not set, cannot find skeleton directory"
	exit 1
fi

if [ -z "$WEB_CONFIG" ]
then
	echo "WEB_CONFIG env variable is not set, cannot find web config file"
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
	testFolder=$(< /dev/urandom tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
	printf "creating folder $testFolder for Test infrastructure setup\n\n"
	export TEST_INFRA_DIRECTORY=$testFolder
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

cp -r "$WEB_PATH"/tests ./"$testFolder"

export SERVER_HOST=${SERVER_HOST:-https://localhost:9200}
export BACKEND_HOST=${BACKEND_HOST:-https://localhost:9200}
export RUN_ON_OCIS='true'
export TEST_TAGS=${TEST_TAGS:-"not @skip"}

yarn run acceptance-tests "$1"

status=$?
exit $status
