* * *

title: "Configuration"
date: "2020-07-01T11:10:52+0200"
weight: 20
geekdocRepo: <https://github.com/owncloud/ocis-reva>
geekdocEditPath: edit/master/docs

## geekdocFilePath: configuration.md

{{&lt; toc >}}

## Configuration

oCIS Single Binary is not responsible for configuring extensions. Instead, each extension could either be configured by environment variables, cli flags or config files.

Each extension has its dedicated documentation page (e.g. <https://owncloud.github.io/extensions/ocis_proxy/configuration>) which lists all possible configurations. Config files and environment variables are picked up if you use the `./bin/ocis server` command within the oCIS single binary. Command line flags must be set explicitly on the extensions subcommands.

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. _i.e: ocis-proxy reads `proxy.json | yaml | toml ...`_.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

### ocis-reva server

Start ocis accounts service

Usage: `ocis-reva server [command options] [arguments...]`

\--http-namespace | $ACCOUNTS_HTTP_NAMESPACE  
: Set the base namespace for the http namespace. Default: `com.owncloud.web`.

\--http-addr | $ACCOUNTS_HTTP_ADDR  
: Address to bind http server. Default: `localhost:9181`.

\--http-root | $ACCOUNTS_HTTP_ROOT  
: Root path of http server. Default: `/`.

\--grpc-namespace | $ACCOUNTS_GRPC_NAMESPACE  
: Set the base namespace for the grpc namespace. Default: `com.owncloud.api`.

\--grpc-addr | $ACCOUNTS_GRPC_ADDR  
: Address to bind grpc server. Default: `localhost:9180`.

\--name | $ACCOUNTS_NAME  
: service name. Default: `accounts`.

\--accounts-data-path | $ACCOUNTS_DATA_PATH  
: accounts folder. Default: `/var/tmp/ocis-accounts`.

\--asset-path | $HELLO_ASSET_PATH  
: Path to custom assets.

### ocis-reva ocis-accounts

Provide accounts and groups for oCIS

Usage: `ocis-reva ocis-accounts [command options] [arguments...]`

\--log-level | $ACCOUNTS_LOG_LEVEL  
: Set logging level. Default: `info`.

\--log-pretty | $ACCOUNTS_LOG_PRETTY  
: Enable pretty logging. Default: `true`.

\--log-color | $ACCOUNTS_LOG_COLOR  
: Enable colored logging. Default: `true`.
