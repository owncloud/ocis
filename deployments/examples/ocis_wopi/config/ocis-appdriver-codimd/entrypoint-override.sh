#!/bin/sh
set -e

apk add curl

#TODO: app driver itself should try again until CodiMD is up...

retries=10
while [[ $retries -gt 0 ]]; do
    if curl --silent --show-error --fail http://codimd:3000 > /dev/null; then
        ocis storage-app-provider server
    else
        echo "CodiMD is not yet available, trying again in 10 seconds"
        sleep 10
        retries=$((retries - 1))
    fi
done
echo 'CodiMD was not available after 100 seconds'
exit 1
