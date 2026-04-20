#!/bin/sh
# DRONE: Serves a fake WOPI hosting discovery XML response on port 8080.
# Hardcodes /drone/src as the repo root (Drone workspace mount point).
# Used by .drone.star wopiValidatorTests pipeline step for FakeOffice WOPI testing.
# When migrating: replace /drone/src with the appropriate path or an env var.

while true; do
  echo -e "HTTP/1.1 200 OK\n\n$(cat /drone/src/tests/config/drone/hosting-discovery.xml)" | nc -l -k -p 8080
done
