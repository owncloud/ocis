---
title: "Tracing"
date: 2020-05-13T12:09:00+01:00
weight: 55
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: tracing.md
---

{{< toc >}}

By default, we use [Jaeger](https://www.jaegertracing.io) for request tracing within oCIS. You can follow these steps
to get started:

1. Start Jaeger by using the all-in-one docker image:
   ```console
   docker run -d --name jaeger \
     -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
     -p 5775:5775/udp \
     -p 6831:6831/udp \
     -p 6832:6832/udp \
     -p 5778:5778 \
     -p 16686:16686 \
     -p 14268:14268 \
     -p 14250:14250 \
     -p 9411:9411 \
     jaegertracing/all-in-one:1.17
   ```
2. Every single oCIS service has its own environment variables for enabling and configuring tracing.
   1. You can enable and configure tracing on each service individually. For example, enable tracing
      in Reva when starting the oCIS single binary like this:
      ```console
      REVA_TRACING_ENABLED=true \
      REVA_TRACING_ENDPOINT=localhost:6831 \
      REVA_TRACING_COLLECTOR=http://localhost:14268/api/traces \
      ./bin/ocis server
      ```
   2. Enabling and configuring tracing on oCIS itself will forward the configuration to all services:
      ```console
      OCIS_TRACING_ENABLED=true \
      OCIS_TRACING_ENDPOINT=localhost:6831 \
      OCIS_TRACING_COLLECTOR=http://localhost:14268/api/traces \
      ./bin/ocis server
      ```
      If you want to set individual tracing configuration for each service, make sure to set
      `OCIS_TRACING_ENABLED=false`.
3. Make the actual request that you want to trace.
4. Open up the [Jaeger UI](http://localhost:16686) to analyze request traces.

For more information on Jaeger, please refer to their [Documentation](https://www.jaegertracing.io/docs/1.17/).
