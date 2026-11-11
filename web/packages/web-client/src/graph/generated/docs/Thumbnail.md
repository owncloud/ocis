
# Thumbnail

The thumbnail resource type represents a thumbnail for an image, video, document, or any item that has a bitmap representation. 

## Properties

Name | Type
------------ | -------------
`content` | string
`height` | number
`sourceItemId` | string
`url` | string
`width` | number

## Example

```typescript
import type { Thumbnail } from ''

// TODO: Update the object below with actual values
const example = {
  "content": null,
  "height": null,
  "sourceItemId": null,
  "url": null,
  "width": null,
} satisfies Thumbnail

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Thumbnail
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


