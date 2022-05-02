#!/bin/sh
set -e

ocis init || true # will only initialize once
ocis server
