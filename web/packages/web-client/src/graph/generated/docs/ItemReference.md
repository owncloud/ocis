# ItemReference


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**driveId** | **string** | Unique identifier of the drive instance that contains the item. Read-only. | [optional] [readonly] [default to undefined]
**driveType** | **string** | Identifies the type of drive. See [drive][] resource for values. Read-only. | [optional] [readonly] [default to undefined]
**id** | **string** | Unique identifier of the item in the drive. Read-only. | [optional] [readonly] [default to undefined]
**name** | **string** | The name of the item being referenced. Read-only. | [optional] [readonly] [default to undefined]
**path** | **string** | Path that can be used to navigate to the item. Read-only. | [optional] [readonly] [default to undefined]

## Example

```typescript
import { ItemReference } from './api';

const instance: ItemReference = {
    driveId,
    driveType,
    id,
    name,
    path,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
