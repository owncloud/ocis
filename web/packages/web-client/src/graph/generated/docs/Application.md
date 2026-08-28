
# Application


## Properties

Name | Type
------------ | -------------
`id` | string
`appRoles` | [Array&lt;AppRole&gt;](AppRole.md)
`displayName` | string

## Example

```typescript
import type { Application } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "appRoles": null,
  "displayName": null,
} satisfies Application

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Application
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


