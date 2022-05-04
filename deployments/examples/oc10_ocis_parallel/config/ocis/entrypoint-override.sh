#!/bin/sh
set -e

ocis init || true # will only initialize once

#chmod 744 -R /etc/ocis
#setpriv --reuid=33 --regid=33 --clear-groups
ocis server
