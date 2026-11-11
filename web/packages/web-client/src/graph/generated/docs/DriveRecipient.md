
# DriveRecipient

Represents a person, group, or other recipient to share a drive item with using the invite action.  When using invite to add permissions, the `driveRecipient` object would specify the `email`, `alias`, or `objectId` of the recipient. Only one of these values is required; multiple values are not accepted. 

## Properties

Name | Type
------------ | -------------
`objectId` | string
`atLibreGraphRecipientType` | string

## Example

```typescript
import type { DriveRecipient } from ''

// TODO: Update the object below with actual values
const example = {
  "objectId": null,
  "atLibreGraphRecipientType": null,
} satisfies DriveRecipient

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DriveRecipient
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


