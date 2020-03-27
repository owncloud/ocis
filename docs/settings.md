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

### Settings bundle definition
A bundle has several properties:
- a name that uniquely identifies the bundle
- a display name that is shown to administrators
- the service the bundle belongs to
- a set of settings (see below)

Each setting is specified as follows:
- a name that uniquely identifies the setting within the settings bundle
- a display name that is shown to users
- a description (optional)
- a set of values (form definition):
  - a default value or function constructing the default value
  - type (required)
    - string
    - integer
    - boolean
    - list (single select / radiobutton)
    - list (multi select / checkbox)
  - set of validation rules
    - email
    - password
    - required
    - min
    - max
  - stepping
  - placeholder
  - options

### Example: `User Profile`
The following JSON defines a Settings Bundle `User Profile` with two settings `Email Address` and `Timezone`.

```json
{
  "name": "user-profile",
  "displayName": "User Profile",
  "extension": "account",
  "settings": [
    {
      "name": "email",
      "displayName": "Email Address",
      "description": null,
      "values": [
        {
          "type": "string",
          "default": null,
          "validation": ["email", "required"],
          "placeholder": "Provide an email address"
        }
      ]
    },
    {
      "name": "timezone",
      "displayName": "Timezone",
      "description": null,
      "values": [
        {
          "type": "list",
          "validation": ["required"],
          "options": [
            {"value": 0, "label": "unknown"},
            {"value": 1, "label": "Europe/Berlin", "default": true},
            {"value": 2, "label": "Europe/Amsterdam"}
          ]
        }
      ]
    }
  ]
}
```

### Registering Settings Bundles

When registering a setting with the settings service, every single setting automatically creates three permissions:
- `read` to allow reading the setting value from the settings service
- `write` to allow writing the setting value to the settings service
- `display` to allow seeing or listing a settings value in the settings ui

A user that has the `read` permission, the `display` permission but no `write` permission could see his timezone in the
settings ui but not change it.

A user that has the `read` permission but no `display` permission should not even see his timezone in the settings ui.

Referencing a specific permission when querying the settings service works by concatenating the extension,
settings bundle, setting and permission, e.g. `account:user-profile:email:read`. When performing an operation against the
settings service, the permission can be omitted, as it is implicitly clear from the operation.

### Additional permissions
TODO: we will probably need additional permissions. Also it might be a good idea to split `write` into `create` and
`update`. We have to decide whether or not it is a good idea to define additional permissions within the settings bundles
or in a separate way.

## Roles

- Every user has roles that are tied to his account.
- Every role has a set of permissions.
- Every permission evaluates to a boolean for a certain scope.

The two roles `user` and `admin` get created on first start of the settings service with the characteristics that the
names imply.

```json
{
    "roles": [
        {
          "name": "admin",
          "displayName": "Admin",
          "scope": {"type": "system"},
          "permissions": [
            {"name": "account:user-profile:email:read", "scope": {"type": "user", "value": ["all"]}},
            {"name": "account:user-profile:email:write", "scope": {"type": "user", "value": ["all"]}}
          ]
        },
        {
          "name": "site-admin",
          "displayName": "Group Admin",
          "scope": {"type": "group", "value": ["group-uid-1", "group-ui-2"]},
          "permissions": [
            {"name":"account:user-profile:email:read", "scope": {"type": "group", "value": ["all"]}},
            {"name":"account:user-profile:email:write", "scope": {"type": "group", "value": ["all"]}},
            {"name":"account:user-profile:email:read", "scope": {"type": "user", "value": ["me"]}},
            {"name":"account:user-profile:email:write", "scope": {"type": "user", "value": ["me"]}}
          ]
        },
        {
          "name": "user",
          "displayName": "User",
          "scope": {"type": "system"},
          "permissions": [
            {"name":"account:user-profile:email:read", "scope": {"type": "user", "value": ["me"]}},
            {"name":"account:user-profile:email:write", "scope": {"type": "user", "value": ["me"]}},
            {"name":"sharing:shares:public-link:create", "scope": {"type": "user", "value": ["me"]}},
            {"name":"files:files:file:read", "scope": {"type":  "user", "value": ["me"]}},
            {"name":"files:files:file:create", "scope": {"type":  "user", "value": ["me"]}},
            {"name":"files:files:file:update", "scope": {"type":  "user", "value": ["me"]}},
            {"name":"files:files:folder:read", "scope": {"type":  "user", "value": ["me"]}},
            {"name":"files:files:folder:create", "scope": {"type":  "user", "value": ["me"]}},
            {"name":"files:files:folder:update", "scope": {"type":  "user", "value": ["me"]}}
          ]
        },
        {
          "name": "share-uid-xyz-member",
          "displayName": "All receivers of share with uid xyz",
          "scope": {"type": "share", "value": ["uid-xyz"]},
          "permissions": [
            {"name": "sharing:shares:file-share:read", "scope": {"type": "user", "value": ["me"]}}
          ]
        },
        {
          "name": "guest",
          "displayName": "Guest",
          "scope": {"type": "system"},
          "permissions": [
            {"name":"files:files:file:read", "scope": {"type":  "user", "value": ["me"]}},
            {"name":"files:files:folder:read", "scope": {"type":  "user", "value": ["me"]}}
          ]
        }
    ]
}

```

### Scopes for roles
We distinguish three different scopes for roles.
1. The `system` scope defines that a role can only exist once. The user role being in the system scope means, that
there can't be another set of permissions that is applied to all users.
2. A role with `group` scope only exists for the group ids set in the scope.
3. A role with `share` scope only exists for the share ids set in the scope.

## Workflow of the settings service
On first start, the settings service creates the roles `admin` and `user` as archetypes. Those two roles
cannot be removed. Every settings bundle that is subsequently registered assigns permissions to these two roles in a way
that they evaluate to `true` for all users for the admin role and evaluate to `true` only for `me` for the user role
unless defined otherwise in the `userPermissions` map of the setting definition.

TODO: the service owner probably has to decide about initial permissions for existing roles when a new service registers
settings bundles. We need a clever way to not let this turn into a configuration mess. Probably a good idea to define
archetype permission mappings in settings.

TODO: passwords
- sometimes we need to store them but only certain extensions should be allowed to read them, eg. the wnd app
- they should be encrypted anyway
- but the settings bundle should have a list of extensions that are allowed to read the setting to narrow down who can access
- a permission to only allow the owner access?

TODO: we should have a json-schema for settings bundle validation

## TODO: Data model for saved values in the settings service
The settings service has to be able to save multiple instances of a setting, coming from different roles or services,
allowing to save a value and a default for each. This way we can override values and provide defaults in a hierarchy.
This also means that the roles need to have a hierarchy (probably a tree).
The language of a user for example could come as a default from the browser or IDP, can be set by the user and could be
overwritten by an admin. There could also be a default defined by the admin, that is only used if the user didn't define
a language on their own. Defaults should not override other defaults - only if a default is missing but needed, we look
in the next higher role for a default. Values on the other hand should override. If a higher role sets a value for a
certain user it must have precedence to the value set by the user.

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

- System configuration
- e.g. service host names and ports
- Changes need to be propagated to other services
- Typically modified on the CLI

**Settings**

- Application level settings
- e.g. default language
- Can be modified at runtime without restarting the service
- Typically modified in the UI

**Preferences**

- User settings
- Subset of "Settings"
- Preferred language
