# Audio

The Audio resource groups audio-related properties on an item into a single structure.  If a DriveItem has a non-null audio facet, the item represents an audio file. The properties of the Audio resource are populated by extracting metadata from the file. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**album** | **string** | The title of the album for this audio file. | [optional] [default to undefined]
**albumArtist** | **string** | The artist named on the album for the audio file. | [optional] [default to undefined]
**artist** | **string** | The performing artist for the audio file. | [optional] [default to undefined]
**bitrate** | **number** | Bitrate expressed in kbps. | [optional] [default to undefined]
**composers** | **string** | The name of the composer of the audio file. | [optional] [default to undefined]
**copyright** | **string** | Copyright information for the audio file. | [optional] [default to undefined]
**disc** | **number** | The number of the disc this audio file came from. | [optional] [default to undefined]
**discCount** | **number** | The total number of discs in this album. | [optional] [default to undefined]
**duration** | **number** | Duration of the audio file, expressed in milliseconds | [optional] [default to undefined]
**genre** | **string** | The genre of this audio file. | [optional] [default to undefined]
**hasDrm** | **boolean** | Indicates if the file is protected with digital rights management. | [optional] [default to undefined]
**isVariableBitrate** | **boolean** | Indicates if the file is encoded with a variable bitrate. | [optional] [default to undefined]
**title** | **string** | The title of the audio file. | [optional] [default to undefined]
**track** | **number** | The number of the track on the original disc for this audio file. | [optional] [default to undefined]
**trackCount** | **number** | The total number of tracks on the original disc for this audio file. | [optional] [default to undefined]
**year** | **number** | The year the audio file was recorded. | [optional] [default to undefined]

## Example

```typescript
import { Audio } from './api';

const instance: Audio = {
    album,
    albumArtist,
    artist,
    bitrate,
    composers,
    copyright,
    disc,
    discCount,
    duration,
    genre,
    hasDrm,
    isVariableBitrate,
    title,
    track,
    trackCount,
    year,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
