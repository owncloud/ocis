---
title: Notification
date: 2023-04-20T12:03:29.962855112Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/notifications
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The notification service is responsible for sending emails to users informing them about events that happened. To do this, it hooks into the event system and listens for certain events that the users need to be informed about.

## Table of Contents

* [Email Notification Templates](#email-notification-templates)
* [Translations](#translations)
  * [Translation Rules](#translation-rules)
* [Example Yaml Config](#example-yaml-config)

## Email Notification Templates

The `notifications` service has embedded email body templates. Email templates can use the placeholders `{{ .Greeting }}`, `{{ .MessageBody }}` and `{{ .CallToAction }}` which are replaced with translations when sent, see the [Translations](#translations) section for more details. Depending on the email purpose, placeholders will contain different strings. An individual translatable string is available for each purpose, finally resolved by the placeholder. Though the email subject is also part of translations, it has no placeholder as it is a mandatory email component. The embedded templates are available for all deployment scenarios.
```text
template 
  placeholders
    translated strings <-- source strings <-- purpose
final output
```
In addition, the notifications service supports custom templates. Custom email templates take precedence over the embedded ones. If a custom email template exists, the embedded templates are not used. To configure custom email templates, the `NOTIFICATIONS_EMAIL_TEMPLATE_PATH` environment variable needs to point to a base folder that will contain the email templates. This path must be available from all instances of the notifications service, a shared storage is recommended. The source templates provided by ocis you can derive from are located in following base folder [https://github.com/owncloud/ocis/tree/master/services/notifications/pkg/email/templates](https://github.com/owncloud/ocis/tree/master/services/notifications/pkg/email/templates) with subfolders `shares` and `spaces`.
-   [shares/shareCreated.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/shares/shareCreated.email.body.tmpl)
-   [shares/shareExpired.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/shares/shareExpired.email.body.tmpl)
-   [spaces/membershipExpired.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/spaces/membershipExpired.email.body.tmpl)
-   [spaces/sharedSpace.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/spaces/sharedSpace.email.body.tmpl)
-   [spaces/unsharedSpace.email.body.tmpl](https://github.com/owncloud/ocis/blob/master/services/notifications/pkg/email/templates/spaces/unsharedSpace.email.body.tmpl)
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
Custom email templates referenced via `NOTIFICATIONS_EMAIL_TEMPLATE_PATH` must also be located in subfolders `shares` and `spaces` and must have the same names as the embedded templates. It is important that the names of these files and  folders match the embedded ones.

## Translations

The `notifications` service has embedded translations sourced via transifex to provide a basic set of translated languages. These embedded translations are available for all deployment scenarios.
In addition, the service supports custom translations, though it is currently not possible to just add custom translations to embedded ones. If custom translations are configured, the embedded ones are not used. To configure custom translations,
the `NOTIFICATIONS_TRANSLATION_PATH` environment variable needs to point to a base folder that will contain the translation files. This path must be available from all instances of the notifications service, a shared storage is recommended. Translation files must be of type  [.po](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html#PO-Files) or [.mo](https://www.gnu.org/software/gettext/manual/html_node/Binaries.html). For each language, the filename needs to be `translations.po` (or `translations.mo`) and stored in a folder structure defining the language code. In general the path/name pattern for a translation file needs to be:
```text
{NOTIFICATIONS_TRANSLATION_PATH}/{language-code}/LC_MESSAGES/translations.po
```
The language code pattern is composed of `language[_territory]` where  `language` is the base language and `_territory` is optional and defines a country.
For example, for the language `de_DE`, one needs to place the corresponding translation files to `{NOTIFICATIONS_TRANSLATION_PATH}/de_DE/LC_MESSAGES/translations.po`.

### Translation Rules

*   If a requested language code is not available, the service tries to fall back to the base language if available.
For example, if `de_DE` is not available, the service tries to fall back to translations in the `de` folder.
*   If the base language `de` is also not available, the service falls back to the system's default English (`en`),
which is the source of the texts provided by the code.

## Example Yaml Config

{{< include file="services/_includes/notifications-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/notifications_configvars.md" >}}

