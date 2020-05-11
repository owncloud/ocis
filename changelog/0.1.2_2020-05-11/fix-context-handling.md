Bugfix: Fix usage of context.Context 

The context was filled with a key defined in the package service but read with a key from the package imgsource.
Since `service.key` and `imgsource.key` are different types imgsource could not read the value provided by service.

https://github.com/owncloud/ocis-thumbnails/issues/18
