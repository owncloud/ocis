# Notification

The notification service is responsible for sending emails to users informing them about events that happened. To do this it hooks into the event system and listens for certain events that the users need to be informed about.

## Email Notification Templates

The `notifications` service has embedded email body templates.
The email templates contain placeholders `{{ .Greeting }}`, `{{ .MessageBody }}`, `{{ .CallToAction }}` that are
replaced with translations, see the [Translations](#translations) section.
These embedded templates are available for all deployment scenarios. In addition, the service supports custom
templates.
The custom email template takes precedence over the embedded one. If a custom email template exists, the embedded ones
are not used. To configure custom email templates,
the `NOTIFICATIONS_EMAIL_TEMPLATE_PATH` environment variable needs to point to a base folder that will contain the email
templates. The source template files provided by ocis are located
in [https://github.com/owncloud/ocis/tree/master/services/notifications/pkg/email/templates](https://github.com/owncloud/ocis/tree/master/services/notifications/pkg/email/templates) in the `shares`
and `spaces` subfolders:
[shares/shareCreated.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/shares/shareCreated.email.body.tmpl)
[shares/shareExpired.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/shares/shareExpired.email.body.tmpl)
[spaces/membershipExpired.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/spaces/membershipExpired.email.body.tmpl)
[spaces/sharedSpace.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/spaces/sharedSpace.email.body.tmpl)
[spaces/unsharedSpace.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/spaces/unsharedSpace.email.body.tmpl)

Custom Email templates referenced via `NOTIFICATIONS_EMAIL_TEMPLATE_PATH` must be located in subfolders `shares`
and `spaces` and have the same names as the embedded templates. This naming must match the embedded ones.
```text
templates
│
└───shares
│   │   shareCreated.email.body.tmpl
│   │   shareExpired.email.body.tmpl
│
└───spaces
    │   membershipExpired.email.body.tmpl
    │   sharedSpace.email.body.tmpl
    │   unsharedSpace.email.body.tmpl
```

## Translations

The `notifications` service has embedded translations sourced via transifex to provide a basic set of translated languages.
These embedded translations are available for all deployment scenarios. In addition, the service supports custom
translations, though it is currently not possible to just add custom translations to embedded ones. If custom
translations are configured, the embedded ones are not used. To configure custom translations,
the `NOTIFICATIONS_TRANSLATION_PATH` environment variable needs to point to a base folder that will contain the translation
files. This path must be available from all instances of the notifications service, a shared storage is recommended.
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
