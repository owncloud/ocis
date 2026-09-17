
# Hashes

Hashes of the file\'s binary content, if available. Read-only.

## Properties

Name | Type
------------ | -------------
`crc32Hash` | string
`quickXorHash` | string
`sha1Hash` | string
`sha256Hash` | string

## Example

```typescript
import type { Hashes } from ''

// TODO: Update the object below with actual values
const example = {
  "crc32Hash": null,
  "quickXorHash": null,
  "sha1Hash": null,
  "sha256Hash": null,
} satisfies Hashes

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Hashes
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


