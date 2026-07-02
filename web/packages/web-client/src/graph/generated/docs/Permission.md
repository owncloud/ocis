# Permission

The Permission resource provides information about a sharing permission granted for a DriveItem resource.  ### Remarks  The Permission resource uses *facets* to provide information about the kind of permission represented by the resource.  Permissions with a `link` facet represent sharing links created on the item. Sharing links contain a unique token that provides access to the item for anyone with the link.  Permissions with a `invitation` facet represent permissions added by inviting specific users or groups to have access to the file. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | The unique identifier of the permission among all permissions on the item. Read-only. | [optional] [readonly] [default to undefined]
**hasPassword** | **boolean** | Indicates whether the password is set for this permission. This property only appears in the response. Optional. Read-only.  | [optional] [readonly] [default to undefined]
**expirationDateTime** | **string** | An optional expiration date which limits the permission in time. | [optional] [default to undefined]
**createdDateTime** | **string** | An optional creation date. Libregraph only. | [optional] [default to undefined]
**grantedToV2** | [**SharePointIdentitySet**](SharePointIdentitySet.md) |  | [optional] [default to undefined]
**link** | [**SharingLink**](SharingLink.md) |  | [optional] [default to undefined]
**roles** | **Array&lt;string&gt;** |  | [optional] [default to undefined]
**grantedToIdentities** | [**Array&lt;IdentitySet&gt;**](IdentitySet.md) | For link type permissions, the details of the identity to whom permission was granted. This could be used to grant access to a an external user that can be identified by email, aka guest accounts. | [optional] [default to undefined]
**libre_graph_permissions_actions** | **Array&lt;string&gt;** | Use this to create a permission with custom actions. | [optional] [default to undefined]
**invitation** | [**SharingInvitation**](SharingInvitation.md) |  | [optional] [default to undefined]

## Example

```typescript
import { Permission } from './api';

const instance: Permission = {
    id,
    hasPassword,
    expirationDateTime,
    createdDateTime,
    grantedToV2,
    link,
    roles,
    grantedToIdentities,
    libre_graph_permissions_actions,
    invitation,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
