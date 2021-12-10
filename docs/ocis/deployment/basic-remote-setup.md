---
title: "Basic Remote Setup"
date: 2020-02-27T20:35:00+01:00
weight: 16
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: basic-remote-setup.md
---

{{< toc >}}

The default configuration of the oCIS binary and the `owncloud/ocis` docker image assume, that you access oCIS on `localhost`. This enables you to do quick testing and development without any configuration.

If you need to access oCIS running in a docker container, on a VM or a remote machine via an other hostname than `localhost`, you need to configure this hostname in oCIS. The same applies if you are not using hostnames but instead an IP (eg. `192.168.178.25`).

## Start the oCIS fullstack server from binary

Upon first start of the oCIS fullstack server with `./bin/ocis server` it will generate a file `identifier-registration.yml` in the config folder in your current working directory. This file is used to configure the built-in identity provider and therefore contains the OpenID Connect issuer and also information about relying parties, for example ownCloud Web and our desktop and mobile applications.

{{< hint warning >}}
The `identifier-registration.yml` file will only be generated if it does not exist yet. If you want to change certain environment variables like `OCIS_URL`, please delete this file first before doing so. Otherwise your changes will not be applied correctly and you will run into errors.
{{< /hint >}}

{{< hint warning >}}
oCIS is currently in a Tech Preview state and is shipped with demo users. In order to secure your oCIS instances please follow following guide: [secure an oCIS instance]({{< ref "./#secure-an-ocis-instance" >}})
{{< /hint >}}

For the following examples you need to have the oCIS binary in your current working directory, we assume it is named `ocis` and it needs to be marked as executable. See [Getting Started]({{< ref "../getting-started/#binaries" >}}) for where to get the binary from.

### Using automatically generated certificates

In order to run oCIS with automatically generated and self signed certificates please execute following command. You need to replace `your-host` with an IP or hostname. Since you have only self signed certificates you need to have `OCIS_INSECURE` set to `true`.

```bash
OCIS_INSECURE=true \
PROXY_HTTP_ADDR=0.0.0.0:9200 \
OCIS_URL=https://your-host:9200 \
./ocis server
```

### Using already present certificates

If you have your own certificates already in place, you may want to make oCIS use them:

```bash
OCIS_INSECURE=false \
PROXY_HTTP_ADDR=0.0.0.0:9200 \
OCIS_URL=https://your-host:9200 \
PROXY_TRANSPORT_TLS_KEY=./certs/your-host.key \
PROXY_TRANSPORT_TLS_CERT=./certs/your-host.crt \
./ocis server
```

If you generated these certificates on your own, you might need to set `OCIS_INSECURE` to `true`.

For more configuration options check the configuration section in [oCIS]({{< ref "../config" >}}) and the oCIS extensions.

## Start the oCIS fullstack server with Docker Compose

Please have a look at our other [deployment examples]({{< ref "./" >}}).
