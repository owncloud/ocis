# Photo

The photo resource provides photo and camera properties, for example, EXIF metadata, on a driveItem. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**cameraMake** | **string** | Camera manufacturer. Read-only. | [optional] [default to undefined]
**cameraModel** | **string** | Camera model. Read-only. | [optional] [default to undefined]
**exposureDenominator** | **number** | The denominator for the exposure time fraction from the camera. Read-only. | [optional] [default to undefined]
**exposureNumerator** | **number** | The numerator for the exposure time fraction from the camera. Read-only. | [optional] [default to undefined]
**fNumber** | **number** | The F-stop value from the camera. Read-only. | [optional] [default to undefined]
**focalLength** | **number** | The focal length from the camera. Read-only. | [optional] [default to undefined]
**iso** | **number** | The ISO value from the camera. Read-only. | [optional] [default to undefined]
**orientation** | **number** | The orientation value from the camera. Read-only. | [optional] [default to undefined]
**takenDateTime** | **string** | Represents the date and time the photo was taken. Read-only. | [optional] [default to undefined]

## Example

```typescript
import { Photo } from './api';

const instance: Photo = {
    cameraMake,
    cameraModel,
    exposureDenominator,
    exposureNumerator,
    fNumber,
    focalLength,
    iso,
    orientation,
    takenDateTime,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
