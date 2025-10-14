---
title: "Glossary"
date: 2020-05-04T12:35:00+01:00
weight: 80
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/settings
geekdocFilePath: glossary.md
---

In the context of this extension and oCIS in general, we are using the following terminology.

### Configuration

- System configuration
- e.g. service host names and ports
- Changes need to be propagated to other services
- Typically modified on the CLI

### Settings

- Application level settings
- e.g. default language
- Can be modified at runtime without restarting the service

### Preferences

- User settings
- Subset of "Settings"
- e.g. preferred language of a user

### Settings Bundle

- Collection of related settings
- Registered by an oCIS extension

### Settings Value

- Manifestation of a setting for a specific user
- E.g. used for customization (at runtime) in `ocis-web`
- Can be queried and modified by other oCIS services
