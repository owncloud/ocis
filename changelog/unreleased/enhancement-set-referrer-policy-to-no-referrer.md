Enhancement: Set Referrer-Policy to no-referrer

Change the Referrer-Policy from 'strict-origin-when-cross-origin'
to 'no-referrer' to enhance user privacy and security.

Previously, the origin was sent on cross-origin requests. This change
completely removes the Referrer header from all outgoing requests,
preventing any potential leakage of browsing information to third parties.
This is a more robust approach to protecting user privacy.

https://github.com/owncloud/ocis/pull/11722
