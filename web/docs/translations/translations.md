---
title: "Working With Translations"
date: 2025-05-21T00:00:00+00:00
weight: 55
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/translations
geekdocFilePath: translations.md
geekdocCollapseSection: true
---

{{< toc >}}

## Translation Setup

The translation setup is managed within each Packages folder that contains a translation source. These directories contain a structure relevant to translation. The structure is shown below using the `web-runtime` folder as an example. When a new app (parent directory) is created, this structure must be applied.

```
l10n			(dir)
  locale		(dir)
  template.pot		(file)
  translations.json	(file)
  .tx			(dir)
    config		(file)
```

The following files are managed via `make` commands in all folders containing translations, see section [Translation Process](#translation-process) for more details.

`make l10n-read` --> `template.pot` and\
`make l10n-write` --> `translations.json`,

The `config` files must be created manually and define the Transifex setup for each folder. Here is an example how this file looks like:

```
[main]
host = https://www.transifex.com

[o:owncloud-org:p:owncloud-web:r:core]
file_filter = locale/<lang>/app.po
minimum_perc = 0
resource_name = web-runtime
source_file = template.pot
source_lang = en
type = PO
```

## Translation Definitions

Translations are defined in 4 levels:

1. The source text such as `$gettext('Administration Settings')`.
2. For each language, the minimum translation completion percentage required for processing.\
This can either be the individual setting in the config file or the global setting behind `make l10n-pull`, the latter takes precedence.
3. Languages that will be processed by the make commands are defined in `/gettext.config.cjs`.
4. Languages shown in the webUI including their user-facing text are defined in\
`/packages/web-runtime/src/defaults/languages.ts`.

Note 1: Languages defined in 3. and 4. should match.\
Note 2: The source string is used for any language that does not meet 2.) **OR** is not listed in 3.)

## Translation Process

A nightly sync automatically extracts and synchronizes data, see the [Add Translations](https://owncloud.dev/services/general-info/add-translations/) documentation for more details, but it can also be triggered manually. Here are the steps for a manual process, run the commands from the repos root:

```
make l10n-read         --> Extract source strings from gettext
make l10n-push         --> Push the translation sources to Transifex
make l10n-pull         --> Get the translations down from Transifex
make l10n-write        --> Apply translations to be processable by the webUI
```

Once the final step has been processed, you can build the webUI, which will include the defined and created translations.

{{< hint warning >}}
'make l10n-push' will only work if you have proper write permissions to Transifex.\
You will get a '403, permission_denied: You do not have permission to perform this action.' otherwise. If you need to push to Transifex manually but lack the proper write permissions, request or trigger a sync which is described via the above link.
{{< /hint >}}

Note that only the command `make l10n-write` will create changes that can be committed, following files are excluded by `.gitignore`:
```
**/l10n/locale
**/l10n/template.pot
```

If you want to start with a clean translation base by removing all extractions and translations, use the following command **before** starting with `make l10n-read`:

```
make l10n-clean        --> Delete all translation template.pot and locale files
``` 
