
# FileSystemInfo

File system information on client. Read-write.

## Properties

Name | Type
------------ | -------------
`createdDateTime` | Date
`lastAccessedDateTime` | Date
`lastModifiedDateTime` | Date

## Example

```typescript
import type { FileSystemInfo } from ''

// TODO: Update the object below with actual values
const example = {
  "createdDateTime": null,
  "lastAccessedDateTime": null,
  "lastModifiedDateTime": null,
} satisfies FileSystemInfo

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as FileSystemInfo
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


