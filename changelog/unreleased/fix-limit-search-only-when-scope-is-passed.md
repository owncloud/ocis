Bugfix: Limit search only when scope is passed

Previously, the search service would limit the search to the according space when searching `/dav/spaces/`.
This was not correct, as the search should be limited to the according space when a `scope` is passed in the search pattern instead.

https://github.com/owncloud/ocis/pull/11664
