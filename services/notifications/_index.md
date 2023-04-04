---
title: Notification service
date: 2023-04-04T08:47:13.617764295Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/notifications
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The notification service is responsible for sending emails to users informing them about events that happened. To do this it hooks into the event system and listens for certain events that the users need to be informed about.

## Table of Contents

* [Example Yaml Config](#example-yaml-config)

## Example Yaml Config

{{< include file="services/_includes/notifications-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/notifications_configvars.md" >}}

