---
title: "GRPC API"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: grpc.md
---

{{< toc >}}

## settings.proto

### AddSettingToBundleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle_id | [string](#string) |  |  |
| setting | [Setting](#setting) |  |  |

### AddSettingToBundleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| setting | [Setting](#setting) |  |  |

### AssignRoleToUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_uuid | [string](#string) |  |  |
| role_id | [string](#string) |  | the role_id is a bundle_id internally |

### AssignRoleToUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignment | [UserRoleAssignment](#userroleassignment) |  |  |

### Bool



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default | [bool](#bool) |  |  |
| label | [string](#string) |  |  |

### Bundle



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| type | [Bundle.Type](#bundletype) |  |  |
| extension | [string](#string) |  |  |
| display_name | [string](#string) |  |  |
| settings | [Setting](#setting) | repeated |  |
| resource | [Resource](#resource) |  |  |

### GetBundleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle_id | [string](#string) |  |  |

### GetBundleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle | [Bundle](#bundle) |  |  |

### GetPermissionByIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| permission_id | [string](#string) |  |  |

### GetPermissionByIDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| permission | [Permission](#permission) |  |  |

### GetValueByUniqueIdentifiersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_uuid | [string](#string) |  |  |
| setting_id | [string](#string) |  |  |

### GetValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### GetValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [ValueWithIdentifier](#valuewithidentifier) |  |  |

### Identifier



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| extension | [string](#string) |  |  |
| bundle | [string](#string) |  |  |
| setting | [string](#string) |  |  |

### Int



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default | [int64](#int64) |  |  |
| min | [int64](#int64) |  |  |
| max | [int64](#int64) |  |  |
| step | [int64](#int64) |  |  |
| placeholder | [string](#string) |  |  |

### ListBundlesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle_ids | [string](#string) | repeated |  |

### ListBundlesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundles | [Bundle](#bundle) | repeated |  |

### ListOption



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [ListOptionValue](#listoptionvalue) |  |  |
| default | [bool](#bool) |  |  |
| display_value | [string](#string) |  |  |

### ListOptionValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| string_value | [string](#string) |  |  |
| int_value | [int64](#int64) |  |  |

### ListPermissionsByResourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource | [Resource](#resource) |  |  |

### ListPermissionsByResourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| permissions | [Permission](#permission) | repeated |  |

### ListRoleAssignmentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account_uuid | [string](#string) |  |  |

### ListRoleAssignmentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignments | [UserRoleAssignment](#userroleassignment) | repeated |  |

### ListValue



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [ListOptionValue](#listoptionvalue) | repeated |  |

### ListValuesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle_id | [string](#string) |  |  |
| account_uuid | [string](#string) |  |  |

### ListValuesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [ValueWithIdentifier](#valuewithidentifier) | repeated |  |

### MultiChoiceList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| options | [ListOption](#listoption) | repeated |  |

### Permission



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [Permission.Operation](#permissionoperation) |  |  |
| constraint | [Permission.Constraint](#permissionconstraint) |  |  |

### RemoveRoleFromUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### RemoveSettingFromBundleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle_id | [string](#string) |  |  |
| setting_id | [string](#string) |  |  |

### Resource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [Resource.Type](#resourcetype) |  |  |
| id | [string](#string) |  |  |

### SaveBundleRequest

---
requests and responses for settings bundles
---

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle | [Bundle](#bundle) |  |  |

### SaveBundleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bundle | [Bundle](#bundle) |  |  |

### SaveValueRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [Value](#value) |  |  |

### SaveValueResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [ValueWithIdentifier](#valuewithidentifier) |  |  |

### Setting



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| display_name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| int_value | [Int](#int) |  |  |
| string_value | [String](#string) |  |  |
| bool_value | [Bool](#bool) |  |  |
| single_choice_value | [SingleChoiceList](#singlechoicelist) |  |  |
| multi_choice_value | [MultiChoiceList](#multichoicelist) |  |  |
| permission_value | [Permission](#permission) |  |  |
| resource | [Resource](#resource) |  |  |

### SingleChoiceList



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| options | [ListOption](#listoption) | repeated |  |

### String



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default | [string](#string) |  |  |
| required | [bool](#bool) |  |  |
| min_length | [int32](#int32) |  |  |
| max_length | [int32](#int32) |  |  |
| placeholder | [string](#string) |  |  |

### UserRoleAssignment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is generated upon saving the assignment |
| account_uuid | [string](#string) |  |  |
| role_id | [string](#string) |  | the role_id is a bundle_id internally |

### Value



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id is the id of the Value. It is generated on saving it. |
| bundle_id | [string](#string) |  |  |
| setting_id | [string](#string) |  | setting_id is the id of the setting from within its bundle. |
| account_uuid | [string](#string) |  |  |
| resource | [Resource](#resource) |  |  |
| bool_value | [bool](#bool) |  |  |
| int_value | [int64](#int64) |  |  |
| string_value | [string](#string) |  |  |
| list_value | [ListValue](#listvalue) |  |  |

### ValueWithIdentifier



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| identifier | [Identifier](#identifier) |  |  |
| value | [Value](#value) |  |  |

### Bundle.Type



| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNKNOWN | 0 |  |
| TYPE_DEFAULT | 1 |  |
| TYPE_ROLE | 2 |  |
### Permission.Constraint



| Name | Number | Description |
| ---- | ------ | ----------- |
| CONSTRAINT_UNKNOWN | 0 |  |
| CONSTRAINT_OWN | 1 |  |
| CONSTRAINT_SHARED | 2 |  |
| CONSTRAINT_ALL | 3 |  |
### Permission.Operation



| Name | Number | Description |
| ---- | ------ | ----------- |
| OPERATION_UNKNOWN | 0 |  |
| OPERATION_CREATE | 1 |  |
| OPERATION_READ | 2 |  |
| OPERATION_UPDATE | 3 |  |
| OPERATION_DELETE | 4 |  |
| OPERATION_WRITE | 5 | WRITE is a combination of CREATE and UPDATE |
| OPERATION_READWRITE | 6 | READWRITE is a combination of READ and WRITE |
### Resource.Type



| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNKNOWN | 0 |  |
| TYPE_SYSTEM | 1 |  |
| TYPE_FILE | 2 |  |
| TYPE_SHARE | 3 |  |
| TYPE_SETTING | 4 |  |
| TYPE_BUNDLE | 5 |  |
| TYPE_USER | 6 |  |
| TYPE_GROUP | 7 |  |

### BundleService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SaveBundle | [SaveBundleRequest](#savebundlerequest) | [SaveBundleResponse](#savebundleresponse) |  |
| GetBundle | [GetBundleRequest](#getbundlerequest) | [GetBundleResponse](#getbundleresponse) |  |
| ListBundles | [ListBundlesRequest](#listbundlesrequest) | [ListBundlesResponse](#listbundlesresponse) |  |
| AddSettingToBundle | [AddSettingToBundleRequest](#addsettingtobundlerequest) | [AddSettingToBundleResponse](#addsettingtobundleresponse) |  |
| RemoveSettingFromBundle | [RemoveSettingFromBundleRequest](#removesettingfrombundlerequest) | [.google.protobuf.Empty](#googleprotobufempty) |  |

### PermissionService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListPermissionsByResource | [ListPermissionsByResourceRequest](#listpermissionsbyresourcerequest) | [ListPermissionsByResourceResponse](#listpermissionsbyresourceresponse) |  |
| GetPermissionByID | [GetPermissionByIDRequest](#getpermissionbyidrequest) | [GetPermissionByIDResponse](#getpermissionbyidresponse) |  |

### RoleService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListRoles | [ListBundlesRequest](#listbundlesrequest) | [ListBundlesResponse](#listbundlesresponse) |  |
| ListRoleAssignments | [ListRoleAssignmentsRequest](#listroleassignmentsrequest) | [ListRoleAssignmentsResponse](#listroleassignmentsresponse) |  |
| AssignRoleToUser | [AssignRoleToUserRequest](#assignroletouserrequest) | [AssignRoleToUserResponse](#assignroletouserresponse) |  |
| RemoveRoleFromUser | [RemoveRoleFromUserRequest](#removerolefromuserrequest) | [.google.protobuf.Empty](#googleprotobufempty) |  |

### ValueService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SaveValue | [SaveValueRequest](#savevaluerequest) | [SaveValueResponse](#savevalueresponse) |  |
| GetValue | [GetValueRequest](#getvaluerequest) | [GetValueResponse](#getvalueresponse) |  |
| ListValues | [ListValuesRequest](#listvaluesrequest) | [ListValuesResponse](#listvaluesresponse) |  |
| GetValueByUniqueIdentifiers | [GetValueByUniqueIdentifiersRequest](#getvaluebyuniqueidentifiersrequest) | [GetValueResponse](#getvalueresponse) |  |

## Scalar Value Types

| .proto Type | Notes | C++ | Java |
| ----------- | ----- | --- | ---- |
| {{< div id="double" content="double" >}} |  | double | double |
| {{< div id="float" content="float" >}} |  | float | float |
| {{< div id="int32" content="int32" >}} | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int |
| {{< div id="int64" content="int64" >}} | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long |
| {{< div id="uint32" content="uint32" >}} | Uses variable-length encoding. | uint32 | int |
| {{< div id="uint64" content="uint64" >}} | Uses variable-length encoding. | uint64 | long |
| {{< div id="sint32" content="sint32" >}} | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int |
| {{< div id="sint64" content="sint64" >}} | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long |
| {{< div id="fixed32" content="fixed32" >}} | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int |
| {{< div id="fixed64" content="fixed64" >}} | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long |
| {{< div id="sfixed32" content="sfixed32" >}} | Always four bytes. | int32 | int |
| {{< div id="sfixed64" content="sfixed64" >}} | Always eight bytes. | int64 | long |
| {{< div id="bool" content="bool" >}} |  | bool | boolean |
| {{< div id="string" content="string" >}} | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String |
| {{< div id="bytes" content="bytes" >}} | May contain any arbitrary sequence of bytes. | string | ByteString |
