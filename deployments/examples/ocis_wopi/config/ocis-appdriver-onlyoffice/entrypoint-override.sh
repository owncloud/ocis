#!/bin/sh
set -e

sleep 120 #TODO: app driver should try again until onlyoffice is up...

ocis storage-app-provider server
