# EducationClassApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addUserToClass**](#addusertoclass) | **POST** /v1.0/education/classes/{class-id}/members/$ref | Assign a user to a class|
|[**createClass**](#createclass) | **POST** /v1.0/education/classes | Add new education class|
|[**deleteClass**](#deleteclass) | **DELETE** /v1.0/education/classes/{class-id} | Delete education class|
|[**deleteUserFromClass**](#deleteuserfromclass) | **DELETE** /v1.0/education/classes/{class-id}/members/{user-id}/$ref | Unassign user from a class|
|[**getClass**](#getclass) | **GET** /v1.0/education/classes/{class-id} | Get class by key|
|[**listClassMembers**](#listclassmembers) | **GET** /v1.0/education/classes/{class-id}/members | Get the educationClass resources owned by an educationSchool|
|[**listClasses**](#listclasses) | **GET** /v1.0/education/classes | list education classes|
|[**updateClass**](#updateclass) | **PATCH** /v1.0/education/classes/{class-id} | Update properties of a education class|

# **addUserToClass**
> addUserToClass(classMemberReference)


### Example

```typescript
import {
    EducationClassApi,
    Configuration,
    ClassMemberReference
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)
let classMemberReference: ClassMemberReference; //educationUser to be added as member

const { status, data } = await apiInstance.addUserToClass(
    classId,
    classMemberReference
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classMemberReference** | **ClassMemberReference**| educationUser to be added as member | |
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createClass**
> EducationClass createClass(educationClass)


### Example

```typescript
import {
    EducationClassApi,
    Configuration,
    EducationClass
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let educationClass: EducationClass; //New entity

const { status, data } = await apiInstance.createClass(
    educationClass
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationClass** | **EducationClass**| New entity | |


### Return type

**EducationClass**

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

# **deleteClass**
> deleteClass()


### Example

```typescript
import {
    EducationClassApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)

const { status, data } = await apiInstance.deleteClass(
    classId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|


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

# **deleteUserFromClass**
> deleteUserFromClass()


### Example

```typescript
import {
    EducationClassApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)
let userId: string; //key: id or username of the user to unassign from class (default to undefined)

const { status, data } = await apiInstance.deleteUserFromClass(
    classId,
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|
| **userId** | [**string**] | key: id or username of the user to unassign from class | defaults to undefined|


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

# **getClass**
> EducationClass getClass()


### Example

```typescript
import {
    EducationClassApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)

const { status, data } = await apiInstance.getClass(
    classId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|


### Return type

**EducationClass**

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

# **listClassMembers**
> CollectionOfEducationUser listClassMembers()


### Example

```typescript
import {
    EducationClassApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)

const { status, data } = await apiInstance.listClassMembers(
    classId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|


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
|**200** | Retrieved class members |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listClasses**
> CollectionOfClass listClasses()


### Example

```typescript
import {
    EducationClassApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

const { status, data } = await apiInstance.listClasses();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CollectionOfClass**

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

# **updateClass**
> EducationClass updateClass(educationClass)


### Example

```typescript
import {
    EducationClassApi,
    Configuration,
    EducationClass
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)
let educationClass: EducationClass; //New property values

const { status, data } = await apiInstance.updateClass(
    classId,
    educationClass
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationClass** | **EducationClass**| New property values | |
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|


### Return type

**EducationClass**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | New property values |  -  |
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

