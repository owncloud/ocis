# Application


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | The unique identifier for the object. 12345678-9abc-def0-1234-56789abcde. The value of the ID property is often, but not exclusively, in the form of a GUID. The value should be treated as an opaque identifier and not based in being a GUID. Null values are not allowed. Read-only. | [readonly] [default to undefined]
**appRoles** | [**Array&lt;AppRole&gt;**](AppRole.md) | The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable. | [optional] [default to undefined]
**displayName** | **string** | The display name for the application. | [optional] [default to undefined]

## Example

```typescript
import { Application } from './api';

const instance: Application = {
    id,
    appRoles,
    displayName,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
