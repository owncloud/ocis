#!/bin/sh
set -e

mkdir -p /var/tmp/ocis/.config/
cp /config/proxy-config.dist.json /var/tmp/ocis/.config/proxy-config.json
# TODO: remove replace logic when log level configuration is fixed
sed -i 's/PROXY_LOG_LEVEL/${PROXY_LOG_LEVEL}/g' /var/tmp/ocis/.config/proxy-config.json

ocis server &
sleep 10

# idp, glauth and accounts are not needed -> replaced by Keycloak and OpenLDAP
ocis kill idp
ocis kill glauth
ocis kill accounts

# workaround for loading proxy configuration
ocis kill proxy
sleep 10
ocis proxy server &

wait
