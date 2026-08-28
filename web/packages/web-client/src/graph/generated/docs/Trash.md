
# Trash

Metadata for trashed drive Items

## Properties

Name | Type
------------ | -------------
`trashedBy` | [IdentitySet](IdentitySet.md)
`trashedDateTime` | string

## Example

```typescript
import type { Trash } from ''

// TODO: Update the object below with actual values
const example = {
  "trashedBy": null,
  "trashedDateTime": null,
} satisfies Trash

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Trash
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


