# Drive

The drive represents a space on the storage.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | The unique identifier for this drive. | [optional] [readonly] [default to undefined]
**createdBy** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**createdDateTime** | **string** | Date and time of item creation. Read-only. | [optional] [readonly] [default to undefined]
**description** | **string** | Provides a user-visible description of the item. Optional. | [optional] [default to undefined]
**eTag** | **string** | ETag for the item. Read-only. | [optional] [readonly] [default to undefined]
**lastModifiedBy** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**lastModifiedDateTime** | **string** | Date and time the item was last modified. Read-only. | [optional] [readonly] [default to undefined]
**name** | **string** | The name of the item. Read-write. | [default to undefined]
**parentReference** | [**ItemReference**](ItemReference.md) |  | [optional] [default to undefined]
**webUrl** | **string** | URL that displays the resource in the browser. Read-only. | [optional] [readonly] [default to undefined]
**driveType** | **string** | Describes the type of drive represented by this resource. Values are \&quot;personal\&quot; for users home spaces, \&quot;project\&quot;, \&quot;virtual\&quot; or \&quot;share\&quot;. Read-only. | [optional] [readonly] [default to undefined]
**driveAlias** | **string** | The drive alias can be used in clients to make the urls user friendly. Example: \&#39;personal/einstein\&#39;. This will be used to resolve to the correct driveID. | [optional] [default to undefined]
**owner** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**quota** | [**Quota**](Quota.md) |  | [optional] [default to undefined]
**items** | [**Array&lt;DriveItem&gt;**](DriveItem.md) | All items contained in the drive. Read-only. Nullable. | [optional] [readonly] [default to undefined]
**root** | [**DriveItem**](DriveItem.md) |  | [optional] [default to undefined]
**special** | [**Array&lt;DriveItem&gt;**](DriveItem.md) | A collection of special drive resources. | [optional] [default to undefined]

## Example

```typescript
import { Drive } from './api';

const instance: Drive = {
    id,
    createdBy,
    createdDateTime,
    description,
    eTag,
    lastModifiedBy,
    lastModifiedDateTime,
    name,
    parentReference,
    webUrl,
    driveType,
    driveAlias,
    owner,
    quota,
    items,
    root,
    special,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
