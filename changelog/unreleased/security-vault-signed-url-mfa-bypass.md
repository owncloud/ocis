Security: Require MFA-bound signing key for Safe (vault) signed URLs

With the Safe and multi-factor authentication enabled, the proxy remembered a
user's proven second factor per user rather than per session and applied it to
any request that carried no login token. A signed URL carries no token, so a
session that had only passed the first factor could fetch the user's URL
signing key and send signed requests that the server treated as if the second
factor had been given, letting it list and download the user's Safe while any
other session of that user had MFA active.

The signing-key endpoint now issues a separate, MFA-bound signing key only to
requests that already carry MFA, and the proxy grants the MFA status to a
signed URL only when the URL was signed with that MFA-bound key and the user's
MFA session is still valid. First-factor-only sessions, basic auth and app
tokens no longer inherit MFA. Signed downloads in the regular Drive are
unaffected.

https://github.com/owncloud/ocis/pull/13485
