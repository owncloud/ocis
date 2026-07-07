# DriveItemInvite


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**recipients** | [**Array&lt;DriveRecipient&gt;**](DriveRecipient.md) | A collection of recipients who will receive access and the sharing invitation. Currently, only internal users or groups are supported. | [optional] [default to undefined]
**roles** | **Array&lt;string&gt;** | Specifies the roles that are to be granted to the recipients of the sharing invitation. | [optional] [default to undefined]
**libre_graph_permissions_actions** | **Array&lt;string&gt;** | Specifies the actions that are to be granted to the recipients of the sharing invitation, in effect creating a custom role. | [optional] [default to undefined]
**expirationDateTime** | **string** | Specifies the dateTime after which the permission expires. | [optional] [default to undefined]

## Example

```typescript
import { DriveItemInvite } from './api';

const instance: DriveItemInvite = {
    recipients,
    roles,
    libre_graph_permissions_actions,
    expirationDateTime,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
