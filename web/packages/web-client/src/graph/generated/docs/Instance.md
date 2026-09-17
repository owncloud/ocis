
# Instance

An oCIS instance that the user is either a member or a guest of.

## Properties

Name | Type
------------ | -------------
`url` | string
`primary` | boolean

## Example

```typescript
import type { Instance } from ''

// TODO: Update the object below with actual values
const example = {
  "url": null,
  "primary": null,
} satisfies Instance

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Instance
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


