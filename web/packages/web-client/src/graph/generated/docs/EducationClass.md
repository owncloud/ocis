
# EducationClass

And extension of group representing a class or course

## Properties

Name | Type
------------ | -------------
`id` | string
`description` | string
`displayName` | string
`members` | [Array&lt;User&gt;](User.md)
`membersodataBind` | Set&lt;string&gt;
`classification` | string
`externalId` | string

## Example

```typescript
import type { EducationClass } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "description": null,
  "displayName": null,
  "members": null,
  "membersodataBind": null,
  "classification": null,
  "externalId": null,
} satisfies EducationClass

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as EducationClass
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


