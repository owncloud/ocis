# ObjectIdentity

Represents an identity used to sign in to a user account

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**issuer** | **string** | domain of the Provider issuing the identity | [optional] [default to undefined]
**issuerAssignedId** | **string** | The unique id assigned by the issuer to the account | [optional] [default to undefined]

## Example

```typescript
import { ObjectIdentity } from './api';

const instance: ObjectIdentity = {
    issuer,
    issuerAssignedId,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
