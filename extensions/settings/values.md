---
title: "Settings Values"
date: 2020-05-04T00:00:00+00:00
weight: 51
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/settings
geekdocFilePath: values.md
---

A **Settings Value** is the value an authenticated user has chosen for a specific setting, defined in a
*settings bundle*. For choosing settings values as a user the sole entry point is the ocis-web extension
provided by this service.

## Identifying settings values

A *settings value* is uniquely identified by four attributes. Three of them are coming from the definition of
the setting within it's settings bundle (see [Settings Bundles]({{< ref "bundles" >}})
for an example). The fourth identifies the user.
- extension: Key of the extension that registered the settings bundle,
- bundleKey: Key of the settings bundle,
- settingKey: Key of the setting as defined within the bundle,
- accountUuid: The UUID of the authenticated user who has saved the setting.

{{< hint info >}}
When requests are going through `ocis-proxy`, the accountUuid attribute can be set to the static keyword `me`
instead of using a real UUID. `ocis-proxy` will take care of minting the UUID of the authenticated user into
a JWT, providing it in the HTTP header as `x-access-token`. That UUID is then used in this service, to replace
`me` with the actual UUID of the authenticated user.
{{< /hint >}}

## Example of stored settings values

```json
{
  "values": {
    "language": {
      "identifier": {
        "extension": "ocis-accounts",
        "bundleKey": "profile",
        "settingKey": "language",
        "accountUuid": "5681371f-4a6e-43bc-8bb5-9c9237fa9c58"
      },
      "listValue": {
        "values": [
          {
            "stringValue": "de"
          }
        ]
      }
    },
    "timezone": {
      "identifier": {
        "extension": "ocis-accounts",
        "bundleKey": "profile",
        "settingKey": "timezone",
        "accountUuid": "5681371f-4a6e-43bc-8bb5-9c9237fa9c58"
      },
      "listValue": {
        "values": [
          {
            "stringValue": "Europe/Berlin"
          }
        ]
      }
    }
  }
}
```

## gRPC endpoints
The obvious way of modifying settings is the ocis-web extension, as described earlier. However, services can
use the respective gRPC endpoints of the `ValueService` to query and modify *settings values* as well.
The gRPC endpoints require the same identifier attributes as described above, so for making a request to
the `ValueService` you will have to make sure that the accountUuid of the authenticated user is available in
your service at the time of the request.
