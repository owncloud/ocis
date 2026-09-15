
# ThumbnailSet

The ThumbnailSet resource is a keyed collection of thumbnail resources. It\'s used to represent a set of thumbnails associated with a DriveItem. 

## Properties

Name | Type
------------ | -------------
`id` | string
`large` | [Thumbnail](Thumbnail.md)
`medium` | [Thumbnail](Thumbnail.md)
`small` | [Thumbnail](Thumbnail.md)
`source` | [Thumbnail](Thumbnail.md)

## Example

```typescript
import type { ThumbnailSet } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "large": null,
  "medium": null,
  "small": null,
  "source": null,
} satisfies ThumbnailSet

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ThumbnailSet
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


