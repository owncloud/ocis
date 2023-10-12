# Settings

The `settings` service provides functionality for other services to register new settings as well as storing and retrieving the respective settings' values.

## Settings Managed

The settings service is currently used for managing the:

*   users' `profile` settings like the language and the email notification settings,
*   possible user roles and their respective permissions,
*   assignment of roles to users.

As an example, user profile settings that can be changed in the Web UI must be persistent.

The settings service supports two different backends for persisting the data. The backend can be set via the `SETTINGS_STORE_TYPE` environment variable. Supported values are:

*   `metadata`: The default. This backend persists the settings data via the `storage-system` service.
*   `filesystem`: This backend persists the settings data in a directory on the local filesystem.
  The directory can be configured with `SETTINGS_DATA_PATH`. This backend is **not** suitable for running
  multiple intances of the `settings` service in a scale-out deployment and should be therefore considered
  deprecated.

<!--- Note: The diagramm is outdate, leaving it here for a future rework
The diagram shows how the settings service integrates into oCIS:

The diagram shows how the settings service integrates into oCIS:

```mermaid
graph TD
    ows ---|"listSettingsBundles(),<br>saveSettingsValue(value)"| os[ocis-settings]
    owc ---|"listSettingsValues()"| sdk[oC SDK]
    sdk --- sdks{ocis-settings<br>available?}
    sdks ---|"yes"| os
    sdks ---|"no"| defaults[Use set of<br>default values]
    oa[oCIS services<br>e.g. ocis-accounts] ---|"saveSettingsBundle(bundle)"| os
```
-->

## Caching

When using `SETTINGS_STORE_TYPE=metadata`, the `settings` service caches the results of queries against the storage backend to provide faster responses. The content of this cache is independent of the cache used in the `storage-system` service as it caches directory listing and settings content stored in files.

The store used for the cache can be configured using the `SETTINGS_CACHE_STORE` environment variable. Possible stores are:

  -   `memory`: Basic in-memory store and the default.
  -   `redis`: Stores metadata in a configured Redis cluster.
  -   `redis-sentinel`: Stores metadata in a configured Redis Sentinel cluster.
  -   `etcd`: Stores metadata in a configured etcd cluster.
  -   `nats-js`: Stores metadata using the key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in production environments.

1.  Note that in-memory stores are by nature not reboot-persistent.
2.  Though usually not necessary, a database name can be configured for event
    stores if the event store supports this. Generally not applicable for
    stores of type `in-memory`, `redis` and `redis-sentinel`. These settings
    are blank by default which means that the standard settings of the
    configured store apply.
3.  The `settings` service can be scaled if not using `in-memory` stores and
    the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via
    `SETTINGS_CACHE_STORE_NODES` in the form of
    `<sentinel-host>:<sentinel-port>/<redis-master>` like
    `10.10.0.200:26379/mymaster`.

## Settings Management

Infinite Scale services can register *settings bundles* with the settings service.

## Settings Usage

Services can set or query Infinite Scale *setting values* of a user from settings bundles.

## Service Accounts

The settings service needs to know the ID's of service accounts but it doesn't need their secrets. Currently only one service account can be configured which has the admin role. This can be set with the `SETTINGS_SERVICE_ACCOUNT_ID_ADMIN` envvar, but it will also pick up the global `OCIS_SERVICE_ACCOUNT_ID` environment variable. Also see the 'auth-service' service description for additional details.

## Default Language

The default language can be defined via `OCIS_DEFAULT_LANGUAGE` environment variable. If this variable is not defined, English will be used as default. The value has the ISO 639-1 format ("de", "en", etc.) and is limited by the list supported languages. This setting can be used to set the default language for notification and invitation emails.

Important developer note: the list of supported languages is at the moment not easy defineable, as it is the minimum intersection of languages shown in the WebUI and languages defined in the ocis code for the use of notifications. Even more, not all languages where there are translations available on transifex, are available in the WebUI respectively for ocis notifications, and the translation rate for existing languages is partially not that high. You will see therefore quite often English default strings though a supported language may exist and was selected.
