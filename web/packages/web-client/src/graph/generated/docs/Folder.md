
# Folder

Folder metadata, if the item is a folder. Read-only.

## Properties

Name | Type
------------ | -------------
`childCount` | number
`view` | [FolderView](FolderView.md)

## Example

```typescript
import type { Folder } from ''

// TODO: Update the object below with actual values
const example = {
  "childCount": null,
  "view": null,
} satisfies Folder

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Folder
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


