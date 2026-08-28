
# AppRoleAssignment


## Properties

Name | Type
------------ | -------------
`id` | string
`deletedDateTime` | string
`appRoleId` | string
`createdDateTime` | string
`principalDisplayName` | string
`principalId` | string
`principalType` | string
`resourceDisplayName` | string
`resourceId` | string

## Example

```typescript
import type { AppRoleAssignment } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "deletedDateTime": null,
  "appRoleId": null,
  "createdDateTime": null,
  "principalDisplayName": null,
  "principalId": null,
  "principalType": null,
  "resourceDisplayName": null,
  "resourceId": null,
} satisfies AppRoleAssignment

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as AppRoleAssignment
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


