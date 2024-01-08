#!/bin/sh
set -e

BASE_URL="$1"

apk add curl

#TODO: app driver itself should try again until OnlyOffice/Collabora is up...

retries=10
while [[ $retries -gt 0 ]]; do
    if curl --insecure --silent --show-error --fail "$BASE_URL"/hosting/discovery > /dev/null; then
        ocis app-provider server
    else
        echo "Office app is not yet available, trying again in 10 seconds"
        sleep 10
        retries=$((retries - 1))
    fi
done
echo 'Office app was not available after 100 seconds'
exit 1
