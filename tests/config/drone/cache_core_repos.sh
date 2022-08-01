#!/bin/bash
source .drone.env

echo "Checking core version - $CORE_COMMITID in cache"

URL="https://cache.owncloud.com/owncloud/ocis/oc10-test-runner/$CORE_COMMITID"

if curl --output /dev/null --silent --head --fail "$URL"; then
	echo "Web cache for ${CORE_COMMITID} already available in cache"
	exit 0
else
	# cache using the minio/mc client to the 'owncloud' bucket (long term bucket)
	mc alias set s3 "${MC_HOST}" "${AWS_ACCESS_KEY_ID}" "${AWS_SECRET_ACCESS_KEY}"
	mc mirror --overwrite --remove --debug oc10TestRunner "s3/owncloud/ocis/oc10-test-runner/""${CORE_COMMITID}"
	mc mirror --overwrite --remove --debug testingApp s3/owncloud/ocis/oc10-testing-app
fi
