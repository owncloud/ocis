Enhancement: Add RED metrics to the metrics endpoint

We added three new metrics to the metrics endpoint to support the RED method for monitoring microservices.

- Request Rate: The number of requests per second. The total count of requests is available under `ocis_proxy_requests_total`.
- Error Rate: The number of failed requests per second. The total count of failed requests is available under `ocis_proxy_errors_total`.
- Duration: The amount of time each request takes. The duration of all requests is available under `ocis_proxy_request_duration_seconds`. This is a histogram metric, so it also provides information about the distribution of request durations.

The metrics are available under the following paths: `PROXY_DEBUG_ADDR/metrics` in a prometheus compatible format and maybe secured by `PROXY_DEBUG_TOKEN`.

https://github.com/owncloud/ocis/pull/7994
