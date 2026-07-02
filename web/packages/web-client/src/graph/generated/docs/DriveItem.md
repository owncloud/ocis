# DriveItem

Represents a resource inside a drive. Read-only.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Read-only. | [optional] [readonly] [default to undefined]
**createdBy** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**createdDateTime** | **string** | Date and time of item creation. Read-only. | [optional] [readonly] [default to undefined]
**description** | **string** | Provides a user-visible description of the item. Optional. | [optional] [default to undefined]
**eTag** | **string** | ETag for the item. Read-only. | [optional] [readonly] [default to undefined]
**lastModifiedBy** | [**IdentitySet**](IdentitySet.md) |  | [optional] [default to undefined]
**lastModifiedDateTime** | **string** | Date and time the item was last modified. Read-only. | [optional] [readonly] [default to undefined]
**name** | **string** | The name of the item. Read-write. | [optional] [default to undefined]
**parentReference** | [**ItemReference**](ItemReference.md) |  | [optional] [default to undefined]
**webUrl** | **string** | URL that displays the resource in the browser. Read-only. | [optional] [readonly] [default to undefined]
**content** | **string** | The content stream, if the item represents a file. | [optional] [default to undefined]
**cTag** | **string** | An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Note This property is not returned if the item is a folder. Read-only. | [optional] [readonly] [default to undefined]
**deleted** | [**Deleted**](Deleted.md) |  | [optional] [default to undefined]
**file** | [**OpenGraphFile**](OpenGraphFile.md) |  | [optional] [default to undefined]
**fileSystemInfo** | [**FileSystemInfo**](FileSystemInfo.md) |  | [optional] [default to undefined]
**folder** | [**Folder**](Folder.md) |  | [optional] [default to undefined]
**image** | [**Image**](Image.md) |  | [optional] [default to undefined]
**photo** | [**Photo**](Photo.md) |  | [optional] [default to undefined]
**location** | [**GeoCoordinates**](GeoCoordinates.md) |  | [optional] [default to undefined]
**thumbnails** | [**Array&lt;ThumbnailSet&gt;**](ThumbnailSet.md) | Collection containing ThumbnailSet objects associated with the item. Read-only. Nullable. | [optional] [default to undefined]
**root** | **object** | If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive. | [optional] [default to undefined]
**trash** | [**Trash**](Trash.md) |  | [optional] [default to undefined]
**specialFolder** | [**SpecialFolder**](SpecialFolder.md) |  | [optional] [default to undefined]
**remoteItem** | [**RemoteItem**](RemoteItem.md) |  | [optional] [default to undefined]
**size** | **number** | Size of the item in bytes. Read-only. | [optional] [readonly] [default to undefined]
**webDavUrl** | **string** | WebDAV compatible URL for the item. Read-only. | [optional] [readonly] [default to undefined]
**children** | [**Array&lt;DriveItem&gt;**](DriveItem.md) | Collection containing Item objects for the immediate children of Item. Only items representing folders have children. Read-only. Nullable. | [optional] [readonly] [default to undefined]
**permissions** | [**Array&lt;Permission&gt;**](Permission.md) | The set of permissions for the item. Read-only. Nullable. | [optional] [readonly] [default to undefined]
**audio** | [**Audio**](Audio.md) |  | [optional] [default to undefined]
**video** | [**Video**](Video.md) |  | [optional] [default to undefined]
**client_synchronize** | **boolean** | Indicates if the item is synchronized with the underlying storage provider. Read-only. | [optional] [default to undefined]
**UI_Hidden** | **boolean** | Properties or facets (see UI.Facet) annotated with this term will not be rendered if the annotation evaluates to true. Users can set this to hide permissions. | [optional] [default to undefined]

## Example

```typescript
import { DriveItem } from './api';

const instance: DriveItem = {
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
    content,
    cTag,
    deleted,
    file,
    fileSystemInfo,
    folder,
    image,
    photo,
    location,
    thumbnails,
    root,
    trash,
    specialFolder,
    remoteItem,
    size,
    webDavUrl,
    children,
    permissions,
    audio,
    video,
    client_synchronize,
    UI_Hidden,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
