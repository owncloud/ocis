#!/bin/bash

set -e

HOST_IP=$(hostname -i | awk '{print $1}')

NAMESPACE="ocis-server"

while [[ $# -gt 0 ]]; do
  case "$1" in
    -n|--namespace)
      NAMESPACE="$2"
      shift 2
      ;;
    *)
      services="$1"
      break
      ;;
  esac
done

if [ -z "$services" ]; then
  echo "[ERR] No services provided."
  echo "Usage: $0 [-n <namespace>] <svc1:port1,svc2:port2,...>"
  exit 1
fi

function expose_svc() {
  local k8s_endpoint k8s_service svc port server_ip
  svc=$1
  port=$2
  server_ip=$3

  k8s_endpoint=$(cat <<EOF
apiVersion: v1
kind: Endpoints
metadata:
  name: $svc
  namespace: $NAMESPACE
subsets:
- addresses:
  - ip: $server_ip
  ports:
  - port: $port
EOF
)
  echo -e "$k8s_endpoint" | kubectl apply -f -

  k8s_service=$(cat <<EOF
apiVersion: v1
kind: Service
metadata:
  name: $svc
  namespace: $NAMESPACE
spec:
  ports:
  - port: $port
    targetPort: $port
EOF
)
  echo -e "$k8s_service" | kubectl apply -f -
}

IFS=',' read -ra ADDR <<< "$services"
for service in "${ADDR[@]}"; do
  IFS=':' read -ra SVC <<< "$service"
  if [ ${#SVC[@]} -ne 2 ]; then
    echo "[ERR] Invalid service format (<svc>:<port>): $service"
    exit 1
  fi
  SERVER_IP=$(getent hosts "${SVC[0]}" | awk '{print $1}')
  if [ -z "$SERVER_IP" ]; then
    echo "[ERR] Could not resolve IP for service ${SVC[0]}"
    exit 1
  fi
  # for localhost services, use the host IP
  if [[ "$SERVER_IP" == "127.0.0.1" ]]; then
    SERVER_IP="$HOST_IP"
  fi
  expose_svc "${SVC[0]}" "${SVC[1]}" "$SERVER_IP"
done
