
# UnifiedRolePermission

Represents a collection of allowed resource actions and the conditions that must be met for the action to be allowed. Resource actions are tasks that can be performed on a resource. For example, an application resource may support create, update, delete, and reset password actions. 

## Properties

Name | Type
------------ | -------------
`allowedResourceActions` | Array&lt;string&gt;
`condition` | string

## Example

```typescript
import type { UnifiedRolePermission } from ''

// TODO: Update the object below with actual values
const example = {
  "allowedResourceActions": null,
  "condition": null,
} satisfies UnifiedRolePermission

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as UnifiedRolePermission
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


