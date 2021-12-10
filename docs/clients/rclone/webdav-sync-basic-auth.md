---
title: WebDAV with Basic Authentication
date: 2021-11-17T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/clients/rclone
geekdocFilePath: webdav-sync-basic-auth.md
geekdocCollapseSection: true
---


## WebDAV with Basic Authentication

{{< hint danger >}}
Basic Authentication is disabled by default in oCIS because of security considerations. In order to make the following Rclone commands work the oCIS administrator needs to enable Basic Authentication eg. by setting the the environment variable `PROXY_ENABLE_BASIC_AUTH` to `true`. 

Please consider to use [Rclone with OpenID Connect]({{< ref "webdav-sync-oidc.md" >}}) instead.
{{< /hint >}}

For the usage of a WebDAV remote with Rclone see also the [Rclone documentation](https://rclone.org/webdav/)

## Configure the WebDAV remote

First of all we need to set up our credentials and the WebDAV remote for Rclone. In this example we do this by setting environment variables. You might also set up a named remote or use command line options to achieve the same.

``` bash
export RCLONE_WEBDAV_VENDOR=owncloud
export RCLONE_WEBDAV_URL=https://ocis.owncloud.test/remote.php/webdav/
export RCLONE_WEBDAV_USER=einstein
export RCLONE_WEBDAV_PASS=$(rclone obscure relativity)
```

{{< hint info >}}
Please note that `RCLONE_WEBDAV_PASS` is not set to the actual password, but to the value returned by `rclone obscure <password>`.
{{< /hint >}}

We now can use Rclone to sync the local folder `/tmp/test` to `/test` in your oCIS home folder.


### Sync to the WebDAV remote

``` bash
rclone sync :local:/tmp :webdav:/test
```

If your oCIS doesn't use valid SSL certificates, you may need to use `rclone --no-check-certificate sync ...`. 
