
# DriveUpdate

The drive represents an update to a space on the storage.

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
`driveType` | string
`driveAlias` | string
`owner` | [IdentitySet](IdentitySet.md)
`quota` | [Quota](Quota.md)
`items` | [Array&lt;DriveItem&gt;](DriveItem.md)
`root` | [DriveItem](DriveItem.md)
`special` | [Array&lt;DriveItem&gt;](DriveItem.md)

## Example

```typescript
import type { DriveUpdate } from ''

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
  "driveType": null,
  "driveAlias": null,
  "owner": null,
  "quota": null,
  "items": null,
  "root": null,
  "special": null,
} satisfies DriveUpdate

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DriveUpdate
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


