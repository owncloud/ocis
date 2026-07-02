# EducationClass

And extension of group representing a class or course

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Read-only. | [optional] [readonly] [default to undefined]
**description** | **string** | An optional description for the group. Returned by default. | [optional] [default to undefined]
**displayName** | **string** | The display name for the group. This property is required when a group is created and cannot be cleared during updates. Returned by default. Supports $search and $orderBy. | [optional] [default to undefined]
**members** | [**Array&lt;User&gt;**](User.md) | Users and groups that are members of this group. HTTP Methods: GET (supported for all groups), Nullable. Supports $expand. | [optional] [default to undefined]
**membersodata_bind** | **Set&lt;string&gt;** | A list of member references to the members to be added. Up to 20 members can be added with a single request | [optional] [default to undefined]
**classification** | **string** | Classification of the group, i.e. \&quot;class\&quot; or \&quot;course\&quot; | [optional] [default to undefined]
**externalId** | **string** | An external unique ID for the class | [optional] [default to undefined]

## Example

```typescript
import { EducationClass } from './api';

const instance: EducationClass = {
    id,
    description,
    displayName,
    members,
    membersodata_bind,
    classification,
    externalId,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
