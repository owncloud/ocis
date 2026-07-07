# OpenGraphFile

File metadata, if the item is a file. Read-only.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**hashes** | [**Hashes**](Hashes.md) |  | [optional] [default to undefined]
**mimeType** | **string** | The MIME type for the file. This is determined by logic on the server and might not be the value provided when the file was uploaded. Read-only. | [optional] [readonly] [default to undefined]
**processingMetadata** | **boolean** |  | [optional] [default to undefined]

## Example

```typescript
import { OpenGraphFile } from './api';

const instance: OpenGraphFile = {
    hashes,
    mimeType,
    processingMetadata,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
