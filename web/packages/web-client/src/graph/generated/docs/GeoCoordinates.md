
# GeoCoordinates

The GeoCoordinates resource provides geographic coordinates and elevation of a location based on metadata contained within the file. If a DriveItem has a non-null location facet, the item represents a file with a known location associated with it. 

## Properties

Name | Type
------------ | -------------
`altitude` | number
`latitude` | number
`longitude` | number

## Example

```typescript
import type { GeoCoordinates } from ''

// TODO: Update the object below with actual values
const example = {
  "altitude": null,
  "latitude": null,
  "longitude": null,
} satisfies GeoCoordinates

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as GeoCoordinates
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


