---
title: Frontend
date: 2022-03-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/frontend
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The frontend service provides multiple HTTP endpoints to translate OCS, archiver and approvider requests into CS3 requests.

## Table of Contents

{{< toc-tree >}}

## OCS

The OCS endpoint implements the open collaboration services API in a backwards compatible manner.

### Sharing

Aggregating share information is one of the most time consuming operations in OCIS. The service fetches a list of either received or created shares and has to stat every resource individually. While stats are fast, the default behavior scales linearly with the number of shares.

To save network trips the sharing implementation can cache the stat requests with an in memory cache or in redis. It will shorten the response time by the network rountrip overhead at the cost of the API only eventually being updated. 

Setting `FRONTEND_OCS_RESOURCE_INFO_CACHE_TTL=60` would cache the stat info for 60 seconds. Increasing this value makes sense for large deployments with thousands of active users that keep the cache up to date. Low frequency usage scenarios should not expect a noticeable improvement.

## Archiver

The archiver endpoint provides bundled downloads of multiple files and folders.

## Appprovider

The appprovider endpoint is used to manage available apps that can be used to open different file types.