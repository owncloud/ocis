Enhancement: Add `--diff` to the `ocis init` command

We have added a new flag `--diff` to the `ocis init` command to show the diff of the configuration files.
This is useful to see what has changed in the configuration files when you run the `ocis init` command.
The diff is stored to the ocispath in the config folder as ocis.config.patch and can be applied using the
linux `patch` command.

https://github.com/owncloud/ocis/pull/9693
https://github.com/owncloud/ocis/issues/3645
