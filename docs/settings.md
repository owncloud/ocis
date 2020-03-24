---
title: "Settings"
date: 2020-02-27T20:35:00+01:00
weight: 45
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: settings.md
---



## The Settings service

Extensions can register a settings bundle with the settings service.

A bundle has several properties:
- an id that uniquely identifies the bundle
- a display name that is show to administrators
- the service the bundle belongs to

These bundles consist of a set of settings, each with:
- an id that uniquely identifies it
- a display name that is shown to users
- a description (optional) 
- a placeholder
- a default value (type, min, max, stepping ...)
  - string
  - integer (min, max)
  - checkbox (bool)
  - list
- a permission that is allowed to lock down the value
- a permission that is allowed to override the default
- a permission that is allowed to set the value

### Examples
#### Timezone setting
```json
{ // our bundle
  "id": "d6d74cd9-be10-44cc-91ea-0e0892a0d162",
  "name": "Timezone Settings",
  "extension": {
    "name": "Calendar",
    "id": "448635a7-b145-4455-8aba-341c12c472ae"
  },
  // every bundle has a list of settings
  "settings": [
    {
      "id": "b8e72ade-c963-48c4-8846-bf626f3e0257",
      "name": "Timezone",
      "description": null,
      "scope": "user",
      "placeholder": "Please select a timezone", // TODO needs translation urgh
      "values": {
        "type": "list",
        "options": [
          {"value":0, "label":"unknown"}
          {"value":1, "label":"Europe/Berlin", "default":true}
          {"value":2, "label":"Europe/Amsterdam"}
          ...
        ]
      },
    },
  ]
}
```

Every setting automatically creates three permissions:
- `read` to allow reading the setting value from the settings service
- `write` to allow writing the setting value to the settings service
- `display` to allow seeing or listing a settings value in the ui

A user that has the `read` permission, the `display` permission but no `write` permission could see his timezone in the ui but not change it.
A user that has the `read` permission but no `display` permission should not even see his timezone in the settings ui.

## Roles

Every user has roles that are tied to his account.
Every role has a set of permissions.
Every permission is a boolean flag.
```json
{
  "name": "Admin",
  "permissions": [
    {"id":"u-u-i-dr", "scope":{"type":"user","value":["all"]}},
    {"id":"u-u-i-dw", "scope":{"type":"user","value":["all"]}},
    ...
  ]
},
{
  "name": "Site Admin",
  "permissions": [
    {"id":"u-u-i-dr", "scope":{"type":"group","value":["u-u-i-d"]}},
    {"id":"u-u-i-dw", "scope":{"type":"group","value":["u-u-i-d"]}},
    {"id":"u-u-i-dr", "scope":{"type":"user","value":["me"]}},
    {"id":"u-u-i-dw", "scope":{"type":"user","value":["me"]}},
    ...
  ]
},
{
  "name": "User",
  "permissions": [
    {"id":"u-u-i-dr", "scope":{"type":"user","value":["me"]}},
  ]
}

{
  "id":"u-u-i-dr", 
  "name":"read Timezone"
},
{
  "id":"u-u-i-dw", 
  "name":"write Timezone"
}

```

TODO: passwords
- sometimes we need to store them but only certain extensions should be allowed to read them, eg. the wnd app
- they should be encrypted anyway
- but the settings bundle should have a list of extensions that are allowed to read the setting to narrow down who can access 
- a permission to only allow the owner access?

### Examples

### Timezone
Idea: User should not be able to mess with his timezone

1. Admin
- has the "Edit Timezone" permission

2. Site Admin
- has the "Edit Timezone" permission

3. User
- does NOT have the "Edit Timezone" permission

#### Quota
Idea: user should not be able to change his quota.

1. Admin
- has permission to change default quota value
- has permission to change any users quota value

2. Site Admin


## Glossary

**Configuration**

- System settings
- Service host names and ports
- Changes need to be propagated to other services
- Typically modified on the CLI

**Settings**

- Application level settings
- Can be modified at runtime without restarting the service
- Typically modified in the UI

**Preferences**

- User settings
- Subset of "Settings"
- Preferred language