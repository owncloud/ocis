#!/bin/sh

set -e

cp /config/proxy-config.dist.json /config/proxy-config.json
# TODO: remove replace logic when log level configuration is fixed
sed -i 's/PROXY_LOG_LEVEL/${PROXY_LOG_LEVEL}/g' /config/proxy-config.json

# start everything except glauth and idp https://github.com/owncloud/ocis/pull/2229
#ocis server --extensions="accounts, graph, graph-explorer, ocs, onlyoffice, proxy, settings, storage-authbasic, storage-authbearer, storage-frontend, storage-gateway, storage-groupsprovider, storage-home, storage-metadata, storage-public-link, storage-sharing, storage-users, storage-users-provider, store, thumbnails, web, webdav"

ocis server &
sleep 10

# idp and glauth are not needed -> replaced by Keycloak and OpenLDAP
ocis kill idp
ocis kill glauth

# workaround for loading proxy configuration
ocis kill proxy
sleep 10
ocis proxy server &


ocis list
wait
