#!/bin/sh
set -e
cp /tmp/only-office.json /etc/onlyoffice/documentserver/local.json
openssl req -x509 -newkey rsa:4096 -keyout onlyoffice.key -out onlyoffice.crt -sha256 -days 365 -batch -nodes
mkdir -p /var/www/onlyoffice/Data/certs
cp onlyoffice.key /var/www/onlyoffice/Data/certs/
cp onlyoffice.crt /var/www/onlyoffice/Data/certs/
chmod 400 /var/www/onlyoffice/Data/certs/onlyoffice.key
/app/ds/run-document-server.sh
