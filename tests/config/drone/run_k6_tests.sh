#!/bin/sh

set -e

if [ "$1" = "--ocis-log" ]; then
    sshpass -p "$SSH_OCIS_PASSWORD" ssh -o StrictHostKeyChecking=no "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" "bash ~/scripts/ocis.sh log"
    exit 0
fi

# start ocis server
sshpass -p "$SSH_OCIS_PASSWORD" ssh -o StrictHostKeyChecking=no "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" \
    "OCIS_URL=${TEST_SERVER_URL} \
    OCIS_COMMIT_ID=${DRONE_COMMIT} \
    bash ~/scripts/ocis.sh start"

# start k6 tests
sshpass -p "$SSH_K6_PASSWORD" ssh -o StrictHostKeyChecking=no "$SSH_K6_USERNAME@$SSH_K6_REMOTE" \
    "TEST_SERVER_URL=${TEST_SERVER_URL} \
    bash ~/scripts/k6-tests.sh"

# stop ocis server
sshpass -p "$SSH_OCIS_PASSWORD" ssh -o StrictHostKeyChecking=no "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" "bash ~/scripts/ocis.sh stop"
