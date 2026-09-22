Bugfix: Stop endless redirect loop when vault MFA step-up fails

Web requested the vault MFA level via `acr_values`, which is a voluntary claim.
IdPs like Keycloak return a token with a lower `acr` instead of an error when the
user can't complete the second factor (none enrolled, device not at hand). Web
then started the step-up again on every return, trapping the user in a redirect
loop. We now remember a pending step-up for the current tab and, if the IdP
returns without the required `acr`, stop retrying, send the user back to the
default view and show an error message.

https://github.com/owncloud/ocis/issue/12984
