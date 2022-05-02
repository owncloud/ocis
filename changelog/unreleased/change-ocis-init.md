Change: Introduce `ocis init` and remove all default secrets

We've removed all default secrets. This means you can't start oCIS any longer
without setting these via environment variable or configuration file.

In order to make this easy for you, we introduced a new command: `ocis init`.
You can run this command before starting oCIS with `ocis server` and it will
bootstrap you a configuration file for a secure oCIS instance.

https://github.com/owncloud/ocis/pull/3551
