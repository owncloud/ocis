#!/bin/sh

while true; do
  echo -e "HTTP/1.1 200 OK\n\n$(cat /ocis/tests/config/ci/hosting-discovery.xml)" | nc -l -k -p 8080
done
