
# SharingLinkType

The type of the link created.  | Value          | Display name      | Description                                                     | | -------------- | ----------------- | --------------------------------------------------------------- | | internal       | Internal          | Creates an internal link without any permissions.               | | view           | View              | Creates a read-only link to the driveItem.                      | | upload         | Upload            | Creates a read-write link to the folder driveItem.              | | edit           | Edit              | Creates a read-write link to the driveItem.                     | | createOnly     | File Drop         | Creates an upload-only link to the folder driveItem.            | | blocksDownload | Secure View       | Creates a read-only link that blocks download to the driveItem. | 

## Properties

Name | Type
------------ | -------------

## Example

```typescript
import type { SharingLinkType } from ''

// TODO: Update the object below with actual values
const example = {
} satisfies SharingLinkType

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SharingLinkType
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


