Security: Update the saml and ws dependencies

We updated two indirect Go dependencies to fix known vulnerabilities.
The module github.com/crewjam/saml is now v0.5.2.
The module github.com/gobwas/ws is now v1.5.0.
The module github.com/russellhaering/goxmldsig is now v1.6.1, as the new saml
version requires it.

https://github.com/owncloud/ocis/pull/13575
