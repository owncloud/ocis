---
title: WebDAV with OpenID Connect
date: 2021-11-17T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/clients/rclone
geekdocFilePath: webdav-sync-oidc.md
geekdocCollapseSection: true
---


## WebDAV with OpenID Connect

Rclone itself is not able to open and maintain an OpenID Connect session. But it is able to still use OpenID Connect for authentication by leveraging a so called OIDC-agent.

### Setting up the OIDC-agent

You need to install the [OIDC-agent](https://github.com/indigo-dc/oidc-agent) from your OS' package repository (eg. [Debian](https://github.com/indigo-dc/oidc-agent#debian-packages) or [MacOS](https://github.com/indigo-dc/oidc-agent#debian-packages)).


### Configuring the the OIDC-agent

Run the following command to add a OpenID Connect profile to your OIDC-agent. It will open the login page of OpenID Connect identity provider where you need to log in if you don't have an active session.

``` bash
oidc-gen \
 --client-id=oidc-agent \
 --client-secret="" \
 --pub \
 --issuer https://ocis.owncloud.test \
 --redirect-uri=http://localhost:12345 \
 --scope max \
 einstein-ocis-owncloud-test
```

If you have dynamic client registration enabled on your OpenID Connect identity provider, you can skip the `--client-id`,  `--client-secret` and `--pub` options.

If your're using a dedicated OpenID Connect client for the OIDC-agent, we recommend a public one with the following two redirect URIs: `http://127.0.0.1:*` and `http://localhost:*`. Alternatively you also may use the already existing OIDC client of the ownCloud Desktop Client (`--client-id=xdXOt13JKxym1B1QcEncf2XDkLAexMBFwiT9j6EfhhHFJhs2KM9jbjTmf8JBXE69` and `--client-secret=UBntmLjC2yYCeHwsyj73Uwo9TAaecAetRwMw0xYcvNL9yRdLSUi0hUAHfvCHFeFh`, no `--pub` set)

Please also note that the OIDC-agent will listen on your localhost interface on port 12345 for the time of the intial authentication. If that port is already occupied on your machine, you can easily change that by setting the `--redirect-uri` parameter to a different value.

After a successful login or an already existing session you will be redirected to success page of the OIDC-agent.
You will now be asked for a password for your account configuration, so that your OIDC session is secured and cannot be used by other people with access to your computer.



## Configure the WebDAV remote

First of all we need to set up our credentials and the WebDAV remote for Rclone. In this example we do this by setting environment variables. You might also set up a named remote or use command line options to achieve the same.

``` bash
export RCLONE_WEBDAV_VENDOR=owncloud
export RCLONE_WEBDAV_URL=https://ocis.owncloud.test/remote.php/webdav/
export RCLONE_WEBDAV_BEARER_TOKEN_COMMAND="oidc-token einstein-ocis-owncloud-test"
```


### Sync to the WebDAV remote

We now can use Rclone to sync the local folder `/tmp/test` to `/test` in your oCIS home folder.

``` bash
rclone sync :local:/tmp :webdav:/test
```

If your oCIS doesn't use valid SSL certificates, you may need to use `rclone --no-check-certificate sync ...`. 
