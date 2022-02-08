#!/usr/bin/env bash

# enable testing app
echo "Cloning and enabling testing app..."
git clone --depth 1 https://github.com/owncloud/testing.git /var/www/owncloud/apps/testing
occ app:enable testing

true
