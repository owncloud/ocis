
# OpenGraphFile

File metadata, if the item is a file. Read-only.

## Properties

Name | Type
------------ | -------------
`hashes` | [Hashes](Hashes.md)
`mimeType` | string
`processingMetadata` | boolean

## Example

```typescript
import type { OpenGraphFile } from ''

// TODO: Update the object below with actual values
const example = {
  "hashes": null,
  "mimeType": null,
  "processingMetadata": null,
} satisfies OpenGraphFile

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as OpenGraphFile
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


