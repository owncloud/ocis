
# DriveItemInvite


## Properties

Name | Type
------------ | -------------
`recipients` | [Array&lt;DriveRecipient&gt;](DriveRecipient.md)
`roles` | Array&lt;string&gt;
`atLibreGraphPermissionsActions` | Array&lt;string&gt;
`expirationDateTime` | string

## Example

```typescript
import type { DriveItemInvite } from ''

// TODO: Update the object below with actual values
const example = {
  "recipients": null,
  "roles": null,
  "atLibreGraphPermissionsActions": null,
  "expirationDateTime": null,
} satisfies DriveItemInvite

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DriveItemInvite
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


