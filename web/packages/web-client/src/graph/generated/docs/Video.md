
# Video

The video resource groups video-related data items into a single structure.  If a driveItem has a non-null video facet, the item represents a video file. The properties of the video resource are populated by extracting metadata from the file. 

## Properties

Name | Type
------------ | -------------
`audioBitsPerSample` | number
`audioChannels` | number
`audioFormat` | string
`audioSamplesPerSecond` | number
`bitrate` | number
`duration` | number
`fourCC` | string
`frameRate` | number
`height` | number
`width` | number

## Example

```typescript
import type { Video } from ''

// TODO: Update the object below with actual values
const example = {
  "audioBitsPerSample": null,
  "audioChannels": null,
  "audioFormat": null,
  "audioSamplesPerSecond": null,
  "bitrate": null,
  "duration": null,
  "fourCC": null,
  "frameRate": null,
  "height": null,
  "width": null,
} satisfies Video

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Video
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


