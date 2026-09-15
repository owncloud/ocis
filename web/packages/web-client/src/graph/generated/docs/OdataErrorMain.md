
# OdataErrorMain


## Properties

Name | Type
------------ | -------------
`code` | string
`message` | string
`target` | string
`details` | [Array&lt;OdataErrorDetail&gt;](OdataErrorDetail.md)
`innererror` | object

## Example

```typescript
import type { OdataErrorMain } from ''

// TODO: Update the object below with actual values
const example = {
  "code": null,
  "message": null,
  "target": null,
  "details": null,
  "innererror": null,
} satisfies OdataErrorMain

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as OdataErrorMain
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


