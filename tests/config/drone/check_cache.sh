#!/bin/bash

# if no $1 is supplied end the script
# Can be web, acceptance or e2e
if [ -z "$1" ]; then
	echo "No cache item is supplied."
	exit 1
fi

URL="$CACHE_ENDPOINT/$CACHE_BUCKET/ocis/$1"

echo "Checking cache at '$URL'"

if curl --output /dev/null --silent --head --fail "$URL"; then
	echo "'$1' cache item already available."
	# https://discourse.drone.io/t/how-to-exit-a-pipeline-early-without-failing/3951
	# exit a Pipeline early without failing
	exit 78
else
	echo "'$1' cache item was not available."
fi
