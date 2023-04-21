---
title: "oCIS and Containers"
date: 2022-06-14T16:00:00+02:00
weight: 5
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/guides
geekdocFilePath: ocis-and-containers.md
geekdocCollapseSection: true
---

## Cloud Native

Why do we recommend to work with containers?

{{< columns >}}

### {{< icon "scale-balanced" >}} &nbsp; Lightweight

Containers are more lightweight than VMs. It is easier to work with shared volumes and networks because they are isolated from the host system.
<--->

### {{< icon "shield-halved" >}} &nbsp; Dependencies

The container images have all dependencies installed and the maintainer takes care for keeping them up-to-date.

<--->

### {{< icon "gauge-high" >}} &nbsp; Scaling

In addition to that, containers help with scaling. You can run multiple instances of one container and distribute them across hosts.

{{< /columns >}}

## Docker compose

For oCIS deployments you often need multiple services. These services need to share resources like volumes and networks. If you do not use any orchestration tool, you would end up writing bash scripts to create and update containers and volumes and connect them via networks. This is what orchestration tools like docker compose can do for you. You define a service mesh using .yaml files and the tool tries to run and maintain that. You gain more value and a version history by using a version control system. Your deployment configuration is fully written down as a spec and you will never touch any system directly and change the config manually.

## Kubernetes

Containers are also used in [kubernetes](https://kubernetes.io/). Kubernetes is part of a huge ecosystem and is founded on best-of-breed practises to orchestrate large scale container applications and services.

## oCIS and Containers

oCIS was developed as microservices. We do not scale the whole system as a monolith but we scale the individual services.
