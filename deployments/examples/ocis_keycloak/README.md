---
document this deployment example in: docs/ocis/deployment/ocis_keycloak.md
---

Please refer to [our documentation](https://owncloud.dev/ocis/deployment/ocis_keycloak/)
for instructions on how to deploy this scenario.

## Branding / custom login theme

This example ships a custom Keycloak login theme named `owncloud` in
[`themes/owncloud/login/`](themes/owncloud/login/). It applies the ownCloud
branding background as a full-viewport image and styles the login form to
match the builtin oCIS IDP login screen. The theme is:

- mounted into the Keycloak container via `docker-compose.yml`
  (`./themes/owncloud:/opt/keycloak/themes/owncloud:ro`)
- enabled for the imported realm via `"loginTheme": "owncloud"` in
  [`config/keycloak/ocis-realm.dist.json`](config/keycloak/ocis-realm.dist.json)

It inherits everything from the standard `keycloak` login theme
(`parent=keycloak` in `theme.properties`) and only overrides CSS and images.

### Replacing the background image

Overwrite [`themes/owncloud/login/resources/img/owncloud-background.png`](themes/owncloud/login/resources/img/owncloud-background.png)
with your own image (PNG/JPEG/WebP, recommended minimum resolution 1920×1080)
and restart the Keycloak container. If you use a different file name or
format, adjust the `--oc-login-background-image` custom property at the top of
[`themes/owncloud/login/resources/css/owncloud.css`](themes/owncloud/login/resources/css/owncloud.css).
The browser favicon comes from `resources/img/favicon.ico` in the same theme.
There is intentionally no logo above the login card since the background
image already contains the ownCloud logo; to add one back, style
`#kc-header-wrapper` in `owncloud.css`.

Note: Keycloak caches themes. When iterating on the theme, either restart the
container or start Keycloak with `--spi-theme-static-max-age=-1
--spi-theme-cache-themes=false --spi-theme-cache-templates=false`.

### Background of the builtin oCIS IDP (when not using Keycloak)

Deployments that use the builtin oCIS IDP (LibreGraph Connect) instead of
Keycloak show the same branding background by default: if the environment
variable `IDP_LOGIN_BACKGROUND_URL` is unset or empty, the IDP falls back to
the embedded asset served from `/signin/v1/static/owncloud-background.jpg`.
Setting `IDP_LOGIN_BACKGROUND_URL` to any URL overrides the background as
before. The embedded asset itself can be replaced on the filesystem through
`IDP_ASSET_PATH`.
