# RemoteItem

Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdBy** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**createdDateTime** | **string** | Date and time of item creation. Read-only. | [optional] [default to undefined]
**file** | [**OpenGraphFile**](OpenGraphFile.md) |  | [optional] [default to undefined]
**fileSystemInfo** | [**FileSystemInfo**](FileSystemInfo.md) |  | [optional] [default to undefined]
**folder** | [**Folder**](Folder.md) |  | [optional] [default to undefined]
**driveAlias** | **string** | The drive alias can be used in clients to make the urls user friendly. Example: \&#39;personal/einstein\&#39;. This will be used to resolve to the correct driveID. | [optional] [default to undefined]
**path** | **string** | The relative path of the item in relation to its drive root. | [optional] [default to undefined]
**rootId** | **string** | Unique identifier for the drive root of this item. Read-only. | [optional] [default to undefined]
**id** | **string** | Unique identifier for the remote item in its drive. Read-only. | [optional] [default to undefined]
**image** | [**Image**](Image.md) |  | [optional] [default to undefined]
**lastModifiedBy** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**lastModifiedDateTime** | **string** | Date and time the item was last modified. Read-only. | [optional] [default to undefined]
**name** | **string** | Optional. Filename of the remote item. Read-only. | [optional] [default to undefined]
**eTag** | **string** | ETag for the item. Read-only. | [optional] [readonly] [default to undefined]
**cTag** | **string** | An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Note This property is not returned if the item is a folder. Read-only. | [optional] [readonly] [default to undefined]
**parentReference** | [**ItemReference**](ItemReference.md) |  | [optional] [default to undefined]
**permissions** | [**Array&lt;Permission&gt;**](Permission.md) | The set of permissions for the item. Read-only. Nullable. | [optional] [readonly] [default to undefined]
**size** | **number** | Size of the remote item. Read-only. | [optional] [default to undefined]
**specialFolder** | [**SpecialFolder**](SpecialFolder.md) |  | [optional] [default to undefined]
**webDavUrl** | **string** | DAV compatible URL for the item. | [optional] [default to undefined]
**webUrl** | **string** | URL that displays the resource in the browser. Read-only. | [optional] [default to undefined]
**spaceId** | **string** | The UUID of the space that contains the item. | [optional] [default to undefined]

## Example

```typescript
import { RemoteItem } from './api';

const instance: RemoteItem = {
    createdBy,
    createdDateTime,
    file,
    fileSystemInfo,
    folder,
    driveAlias,
    path,
    rootId,
    id,
    image,
    lastModifiedBy,
    lastModifiedDateTime,
    name,
    eTag,
    cTag,
    parentReference,
    permissions,
    size,
    specialFolder,
    webDavUrl,
    webUrl,
    spaceId,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
