---
title: "Minimalistic Evaluation Guide for oCIS with Docker"
date: 2025-05-08T16:00:00+02:00
weight: 7
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/guides
geekdocFilePath: ocis-mini-eval.md
geekdocCollapseSection: true
---

{{< toc >}}

**IMPORTANT** This type of installation is **ONLY** intended for a quick evaluation and to access the instance from your local machine using internal certificates only. We recommend using test data only for this evaluation.

## Prerequisites

As a prerequisite, we assume you have Docker installed locally on a machine that has a GUI. A headless one will not work. You know how to use Docker commands and what they mean. This guide is not intended to be a detailed explanation.

## What you can Expect

By passing the commands as described, you will be able to get a very first look and feel of Infinite Scale, which can be accessed from your browser using `localhost`.

## Setup Procedure

The setup process is quite simple and is done from a terminal.

### Prepare Paths

Create directories if not exists:

```bash
mkdir -p $HOME/ocis/ocis-config \
mkdir -p $HOME/ocis/ocis-data
```

Set the user for the directories to be the same as the user inside the container:

```bash
sudo chown -Rfv 1000:1000 $HOME/ocis/
```

### Pull the Image

```bash
docker pull owncloud/ocis
```

### First Time Initialisation

```bash
docker run --rm -it \
    --mount type=bind,source=$HOME/ocis/ocis-config,target=/etc/ocis \
    --mount type=bind,source=$HOME/ocis/ocis-data,target=/var/lib/ocis \
    owncloud/ocis init --insecure yes
```

You will get an output like the following:

```txt {hl_lines=[6]}
=========================================
 generated OCIS Config
=========================================
 configpath : /etc/ocis/ocis.yaml
 user       : admin
 password   : t3p4N0jJ47LbhpQ04s9W%u1$d2uE3Y.3
```

Note that the password displayed is the one that will be used when you first log in or until it is changed.

### Recurring Start of Infinite Scale

```bash
docker run \
    --name ocis_runtime \
    --rm \
    -it \
    -p 9200:9200 \
    --mount type=bind,source=$HOME/ocis/ocis-config,target=/etc/ocis \
    --mount type=bind,source=$HOME/ocis/ocis-data,target=/var/lib/ocis \
    -e OCIS_INSECURE=true \
    -e PROXY_HTTP_ADDR=0.0.0.0:9200 \
    -e OCIS_URL=https://localhost:9200 \
    owncloud/ocis
```

## Access Infinite Scale with the Browser

To access Infinite Scale, open your browser and type as url:

```
https://localhost:9200
```

## Remove the Evaluation

```bash
sudo docker rmi owncloud/ocis \
sudo rm -r $HOME/ocis
```

## Next Steps

After evaluation, we strongly recommend that you use the [Install Infinite Scale on a Server](https://doc.owncloud.com/ocis/next/depl-examples/ubuntu-compose/ubuntu-compose-prod.html) documentation which is ready for production and remove this evaluation setup completely.
