
# EducationUser

An extension of user with education-specific attributes

## Properties

Name | Type
------------ | -------------
`id` | string
`accountEnabled` | boolean
`displayName` | string
`drives` | [Array&lt;Drive&gt;](Drive.md)
`drive` | [Drive](Drive.md)
`identities` | [Array&lt;ObjectIdentity&gt;](ObjectIdentity.md)
`mail` | string
`memberOf` | [Array&lt;Group&gt;](Group.md)
`onPremisesSamAccountName` | string
`passwordProfile` | [PasswordProfile](PasswordProfile.md)
`surname` | string
`givenName` | string
`primaryRole` | string
`userType` | string
`externalID` | string

## Example

```typescript
import type { EducationUser } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "accountEnabled": null,
  "displayName": null,
  "drives": null,
  "drive": null,
  "identities": null,
  "mail": null,
  "memberOf": null,
  "onPremisesSamAccountName": null,
  "passwordProfile": null,
  "surname": null,
  "givenName": null,
  "primaryRole": null,
  "userType": null,
  "externalID": null,
} satisfies EducationUser

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as EducationUser
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


