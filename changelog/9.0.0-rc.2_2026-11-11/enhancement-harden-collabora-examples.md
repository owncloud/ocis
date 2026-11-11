Enhancement: Harden the Collabora deployment examples

The `ocis_full` and `ocis_ocm` deployment examples now ship a hardened
Collabora configuration. Collabora is restricted to the hosts it needs, runs on
its own internal network that only the proxy and the collaboration service
share, and the admin console requires explicitly configured credentials
(`COLLABORA_ADMIN_USER` and `COLLABORA_ADMIN_PASSWORD`). If you copy these
examples, make sure to set the credentials in your `.env` file.

https://github.com/owncloud/ocis/pull/13395
