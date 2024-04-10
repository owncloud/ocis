Enhancement: Make IDP cookies same site strict

To enhance the security of our application and prevent Cross-Site Request Forgery (CSRF) attacks, we have updated the
SameSite attribute of the build in Identity Provider (IDP) cookies to Strict.

This change restricts the browser from sending these cookies with any cross-site requests,
thereby limiting the exposure of the user's session to potential threats.

This update does not impact the existing functionality of the application but provides an additional layer of security
where needed.

https://github.com/owncloud/ocis/pull/8716
