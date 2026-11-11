
# PasswordProfile

Password Profile associated with a user

## Properties

Name | Type
------------ | -------------
`forceChangePasswordNextSignIn` | boolean
`password` | string

## Example

```typescript
import type { PasswordProfile } from ''

// TODO: Update the object below with actual values
const example = {
  "forceChangePasswordNextSignIn": null,
  "password": null,
} satisfies PasswordProfile

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as PasswordProfile
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


