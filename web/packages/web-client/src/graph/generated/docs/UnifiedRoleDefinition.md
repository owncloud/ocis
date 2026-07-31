# UnifiedRoleDefinition

A role definition is a collection of permissions in libre graph listing the operations that can be performed and the resources against which they can performed. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **string** | The description for the unifiedRoleDefinition. | [optional] [default to undefined]
**displayName** | **string** | The display name for the unifiedRoleDefinition. Required. Supports $filter (&#x60;eq&#x60;, &#x60;in&#x60;). | [optional] [default to undefined]
**id** | **string** | The unique identifier for the role definition. Key, not nullable, Read-only. Inherited from entity. Supports $filter (&#x60;eq&#x60;, &#x60;in&#x60;). | [optional] [default to undefined]
**rolePermissions** | [**Array&lt;UnifiedRolePermission&gt;**](UnifiedRolePermission.md) | List of permissions included in the role. | [optional] [default to undefined]
**libre_graph_weight** | **number** | When presenting a list of roles the weight can be used to order them in a meaningful way. Lower weight gets higher precedence. So content with lower weight will come first. If set, weights should be non-zero, as 0 is interpreted as an unset weight.  | [optional] [default to undefined]

## Example

```typescript
import { UnifiedRoleDefinition } from './api';

const instance: UnifiedRoleDefinition = {
    description,
    displayName,
    id,
    rolePermissions,
    libre_graph_weight,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
