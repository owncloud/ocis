# Instance

An oCIS instance that the user is either a member or a guest of.

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**url** | **string** | The URL of the oCIS instance. | [optional] [default to undefined]
**primary** | **boolean** | Whether the instance is the user\&#39;s primary instance. | [optional] [default to undefined]

## Example

```typescript
import { Instance } from './api';

const instance: Instance = {
    url,
    primary,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
