#!/usr/bin/env bash
echo "Writing custom config files..."

# user LDAP
gomplate \
  -f /etc/templates/ldap-config.tmpl.json \
  -o ${OWNCLOUD_VOLUME_CONFIG}/ldap-config.json

CONFIG=$(cat ${OWNCLOUD_VOLUME_CONFIG}/ldap-config.json)
occ config:import <<< $CONFIG

occ ldap:test-config "s01"
occ app:enable user_ldap
/bin/bash -c 'occ user:sync "OCA\User_LDAP\User_Proxy" -r -m remove'

# enable testing app
echo "Cloning and enabling testing app..."
git clone --depth 1 https://github.com/owncloud/testing.git /var/www/owncloud/apps/testing
occ app:enable testing

true
