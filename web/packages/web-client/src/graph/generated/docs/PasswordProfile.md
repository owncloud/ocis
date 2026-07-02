# PasswordProfile

Password Profile associated with a user

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**forceChangePasswordNextSignIn** | **boolean** | If true the user is required to change their password upon the next login | [optional] [default to false]
**password** | **string** | The user\&#39;s password | [optional] [default to undefined]

## Example

```typescript
import { PasswordProfile } from './api';

const instance: PasswordProfile = {
    forceChangePasswordNextSignIn,
    password,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
