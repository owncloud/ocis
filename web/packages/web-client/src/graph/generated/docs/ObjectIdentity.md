
# ObjectIdentity

Represents an identity used to sign in to a user account

## Properties

Name | Type
------------ | -------------
`issuer` | string
`issuerAssignedId` | string

## Example

```typescript
import type { ObjectIdentity } from ''

// TODO: Update the object below with actual values
const example = {
  "issuer": null,
  "issuerAssignedId": null,
} satisfies ObjectIdentity

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ObjectIdentity
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


