
# Quota

Optional. Information about the drive\'s storage space quota. Read-only.

## Properties

Name | Type
------------ | -------------
`deleted` | number
`remaining` | number
`state` | string
`total` | number
`used` | number

## Example

```typescript
import type { Quota } from ''

// TODO: Update the object below with actual values
const example = {
  "deleted": null,
  "remaining": null,
  "state": null,
  "total": null,
  "used": null,
} satisfies Quota

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Quota
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


