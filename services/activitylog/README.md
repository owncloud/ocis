# Activitylog

The `activitylog` service is responsible for storing events (activities) per resource.

## The Log Service Ecosystem

Log services like the `activitylog`, `userlog`, `clientlog` and `sse` are responsible for composing notifications for a specific audience.
  -   The `userlog` service translates and adjusts messages to be human readable.
  -   The `clientlog` service composes machine readable messages, so clients can act without the need to query the server.
  -   The `sse` service is only responsible for sending these messages. It does not care about their form or language.
  -   The `activitylog` service stores events per resource. These can be retrieved to show item activities

## Activitylog Store

The `activitylog` stores activities for each resource. It works in conjunction with the `eventhistory` service to keep the data it needs to store to a minimum.

## Translations

The `activitylog` service has embedded translations sourced via transifex to provide a basic set of translated languages. These embedded translations are available for all deployment scenarios. In addition, the service supports custom translations, though it is currently not possible to just add custom translations to embedded ones. If custom translations are configured, the embedded ones are not used. To configure custom translations, the `ACTIVITYLOG_TRANSLATION_PATH` environment variable needs to point to a base folder that will contain the translation files. This path must be available from all instances of the activitylog service, a shared storage is recommended. Translation files must be of type  [.po](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html#PO-Files) or [.mo](https://www.gnu.org/software/gettext/manual/html_node/Binaries.html). For each language, the filename needs to be `activitylog.po` (or `activitylog.mo`) and stored in a folder structure defining the language code. In general the path/name pattern for a translation file needs to be:

```text
{ACTIVITYLOG_TRANSLATION_PATH}/{language-code}/LC_MESSAGES/activitylog.po
```

The language code pattern is composed of `language[_territory]` where  `language` is the base language and `_territory` is optional and defines a country.

For example, for the language `de`, one needs to place the corresponding translation files to `{ACTIVITYLOG_TRANSLATION_PATH}/de_DE/LC_MESSAGES/activitylog.po`.

<!-- also see the notifications readme -->

Important: For the time being, the embedded ownCloud Web frontend only supports the main language code but does not handle any territory. When strings are available in the language code `language_territory`, the web frontend does not see it as it only requests `language`. In consequence, any translations made must exist in the requested `language` to avoid a fallback to the default.

### Translation Rules

*   If a requested language code is not available, the service tries to fall back to the base language if available. For example, if the requested language-code `de_DE` is not available, the service tries to fall back to translations in the `de` folder.
*   If the base language `de` is also not available, the service falls back to the system's default English (`en`),
which is the source of the texts provided by the code.

## Default Language

The default language can be defined via the `OCIS_DEFAULT_LANGUAGE` environment variable. See the `settings` service for a detailed description.
