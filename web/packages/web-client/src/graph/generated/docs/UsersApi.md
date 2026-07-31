# UsersApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createUser**](#createuser) | **POST** /v1.0/users | Add new entity to users|
|[**listUsers**](#listusers) | **GET** /v1.0/users | Get entities from users|

# **createUser**
> User createUser(user)


### Example

```typescript
import {
    UsersApi,
    Configuration,
    User
} from './api';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let user: User; //New entity

const { status, data } = await apiInstance.createUser(
    user
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **user** | **User**| New entity | |


### Return type

**User**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listUsers**
> CollectionOfUser listUsers()


### Example

```typescript
import {
    UsersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let $search: string; //Search items by search phrases (optional) (default to undefined)
let $filter: string; //Filter users by property values and relationship attributes (optional) (default to undefined)
let $orderby: Set<'displayName' | 'displayName desc' | 'mail' | 'mail desc' | 'onPremisesSamAccountName' | 'onPremisesSamAccountName desc'>; //Order items by property values (optional) (default to undefined)
let $select: Set<'id' | 'displayName' | 'mail' | 'memberOf' | 'onPremisesSamAccountName' | 'surname'>; //Select properties to be returned (optional) (default to undefined)
let $expand: Set<'drive' | 'drives' | 'memberOf' | 'appRoleAssignments'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.listUsers(
    $search,
    $filter,
    $orderby,
    $select,
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **$search** | [**string**] | Search items by search phrases | (optional) defaults to undefined|
| **$filter** | [**string**] | Filter users by property values and relationship attributes | (optional) defaults to undefined|
| **$orderby** | **Array<&#39;displayName&#39; &#124; &#39;displayName desc&#39; &#124; &#39;mail&#39; &#124; &#39;mail desc&#39; &#124; &#39;onPremisesSamAccountName&#39; &#124; &#39;onPremisesSamAccountName desc&#39;>** | Order items by property values | (optional) defaults to undefined|
| **$select** | **Array<&#39;id&#39; &#124; &#39;displayName&#39; &#124; &#39;mail&#39; &#124; &#39;memberOf&#39; &#124; &#39;onPremisesSamAccountName&#39; &#124; &#39;surname&#39;>** | Select properties to be returned | (optional) defaults to undefined|
| **$expand** | **Array<&#39;drive&#39; &#124; &#39;drives&#39; &#124; &#39;memberOf&#39; &#124; &#39;appRoleAssignments&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**CollectionOfUser**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entities |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

