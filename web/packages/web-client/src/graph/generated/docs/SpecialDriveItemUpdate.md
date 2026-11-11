
# SpecialDriveItemUpdate

References an existing item inside a drive to be assigned as one of the drive\'s special resources, such as its image or readme.

## Properties

Name | Type
------------ | -------------
`id` | string
`specialFolder` | [SpecialFolder](SpecialFolder.md)

## Example

```typescript
import type { SpecialDriveItemUpdate } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "specialFolder": null,
} satisfies SpecialDriveItemUpdate

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SpecialDriveItemUpdate
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


