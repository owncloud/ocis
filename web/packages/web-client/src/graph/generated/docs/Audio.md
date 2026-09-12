
# Audio

The Audio resource groups audio-related properties on an item into a single structure.  If a DriveItem has a non-null audio facet, the item represents an audio file. The properties of the Audio resource are populated by extracting metadata from the file. 

## Properties

Name | Type
------------ | -------------
`album` | string
`albumArtist` | string
`artist` | string
`bitrate` | number
`composers` | string
`copyright` | string
`disc` | number
`discCount` | number
`duration` | number
`genre` | string
`hasDrm` | boolean
`isVariableBitrate` | boolean
`title` | string
`track` | number
`trackCount` | number
`year` | number

## Example

```typescript
import type { Audio } from ''

// TODO: Update the object below with actual values
const example = {
  "album": null,
  "albumArtist": null,
  "artist": null,
  "bitrate": null,
  "composers": null,
  "copyright": null,
  "disc": null,
  "discCount": null,
  "duration": null,
  "genre": null,
  "hasDrm": null,
  "isVariableBitrate": null,
  "title": null,
  "track": null,
  "trackCount": null,
  "year": null,
} satisfies Audio

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Audio
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


