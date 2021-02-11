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

echo "default secrets changed"
echo "##################################################"

echo "##################################################"
echo "delete demo users:" # demo users are provided by keycloak

set +e # accounts can only delete once, so it will fail the second time
ocis accounts remove 4c510ada-c86b-4815-8820-42cdf82c3d51
ocis accounts remove ddc2004c-0977-11eb-9d3f-a793888cd0f8
ocis accounts remove 932b4540-8d16-481e-8ef4-588e4b6b151c
ocis accounts remove 058bff95-6708-4fe5-91e4-9ea3d377588b
ocis accounts remove f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c
set -e

echo "##################################################"

killall ocis

ocis server
