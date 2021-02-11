#!/bin/sh

set -e

ocis server&
sleep 10

echo "##################################################"
echo "change default secrets:"

# IDP
IDP_USER_UUID=$(ocis accounts list | grep "| Kopano IDP " | egrep '[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}' -o)
echo "  IDP user UUID: $IDP_USER_UUID"
ocis accounts update --password $IDP_LDAP_BIND_PASSWORD $IDP_USER_UUID

# REVA
REVA_USER_UUID=$(ocis accounts list | grep " | Reva Inter " | egrep '[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}' -o)
echo "  Reva user UUID: $REVA_USER_UUID"
ocis accounts update --password $STORAGE_LDAP_BIND_PASSWORD $REVA_USER_UUID

killall ocis
echo "default secrets changed"
echo "##################################################"

ocis server
