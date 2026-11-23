#!/bin/bash
printenv
# replace oCIS domain and the invitations client secret in keycloak realm import
mkdir /opt/keycloak/data/import
sed -e "s/ocis.owncloud.test/${OCIS_DOMAIN}/g" \
    -e "s/__INVITATIONS_CLIENT_SECRET__/${INVITATIONS_CLIENT_SECRET}/g" \
    /opt/keycloak/data/import-dist/ocis-realm.json > /opt/keycloak/data/import/oCIS-realm.json

# run original docker-entrypoint
/opt/keycloak/bin/kc.sh "$@"
