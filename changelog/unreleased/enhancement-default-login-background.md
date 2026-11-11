Enhancement: Ship ownCloud branding as the default login background

The builtin IDP login page now ships the ownCloud branding image as its default
background. The bundled image (services/idp/assets/identifier/static/media,
referenced from theme.css) is anchored to the top right so the ownCloud logo
stays fully visible at every viewport size instead of being cropped. Operators
can still point the login background at a different image via
IDP_LOGIN_BACKGROUND_URL.

The same branding background is also embedded in the web service theme assets
(services/web/assets/themes/owncloud/assets/loginBackground.jpg), so the
ownCloud Web login page shows it out of the box, taking precedence over the
background bundled with the vendored web release.

Additionally the ocis_keycloak deployment example now ships a custom Keycloak
login theme ("owncloud") that applies the same branding background and styles
the login form to match the IDP login screen.

https://github.com/owncloud/ocis/pull/12406
