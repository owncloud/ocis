Security: Harden public link brute force protection

We've hardened the brute force protection for password protected public links so
that failed password attempts are reliably counted and the rate limit is
consistently enforced.

https://github.com/owncloud/ocis/pull/12945
