---
title: "Kubernetes"
date: 2021-10-14T11:04:00+01:00
weight: 25
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: kubernetes.md
---

{{< toc >}}

## What is Kubernetes

Formally described as:

> Kubernetes is a portable, extensible, open-source platform for managing containerized workloads and services, that facilitates both declarative configuration and automation.

_[source](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/)_

Without getting too deep in definitions, and for the purpose of compactness, Kubernetes can be summarized as a way of managing containers that run applications to ensure that there is no downtime and a optimal usage of resources. It provides with a framework in which to run distributed systems.

Kubernetes provides you with:
- **Service discovery and load balancing**: Kubernetes can expose a container using the DNS name or using their own IP address. If traffic to a container is high, Kubernetes is able to load balance and distribute the network traffic so that the deployment is stable.
- **Storage orchestration**: Kubernetes allows you to automatically mount a storage system of your choice, such as local storages, public cloud providers, and more.
- **Automated rollouts and rollbacks**: You can describe the desired state for your deployed containers using Kubernetes, and it can change the actual state to the desired state at a controlled rate. For example, you can automate Kubernetes to create new containers for your deployment, remove existing containers and adopt all their resources to the new container.
- **Automatic bin packing**: You provide Kubernetes with a cluster of nodes that it can use to run containerized tasks. You tell Kubernetes how much CPU and memory (RAM) each container needs. Kubernetes can fit containers onto your nodes to make the best use of your resources.
- **Self-healing**: Kubernetes restarts containers that fail, replaces containers, kills containers that don't respond to your user-defined health check, and doesn't advertise them to clients until they are ready to serve.
- **Secret and configuration management**: Kubernetes lets you store and manage sensitive information, such as passwords, OAuth tokens, and SSH keys. You can deploy and update secrets and application configuration without rebuilding your container images, and without exposing secrets in your stack configuration.

_[extracted from k8s docs](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/#why-you-need-kubernetes-and-what-can-it-do)_

If that is still too abstract, [here is an ELI5 writeup](https://dev.to/miguelmota/comment/filh).

### How does oCIS fit in the Kubernetes model

oCIS was designed with running on Kubernetes in mind. We set up to adopt the [Twelve-Factor App](https://12factor.net/) principles regarding configuration, with almost every aspect of oCIS being modifiable via environment variables. This comes in handy when you especially have a look at how a helm chart's (we will introduce this concept shortly) [list of values](https://github.com/refs/ocis-charts/blob/d8735e3222d2050504303851d3461909c86fcc89/ocis/values.yaml) looks like.

## What is Minikube

[Minikube](https://minikube.sigs.k8s.io/docs/) lets you run a Kubernetes cluster locally. It is the most approachable way to test a deployment. It requires no extra configuration on any cloud platform, as everything runs on your local machine. For the purpose of these docs, this is the first approach we chose to run oCIS and will develop on how to set it up.

## What is `kubectl`

[kubectl](https://kubernetes.io/docs/tasks/tools/) is the command-line tool for Kubernetes. It allows users to run commands against a k8s cluster the user has access to. It supports for having multiple contexts for as many clusters as you have access to. In these docs we will setup 2 contexts, a minikube and a GCP context.

## What are Helm Charts, and why they are useful for oCIS

[Helm](https://helm.sh/) is the equivalent of a package manager for Kubernetes. It can be described as a layer on top of how you would write pods, deployments or any other k8s resource declaration.

### Installing Helm

[Follow the official installation guide](https://helm.sh/docs/intro/install/).

## Setting up Minikube

For a guide on how to set minikube up follow the [official minikube start guide](https://minikube.sigs.k8s.io/docs/start/) for your specific OS.

### Start minikube

First off, verify your installation is correct:

```console
~/code/refs/ocis-charts
❯ minikube status
minikube
type: Control Plane
host: Stopped
kubelet: Stopped
apiserver: Stopped
kubeconfig: Stopped
```

After that, start the cluster:

```console
~/code/refs/ocis-charts
❯ minikube start
😄  minikube v1.23.0 on Darwin 11.4
✨  Using the docker driver based on existing profile
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
🔄  Restarting existing docker container for "minikube" ...
🐳  Preparing Kubernetes v1.22.1 on Docker 20.10.8 ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

_On these docs, we are using the Docker driver on Mac._

## Run a chart

The easiest way to run the entire package is by using the available charts on https://github.com/refs/ocis-charts. It is not the purpose of this guide to explain the inner working of Kubernetes or its resources, as Helm builds an abstraction oon top of it, letting you interact with a refined interface that roughly translates as "helm install" and "helm uninstall".

In order to host charts one can create a [charts repository](https://helm.sh/docs/topics/chart_repository/), but this is outside the scope of this documentation. Having said that, we will assume you have access to a cli and git.

### Requirements

1. minikube up and running.
2. `kubectl` installed. By [default you should be able to access the minikube's cluster](https://minikube.sigs.k8s.io/docs/handbook/kubectl/). If you chose not to install `kubectl`, minikube wraps `kubectl` as `minikube kubectl`.
3. helm cli installed.
4. git installed.

### Setup

1. clone the charts: `git clone https://github.com/refs/ocis-charts.git /var/tmp/ocis-charts`
2. cd into the charts root: `cd /var/tmp/ocis-charts/ocis`
3. install the package: `helm install ocis .`
4. verify the application is running in the cluster: `kubectl get pods`

```console
❯ kubectl get pods
NAME                          READY   STATUS    RESTARTS         AGE
glauth-5fb678b9cb-zs5qh       1/1     Running   3 (10m ago)      3h33m
ocis-proxy-848f988687-g7fmb   1/1     Running   2 (10m ago)      130m
ocs-6bb8896dd6-t4bkx          1/1     Running   3 (10m ago)      3h33m
settings-6bf77f978d-27rdf     1/1     Running   3 (10m ago)      3h33m
storages-6b45f9c4-2j696       10/10   Running   23 (4m43s ago)   112m
store-cf79db94d-hvb7z         1/1     Running   3 (10m ago)      3h33m
web-8685fdd574-tmkfb          1/1     Running   2 (10m ago)      157m
webdav-f8d4dd7c6-vv4n7        1/1     Running   3 (10m ago)      3h33m
```

5. expose the proxy as a service to the host

```console
~/code/refs/ocis-charts
❯ minikube service proxy-service --url
🏃  Starting tunnel for service proxy-service.
|-----------|---------------|-------------|------------------------|
| NAMESPACE |     NAME      | TARGET PORT |          URL           |
|-----------|---------------|-------------|------------------------|
| default   | proxy-service |             | http://127.0.0.1:63633 |
|-----------|---------------|-------------|------------------------|
http://127.0.0.1:63633
❗  Because you are using a Docker driver on darwin, the terminal needs to be open to run it.
```

6. attempt a `PROPFIND` WebDAV request to the storage: `curl -v -k -u einstein:relativity -H "depth: 0" -X PROPFIND https://127.0.0.1:63633/remote.php/dav/files/ | xmllint --format -`

If all is correctly setup, you should expect a response back:

```xml
<?xml version="1.0" encoding="utf-8"?>
<d:multistatus xmlns:d="DAV:" xmlns:s="http://sabredav.org/ns" xmlns:oc="http://owncloud.org/ns">
  <d:response>
    <d:href>/remote.php/dav/files/einstein/</d:href>
    <d:propstat>
      <d:prop>
        <oc:id>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OjZlMWIyMjdmLWZmYTQtNDU4Ny1iNjQ5LWE1YjBlYzFkMTNmYw==</oc:id>
        <oc:fileid>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OjZlMWIyMjdmLWZmYTQtNDU4Ny1iNjQ5LWE1YjBlYzFkMTNmYw==</oc:fileid>
        <d:getetag>"92cc7f069c8496ee2ce33ad4f29de763"</d:getetag>
        <oc:permissions>WCKDNVR</oc:permissions>
        <d:resourcetype>
          <d:collection/>
        </d:resourcetype>
        <d:getcontenttype>httpd/unix-directory</d:getcontenttype>
        <oc:size>4096</oc:size>
        <d:getlastmodified>Tue, 14 Sep 2021 12:45:29 +0000</d:getlastmodified>
        <oc:favorite>0</oc:favorite>
      </d:prop>
      <d:status>HTTP/1.1 200 OK</d:status>
    </d:propstat>
  </d:response>
</d:multistatus>
```

## Setting up an external identity provider

The previous setup works because the proxy is configured to run using basic auth, but if we want to actually use the WebUI we will need an external identity provider. From here on the setup is composed of:

- keycloak
  - traefik
  - postgresql

Running on i.e: `https://keycloak.owncloud.works`. Because of this we have to adjust some of `values.yaml` key / values to:

```diff
diff --git a/ocis/values.yaml b/ocis/values.yaml
index fbc229c..5b36fbd 100644
--- a/ocis/values.yaml
+++ b/ocis/values.yaml
@@ -1,9 +1,9 @@
 # when in local tunnel mode, ingressDomain is the proxy address.
 # sadly when in combination with --set, anchors are lost.
-ingressDomain: &ingressDomain "https://stale-wasp-86.loca.lt"
+ingressDomain: &ingressDomain "https://keycloak.owncloud.works"

 # base ocis image
-image: owncloud/ocis:1.0.0-rc8-linux-amd64
+image: owncloud/ocis:1.11.0-linux-amd64

 # set of ocis services to create deployments objects.
 services:
@@ -22,6 +22,8 @@ services:
       value: "debug"
     - name: "PROXY_REVA_GATEWAY_ADDR"
       value: "storages-service:9142"
+    - name: "PROXY_OIDC_ISSUER"
+      value: "https://keycloak.ocis-keycloak.released.owncloud.works/auth/realms/oCIS"
     - name: "PROXY_ENABLE_BASIC_AUTH"
       value: "'true'" # see https://stackoverflow.com/a/44692213/2295410
     volumeMounts:
@@ -81,34 +85,6 @@ services:
     labels:
       app: "glauth"
     args: ["glauth"]
   settings:
     metadata:
       name: "settings"
@@ -135,11 +111,11 @@ services:
     args: ["web"]
     env:
     - name: "WEB_UI_CONFIG_SERVER"
-      value: *ingressDomain
+      value: "https://127.0.0.1:51559/"
     - name: "WEB_OIDC_METADATA_URL"
-      value: *ingressDomain
+      value: "https://keycloak.owncloud.works/auth/realms/oCIS/.well-known/openid-configuration"
     - name: "WEB_OIDC_AUTHORITY"
-      value: *ingressDomain
+      value: "https://keycloak.owncloud.works/auth/realms/oCIS/.well-known/openid-configuration"
     ports:
       values:
       - name: "http"
@@ -231,4 +207,4 @@ kubeServices:
       - protocol: TCP
         port: 9100
         targetPort: 9100
```

NOTE: the IDP has to be properly configure with an oCIS realm and a `web` client configured. There are example config file that have to be adjusted depending on your environment on our [docker-compose examples](https://github.com/owncloud/ocis/tree/master/deployments/examples/ocis_keycloak/config/keycloak).

## What is GCP

> Google Cloud Platform (GCP), offered by Google, is a suite of cloud computing services that runs on the same infrastructure that Google uses internally for its end-user products

One of such offered services are [Google Kubernetes Engines (GKE)](https://cloud.google.com/kubernetes-engine).

### Can Helm charts run on GCP?

Yes. The next logical step would be to deploy this charts on GKE. There is a pretty thorough guide [at shippable.com](http://docs.shippable.com/deploy/tutorial/deploy-to-gcp-gke-helm/) that, for the purposes of our docs, we are only interested on step 5, as we already explain the previous concepts, and provide with the Charts.

## TODOs

- while log-in works and creating folders work, uploading fails, most likely a configuration issue that has to be solved.
