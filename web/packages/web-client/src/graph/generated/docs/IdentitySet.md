
# IdentitySet

Optional. User account.

## Properties

Name | Type
------------ | -------------
`application` | [Identity](Identity.md)
`device` | [Identity](Identity.md)
`user` | [Identity](Identity.md)
`group` | [Identity](Identity.md)

## Example

```typescript
import type { IdentitySet } from ''

// TODO: Update the object below with actual values
const example = {
  "application": null,
  "device": null,
  "user": null,
  "group": null,
} satisfies IdentitySet

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as IdentitySet
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


