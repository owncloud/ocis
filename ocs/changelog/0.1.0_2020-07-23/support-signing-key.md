Enhancement: Support signing key

We added support for the `/v[12].php/cloud/user/signing-key` endpoint that is used by the owncloud-sdk to generate signed URLs. This allows directly downloading large files with browsers instead of using `blob://` urls, which eats memory ...

<https://github.com/owncloud/ocis/ocs/pull/18>
<https://github.com/owncloud/ocis-proxy/pull/75>
<https://github.com/owncloud/owncloud-sdk/pull/504>
