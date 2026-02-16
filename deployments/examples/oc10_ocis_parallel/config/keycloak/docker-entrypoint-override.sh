#!/bin/bash
printenv
# replace owncloud domain in keycloak realm import
cp /opt/jboss/keycloak/owncloud-realm.dist.json /opt/jboss/keycloak/owncloud-realm.json
sed -i "s/cloud.owncloud.test/${CLOUD_DOMAIN}/g" /opt/jboss/keycloak/owncloud-realm.json
sed -i "s/oc10-oidc-secret/${OC10_OIDC_CLIENT_SECRET}/g" /opt/jboss/keycloak/owncloud-realm.json
sed -i "s/ldap-bind-credential/${LDAP_ADMIN_PASSWORD}/g" /opt/jboss/keycloak/owncloud-realm.json



# run original docker-entrypoint
/opt/jboss/tools/docker-entrypoint.sh
