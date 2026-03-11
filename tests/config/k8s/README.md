# Running API Tests

## Pre-requisites

Install the following tools:

- [k3d](https://k3d.io/stable/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
- [helm](https://helm.sh/docs/intro/install/)

## Setup

### 1. Create K8s Cluster

To set up the local Kubernetes cluster for running API tests, use the following commands:

```bash
make create-cluster
```

### 2. Setup Charts

```bash
make setup-charts
```

### 3. Deploy oCIS

```bash
make deploy-ocis
```
