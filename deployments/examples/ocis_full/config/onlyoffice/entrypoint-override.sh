#!/bin/sh
set -e

# we can't mount it directly because the run-document-server.sh script wants to move it
cp /etc/onlyoffice/documentserver/local.dist.json /etc/onlyoffice/documentserver/local.json

# Ensure license file has the correct permissions
if [ -f /var/www/onlyoffice/Data/license.lic ]; then
  echo "Fixing permissions on license.lic..."
  chmod 644 /var/www/onlyoffice/Data/license.lic || true
fi

/app/ds/run-document-server.sh
