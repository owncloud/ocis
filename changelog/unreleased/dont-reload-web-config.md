Bugfix: Don't reload web config

When requesting `config.json` file from the server, web service would reload the file if a path is set. This will remove config entries set via Envvar. Since we want to have the possiblity to set configuration from both sources we removed the reading from file. The file will still be loaded on service startup.

https://github.com/owncloud/ocis/pull/7369
