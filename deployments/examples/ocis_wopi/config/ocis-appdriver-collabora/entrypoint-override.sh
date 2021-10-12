#!/bin/sh
set -e

# if Collabora is already up and we have a new oCIS image, this app provider starts up too fast for oCIS
sleep 20

ocis storage-app-provider server
