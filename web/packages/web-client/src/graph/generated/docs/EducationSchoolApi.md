# EducationSchoolApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addClassToSchool**](#addclasstoschool) | **POST** /v1.0/education/schools/{school-id}/classes/$ref | Assign a class to a school|
|[**addUserToSchool**](#addusertoschool) | **POST** /v1.0/education/schools/{school-id}/users/$ref | Assign a user to a school|
|[**createSchool**](#createschool) | **POST** /v1.0/education/schools | Add new school|
|[**deleteClassFromSchool**](#deleteclassfromschool) | **DELETE** /v1.0/education/schools/{school-id}/classes/{class-id}/$ref | Unassign class from a school|
|[**deleteSchool**](#deleteschool) | **DELETE** /v1.0/education/schools/{school-id} | Delete school|
|[**deleteUserFromSchool**](#deleteuserfromschool) | **DELETE** /v1.0/education/schools/{school-id}/users/{user-id}/$ref | Unassign user from a school|
|[**getSchool**](#getschool) | **GET** /v1.0/education/schools/{school-id} | Get the properties of a specific school|
|[**listSchoolClasses**](#listschoolclasses) | **GET** /v1.0/education/schools/{school-id}/classes | Get the educationClass resources owned by an educationSchool|
|[**listSchoolUsers**](#listschoolusers) | **GET** /v1.0/education/schools/{school-id}/users | Get the educationUser resources associated with an educationSchool|
|[**listSchools**](#listschools) | **GET** /v1.0/education/schools | Get a list of schools and their properties|
|[**updateSchool**](#updateschool) | **PATCH** /v1.0/education/schools/{school-id} | Update properties of a school|

# **addClassToSchool**
> addClassToSchool(classReference)


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration,
    ClassReference
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)
let classReference: ClassReference; //educationClass to be added as member

const { status, data } = await apiInstance.addClassToSchool(
    schoolId,
    classReference
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classReference** | **ClassReference**| educationClass to be added as member | |
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


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

# **addUserToSchool**
> addUserToSchool(educationUserReference)


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration,
    EducationUserReference
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)
let educationUserReference: EducationUserReference; //educationUser to be added as member

const { status, data } = await apiInstance.addUserToSchool(
    schoolId,
    educationUserReference
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationUserReference** | **EducationUserReference**| educationUser to be added as member | |
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


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

# **createSchool**
> EducationSchool createSchool(educationSchool)


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration,
    EducationSchool
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let educationSchool: EducationSchool; //New school

const { status, data } = await apiInstance.createSchool(
    educationSchool
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationSchool** | **EducationSchool**| New school | |


### Return type

**EducationSchool**

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

# **deleteClassFromSchool**
> deleteClassFromSchool()


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)
let classId: string; //key: id or externalId of the class to unassign from school (default to undefined)

const { status, data } = await apiInstance.deleteClassFromSchool(
    schoolId,
    classId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|
| **classId** | [**string**] | key: id or externalId of the class to unassign from school | defaults to undefined|


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

# **deleteSchool**
> deleteSchool()

Deletes a school. A school can only be delete if it has the terminationDate property set. And if that termination Date is in the past.

### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)

const { status, data } = await apiInstance.deleteSchool(
    schoolId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


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

# **deleteUserFromSchool**
> deleteUserFromSchool()


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)
let userId: string; //key: id or username of the user to unassign from school (default to undefined)

const { status, data } = await apiInstance.deleteUserFromSchool(
    schoolId,
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|
| **userId** | [**string**] | key: id or username of the user to unassign from school | defaults to undefined|


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

# **getSchool**
> EducationSchool getSchool()


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)

const { status, data } = await apiInstance.getSchool(
    schoolId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


### Return type

**EducationSchool**

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

# **listSchoolClasses**
> CollectionOfEducationClass listSchoolClasses()


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)

const { status, data } = await apiInstance.listSchoolClasses(
    schoolId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


### Return type

**CollectionOfEducationClass**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved classes |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listSchoolUsers**
> CollectionOfEducationUser listSchoolUsers()


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)

const { status, data } = await apiInstance.listSchoolUsers(
    schoolId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


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
|**200** | Retrieved educationUser |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listSchools**
> CollectionOfSchools listSchools()


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

const { status, data } = await apiInstance.listSchools();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CollectionOfSchools**

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

# **updateSchool**
> EducationSchool updateSchool(educationSchool)


### Example

```typescript
import {
    EducationSchoolApi,
    Configuration,
    EducationSchool
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationSchoolApi(configuration);

let schoolId: string; //key: id or schoolNumber of school (default to undefined)
let educationSchool: EducationSchool; //New property values

const { status, data } = await apiInstance.updateSchool(
    schoolId,
    educationSchool
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **educationSchool** | **EducationSchool**| New property values | |
| **schoolId** | [**string**] | key: id or schoolNumber of school | defaults to undefined|


### Return type

**EducationSchool**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

