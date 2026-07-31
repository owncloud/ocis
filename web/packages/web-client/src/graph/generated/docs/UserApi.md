# UserApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**deleteUser**](#deleteuser) | **DELETE** /v1.0/users/{user-id} | Delete entity from users|
|[**exportPersonalData**](#exportpersonaldata) | **POST** /v1.0/users/{user-id}/exportPersonalData | export personal data of a user|
|[**getUser**](#getuser) | **GET** /v1.0/users/{user-id} | Get entity from users by key|
|[**updateUser**](#updateuser) | **PATCH** /v1.0/users/{user-id} | Update entity in users|

# **deleteUser**
> deleteUser()


### Example

```typescript
import {
    UserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let userId: string; //key: id or name of user (default to undefined)
let ifMatch: string; //ETag (optional) (default to undefined)

const { status, data } = await apiInstance.deleteUser(
    userId,
    ifMatch
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**string**] | key: id or name of user | defaults to undefined|
| **ifMatch** | [**string**] | ETag | (optional) defaults to undefined|


### Return type

void (empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **exportPersonalData**
> exportPersonalData()


### Example

```typescript
import {
    UserApi,
    Configuration,
    ExportPersonalDataRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let userId: string; //key: id or name of user (default to undefined)
let exportPersonalDataRequest: ExportPersonalDataRequest; //destination the file should be created at (optional)

const { status, data } = await apiInstance.exportPersonalData(
    userId,
    exportPersonalDataRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **exportPersonalDataRequest** | **ExportPersonalDataRequest**| destination the file should be created at | |
| **userId** | [**string**] | key: id or name of user | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**202** | success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUser**
> User getUser()


### Example

```typescript
import {
    UserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let userId: string; //key: id or name of user (default to undefined)
let $select: Set<'id' | 'displayName' | 'drive' | 'drives' | 'mail' | 'memberOf' | 'onPremisesSamAccountName' | 'surname'>; //Select properties to be returned (optional) (default to undefined)
let $expand: Set<'drive' | 'drives' | 'memberOf' | 'appRoleAssignments'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.getUser(
    userId,
    $select,
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**string**] | key: id or name of user | defaults to undefined|
| **$select** | **Array<&#39;id&#39; &#124; &#39;displayName&#39; &#124; &#39;drive&#39; &#124; &#39;drives&#39; &#124; &#39;mail&#39; &#124; &#39;memberOf&#39; &#124; &#39;onPremisesSamAccountName&#39; &#124; &#39;surname&#39;>** | Select properties to be returned | (optional) defaults to undefined|
| **$expand** | **Array<&#39;drive&#39; &#124; &#39;drives&#39; &#124; &#39;memberOf&#39; &#124; &#39;appRoleAssignments&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**User**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateUser**
> User updateUser(userUpdate)


### Example

```typescript
import {
    UserApi,
    Configuration,
    UserUpdate
} from './api';

const configuration = new Configuration();
const apiInstance = new UserApi(configuration);

let userId: string; //key: id of user (default to undefined)
let userUpdate: UserUpdate; //New property values

const { status, data } = await apiInstance.updateUser(
    userId,
    userUpdate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userUpdate** | **UserUpdate**| New property values | |
| **userId** | [**string**] | key: id of user | defaults to undefined|


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
|**200** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

