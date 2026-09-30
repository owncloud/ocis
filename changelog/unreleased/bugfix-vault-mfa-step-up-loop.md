Bugfix: Stop endless redirect loop when vault MFA step-up fails

Web requested the vault MFA level via `acr_values`, which is a voluntary claim.
IdPs like Keycloak return a token with a lower `acr` instead of an error when the
user can't complete the second factor (none enrolled, device not at hand). Web
then started the step-up again on every return, trapping the user in a redirect
loop. We now remember a pending step-up for the current tab and, if the IdP
returns without the required `acr`, stop retrying, reload the page outside the
vault and show an error message there. Leaving the vault needs a full page load,
because vault mode is fixed for the lifetime of the page.

https://github.com/owncloud/ocis/pull/12985
