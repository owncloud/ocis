# SharePointIdentitySet

This resource is used to represent a set of identities associated with various events for an item, such as created by or last modified by.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**user** | [**Identity**](Identity.md) |  | [optional] [default to undefined]
**group** | [**Identity**](Identity.md) |  | [optional] [default to undefined]

## Example

```typescript
import { SharePointIdentitySet } from './api';

const instance: SharePointIdentitySet = {
    user,
    group,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
