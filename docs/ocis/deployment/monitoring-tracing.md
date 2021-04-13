---
title: "Monitoring & Tracing"
date: 2020-02-27T20:35:00+01:00
weight: 10
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: monitoring_tracing.md
---

{{< toc >}}

Monitoring and tracing gives developers and admin insights into a complex system, in this case oCIS.

If you are a developer and want to trace during developing you should have a look at [example server setup]({{< ref "../development/tracing" >}}).

This documentation describes how to set up a long running monitoring & tracing infrastructure for one or multiple oCIS servers or deployments. After reading this guide, you also should know everything needed to integrate oCIS into your existing monitoring and tracing infrastructure.

# Overview about the proposed solution

{{< svg src="ocis/static/monitoring_tracing_overview.drawio.svg" >}}

## Monitoring & tracing clients

We assume that you already have oCIS deployed on one or multiple servers by using our deployment examples (see rectangle on the left). On these servers our monitoring & tracing clients, namely Telegraf and Jaeger agent, need to be added.

Telegraf will collect host metrics (CPU, RAM, network, processes, ...) and docker metrics (per container CPU, RAM, network, ...). Telegraf is also configured to scrape metrics from Prometheus metric endpoints which oCIS exposes, this is done by the Prometheus input plugin . The metrics from oCIS and all other metrics gathered will be exposed with the Prometheus output plugin and can therefore be scraped by our monitoring & tracing server.

Jaeger agent is is being configured as target for traces in oCIS. It then will receive traces from all oCIS extensions, add some process tags to them and forward them to our Jaeger collector on our monitoring & tracing server.

For more information and how to deploy it, see [monitoring & tracing client](https://github.com/owncloud-devops/monitoring-tracing-client).

## Monitoring & tracing server

The monitoring & tracing server is considered as shared infrastructure and is normally used for different services. This means that oCIS is not the only software whose metrics and traces are available on the monitoring server. It is also possible that data of multiple oCIS instances are available on the monitoring server.

Metrics are scraped, stored and can be queried with Prometheus. For the visualization of these metrics Grafana is used. Because Prometheus is scraping the metrics from the oCIS server (pull model instead of a push model), the Prometheus server must have access to the exposed endpoint of the Telegraf Prometheus output plugin.

Jaeger collector receives traces sent by the Jaeger agent on the oCIS servers and persists them in ElasticSearch. From there the user can query and visualize the traces in Jaeger query or in Grafana. Because Jaeger agent is actively sending traces to the monitoring & tracing server, the server must be reachable from the oCIS server.

For more information and how to deploy it, see [monitoring & tracing server](https://github.com/owncloud-devops/monitoring-tracing-server).
