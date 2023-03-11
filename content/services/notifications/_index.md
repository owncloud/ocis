---
title: Notifications
date: 2022-03-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/notifications
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The notifications extension is responsible for making users aware of changes. It listens on the event bus, filters relevant events, looks up the recipients email address and then queues an email with an external MTA.

## Table of Contents

{{< toc-tree >}}