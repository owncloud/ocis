#!/bin/bash
source .drone.env

echo "Checking core version - $CORE_COMMITID in cache"

URL="https://cache.owncloud.com/owncloud/ocis/oc10-test-runner/$CORE_COMMITID"

if curl --output /dev/null --silent --head --fail "$URL"; then
	echo "Core cache for $CORE_COMMITID already available in cache"
	exit 0
else
	echo "Cache for $CORE_COMMITID not available in cache, cloning..."

	# clone the core repository
	git clone -b "${CORE_BRANCH}" --single-branch --no-tags https://github.com/owncloud/core.git oc10TestRunner

	cd oc10TestRunner && git checkout "${CORE_COMMITID}"
	ls -la .
	cd ..

	# clone the testing app repository
	git clone -b master --single-branch --no-tags https://github.com/owncloud/testing.git testingApp

	ls -la testingApp
fi
