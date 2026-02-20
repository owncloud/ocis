#!/bin/bash

# Cleanup function
cleanup() {
  kill 0
  exit 0
}

# Trap Ctrl+C and termination signals
trap cleanup SIGINT SIGTERM

DEPLOYMENTS=$(kubectl get deployments -n ocis -o jsonpath='{.items[*].metadata.name}')

for APP_NAME in $DEPLOYMENTS; do
  (
    while true; do
    #   echo "[$APP_NAME] Attaching to pods..."

      kubectl logs -f -n "ocis" \
        -l app="$APP_NAME" \
        --all-containers=true \

    done
  ) &
done

wait