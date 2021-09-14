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
~/code/refs/ocis-charts/ocis
❯ minikube status
m01
host: Stopped
kubelet: Stopped
apiserver: Stopped
kubeconfig: Stopped
```

After that, start it:

```console
~/code/refs/ocis-charts/ocis
❯ minikube start
😄  minikube v1.9.2 on Darwin 11.4
✨  Using the hyperkit driver based on existing profile
👍  Starting control plane node m01 in cluster minikube
🔄  Restarting existing hyperkit VM for "minikube" ...
🐳  Preparing Kubernetes v1.18.0 on Docker 19.03.8 ...
🌟  Enabling addons: default-storageclass, storage-provisioner
🏄  Done! kubectl is now configured to use "minikube"
```

## Run a chart

The easiest way to run the entire package is by using the available charts on https://github.com/refs/ocis-charts. It is not the purpose of this guide to explain the inner working of Kubernetes or its resources, as Helm builds an abstraction oon top of it, letting you interact with a simplified UI that roughly translates as "helm install" and "helm uninstall".

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
NAME                          READY   STATUS    RESTARTS   AGE
glauth-67b6d89577-zcf65       1/1     Running   0          23s
konnectd-85b9d6db59-s9wxq     1/1     Running   0          23s
ocis-proxy-6f6667986d-htdgq   1/1     Running   0          23s
ocs-6756757547-vqdb9          1/1     Running   0          23s
settings-9776fd95c-tx7dg      1/1     Running   0          23s
storages-6df6d479-j8t4k       10/10   Running   1          23s
store-85844f776f-fnsb2        1/1     Running   0          23s
web-56cb5c95b5-vr8qf          1/1     Running   0          23s
webdav-785b9f9ccc-4ll5n       1/1     Running   0          23s
```

5. get the exposed port for the kubernetes service: `minikube service list`

```console
|-------------|------------------|--------------|---------------------------|
|  NAMESPACE  |       NAME       | TARGET PORT  |            URL            |
|-------------|------------------|--------------|---------------------------|
| default     | konnectd-service | No node port |
| default     | kubernetes       | No node port |
| default     | ldap-service     | No node port |
| default     | ocs-service      | No node port |
| default     | proxy-service    |         9200 | http://192.168.64.5:30325 |
| default     | settings-service | No node port |
| default     | storages-service | No node port |
| default     | web-service      | No node port |
| kube-system | kube-dns         | No node port |
|-------------|------------------|--------------|---------------------------|
```

6. attempt a `PROPFIND` WebDAV request to the storage: `curl -v -k -u einstein:relativity -H "depth: 0" -X PROPFIND https://192.168.64.5:30325/remote.php/dav/files/ | xmllint --format -`

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

## What is GCP

### Can Helm charts run on GCP?

## Running on GCP (Google Cloud Platform)

## TODO

- setup an external IDP?
- make it work using the WebUI...
