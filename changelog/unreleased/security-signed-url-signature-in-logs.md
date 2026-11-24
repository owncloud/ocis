Security: Stop logging signed-URL signatures in the proxy

When a pre-signed URL failed verification, the proxy logged, at error level,
the signature it expected for that URL together with the full request URL.
A logged URL could still carry a valid signature, for example when only the
HTTP method did not match. Anyone who could read the proxy logs could reuse
these signatures to send signed requests as the user named in the URL.

The proxy now logs only the request path and a generic "signature mismatch",
and compares signatures in constant time.

https://github.com/owncloud/ocis/pull/13569
