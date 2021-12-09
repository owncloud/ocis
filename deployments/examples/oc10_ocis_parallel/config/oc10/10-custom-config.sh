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

cp /tmp/ldap-sync-cron /etc/cron.d
chown root:root /etc/cron.d/ldap-sync-cron

# ownCloud Web
gomplate \
  -f /etc/templates/web.config.php \
  -o ${OWNCLOUD_VOLUME_CONFIG}/web.config.php

gomplate \
  -f /etc/templates/web-config.tmpl.json \
  -o ${OWNCLOUD_VOLUME_CONFIG}/config.json

occ market:upgrade --major web
occ app:enable web

true
