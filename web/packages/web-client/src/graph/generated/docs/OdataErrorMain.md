# OdataErrorMain


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**code** | **string** |  | [default to undefined]
**message** | **string** |  | [default to undefined]
**target** | **string** |  | [optional] [default to undefined]
**details** | [**Array&lt;OdataErrorDetail&gt;**](OdataErrorDetail.md) |  | [optional] [default to undefined]
**innererror** | **object** | The structure of this object is service-specific | [optional] [default to undefined]

## Example

```typescript
import { OdataErrorMain } from './api';

const instance: OdataErrorMain = {
    code,
    message,
    target,
    details,
    innererror,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
