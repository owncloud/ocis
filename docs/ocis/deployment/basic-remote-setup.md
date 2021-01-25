---
title: "Basic Remote Setup"
date: 2020-02-27T20:35:00+01:00
weight: 16
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: basic-remote-setup.md
---

{{< toc >}}

Out of the box the ocis single binary and the `owncloud/ocis` docker image are configured to run on localhost for quick testing and development.

If you need to access ocis on a VM or a remote machine e.g. when testing a mobile client you need to configure ocis to run on a different host.

## Use the binary

If you start the ocis fullstack for the first time with `./bin/ocis server` it will generate a file `identifier-registration.yml` in the config folder relative to its location. This file is used to configure the clients for the built-in Identity Provider.

{{< hint warning >}}
**Outdated version**\
The `identifier-registration.yml` file will only be generated if there is no such file in place. You could miss updates on this file. Run `make clean` to delete the file and keep the development environment tidy otherwise as well.
{{< /hint >}}

### Add your hostname to the idp config

Let us assume `your-host` is your remote domain name or IP address. Add your host to the `identifier-registration.yml` like this:

```yaml {linenos=table,hl_lines=["15-17",21]}
# OpenID Connect client registry.
clients:
  - id: web
    name: ownCloud web app
    application_type: web
    insecure: yes
    trusted: yes
    redirect_uris:
      - http://localhost:9100/
      - http://localhost:9100/oidc-callback.html
      - http://localhost:9100/oidc-silent-redirect.html
      - https://localhost:9200/
      - https://localhost:9200/oidc-callback.html
      - https://localhost:9200/oidc-silent-redirect.html
      - https://your-server:9200/
      - https://your-server:9200/oidc-callback.html
      - https://your-server:9200/oidc-silent-redirect.html
    origins:
      - http://localhost:9100
      - https://localhost:9200
      - https://your-server:9200
```

In this example we do not change the default port (`9200`). But this could be changed to another port.

### Start the ocis fullstack server

You need to configure `your-host` in some services to provide the needed public resources.

This snippet will start the ocis server with auto generated self signed certificates:

```bash
PROXY_HTTP_ADDR=0.0.0.0:9200 \
OCIS_URL=https://your-server:9200 \
KONNECTD_TLS=0 \
./bin/ocis server
```

For more configuration options check the configuration section in [ocis](https://owncloud.github.io/ocis/configuration/) and every ocis extension.

{{< hint info >}}
**TLS Certificate**\
If you have a CA signed certificate for your domain, add the following configurations:
```
PROXY_TRANSPORT_TLS_KEY=./certs/your-host.key \
PROXY_TRANSPORT_TLS_CERT=./certs/your-host.crt \
```
{{< /hint >}}

## Use Docker Compose

We are using our [docker compose playground](https://github.com/owncloud-docker/compose-playground) as a repository to share snippets that make our test setups easier and more aligned.

You can start oCIS with docker very easily on a different host using this snippet.

Let us assume your local IP is `192.168.103.195`

```bash
git clone https://github.com/owncloud-docker/compose-playground.git
cd compose-playground/compose/ocis

sed -i -e 's/your-url/192.168.103.195/g' config/identifier-registration.yml

cat << EOF > .env
OCIS_URL=https://192.168.103.195
OCIS_DOCKER_TAG=latest
EOF

curl -k https://192.168.103.195:9200/status.php
```
