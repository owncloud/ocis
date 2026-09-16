#!/bin/bash

set -e

deploy_name=$1
service_port=$2
node_port=$3
namespace="${TEST_SERVER_DOMAIN:-ocis-server}"

if [[ -z "$deploy_name" ]] || [[ -z "$service_port" ]] || [[ -z "$node_port" ]]; then
  echo "[ERROR] Missing arguments. Usage: $0 <deploy_name> <service_port> <node_port>"
  exit 1
fi

service_name="${deploy_name}-${node_port}-np"

kubectl -n "$namespace" expose deployment "$deploy_name" --type=NodePort --port="$service_port" --name="$service_name"
kubectl -n "$namespace" patch svc "$service_name" -p "{\"spec\":{\"ports\":[{\"port\":${service_port},\"nodePort\":${node_port}}]}}"
