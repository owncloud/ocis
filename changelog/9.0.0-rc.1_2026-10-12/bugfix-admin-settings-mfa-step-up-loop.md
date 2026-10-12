Bugfix: Stop endless redirect loop when admin settings MFA step-up fails

The admin settings route guards required MFA via `requireAcr()`, which
redirected to the IdP whenever the token's `acr` did not match. IdPs like
Keycloak answer an `acr_values` request they can't satisfy with a lower `acr`
instead of an error, so a user without a second factor at hand was trapped in a
redirect loop. `requireAcr()` now uses the same protection as the vault: the
step-up carries its target as OIDC state, and if the IdP answers it without the
required `acr`, `requireAcr()` returns false. The admin settings then redirect
to the default view and show an error message.

https://github.com/owncloud/ocis/issues/13069
https://github.com/owncloud/ocis/pull/13070
