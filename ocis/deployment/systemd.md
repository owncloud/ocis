---
title: "Systemd service"
date: 2020-09-27T06:00:00+01:00
weight: 16
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: systemd.md
---

{{< toc >}}

## Install the oCIS binary

Download the oCIS binary of your preferred version and for your CPU architecture and operating system from [download.owncloud.com](https://download.owncloud.com/ocis/ocis).

Rename the downloaded binary to `ocis` and move it to `/usr/bin/`. As a next step, you need to mark it as executable with `chmod +x /usr/bin/ocis`.

When you now run `ocis help` on your command line, you should see the available options for the oCIS command.

## Systemd service definition

Create the Systemd service definition for oCIS in the file `/etc/systemd/system/ocis.service` with following content:

```systemd
[Unit]
Description=OCIS server

[Service]
Type=simple
User=root
Group=root
EnvironmentFile=/etc/ocis/ocis.env
ExecStart=ocis server
Restart=always

[Install]
WantedBy=multi-user.target
```

{{< hint warning >}}
For reasons of simplicity we are using the root user and group to run oCIS which is not recommended. Please use a non-root user in production environments and modify the oCIS service definition accordingly.
{{< /hint >}}

In the service definition we referenced `/etc/ocis/ocis.env` as our file containing environment variables for the oCIS process.
In order to create the file we need first to create the folder `/etc/ocis/` and then we can add the actual `/etc/ocis/ocis.env` with following content:

```bash
OCIS_URL=https://some-hostname-or-ip:9200
PROXY_HTTP_ADDR=0.0.0.0:9200
OCIS_INSECURE=false

OCIS_LOG_LEVEL=error

OCIS_CONFIG_DIR=/etc/ocis
OCIS_BASE_DATA_PATH=/var/lib/ocis
```

Since we set `OCIS_CONFIG_DIR` to `/etc/ocis` you can also place configuration files in this directory.

Please change your `OCIS_URL` in order to reflect your actual deployment. If you are using self-signed certificates you need to set `OCIS_INSECURE=true` in `/etc/ocis/ocis.env`.

oCIS will store all data in `/var/lib/ocis`, because we configured it so by setting `OCIS_BASE_DATA_PATH`. Therefore you need to create that directory and make it accessible to the user, you use to start oCIS.

## Starting the oCIS service

Initialize the oCIS configuration by running `ocis init --config-path /etc/ocis`.

You can enable oCIS now by running `systemctl enable --now ocis`. It will ensure that oCIS also is restarted after a reboot of the host.

If you need to restart oCIS because of configuration changes in `/etc/ocis/ocis.env`, run `systemctl restart ocis`.

You can have a look at the logs of oCIS by issuing `journalctl -f -u ocis`.
