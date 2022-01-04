#!/usr/bin/env bash
echo "Writing custom config files..."

# openidconnect
gomplate \
  -f /etc/templates/oidc.config.php \
  -o ${OWNCLOUD_VOLUME_CONFIG}/oidc.config.php

# we need at least version 2.1.0 of the oenidconnect app
occ market:upgrade --major openidconnect
occ app:enable openidconnect

# user LDAP
gomplate \
  -f /etc/templates/ldap-config.tmpl.json \
  -o ${OWNCLOUD_VOLUME_CONFIG}/ldap-config.json

CONFIG=$(cat ${OWNCLOUD_VOLUME_CONFIG}/ldap-config.json)
occ config:import <<< $CONFIG

occ ldap:test-config "s01"
occ app:enable user_ldap
/bin/bash -c 'occ user:sync "OCA\User_LDAP\User_Proxy" -r -m remove'

occ market:upgrade --major web
occ app:enable web

# enable testing app
echo "Cloning and enabling testing app..."
git clone --depth 1 https://github.com/owncloud/testing.git /var/www/owncloud/apps/testing
occ app:enable testing

true
