#!/bin/sh

set -e

echo "Check EOS MGM availability"
nc -z -w 3 $EOS_MGM_ALIAS 1094

/setup.sh

sleep 20

echo "----- [ocis] Starting oCIS -----"

ocis server
