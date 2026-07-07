# Identity


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**displayName** | **string** | The identity\&#39;s display name. Note that this may not always be available or up to date. For example, if a user changes their display name, the API may show the new value in a future response, but the items associated with the user won\&#39;t show up as having changed when using delta. | [default to undefined]
**id** | **string** | Unique identifier for the identity. | [optional] [default to undefined]
**libre_graph_userType** | **string** | The type of the identity. This can be either \&quot;Member\&quot; for regular user, \&quot;Guest\&quot; for guest users or \&quot;Federated\&quot; for users imported from a federated instance. Can be used by clients to indicate the type of user. For more details, clients should look up and cache the user at the /users endpoint. | [optional] [default to undefined]

## Example

```typescript
import { Identity } from './api';

const instance: Identity = {
    displayName,
    id,
    libre_graph_userType,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
