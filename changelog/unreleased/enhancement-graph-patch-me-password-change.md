Security: Block password and account changes via PATCH /graph/v1.0/me

`PATCH /graph/v1.0/me` no longer accepts `passwordProfile`, `accountEnabled`, or `onPremisesSamAccountName`. Use `POST /graph/v1.0/me/changePassword` to change your password.

https://github.com/owncloud/ocis/pull/12493