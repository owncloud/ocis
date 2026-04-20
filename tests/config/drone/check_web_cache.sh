#!/bin/bash
# DRONE: Checks S3 MinIO cache for web test runner artifacts (acceptance/e2e/web).
# Sources .drone.env to get WEB_COMMITID used as the cache key.
# Exit code 78 is the Drone "skip pipeline without failure" signal.
# Used by .drone.star build-web-cache pipeline step.
source .drone.env

# if no $1 is supplied end the script
# Can be web, acceptance or e2e
if [ -z "$1" ]; then
  echo "No cache item is supplied."
  exit 1
fi

echo "Checking web version - $WEB_COMMITID in cache"
web_cache=$(mc find s3/$CACHE_BUCKET/ocis/web-test-runner/$WEB_COMMITID/$1 2>&1 | grep 'Object does not exist')

if [[ -z "$web_cache" ]]; then
  echo "$1 cache with commit id $WEB_COMMITID already available."
  # https://discourse.drone.io/t/how-to-exit-a-pipeline-early-without-failing/3951
  # exit a Pipeline early without failing
  exit 78
else
  echo "$1 cache with commit id $WEB_COMMITID was not available."
fi
