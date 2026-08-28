
# Group


## Properties

Name | Type
------------ | -------------
`id` | string
`description` | string
`displayName` | string
`groupTypes` | Array&lt;string&gt;
`members` | [Array&lt;User&gt;](User.md)
`membersodataBind` | Set&lt;string&gt;

## Example

```typescript
import type { Group } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "description": null,
  "displayName": null,
  "groupTypes": null,
  "members": null,
  "membersodataBind": null,
} satisfies Group

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Group
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


