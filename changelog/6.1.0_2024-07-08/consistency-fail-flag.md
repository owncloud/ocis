Enhancement: Add fail flag to consistency check

We added a `--fail` flag to the `ocis backup consistency` command. If set to true, the command will return a non-zero exit code if any inconsistencies are found. This allows you to use the command in scripts and CI/CD pipelines to ensure that backups are consistent.

https://github.com/owncloud/ocis/pull/9447
