#!/bin/bash
source .drone.env

echo "Checking web version - $WEB_COMMITID in cache"

URL="https://cache.owncloud.com/owncloud/ocis/web-test-runner/$WEB_COMMITID"

if curl --output /dev/null --silent --head --fail "$URL"; then
	echo "web cache for $WEB_COMMITID already available in cache"
	exit 0
else
	echo "Cache for $WEB_COMMITID not available in cache, cloning..."
	# clone the "owncloud/web" repository
	git clone -b "${WEB_BRANCH}" --single-branch --no-tags https://github.com/owncloud/web.git webTestRunner
	cd webTestRunner && git checkout "${WEB_COMMITID}"

	ls -la .
fi
