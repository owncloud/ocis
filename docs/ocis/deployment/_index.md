---
title: "Deployment"
date: 2020-10-01T20:35:00+01:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/deployment
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc >}}

## Deployments scenarios and examples
This section handles deployments and operations for admins and people who are interested in how versatile oCIS is. If you want to just try oCIS you may also follow [Getting started]({{< ref "../getting-started" >}}).

### Setup oCIS on your server
oCIS deployments are super simple, yet there are many configurations possible for advanced setups.

- [Basic oCIS setup]({{< ref "basic-remote-setup" >}}) - configure domain, certificates and port
- [oCIS setup with Traefik for SSL termination]({{< ref "ocis_traefik" >}})
- [oCIS setup with Keycloak as identity provider]({{< ref "ocis_keycloak" >}})
- [oCIS setup with WOPI server to open office documents in your browser]({{< ref "ocis_wopi" >}})
- [oCIS with S3 storage backend (MinIO)]({{< ref "ocis_s3" >}})
- [oCIS with the Hello extension example]({{< ref "ocis_hello" >}})


## Secure an oCIS instance

### Change default secrets
oCIS uses two system users which are needed for being operational:
- Reva Inter Operability Platform (bc596f3c-c955-4328-80a0-60d018b4ad57)
- Kopano IDP (820ba2a1-3f54-4538-80a4-2d73007e30bf)

Both have simple default passwords which need to be changed. Currently, changing a password is only possible on the command line. You need to run `ocis accounts update --password <new-password> <id>` for both users.

The new password for the Reva Inter Operability Platform user must be made available to oCIS by using the environment variable `STORAGE_LDAP_BIND_PASSWORD`. The same applies to the new Kopano IDP user password, which needs do be made available to oCIS in `IDP_LDAP_BIND_PASSWORD`.

Furthermore, oCIS uses a shared secret to sign JWT tokens for inter service authorization, which also needs to be changed by the user.
You can change it by setting the `OCIS_JWT_SECRET` environment variable for oCIS to a random string.

Another is used secret for singing JWT tokens for uploads and downloads, which also needs to be changed by the user.
You can change it by setting the `STORAGE_TRANSFER_SECRET` environment variable for oCIS to a random string.

### Delete demo users

{{< hint info >}}
Before deleting the demo users mentioned below, you must create a new account for yourself and assign it to the administrator role.
{{< /hint >}}

oCIS ships with a few demo users besides the system users:
- Admin (ddc2004c-0977-11eb-9d3f-a793888cd0f8)
- Albert Einstein (4c510ada-c86b-4815-8820-42cdf82c3d51)
- Richard Feynman (932b4540-8d16-481e-8ef4-588e4b6b151c)
- Maurice Moss (058bff95-6708-4fe5-91e4-9ea3d377588b)
- Marie Curie (f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c)

You can view them in ownCloud Web if you log in as Admin user or list them by running `ocis accounts list`.
After adding your own user it is safe to delete the demo users in the web UI or with the command `ocis accounts remove <id>`. Please do not delete the system users (see [change default secrets]({{< ref "./#change-default-secrets" >}})) or oCIS will not function properly anymore.
