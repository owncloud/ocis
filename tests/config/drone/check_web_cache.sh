#!/bin/bash
source .drone.env

echo "Checking web version - $WEB_COMMITID in cache"

URL="$CACHE_ENDPOINT/$CACHE_BUCKET/ocis/web-test-runner/$WEB_COMMITID/README.md"

echo "Checking for the web cache at '$URL'."

if curl --output /dev/null --silent --head --fail "$URL"
then
	echo "Web with commit id $WEB_COMMITID already available in cache"
	# https://discourse.drone.io/t/how-to-exit-a-pipeline-early-without-failing/3951
	# exit a Pipeline early without failing
	exit 78
else
	echo "Web with $WEB_COMMITID not available in cache"
fi
