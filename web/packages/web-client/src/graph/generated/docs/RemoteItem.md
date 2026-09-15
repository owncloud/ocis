
# RemoteItem

Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.

## Properties

Name | Type
------------ | -------------
`createdBy` | [IdentitySet](IdentitySet.md)
`createdDateTime` | string
`file` | [OpenGraphFile](OpenGraphFile.md)
`fileSystemInfo` | [FileSystemInfo](FileSystemInfo.md)
`folder` | [Folder](Folder.md)
`driveAlias` | string
`path` | string
`rootId` | string
`id` | string
`image` | [Image](Image.md)
`lastModifiedBy` | [IdentitySet](IdentitySet.md)
`lastModifiedDateTime` | string
`name` | string
`eTag` | string
`cTag` | string
`parentReference` | [ItemReference](ItemReference.md)
`permissions` | [Array&lt;Permission&gt;](Permission.md)
`size` | number
`specialFolder` | [SpecialFolder](SpecialFolder.md)
`webDavUrl` | string
`webUrl` | string
`spaceId` | string

## Example

```typescript
import type { RemoteItem } from ''

// TODO: Update the object below with actual values
const example = {
  "createdBy": null,
  "createdDateTime": null,
  "file": null,
  "fileSystemInfo": null,
  "folder": null,
  "driveAlias": null,
  "path": null,
  "rootId": null,
  "id": null,
  "image": null,
  "lastModifiedBy": null,
  "lastModifiedDateTime": null,
  "name": null,
  "eTag": null,
  "cTag": null,
  "parentReference": null,
  "permissions": null,
  "size": null,
  "specialFolder": null,
  "webDavUrl": null,
  "webUrl": null,
  "spaceId": null,
} satisfies RemoteItem

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RemoteItem
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


