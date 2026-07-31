# FileSystemInfo

File system information on client. Read-write.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdDateTime** | **string** | The UTC date and time the file was created on a client. | [optional] [default to undefined]
**lastAccessedDateTime** | **string** | The UTC date and time the file was last accessed. Available for the recent file list only. | [optional] [default to undefined]
**lastModifiedDateTime** | **string** | The UTC date and time the file was last modified on a client. | [optional] [default to undefined]

## Example

```typescript
import { FileSystemInfo } from './api';

const instance: FileSystemInfo = {
    createdDateTime,
    lastAccessedDateTime,
    lastModifiedDateTime,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
