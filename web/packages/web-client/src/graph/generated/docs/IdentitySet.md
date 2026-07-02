# IdentitySet

Optional. User account.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**application** | [**Identity**](Identity.md) |  | [optional] [default to undefined]
**device** | [**Identity**](Identity.md) |  | [optional] [default to undefined]
**user** | [**Identity**](Identity.md) |  | [optional] [default to undefined]
**group** | [**Identity**](Identity.md) |  | [optional] [default to undefined]

## Example

```typescript
import { IdentitySet } from './api';

const instance: IdentitySet = {
    application,
    device,
    user,
    group,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
