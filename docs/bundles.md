---
title: "Settings Bundles"
date: 2020-05-04T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-settings
geekdocEditPath: edit/master/docs
geekdocFilePath: bundles.md
---

A **Settings Bundle** is a collection of settings, uniquely identified by the key of the
extension registering the bundle and the key of the bundle itself. It's purpose is to let
oCIS extensions define settings and make them available to users. They are dynamically
rendered into forms, available in the frontend.

As of now we support five different types of settings:
- boolean
- integer
- string
- single choice list of integers or strings
- multiple choice list of integers or strings

Each **Setting** is uniquely identified by a key within the bundle. Some attributes
depend on the chosen type of setting. Through the information provided with the
attributes of the setting, the settings frontend dynamically renders form elements,
allowing users to change their settings individually.

## Example

```json
{
  "identifier": {
    "extension": "ocis-accounts",
    "bundleKey": "profile"
  },
  "displayName": "Profile",
  "settings": [
    {
      "settingKey": "lastname",
      "displayName": "Lastname",
      "description": "Input for lastname",
      "stringValue": {
        "placeholder": "Set lastname"
      }
    },
    {
      "settingKey": "age",
      "displayName": "Age",
      "description": "Input for age",
      "intValue": {
        "min": "16",
        "max": "200",
        "step": "2",
        "placeholder": "Set age"
      }
    },
    {
      "settingKey": "timezone",
      "displayName": "Timezone",
      "description": "User timezone",
      "singleChoiceValue": {
        "options": [
          {
            "stringValue": "Europe/Berlin",
            "displayValue": "Europe/Berlin"
          },
          {
            "stringValue": "Asia/Kathmandu",
            "displayValue": "Asia/Kathmandu"
          }
        ]
      }
    }
  ]
}
```
