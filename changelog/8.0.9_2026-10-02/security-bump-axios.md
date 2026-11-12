Security: Bump axios to 1.20.0

We've updated the axios dependency of the idp service to 1.20.0. This fixes
several denial-of-service, request-smuggling and server-side request forgery
vulnerabilities flagged by the release filesystem scan.

https://github.com/owncloud/ocis/pull/13056
