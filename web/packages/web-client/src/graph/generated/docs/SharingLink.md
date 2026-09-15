
# SharingLink

The `SharingLink` resource groups link-related data items into a single structure.  If a `permission` resource has a non-null `sharingLink` facet, the permission represents a sharing link (as opposed to permissions granted to a person or group). 

## Properties

Name | Type
------------ | -------------
`type` | [SharingLinkType](SharingLinkType.md)
`preventsDownload` | boolean
`webUrl` | string
`atLibreGraphDisplayName` | string
`atLibreGraphQuickLink` | boolean

## Example

```typescript
import type { SharingLink } from ''

// TODO: Update the object below with actual values
const example = {
  "type": null,
  "preventsDownload": null,
  "webUrl": null,
  "atLibreGraphDisplayName": null,
  "atLibreGraphQuickLink": null,
} satisfies SharingLink

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SharingLink
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


