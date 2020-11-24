Bugfix: add rate limiting middleware

Tags: proxy

If too many concurrent requests come in some services can be overwhelmed.
By default it's limited to 100 requests per second which is configurable with the PROXY_RATE_LIMIT env variable.