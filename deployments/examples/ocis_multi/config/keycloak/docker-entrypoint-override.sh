#!/bin/bash
printenv
# replace oCIS domain in keycloak realm import
mkdir /opt/keycloak/data/import
sed -e "s/ocis.owncloud.test/${OCIS_DOMAIN}/g" /opt/keycloak/data/import-dist/ocis-realm.json > /opt/keycloak/data/import/oCIS-realm.json
# sed -e "s/ocis.ocm.owncloud.test/${OCIS_OCM_DOMAIN}/g" /opt/keycloak/data/import-dist/ocis-realm.json > /opt/keycloak/data/import/oCIS-realm.json

# run original docker-entrypoint
/opt/keycloak/bin/kc.sh "$@"
