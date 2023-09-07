# OCS Service

The `ocs` service (open collaboration services) serves one purpose: it has an endpoint for signing keys which the web frontend accesses when uploading data.

## Signing-Keys Endpoint

The `ocs` service contains an endpoint `/cloud/user/signing-key` on which a user can GET a signing key. Note, this functionality might be deprecated or moved in the future.
