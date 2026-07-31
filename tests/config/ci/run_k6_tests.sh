#!/bin/sh

set -e

SSH_OPTS="-o StrictHostKeyChecking=no -o ServerAliveInterval=60 -o ServerAliveCountMax=10"

if [ "$1" = "--ocis-log" ]; then
    sshpass -p "$SSH_OCIS_PASSWORD" ssh $SSH_OPTS "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" "bash ~/scripts/ocis.sh log"
    exit 0
fi

# clean up from previous runs
sshpass -p "$SSH_OCIS_PASSWORD" ssh $SSH_OPTS "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" "rm -rf k6-ocis"

# start ocis server
sshpass -p "$SSH_OCIS_PASSWORD" ssh $SSH_OPTS "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" \
    "OCIS_URL=${TEST_SERVER_URL} \
    OCIS_COMMIT_ID=${OCIS_COMMIT_SHA} \
    bash ~/scripts/ocis.sh start"

# wait for ocis to be ready (via SSH since TEST_SERVER_URL is on a private network)
echo "Waiting for OCIS to be ready at ${TEST_SERVER_URL}..."
RETRIES=0
MAX_RETRIES=24
until [ "$RETRIES" -ge "$MAX_RETRIES" ]; do
    STATUS=$(sshpass -p "$SSH_OCIS_PASSWORD" ssh $SSH_OPTS "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" \
        "curl -s -o /dev/null -w '%{http_code}' -k '${TEST_SERVER_URL}/.well-known/openid-configuration'" 2>/dev/null || echo "000")
    if [ "$STATUS" = "200" ]; then
        echo "OCIS is ready (HTTP ${STATUS})"
        break
    fi
    RETRIES=$((RETRIES + 1))
    echo "OCIS not ready (HTTP ${STATUS}), retrying in 5s... (${RETRIES}/${MAX_RETRIES})"
    sleep 5
done
if [ "$RETRIES" -ge "$MAX_RETRIES" ]; then
    echo "ERROR: OCIS did not become ready within $((MAX_RETRIES * 5))s"
    exit 1
fi

# start k6 tests
sshpass -p "$SSH_K6_PASSWORD" ssh $SSH_OPTS "$SSH_K6_USERNAME@$SSH_K6_REMOTE" \
    "TEST_SERVER_URL=${TEST_SERVER_URL} \
    bash ~/scripts/k6-tests.sh"

# stop ocis server
sshpass -p "$SSH_OCIS_PASSWORD" ssh $SSH_OPTS "$SSH_OCIS_USERNAME@$SSH_OCIS_REMOTE" "bash ~/scripts/ocis.sh stop"
