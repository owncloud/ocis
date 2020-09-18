Bugfix: Fix Module Path

The module version must be in the path. See https://github.com/golang/go/wiki/Modules#semantic-import-versioning for more information.
> If the module is version v2 or higher, the major version of the module must be included as a /vN at the end of the module paths used in go.mod files (e.g., module github.com/my/mod/v2, require github.com/my/mod/v2 v2.0.1) and in the package import path (e.g., import "github.com/my/mod/v2/mypkg"). This includes the paths used in go get commands (e.g., go get github.com/my/mod/v2@v2.0.1. Note there is both a /v2 and a @v2.0.1 in that example. One way to think about it is that the module name now includes the /v2, so include /v2 whenever you are using the module name).

https://github.com/owncloud/ocis-pkg/pull/25
