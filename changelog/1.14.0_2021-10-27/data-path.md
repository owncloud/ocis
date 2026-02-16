Change: New default data paths and easier configuration of the data path

We've changed the default data path for our release artifacts:
- oCIS docker images will now store all data in `/var/lib/ocis` instead in `/var/tmp/ocis`
- binary releases will now store all data in `~/.ocis` instead of `/var/tmp/ocis`

Also if you're a developer and you run oCIS from source, it will store all data in `~/.ocis` from now on.

You can now easily change the data path for all extensions by setting the environment variable `OCIS_BASE_DATA_PATH`.

If you want to package oCIS, you also can set the default data path at compile time, eg. by passing `-X "github.com/owncloud/ocis/ocis-pkg/config/defaults.BaseDataPathType=path" -X "github.com/owncloud/ocis/ocis-pkg/config/defaults.BaseDataPathValue=/var/lib/ocis"` to your go build step.

https://github.com/owncloud/ocis/pull/2590
