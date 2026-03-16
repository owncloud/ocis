# Running API Tests in Kubernetes Cluster

## Table of Contents

- [K8s Setup](#k8s-setup)
- [Running API Tests](#running-api-tests)
  - [Run General API tests](#run-general-api-tests)
  - [Run Notification API tests](#run-notification-api-tests)
  - [Run Antivirus API tests](#run-antivirus-api-tests)
  - [Run Full Text Search API tests](#run-full-text-search-api-tests)
- [Cleanup the Setup](#cleanup-the-setup)

## K8s Setup

### Pre-requisites

1. Install the following tools:
   - [k3d](https://k3d.io/stable/)
   - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
   - [helm](https://helm.sh/docs/intro/install/)

2. Add these hosts to your `/etc/hosts` file:

   ```bash
   echo "127.0.0.1       ocis-server federation-ocis-server \
   clamav collabora onlyoffice fakeoffice tika \
   email" | sudo tee -a /etc/hosts
   ```

### Deploy oCIS in K8s

1. Change directory to `<ocis-rooot>/tests/config/k8s`:

   ```bash
   cd <ocis-root>/tests/config/k8s
   ```

2. Create K8s Cluster

   ```bash
   make create-cluster
   ```

3. Prepare Charts

   ```bash
   make prepare-charts
   ```

   ⚠️ NOTE: To run the test suites that require extra services,
   use the following appropriate environment variables:
   - `ENABLE_ANTIVIRUS=true`: Antivirus test suites
   - `ENABLE_EMAIL=true`: Notification test suites
   - `ENABLE_TIKA=true`: Content search test suites
   - `ENABLE_WOPI=true`: WOPI test suites
   - `ENABLE_OCM=true`: OCM test suites
   - `ENABLE_AUTH_APP=true`: auth-app test suites

   Example:

   ```bash
   ENABLE_EMAIL=true make prepare-charts
   ```

4. Deploy oCIS

   ```bash
   make deploy-ocis
   ```

## Running API Tests

Build and run the ociswrapper to be able to run the API tests tagged with `@env-config`.

1. Change directory to the ocis root:

   ```bash
   cd <ocis-root>
   ```

2. Build the ociswrapper:

   ```bash
   make -C tests/ociswrapper build
   ```

3. Run the ociswrapper:

   ```bash
   tests/ociswrapper/bin/ociswrapper serve --url https://ocis-server \
   --admin-username admin --admin-password admin --skip-ocis-run -n ocis-server
   ```

### Run General API tests

```bash
TEST_SERVER_URL=https://ocis-server \
K8S=true \
BEHAT_FEATURE=<test-suites-path>/apiDownloads/download.feature \
make test-acceptance-api
```

### Run Notification API tests

1. Check if setup [step 3](#deploy-ocis-in-k8s) is done correctly. (`ENABLE_EMAIL=true`)
2. Start the email server

   ```bash
   docker run -d -p 1025:1025 -p 8025:8025 axllent/mailpit:v1.22.3
   ```

3. Expose the email server to the cluster

   ```bash
   bash tests/config/k8s/expose-external-svc.sh email:1025
   ```

4. Run the tests

   ```bash
   TEST_SERVER_URL=https://ocis-server \
   EMAIL_HOST=email \
   EMAIL_PORT=8025 \
   K8S=true \
   BEHAT_FEATURE=<test-suites-path>/apiNotification/notification.feature \
   make test-acceptance-api
   ```

### Run Antivirus API tests

1. Check if setup [step 3](#deploy-ocis-in-k8s) is done correctly. (`ENABLE_ANTIVIRUS=true`)
2. Start the antivirus server

   ```bash
   docker run -d -p 3310:3310 owncloudci/clamavd
   ```

3. Expose the antivirus server to the cluster

   ```bash
   bash tests/config/k8s/expose-external-svc.sh clamav:3310
   ```

4. Run the tests

   ```bash
   TEST_SERVER_URL=https://ocis-server \
   K8S=true \
   BEHAT_FEATURE=<test-suites-path>/apiAntivirus/antivirus.feature \
   make test-acceptance-api
   ```

### Run Full Text Search API tests

1. Check if setup [step 3](#deploy-ocis-in-k8s) is done correctly. (`ENABLE_TIKA=true`)
2. Start the tika server

   ```bash
   docker run -d -p 9998:9998 apache/tika:3.2.2.0-full
   ```

3. Expose the tika server to the cluster

   ```bash
   bash tests/config/k8s/expose-external-svc.sh tika:9998
   ```

4. Run the tests

   ```bash
   TEST_SERVER_URL=https://ocis-server \
   K8S=true \
   BEHAT_FEATURE=<test-suites-path>/apiSearchContent/contentSearch.feature \
   make test-acceptance-api
   ```

## Cleanup the Setup

To delete the cluster and all the setup resources, run the following command:

```bash
make -C <ocis-rooot>tests/config/k8s cleanup
```

This will delete the K8s cluster, ocis-charts, and ocis logs directory.
