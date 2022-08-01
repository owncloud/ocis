#!/bin/bash
source .drone.env

echo "Checking web version - $WEB_COMMITID in cache"

URL="https://cache.owncloud.com/owncloud/ocis/web-test-runner/$WEB_COMMITID"

if curl --output /dev/null --silent --head --fail "$URL"; then
	echo "web cache for $WEB_COMMITID already available in cache"
	exit 0
else
	# cache using the minio/mc client to the 'owncloud' bucket (long term bucket)
	mc alias set s3 "${MC_HOST}" "${AWS_ACCESS_KEY_ID}" "${AWS_SECRET_ACCESS_KEY}"
	mc mirror --overwrite --remove --debug webTestRunner "s3/owncloud/ocis/web-test-runner/""${WEB_COMMITID}"
fi
