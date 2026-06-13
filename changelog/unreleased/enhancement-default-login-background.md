Enhancement: Ship ownCloud branding background as default login background

The IDP login page now falls back to the builtin ownCloud branding background
image (served from "/signin/v1/static/owncloud-background.jpg") when
IDP_LOGIN_BACKGROUND_URL is not set. Operators who configure
IDP_LOGIN_BACKGROUND_URL explicitly are unaffected. The background image can
also be overridden on the filesystem via IDP_ASSET_PATH.

The same branding background is now also embedded in the web service theme
assets (services/web/assets/themes/owncloud/assets/loginBackground.jpg), so
the ownCloud Web login page shows it out of the box, taking precedence over
the background bundled with the vendored web release.

Additionally the ocis_keycloak deployment example now ships a custom Keycloak
login theme ("owncloud") that applies the same branding background and styles
the login form to match the IDP login screen.

https://github.com/owncloud/ocis/pull/12406
