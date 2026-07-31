# AppRoleAssignment


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | The unique identifier for the object. 12345678-9abc-def0-1234-56789abcde. The value of the ID property is often, but not exclusively, in the form of a GUID. The value should be treated as an opaque identifier and not based in being a GUID. Null values are not allowed. Read-only. | [optional] [readonly] [default to undefined]
**deletedDateTime** | **string** |  | [optional] [default to undefined]
**appRoleId** | **string** | The identifier (id) for the app role which is assigned to the user. Required on create. | [default to undefined]
**createdDateTime** | **string** | The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. | [optional] [default to undefined]
**principalDisplayName** | **string** | The display name of the user, group, or service principal that was granted the app role assignment. Read-only. | [optional] [default to undefined]
**principalId** | **string** | The unique identifier (id) for the user, security group, or service principal being granted the app role. Security groups with dynamic memberships are supported. Required on create. | [default to undefined]
**principalType** | **string** | The type of the assigned principal. This can either be User, Group, or ServicePrincipal. Read-only. | [optional] [default to undefined]
**resourceDisplayName** | **string** | The display name of the resource app\&#39;s service principal to which the assignment is made. | [optional] [default to undefined]
**resourceId** | **string** | The unique identifier (id) for the resource service principal for which the assignment is made. Required on create. | [default to undefined]

## Example

```typescript
import { AppRoleAssignment } from './api';

const instance: AppRoleAssignment = {
    id,
    deletedDateTime,
    appRoleId,
    createdDateTime,
    principalDisplayName,
    principalId,
    principalType,
    resourceDisplayName,
    resourceId,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
