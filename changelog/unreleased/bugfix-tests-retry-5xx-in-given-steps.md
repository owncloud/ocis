Bugfix: Retry 5xx responses in Given steps after oCIS restart

When a Background step reconfigures oCIS (e.g. enabling async upload with
delayed postprocessing), oCIS restarts. The readiness probe passes but the
storage layer may not yet be fully initialised, so the first WebDAV request
in the subsequent Given step (such as folder creation) returns 502 Bad
Gateway. This caused `apiAsyncUpload/delayPostprocessing.feature` to fail
spuriously in the `natsJwt-events-stores-file` suite.

The existing 5xx retry in `HttpRequestHelper::sendRequest` was guarded behind
`OcisConfigHelper::isK8s()` and did not apply to Docker/ociswrapper
environments. The fix extends the retry to any Given step that receives a 5xx
response (except 507 Insufficient Storage) regardless of deployment type.

