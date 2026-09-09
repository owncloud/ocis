Bugfix: Remove the unmarshalable request body field from graph log messages

We've removed the `body` field from the log messages of the graph service's HTTP
handlers. The field was set from `http.Request.Body`, an `io.ReadCloser`, which
can never be serialized into a useful log value.

The tracing middleware replaces the request body with a wrapper that carries an
exported function field, so serializing it failed outright and the log line was
emitted with `"body": "marshaling error: json: unsupported type: func(int64)"`
instead of the payload. Without the tracing middleware the field serialized to an
empty object. In both cases the intended payload was never logged.

https://github.com/owncloud/ocis/pull/12915