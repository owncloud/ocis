Security: Add X-XSS-Protection header

Added the X-XSS-Protection header set to "0" to explicitly disable the
deprecated browser XSS filter, which can introduce side-channel
vulnerabilities. Modern XSS protection is provided through the
Content-Security-Policy header.

This change addresses security audit findings requiring explicit
configuration of HTTP security headers per OWASP recommendations.

https://github.com/owncloud/ocis/pull/12092
