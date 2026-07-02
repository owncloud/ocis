# FolderView

A collection of properties defining the recommended view for the folder.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sortBy** | **string** | The method by which the folder should be sorted. | [optional] [default to undefined]
**sortOrder** | **string** | If true, indicates that items should be sorted in descending order. Otherwise, items should be sorted ascending. | [optional] [default to undefined]
**viewType** | **string** | The type of view that should be used to represent the folder. | [optional] [default to undefined]

## Example

```typescript
import { FolderView } from './api';

const instance: FolderView = {
    sortBy,
    sortOrder,
    viewType,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
