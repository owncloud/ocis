# EducationClassTeachersApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addTeacherToClass**](#addteachertoclass) | **POST** /v1.0/education/classes/{class-id}/teachers/$ref | Assign a teacher to a class|
|[**deleteTeacherFromClass**](#deleteteacherfromclass) | **DELETE** /v1.0/education/classes/{class-id}/teachers/{user-id}/$ref | Unassign user as teacher of a class|
|[**getTeachers**](#getteachers) | **GET** /v1.0/education/classes/{class-id}/teachers | Get the teachers for a class|

# **addTeacherToClass**
> addTeacherToClass(classTeacherReference)


### Example

```typescript
import {
    EducationClassTeachersApi,
    Configuration,
    ClassTeacherReference
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassTeachersApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)
let classTeacherReference: ClassTeacherReference; //educationUser to be added as teacher

const { status, data } = await apiInstance.addTeacherToClass(
    classId,
    classTeacherReference
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classTeacherReference** | **ClassTeacherReference**| educationUser to be added as teacher | |
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

# **deleteTeacherFromClass**
> deleteTeacherFromClass()


### Example

```typescript
import {
    EducationClassTeachersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassTeachersApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)
let userId: string; //key: id or username of the user to unassign as teacher (default to undefined)

const { status, data } = await apiInstance.deleteTeacherFromClass(
    classId,
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **classId** | [**string**] | key: id or externalId of class | defaults to undefined|
| **userId** | [**string**] | key: id or username of the user to unassign as teacher | defaults to undefined|


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

# **getTeachers**
> CollectionOfEducationUser getTeachers()


### Example

```typescript
import {
    EducationClassTeachersApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new EducationClassTeachersApi(configuration);

let classId: string; //key: id or externalId of class (default to undefined)

const { status, data } = await apiInstance.getTeachers(
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
|**200** | Retrieved class teachers |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

