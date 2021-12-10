Change: OIDC: fallback if IDP doesn't provide "preferred_username" claim

Some IDPs don't add the "preferred_username" claim. Fallback to the "email"
claim in that case

https://github.com/owncloud/ocis/issues/2644
