Enhancement: Add support for proof keys for the collaboration service

Proof keys support will be enabled by default in order to ensure that all
the requests come from a trusted source.
Since proof keys must be set in the WOPI app (OnlyOffice, Collabora...), it's
possible to disable the verification of the proof keys via configuration.

https://github.com/owncloud/ocis/pull/9366
