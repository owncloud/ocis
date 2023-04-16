# Notification

The notification service is responsible for sending emails to users informing them about events that happened. To do this it hooks into the event system and listens for certain events that the users need to be informed about.

## Translations

The `translations` service has embedded translations sourced via transifex to provide a basic set of translated languages.
These embedded translations are available for all deployment scenarios. In addition, the service supports custom
translations, though it is currently not possible to just add custom translations to embedded ones. If custom
translations are configured, the embedded ones are not used. To configure custom translations,
the `NOTIFICATIONS_TRANSLATION_PATH` environment variable needs to point to a base folder that will contain the translation
files. This path must be available from all instances of the translations service, a shared storage is recommended.
Translation files must be of type  [.po](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html#PO-Files)
or [.mo](https://www.gnu.org/software/gettext/manual/html_node/Binaries.html). For each language, the filename needs to
be `translations.po` (or `translations.mo`) and stored in a folder structure defining the language code. In general the path/name
pattern for a translation file needs to be:

```text
{NOTIFICATIONS_TRANSLATION_PATH}/{language-code}/LC_MESSAGES/translations.po
```

The language code pattern is composed of `language[_territory]` where  `language` is the base language and `_territory`
is optional and defines a country.

For example, for the language `de_DE`, one needs to place the corresponding translation files
to `{NOTIFICATIONS_TRANSLATION_PATH}/de_DE/LC_MESSAGES/translations.po`.

### Translation Rules

*   If a requested language code is not available, the service tries to fall back to the base language if available.
For example, if `de_DE` is not available, the service tries to fall back to translations in the `de` folder.
*   If the base language `de` is also not available, the service falls back to the system's default English (`en`),
which is the source of the texts provided by the code.
