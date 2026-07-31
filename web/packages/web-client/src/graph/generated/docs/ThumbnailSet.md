# ThumbnailSet

The ThumbnailSet resource is a keyed collection of thumbnail resources. It\'s used to represent a set of thumbnails associated with a DriveItem. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | The ID within the item. Read-only. | [optional] [default to undefined]
**large** | [**Thumbnail**](Thumbnail.md) |  | [optional] [default to undefined]
**medium** | [**Thumbnail**](Thumbnail.md) |  | [optional] [default to undefined]
**small** | [**Thumbnail**](Thumbnail.md) |  | [optional] [default to undefined]
**source** | [**Thumbnail**](Thumbnail.md) |  | [optional] [default to undefined]

## Example

```typescript
import { ThumbnailSet } from './api';

const instance: ThumbnailSet = {
    id,
    large,
    medium,
    small,
    source,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
