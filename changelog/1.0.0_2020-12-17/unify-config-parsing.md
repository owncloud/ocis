Change: Unify Configuration Parsing

Tags: ocis

- responsibility for config parsing should be on the subcommand
- if there is a config file in the environment location, env var should take precedence
- general rule of thumb: the more explicit the config file is that would be picked up. Order from less to more explicit:
    - config location (/etc/ocis)
    - environment variable
    - cli flag

https://github.com/owncloud/ocis/pull/675
