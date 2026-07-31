Security: Fix signing-key to public share guests

The /ocs/v[12].php/cloud/user/signing-key endpoint was reachable through a
public share session.
The endpoint `public-token` is no longer allowed by the public-share resource scope in reva.

https://github.com/owncloud/ocis/pull/12332
https://github.com/owncloud/reva/pull/608
