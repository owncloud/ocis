
# SharePointIdentitySet

This resource is used to represent a set of identities associated with various events for an item, such as created by or last modified by.

## Properties

Name | Type
------------ | -------------
`user` | [Identity](Identity.md)
`group` | [Identity](Identity.md)

## Example

```typescript
import type { SharePointIdentitySet } from ''

// TODO: Update the object below with actual values
const example = {
  "user": null,
  "group": null,
} satisfies SharePointIdentitySet

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SharePointIdentitySet
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


