#!/bin/bash
printenv
# replace oCIS domain in keycloak realm import
cp /opt/jboss/keycloak/ocis-realm.dist.json /opt/jboss/keycloak/ocis-realm.json
sed -i "s/ocis.owncloud.test/${OCIS_DOMAIN}/g" /opt/jboss/keycloak/ocis-realm.json

# run original docker-entrypoint
/opt/jboss/tools/docker-entrypoint.sh
