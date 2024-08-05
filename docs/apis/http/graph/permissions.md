---
title: Permissions
weight: 50
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http/graph
geekdocFilePath: permissions.md
---

{{< toc >}}

## Permissions API

The Permissions API is implementing a subset of the functionality of the
[MS Graph Permission resource](https://learn.microsoft.com/en-us/graph/api/resources/permission?view=graph-rest-1.0).

### Example Permissions

The JSON representation of a Drive, as handled by the Spaces API, looks like this:
````json
{
  "@libre.graph.permissions.roles.allowedValues": [
    {
      "id": "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5",
      "description": "Allows reading the shared file or folder",
      "displayName": "Viewer",
      "@libre.graph.weight": 1
    },
    {
      "id": "fb6c3e19-e378-47e5-b277-9732f9de6e21",
      "description": "Allows reading and writing the shared file or folder",
      "displayName": "Editor",
      "@libre.graph.weight": 2
    },
    {
      "id": "312c0871-5ef7-4b3a-85b6-0e4074c64049",
      "description": "Allows managing a space",
      "displayName": "Manager",
      "@libre.graph.weight": 3
    },
    {
      "id": "4916f47e-66d5-49bb-9ac9-748ad00334b",
      "description": "Allows creating new files",
      "displayName": "File Drop",
      "@libre.graph.weight": 4
    }
  ],
  "@libre.graph.permissions.actions.allowedValues": [
    "libre.graph/driveItem/basic/read",
    "libre.graph/driveItem/permissions/read",
    "libre.graph/driveItem/upload/create",
    "libre.graph/driveItem/standard/allTasks",
    "libre.graph/driveItem/upload/create"
  ],
  "value": [
    {
      "id": "67445fde-a647-4dd4-b015-fc5dafd2821d",
      "link": {
        "type": "view",
        "webUrl": "https://cloud.example.org/s/fhGBMIkKFEHWysj"
      }
    },
    {
      "id": "34646ab6-be32-43c9-89e6-987e0c237e9b",
      "roles": [
        "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
      ],
      "grantedToV2": [
        {
          "user": {
            "id": "4c510ada-c86b-4815-8820-42cdf82c3d51",
            "displayName": "Albert Einstein"
          }
        }
      ]
    },
    {
      "id": "81d5bad3-3eff-410a-a2ea-eda2d14d4474",
      "roles": [
        "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
      ],
      "grantedToV2": [
        {
          "user": {
            "id": "4c510ada-c86b-4815-8820-42cdf82c3d51",
            "displayName": "Albert Einstein"
          }
        }
      ]
    },
    {
      "id": "b470677e-a7f5-4304-8ef5-f5056a21fff1",
      "roles": [
        "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5"
      ],
      "grantedToV2": [
        {
          "user": {
            "id": "f7fbf8c8-139b-4376-b307-cf0a8c2d0d9c",
            "displayName": "Marie Skłodowska Curie"
          }
        }
      ]
    },
    {
      "id": "453b02be-4ec2-4e7d-b576-09fc153de812",
      "roles": [
        "fb6c3e19-e378-47e5-b277-9732f9de6e21"
      ],
      "grantedToV2": [
        {
          "user": {
            "id": "4c510ada-c86b-4815-8820-42cdf82c3d51",
            "displayName": "Albert Einstein"
          }
        }
      ],
      "expirationDateTime": "2018-07-15T14:00:00.000Z"
    },
    {
      "id": "86765c0d-3905-444a-9b07-76201f8cf7df",
      "roles": [
        "312c0871-5ef7-4b3a-85b6-0e4074c64049"
      ],
      "grantedToV2": [
        {
          "group": {
            "id": "167cbee2-0518-455a-bfb2-031fe0621e5d",
            "displayName": "Philosophy Haters"
          }
        }
      ]
    },
    {
      "id": "c42b5cbd-2d65-42cf-b0b6-fb6d2b762256",
      "grantedToV2": [
        {
          "user": {
            "id": "4c510ada-c86b-4815-8820-42cdf82c3d51",
            "displayName": "Albert Einstein"
          }
        }
      ],
      "@libre.graph.permissions.actions": [
        "libre.graph/driveItem/basic/read",
        "libre.graph/driveItem/path/update"
      ]
    }
  ]
}
````

## Creating Share Invitation / Link

### Create a link share `POST /drives/{drive-id}/items/{item-id}/createLink`

https://owncloud.dev/libre-graph-api/#/drives.permissions/CreateLink

### Create a user/group share `POST /drives/{drive-id}/items/{item-id}/invite`

https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

## Reading Permissions

### List the effective sharing permissions on a driveitem `GET /drives/{drive-id}/items/{item-id}/permissions`

https://owncloud.dev/libre-graph-api/#/drives.permissions/ListPermissions

### List Get sharing permission for a file or folder `GET /drives/{drive-id}/items/{item-id}/permissions/{perm-id}`

https://owncloud.dev/libre-graph-api/#/drives.permissions/GetPermission

## Updating Permissions

### Updating sharing permission `POST /drives/{drive-id}/items/{item-id}/permissions/{perm-id}`

https://owncloud.dev/libre-graph-api/#/drives.permissions/UpdatePermission

### Set password of permission `POST /drives/{drive-id}/items/{item-id}/permissions/{perm-id}/setPassword`

https://owncloud.dev/libre-graph-api/#/drives.permissions/SetPermissionPassword

### Deleting permission `DELETE /drives/{drive-id}/items/{item-id}/permissions/{perm-id}`

https://owncloud.dev/libre-graph-api/#/drives.permissions/DeletePermission