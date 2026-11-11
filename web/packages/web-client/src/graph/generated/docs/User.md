
# User

Represents an Active Directory user object.

## Properties

Name | Type
------------ | -------------
`id` | string
`accountEnabled` | boolean
`appRoleAssignments` | [Array&lt;AppRoleAssignment&gt;](AppRoleAssignment.md)
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
`userType` | string
`preferredLanguage` | string
`signInActivity` | [SignInActivity](SignInActivity.md)
`externalID` | string
`crossInstanceReference` | string
`instances` | [Array&lt;Instance&gt;](Instance.md)
`attributes` | Array&lt;string&gt;

## Example

```typescript
import type { User } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "accountEnabled": null,
  "appRoleAssignments": null,
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
  "userType": null,
  "preferredLanguage": null,
  "signInActivity": null,
  "externalID": null,
  "crossInstanceReference": null,
  "instances": null,
  "attributes": null,
} satisfies User

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as User
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


