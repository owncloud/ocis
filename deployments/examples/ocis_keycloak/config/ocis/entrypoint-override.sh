#!/bin/sh

set -e

ocis server&
sleep 10

# stop builtin IDP since we use Keycloak as a replacement
ocis kill idp

echo "##################################################"
echo "change default secrets:"

ocis accounts update --password $STORAGE_LDAP_BIND_PASSWORD bc596f3c-c955-4328-80a0-60d018b4ad57 # REVA

echo "##################################################"

echo "##################################################"
echo "delete demo users" # users are provided by keycloak

set +e # accounts can only delete once, so it will fail the second time
# only admin, IDP and REVA user will be created because of ACCOUNTS_DEMO_USERS_AND_GROUPS=false
ocis accounts remove  820ba2a1-3f54-4538-80a4-2d73007e30bf # IDP user
ocis accounts remove ddc2004c-0977-11eb-9d3f-a793888cd0f8 # admin
set -e

echo "##################################################"

wait # wait for oCIS to exit
