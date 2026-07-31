# Quota

Optional. Information about the drive\'s storage space quota. Read-only.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**deleted** | **number** | Total space consumed by files in the recycle bin, in bytes. Read-only. | [optional] [readonly] [default to undefined]
**remaining** | **number** | Total space remaining before reaching the quota limit, in bytes. Read-only. | [optional] [readonly] [default to undefined]
**state** | **string** | Enumeration value that indicates the state of the storage space. Either \&quot;normal\&quot;, \&quot;nearing\&quot;, \&quot;critical\&quot; or \&quot;exceeded\&quot;. Read-only. | [optional] [readonly] [default to undefined]
**total** | **number** | Total allowed storage space, in bytes. Read-only. | [optional] [readonly] [default to undefined]
**used** | **number** | Total space used, in bytes. Read-only. | [optional] [readonly] [default to undefined]

## Example

```typescript
import { Quota } from './api';

const instance: Quota = {
    deleted,
    remaining,
    state,
    total,
    used,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
