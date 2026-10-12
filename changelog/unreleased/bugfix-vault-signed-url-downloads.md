Bugfix: Fix file and folder downloads in vault mode

Downloading a file in the vault failed with a 403 error. A download uses a
signed link without a login token. Thus the proxy uses the stored MFA status of
the user. The proxy did not store this status. Now the proxy stores it after it
knows the user. A download is accepted if the user had a request with MFA within
the MFA session duration.

Downloading a folder or several files in the vault failed with a 404 error. In
vault mode, the web UI requested the signing key from a wrong address. Now the
web UI uses the correct address.

https://github.com/owncloud/ocis/pull/13097
