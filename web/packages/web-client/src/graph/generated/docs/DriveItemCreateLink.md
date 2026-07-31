# DriveItemCreateLink


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | [**SharingLinkType**](SharingLinkType.md) |  | [optional] [default to undefined]
**expirationDateTime** | **string** | Optional. A String with format of yyyy-MM-ddTHH:mm:ssZ of DateTime indicates the expiration time of the permission. | [optional] [default to undefined]
**password** | **string** | Optional.The password of the sharing link that is set by the creator. | [optional] [default to undefined]
**displayName** | **string** | Provides a user-visible display name of the link. Optional. Libregraph only. | [optional] [default to undefined]
**libre_graph_quickLink** | **boolean** | The quicklink property can be assigned to only one link per resource. A quicklink can be used in the clients to provide a one-click copy to clipboard action. Optional. Libregraph only. | [optional] [default to undefined]

## Example

```typescript
import { DriveItemCreateLink } from './api';

const instance: DriveItemCreateLink = {
    type,
    expirationDateTime,
    password,
    displayName,
    libre_graph_quickLink,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
