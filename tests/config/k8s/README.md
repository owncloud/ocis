# Running API Tests

## Pre-requisites

1. Install the following tools:
   - [k3d](https://k3d.io/stable/)
   - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
   - [helm](https://helm.sh/docs/intro/install/)

2. Add these hosts to your `/etc/hosts` file:

```bash
echo "127.0.0.1       ocis-server federation-ocis-server clamav collabora onlyoffice fakeoffice tika email" | sudo tee -a /etc/hosts
```

## Setup

### 1. Create K8s Cluster

To set up the local Kubernetes cluster for running API tests, use the following commands:

```bash
make create-cluster
```

### 2. Prepare Charts

```bash
make prepare-charts
```

### 3. Deploy oCIS

```bash
make deploy-ocis
```
