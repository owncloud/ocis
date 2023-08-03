# Userlog

The `userlog` service is a mediator between the `eventhistory` service and clients who want to be informed about user related events. It provides an API to retrieve those.

## Prerequisites

Running the `userlog` service without running the `eventhistory` service is not possible.

## Storing

The `userlog` service persists information via the configured store in `USERLOG_STORE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured Redis cluster.
  -   `redis-sentinel`: Stores data in a configured Redis Sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in production environments.

1.  Note that in-memory stores are by nature not reboot-persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store apply.
3.  The userlog service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `USERLOG_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.

## Configuring

For the time being, the configuration which user related events are of interest is hardcoded and cannot be changed.

## Retrieving

The `userlog` service provides an API to retrieve configured events. For now, this API is mostly following the [oc10 notification GET API](https://doc.owncloud.com/server/next/developer_manual/core/apis/ocs-notification-endpoint-v1.html#get-user-notifications).

## Subscribing

Additionally to the oc10 API, the `userlog` service also provides an `/sse` (Server-Sent Events) endpoint to be informed by the server when an event happens. See [What is Server-Sent Events](https://medium.com/yemeksepeti-teknoloji/what-is-server-sent-events-sse-and-how-to-implement-it-904938bffd73) for a simple introduction and examples of server sent events. The `sse` endpoint will respect language changes of the user without needing to reconnect. Note that SSE has a limitation of six open connections per browser which can be reached if one has opened various tabs of the Web UI pointing to the same Infinite Scale instance.

## Posting

The userlog service is able to store global messages that will be displayed in the Web UI to all users. If a user deletes the message in the Web UI, it reappears on reload. Global messages use the endpoint `/ocs/v2.php/apps/notifications/api/v1/notifications/global` and are activated by sending a `POST` request. Note that sending another `POST` request of the same type overwrites the previous one. For the time being, only the type `deprovision` is supported.

### Authentication

`POST` and `DELETE` endpoints provide notifications to all users. Therefore only certain users can configure them. Two authentication methods for this endpoint are provided. Users with the `admin` role can always access these endpoints. Additionally, a static secret via the `USERLOG_GLOBAL_NOTIFICATIONS_SECRET` can be defined to enable access for users knowing this secret, which has to be sent with the header containing the request.

### Deprovisioning

Deprovision messages announce a deprovision text including a deprovision date of the instance to all users. With this message, users get informed that the instance will be shut down and deprovisioned and no further access to their data is possible past the given date. This implies that users must download their data before the given date. The text shown to users refers to this information. Note that the task to deprovision the instance does not depend on the message. The text of the message can be translated according to the translation settings, see section [Translations](#translations). The endpoint only expects a `deprovision_date` parameter in the `POST` request body as the final text is assembled automatically. The string hast to be in `RFC3339` format, however, this format can be changed by using `deprovision_date_format`. See the [go time formating](https://pkg.go.dev/time#pkg-constants) for more details.

## Deleting

To delete events for an user, use a `DELETE` request to `ocs/v2.php/apps/notifications/api/v1/notifications` containing the IDs to delete.

Sending a `DELETE` request to the `ocs/v2.php/apps/notifications/api/v1/notifications/global` endpoint to remove a global message is a restricted action, see the [Authentication](#authentication) section for more details.)

## Translations

The `userlog` service has embedded translations sourced via transifex to provide a basic set of translated languages. These embedded translations are available for all deployment scenarios. In addition, the service supports custom translations, though it is currently not possible to just add custom translations to embedded ones. If custom translations are configured, the embedded ones are not used. To configure custom translations, the `USERLOG_TRANSLATION_PATH` environment variable needs to point to a base folder that will contain the translation files. This path must be available from all instances of the userlog service, a shared storage is recommended. Translation files must be of type  [.po](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html#PO-Files) or [.mo](https://www.gnu.org/software/gettext/manual/html_node/Binaries.html). For each language, the filename needs to be `userlog.po` (or `userlog.mo`) and stored in a folder structure defining the language code. In general the path/name pattern for a translation file needs to be:

```text
{USERLOG_TRANSLATION_PATH}/{language-code}/LC_MESSAGES/userlog.po
```

The language code pattern is composed of `language[_territory]` where  `language` is the base language and `_territory` is optional and defines a country.

For example, for the language `de`, one needs to place the corresponding translation files to `{USERLOG_TRANSLATION_PATH}/de_DE/LC_MESSAGES/userlog.po`.

<!-- also see the notifications readme -->

Important: For the time being, the embedded ownCloud Web frontend only supports the main language code but does not handle any territory. When strings are available in the language code `language_territory`, the web frontend does not see it as it only requests `language`. In consequence, any translations made must exist in the requested `language` to avoid a fallback to the default.

### Translation Rules

*   If a requested language code is not available, the service tries to fall back to the base language if available. For example, if the requested language-code `de_DE` is not available, the service tries to fall back to translations in the `de` folder.
*   If the base language `de` is also not available, the service falls back to the system's default English (`en`),
which is the source of the texts provided by the code.
