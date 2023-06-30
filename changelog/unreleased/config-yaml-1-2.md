Change: YAML configuration files are restricted to yaml-1.2

For parsing YAML based configuration files we utilize the gookit/config module.
That module has dropped support for older variants of the YAML format. It now
only supports the YAML 1.2 syntax.
If you're using yaml configuration files, please make sure to update your files
accordingly. The most significant change likely is that only the string `true`
and `false` (including `TRUE`,`True`, `FALSE` and `False`) are now parsed as
booleans. `Yes`, `On` and other values are not longer considered valid values
for booleans.

https://github.com/owncloud/ocis/issues/6510
https://github.com/owncloud/ocis/pull/6493
