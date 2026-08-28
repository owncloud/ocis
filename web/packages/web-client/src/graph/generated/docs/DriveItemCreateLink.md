
# DriveItemCreateLink


## Properties

Name | Type
------------ | -------------
`type` | [SharingLinkType](SharingLinkType.md)
`expirationDateTime` | string
`password` | string
`displayName` | string
`atLibreGraphQuickLink` | boolean

## Example

```typescript
import type { DriveItemCreateLink } from ''

// TODO: Update the object below with actual values
const example = {
  "type": null,
  "expirationDateTime": null,
  "password": null,
  "displayName": null,
  "atLibreGraphQuickLink": null,
} satisfies DriveItemCreateLink

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DriveItemCreateLink
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


