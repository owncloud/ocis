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

set -evax

export SERVER_HOST=${SERVER_HOST:-https://localhost:9200}
export BACKEND_HOST=${BACKEND_HOST:-https://localhost:9200}
export TEST_TAGS=${TEST_TAGS:-"not @skip"}

cd ${WEB_PATH}/tests/acceptance/
yarn test:acceptance:external -- --require /drone/src/settings/ui/tests/acceptance/stepDefinitions ${1}

status=$?
exit $status
