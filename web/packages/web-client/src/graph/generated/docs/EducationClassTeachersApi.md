# EducationClassTeachersApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addTeacherToClass**](EducationClassTeachersApi.md#addteachertoclass) | **POST** /v1.0/education/classes/{class-id}/teachers/$ref | Assign a teacher to a class |
| [**deleteTeacherFromClass**](EducationClassTeachersApi.md#deleteteacherfromclass) | **DELETE** /v1.0/education/classes/{class-id}/teachers/{user-id}/$ref | Unassign user as teacher of a class |
| [**getTeachers**](EducationClassTeachersApi.md#getteachers) | **GET** /v1.0/education/classes/{class-id}/teachers | Get the teachers for a class |



## addTeacherToClass

> addTeacherToClass(classId, classTeacherReference)

Assign a teacher to a class

### Example

```ts
import {
  Configuration,
  EducationClassTeachersApi,
} from '';
import type { AddTeacherToClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassTeachersApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
    // ClassTeacherReference | educationUser to be added as teacher
    classTeacherReference: ...,
  } satisfies AddTeacherToClassRequest;

  try {
    const data = await api.addTeacherToClass(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |
| **classTeacherReference** | [ClassTeacherReference](ClassTeacherReference.md) | educationUser to be added as teacher | |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteTeacherFromClass

> deleteTeacherFromClass(classId, userId)

Unassign user as teacher of a class

### Example

```ts
import {
  Configuration,
  EducationClassTeachersApi,
} from '';
import type { DeleteTeacherFromClassRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassTeachersApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: classId_example,
    // string | key: id or username of the user to unassign as teacher
    userId: 90eedea1-dea1-90ee-a1de-ee90a1deee90,
  } satisfies DeleteTeacherFromClassRequest;

  try {
    const data = await api.deleteTeacherFromClass(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |
| **userId** | `string` | key: id or username of the user to unassign as teacher | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getTeachers

> CollectionOfEducationUser getTeachers(classId)

Get the teachers for a class

### Example

```ts
import {
  Configuration,
  EducationClassTeachersApi,
} from '';
import type { GetTeachersRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new EducationClassTeachersApi(config);

  const body = {
    // string | key: id or externalId of class
    classId: 86948e45-96a6-43df-b83d-46e92afd30de,
  } satisfies GetTeachersRequest;

  try {
    const data = await api.getTeachers(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **classId** | `string` | key: id or externalId of class | [Defaults to `undefined`] |

### Return type

[**CollectionOfEducationUser**](CollectionOfEducationUser.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved class teachers |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

