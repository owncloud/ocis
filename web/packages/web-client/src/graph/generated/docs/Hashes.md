# Hashes

Hashes of the file\'s binary content, if available. Read-only.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**crc32Hash** | **string** | The CRC32 value of the file (if available). Read-only. | [optional] [default to undefined]
**quickXorHash** | **string** | A proprietary hash of the file that can be used to determine if the contents of the file have changed (if available). Read-only. | [optional] [default to undefined]
**sha1Hash** | **string** | SHA1 hash for the contents of the file (if available). Read-only. | [optional] [default to undefined]
**sha256Hash** | **string** | SHA256 hash for the contents of the file (if available). Read-only. | [optional] [default to undefined]

## Example

```typescript
import { Hashes } from './api';

const instance: Hashes = {
    crc32Hash,
    quickXorHash,
    sha1Hash,
    sha256Hash,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
