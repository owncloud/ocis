Security: Bump go-jose and ignore CVE-2026-116587 for go-jose/v3

We've bumped go-jose/v4 to v4.1.6 to fix CVE-2026-116587, a denial of service
when parsing JWKs that contain RSA keys.

The indirect dependency go-jose/v3 has no fix for CVE-2026-116587. It comes
from libregraph/lico. Lico reads JWKs only from local files and from
configured OIDC authorities. We've fixed go-jose/v4 and we ignore the go-jose/v3
finding in the Trivy scan.

https://github.com/owncloud/ocis/pull/13341
