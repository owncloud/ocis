# Video

The video resource groups video-related data items into a single structure.  If a driveItem has a non-null video facet, the item represents a video file. The properties of the video resource are populated by extracting metadata from the file. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**audioBitsPerSample** | **number** | Number of audio bits per sample. | [optional] [default to undefined]
**audioChannels** | **number** | Number of audio channels. | [optional] [default to undefined]
**audioFormat** | **string** | Name of the audio format (AAC, MP3, etc.). | [optional] [default to undefined]
**audioSamplesPerSecond** | **number** | Number of audio samples per second. | [optional] [default to undefined]
**bitrate** | **number** | Bit rate of the video in bits per second. | [optional] [default to undefined]
**duration** | **number** | Duration of the file in milliseconds. | [optional] [default to undefined]
**fourCC** | **string** | \\\&quot;Four character code\\\&quot; name of the video format. | [optional] [default to undefined]
**frameRate** | **number** | Frame rate of the video. | [optional] [default to undefined]
**height** | **number** | Height of the video, in pixels. | [optional] [default to undefined]
**width** | **number** | Width of the video, in pixels. | [optional] [default to undefined]

## Example

```typescript
import { Video } from './api';

const instance: Video = {
    audioBitsPerSample,
    audioChannels,
    audioFormat,
    audioSamplesPerSecond,
    bitrate,
    duration,
    fourCC,
    frameRate,
    height,
    width,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
