---
title: "Kubernetes"
date: 2021-09-23T11:04:00+01:00
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

Without getting too deep in definitions, and for the purpose of compactness, Kubernetes can be summarized as a way of managing containers that run applications to ensure that there is no downtime and an optimal usage of resources. It provides with a framework in which to run distributed systems.

Kubernetes provides you with:
- **Service discovery and load balancing**: Kubernetes can expose a container using the DNS name or using their own IP address. If traffic to a container is high, Kubernetes is able to load balance and distribute the network traffic so that the deployment is stable.
- **Storage orchestration**: Kubernetes allows you to automatically mount a storage system of your choice, such as local storages, public cloud providers, and more.
- **Automated rollouts and rollbacks**: You can describe the desired state for your deployed containers using Kubernetes, and it can change the actual state to the desired state at a controlled rate. For example, you can automate Kubernetes to create new containers for your deployment, remove existing containers and adopt all their resources to the new container.
- **Automatic bin packing**: You provide Kubernetes with a cluster of nodes that it can use to run containerized tasks. You tell Kubernetes how much CPU and memory (RAM) each container needs. Kubernetes can fit containers onto your nodes to make the best use of your resources.
- **Self-healing**: Kubernetes restarts containers that fail, replaces containers, kills containers that don't respond to your user-defined health check, and doesn't advertise them to clients until they are ready to serve.
- **Secret and configuration management**: Kubernetes lets you store and manage sensitive information, such as passwords, OAuth tokens, and SSH keys. You can deploy and update secrets and application configuration without rebuilding your container images, and without exposing secrets in your stack configuration.

_[extracted from k8s docs](https://kubernetes.io/docs/concepts/overview/what-is-kubernetes/#why-you-need-kubernetes-and-what-can-it-do)_

If that is still too abstract, [here is an ELI5 writeup](https://dev.to/miguelmota/comment/filh).

### References and further reads

- [Marcel Wunderlich's](https://github.com/Deaddy) [4 series articles](http://deaddy.net/introduction-to-kubernetes-pt-1.html) on Kubernetes clarifying its declarative nature, deep diving into ingress networking, storage and monitoring.

### How does oCIS fit in the Kubernetes model

oCIS was designed with running on Kubernetes in mind. We set up to adopt the [Twelve-Factor App](https://12factor.net/) principles regarding configuration, with almost every aspect of oCIS being modifiable via environment variables. This comes in handy when you especially have a look at how a helm chart's (we will introduce this concept shortly) [list of values](https://github.com/owncloud/ocis-charts/blob/d8735e3222d2050504303851d3461909c86fcc89/ocis/values.yaml) looks like.

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

## oCIS charts

We have not yet published the oCIS Helm charts, therefore you need to clone the git repository manually. It currently also does not support to be run on Kind or Minikube clusters. For known issues and planned features, please have a look at the [GitHub issue tracker](https://github.com/owncloud/ocis-charts/issues).

Configuration options are described [here](https://github.com/owncloud/ocis-charts/tree/master/charts/ocis#configuration).

### Run oCIS

1. clone the charts: `git clone https://github.com/owncloud/ocis-charts.git /var/tmp/ocis-charts`
2. cd into the charts root: `cd /var/tmp/ocis-charts/charts/ocis`
3. install the package: `helm install ocis .` (you need to set configuration values in almost all cases)
4. verify the application is running in the cluster: `kubectl get pods`
