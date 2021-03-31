Enhancement: Runtime Hostname and Port are now configurable

Without any configuration the ocis runtime will start on `localhost:9250` unless specified otherwise. Usage:

- `OCIS_RUNTIME_PORT=6061 bin/ocis server`
  - overrides the oCIS runtime and starts on port 6061
- `OCIS_RUNTIME_PORT=6061 bin/ocis list`
  - lists running extensions for the runtime on `localhost:6061`

All subcommands are updated and expected to work with the following environment variables:

```
OCIS_RUNTIME_HOST
OCIS_RUNTIME_PORT
```

https://github.com/owncloud/ocis/pull/1822
