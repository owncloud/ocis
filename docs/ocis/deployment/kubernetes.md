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

### How does oCIS firs in the Kubernetes model

oCIS was designed with running on Kubernetes in mind. We set up to adopt the [Twelve-Factor App](https://12factor.net/) principles regarding configuration, so almost every aspect of oCIS can be modified using an environment variable. This comes in handy when you especially have a look at how a helm chart's (we will introduce this concept shortly) [list of values](https://github.com/refs/ocis-charts/blob/d8735e3222d2050504303851d3461909c86fcc89/ocis/values.yaml) looks like.

## What is Minikube

## What is `kubectl`

## Setting up Minikube

## Helm Charts, and why they are handy for oCIS

## Run a local set of Helm charts

## What is GCP

### Can Helm charts run on GCP?

## Running on GCP (Google Cloud Platform)
