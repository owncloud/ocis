Security: Harden origin validation for embed mode delegated authentication

Hardened origin validation for the delegated authentication feature in Web's
embed mode. Deployments using `WEB_OPTION_EMBED_DELEGATE_AUTHENTICATION` should
ensure `WEB_OPTION_EMBED_DELEGATE_AUTHENTICATION_ORIGIN` is explicitly
configured.

https://github.com/owncloud/ocis/pull/12572
