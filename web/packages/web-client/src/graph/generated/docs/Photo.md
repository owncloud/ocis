
# Photo

The photo resource provides photo and camera properties, for example, EXIF metadata, on a driveItem. 

## Properties

Name | Type
------------ | -------------
`cameraMake` | string
`cameraModel` | string
`exposureDenominator` | number
`exposureNumerator` | number
`fNumber` | number
`focalLength` | number
`iso` | number
`orientation` | number
`takenDateTime` | Date

## Example

```typescript
import type { Photo } from ''

// TODO: Update the object below with actual values
const example = {
  "cameraMake": null,
  "cameraModel": null,
  "exposureDenominator": null,
  "exposureNumerator": null,
  "fNumber": null,
  "focalLength": null,
  "iso": null,
  "orientation": null,
  "takenDateTime": null,
} satisfies Photo

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Photo
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


