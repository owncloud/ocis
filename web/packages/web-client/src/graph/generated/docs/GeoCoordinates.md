# GeoCoordinates

The GeoCoordinates resource provides geographic coordinates and elevation of a location based on metadata contained within the file. If a DriveItem has a non-null location facet, the item represents a file with a known location associated with it. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**altitude** | **number** | The altitude (height), in feet, above sea level for the item. Read-only. | [optional] [default to undefined]
**latitude** | **number** | The latitude, in decimal, for the item. Read-only. | [optional] [default to undefined]
**longitude** | **number** | The longitude, in decimal, for the item. Read-only. | [optional] [default to undefined]

## Example

```typescript
import { GeoCoordinates } from './api';

const instance: GeoCoordinates = {
    altitude,
    latitude,
    longitude,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
