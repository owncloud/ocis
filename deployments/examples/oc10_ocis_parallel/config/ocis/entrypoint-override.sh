#!/bin/sh
set -e
ocis server &
sleep 10

# idp, glauth and accounts are not needed -> replaced by Keycloak and OpenLDAP
ocis kill idp
ocis kill glauth
ocis kill accounts

wait
