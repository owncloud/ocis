
# UnifiedRoleDefinition

A role definition is a collection of permissions in libre graph listing the operations that can be performed and the resources against which they can performed. 

## Properties

Name | Type
------------ | -------------
`description` | string
`displayName` | string
`id` | string
`rolePermissions` | [Array&lt;UnifiedRolePermission&gt;](UnifiedRolePermission.md)
`atLibreGraphWeight` | number

## Example

```typescript
import type { UnifiedRoleDefinition } from ''

// TODO: Update the object below with actual values
const example = {
  "description": null,
  "displayName": null,
  "id": null,
  "rolePermissions": null,
  "atLibreGraphWeight": null,
} satisfies UnifiedRoleDefinition

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as UnifiedRoleDefinition
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


