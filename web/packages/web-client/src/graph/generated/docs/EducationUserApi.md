# EducationUserApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createEducationUser**](#createeducationuser) | **POST** /v1.0/education/users | Add new education user|
|[**deleteEducationUser**](#deleteeducationuser) | **DELETE** /v1.0/education/users/{user-id} | Delete educationUser|
|[**getEducationUser**](#geteducationuser) | **GET** /v1.0/education/users/{user-id} | Get properties of educationUser|
|[**listEducationUsers**](#listeducationusers) | **GET** /v1.0/education/users | Get entities from education users|
|[**updateEducationUser**](#updateeducationuser) | **PATCH** /v1.0/education/users/{user-id} | Update properties of educationUser|

# **createEducationUser**
> EducationUser createEducationUser(educationUser)


### Example

```typescript
import {
    EducationUserApi,
    Configuration,
    EducationUser
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationUserApi(configuration);

let educationUser: EducationUser; //New entity

const { status, data } = await apiInstance.createEducationUser(
    educationUser
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationUser** | **EducationUser**| New entity | |


### Return type

**EducationUser**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteEducationUser**
> deleteEducationUser()


### Example

```typescript
import {
    EducationUserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationUserApi(configuration);

let userId: string; //key: id or username of user (default to undefined)

const { status, data } = await apiInstance.deleteEducationUser(
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**string**] | key: id or username of user | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getEducationUser**
> EducationUser getEducationUser()


### Example

```typescript
import {
    EducationUserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationUserApi(configuration);

let userId: string; //key: id or username of user (default to undefined)
let $expand: Set<'memberOf'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.getEducationUser(
    userId,
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**string**] | key: id or username of user | defaults to undefined|
| **$expand** | **Array<&#39;memberOf&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**EducationUser**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listEducationUsers**
> CollectionOfEducationUser listEducationUsers()


### Example

```typescript
import {
    EducationUserApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationUserApi(configuration);

let $orderby: Set<'displayName' | 'displayName desc' | 'mail' | 'mail desc' | 'onPremisesSamAccountName' | 'onPremisesSamAccountName desc'>; //Order items by property values (optional) (default to undefined)
let $expand: Set<'memberOf'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.listEducationUsers(
    $orderby,
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **$orderby** | **Array<&#39;displayName&#39; &#124; &#39;displayName desc&#39; &#124; &#39;mail&#39; &#124; &#39;mail desc&#39; &#124; &#39;onPremisesSamAccountName&#39; &#124; &#39;onPremisesSamAccountName desc&#39;>** | Order items by property values | (optional) defaults to undefined|
| **$expand** | **Array<&#39;memberOf&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**CollectionOfEducationUser**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entities |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateEducationUser**
> EducationUser updateEducationUser(educationUser)


### Example

```typescript
import {
    EducationUserApi,
    Configuration,
    EducationUser
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationUserApi(configuration);

let userId: string; //key: id or username of user (default to undefined)
let educationUser: EducationUser; //New property values

const { status, data } = await apiInstance.updateEducationUser(
    userId,
    educationUser
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationUser** | **EducationUser**| New property values | |
| **userId** | [**string**] | key: id or username of user | defaults to undefined|


### Return type

**EducationUser**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

