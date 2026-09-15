
# DriveItem

Represents a resource inside a drive. Read-only.

## Properties

Name | Type
------------ | -------------
`id` | string
`createdBy` | [IdentitySet](IdentitySet.md)
`createdDateTime` | string
`description` | string
`eTag` | string
`lastModifiedBy` | [IdentitySet](IdentitySet.md)
`lastModifiedDateTime` | string
`name` | string
`parentReference` | [ItemReference](ItemReference.md)
`webUrl` | string
`content` | string
`cTag` | string
`deleted` | [Deleted](Deleted.md)
`file` | [OpenGraphFile](OpenGraphFile.md)
`fileSystemInfo` | [FileSystemInfo](FileSystemInfo.md)
`folder` | [Folder](Folder.md)
`image` | [Image](Image.md)
`photo` | [Photo](Photo.md)
`location` | [GeoCoordinates](GeoCoordinates.md)
`thumbnails` | [Array&lt;ThumbnailSet&gt;](ThumbnailSet.md)
`root` | object
`trash` | [Trash](Trash.md)
`specialFolder` | [SpecialFolder](SpecialFolder.md)
`remoteItem` | [RemoteItem](RemoteItem.md)
`size` | number
`webDavUrl` | string
`children` | [Array&lt;DriveItem&gt;](DriveItem.md)
`permissions` | [Array&lt;Permission&gt;](Permission.md)
`audio` | [Audio](Audio.md)
`video` | [Video](Video.md)
`atClientSynchronize` | boolean
`atUIHidden` | boolean

## Example

```typescript
import type { DriveItem } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "createdBy": null,
  "createdDateTime": null,
  "description": null,
  "eTag": null,
  "lastModifiedBy": null,
  "lastModifiedDateTime": null,
  "name": null,
  "parentReference": null,
  "webUrl": null,
  "content": null,
  "cTag": null,
  "deleted": null,
  "file": null,
  "fileSystemInfo": null,
  "folder": null,
  "image": null,
  "photo": null,
  "location": null,
  "thumbnails": null,
  "root": null,
  "trash": null,
  "specialFolder": null,
  "remoteItem": null,
  "size": null,
  "webDavUrl": null,
  "children": null,
  "permissions": null,
  "audio": null,
  "video": null,
  "atClientSynchronize": null,
  "atUIHidden": null,
} satisfies DriveItem

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DriveItem
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


