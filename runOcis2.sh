#!/usr/bin/bash
kubectl create namespace ocis2

kubectl -n kube-system patch configmap coredns-custom \
  --type merge \
  -p '{"data":{"rewritehost.override2":"rewrite name exact ocis2-server host.k3d.internal"}}'

kubectl -n kube-system rollout restart deployment coredns

helm upgrade --install ocis2 ./charts/ocis/ \
  -n ocis2 --create-namespace \
  --values ./charts/ocis/ci/deployment-values-ocis2.yaml \
  --rollback-on-failure --timeout 5m0s

kubectl create configmap coredns-custom --namespace kube-system \
  --from-literal='rewritehost.override=rewrite name exact ocis-server host.k3d.internal,ocis2-server host.k3d.internal'

kubectl -n kube-system rollout restart deployment coredns



kubectl create configmap coredns-custom \
  -n kube-system \
  --from-literal=rewritehost.override=$'rewrite name exact ocis-server host.k3d.internal\nrewrite name exact ocis2-server host.k3d.internal' \
  --dry-run=client -o yaml | kubectl apply -f -
