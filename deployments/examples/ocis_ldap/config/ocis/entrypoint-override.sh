#!/bin/sh

set -e

mkdir -p /var/tmp/ocis/.config/
cp /config/web-config.dist.json /var/tmp/ocis/.config/web-config.json
sed -i 's/ocis.owncloud.test/'${OCIS_DOMAIN:-ocis.owncloud.test}'/g' /var/tmp/ocis/.config/web-config.json

ocis server&
sleep 10

# stop builtin accounts since we use LDAP only
ocis kill accounts
# stop builtin LDAP server since we use external LDAP only
ocis kill glauth

wait # wait for oCIS to exit
