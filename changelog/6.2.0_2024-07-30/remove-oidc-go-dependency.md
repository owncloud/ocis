Enhancement: Remove oidc-go dependency

Removes the kgol/oidc-go dependency because it was flagged by dependabot. Luckily us we only used it for importing the strings "profile" and "email".

https://github.com/owncloud/ocis/pull/9641
