# SharingLink

The `SharingLink` resource groups link-related data items into a single structure.  If a `permission` resource has a non-null `sharingLink` facet, the permission represents a sharing link (as opposed to permissions granted to a person or group). 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | [**SharingLinkType**](SharingLinkType.md) |  | [optional] [default to undefined]
**preventsDownload** | **boolean** | If &#x60;true&#x60; then the user can only use this link to view the item on the web, and cannot use it to download the contents of the item. | [optional] [readonly] [default to undefined]
**webUrl** | **string** | A URL that opens the item in the browser on the website. | [optional] [readonly] [default to undefined]
**libre_graph_displayName** | **string** | Provides a user-visible display name of the link. Optional. Libregraph only. | [optional] [default to undefined]
**libre_graph_quickLink** | **boolean** | The quicklink property can be assigned to only one link per resource. A quicklink can be used in the clients to provide a one-click copy to clipboard action. Optional. Libregraph only. | [optional] [default to undefined]

## Example

```typescript
import { SharingLink } from './api';

const instance: SharingLink = {
    type,
    preventsDownload,
    webUrl,
    libre_graph_displayName,
    libre_graph_quickLink,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
