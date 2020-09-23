* * *

title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: <https://github.com/owncloud/ocis-accounts>
geekdocEditPath: edit/master/docs

## geekdocFilePath: getting-started.md

{{&lt; toc >}}

## Installation

So far we are offering two different variants for the installation. You can choose between [Docker](https://www.docker.com/) or pre-built binaries which are stored on our download mirrors and GitHub releases. Maybe we will also provide system packages for the major distributions later if we see the need for it.

### Docker

TBD

### Binaries

TBD

## Configuration

We provide overall three different variants of configuration. The variant based on environment variables and commandline flags are split up into global values and command-specific values.

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

#### Server

OCIS_ACCOUNTS_NAME
: Name of the accounts service. It will be part of the namespace.

OCIS_ACCOUNTS_NAMESPACE
: Namespace of the accounts service.

OCIS_ACCOUNTS_ADDRESS
: Endpoint for the grpc service endpoint.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below.

### Configuration file

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/accounts/tree/master/pkg/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/accounts.yml`, `${HOME}/.ocis/accounts.yml` or `$(pwd)/config/accounts.yml`.

## Usage

The program provides a few sub-commands on execution. The available configuration methods have already been mentioned above. Generally you can always see a formated help output if you execute the binary via `ocis-accounts --help`.

### Server

The server command is used to start the grpc server. For further help please execute:

{{&lt; highlight txt >}}
ocis-accounts server --help
{{&lt; / highlight >}}
