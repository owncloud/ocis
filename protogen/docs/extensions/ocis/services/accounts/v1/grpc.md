---
title: "GRPC API"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: grpc.md
---

{{< toc >}}



## ocis/services/accounts/v1/accounts.proto

### AddMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_id | [string](#string) |  | The id of the group to add a member to |
| account_id | [string](#string) |  | The account id to add |

### CreateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) |  | The account resource to create |

### CreateGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) |  | The account resource to create |

### DeleteAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### DeleteGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### GetAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### GetGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |

### ListAccountsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_size | [int32](#int32) |  | Optional. The maximum number of accounts to return in the response |
| page_token | [string](#string) |  | Optional. A pagination token returned from a previous call to `Get` that indicates from where search should continue |
| field_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | Optional. Used to specify a subset of fields that should be returned by a get operation or modified by an update operation. |
| query | [string](#string) |  | Optional. Search criteria used to select the accounts to return. If no search criteria is specified then all accounts will be returned

TODO update query language Query expressions can be used to restrict results based upon the account properties where the operators `=`, `NOT`, `AND` and `OR` can be used along with the suffix wildcard symbol `*`.

The string properties in a query expression should use escaped quotes for values that include whitespace to prevent unexpected behavior.

Some example queries are:

* Query `display_name=Th*` returns accounts whose display_name starts with "Th" * Query `email=foo@example.com` returns accounts with `email` set to `foo@example.com` * Query `display_name=\\"Test String\\"` returns accounts with display names that include both "Test" and "String" |

### ListAccountsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| accounts | [ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) | repeated | The field name should match the noun "accounts" in the method name. There will be a maximum number of items returned based on the page_size field in the request |
| next_page_token | [string](#string) |  | Token to retrieve the next page of results, or empty if there are no more results in the list |

### ListGroupsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_size | [int32](#int32) |  | Optional. The maximum number of groups to return in the response |
| page_token | [string](#string) |  | Optional. A pagination token returned from a previous call to `Get` that indicates from where search should continue |
| field_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | Optional. Used to specify a subset of fields that should be returned by a get operation or modified by an update operation. |
| query | [string](#string) |  | Optional. Search criteria used to select the groups to return. If no search criteria is specified then all groups will be returned

TODO update query language Query expressions can be used to restrict results based upon the account properties where the operators `=`, `NOT`, `AND` and `OR` can be used along with the suffix wildcard symbol `*`.

The string properties in a query expression should use escaped quotes for values that include whitespace to prevent unexpected behavior.

Some example queries are:

* Query `display_name=Th*` returns accounts whose display_name starts with "Th" * Query `display_name=\\"Test String\\"` returns groups with display names that include both "Test" and "String" |

### ListGroupsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| groups | [ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) | repeated | The field name should match the noun "group" in the method name. There will be a maximum number of items returned based on the page_size field in the request |
| next_page_token | [string](#string) |  | Token to retrieve the next page of results, or empty if there are no more results in the list |

### ListMembersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_size | [int32](#int32) |  |  |
| page_token | [string](#string) |  | Optional. A pagination token returned from a previous call to `Get` that indicates from where search should continue |
| field_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | Optional. Used to specify a subset of fields that should be returned by a get operation or modified by an update operation. |
| query | [string](#string) |  | Optional. Search criteria used to select the groups to return. If no search criteria is specified then all groups will be returned

TODO update query language Query expressions can be used to restrict results based upon the account properties where the operators `=`, `NOT`, `AND` and `OR` can be used along with the suffix wildcard symbol `*`.

The string properties in a query expression should use escaped quotes for values that include whitespace to prevent unexpected behavior.

Some example queries are:

* Query `display_name=Th*` returns accounts whose display_name starts with "Th" * Query `display_name=\\"Test String\\"` returns groups with display names that include both "Test" and "String" |
| id | [string](#string) |  | The id of the group to list members from |

### ListMembersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| members | [ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) | repeated | The field name should match the noun "members" in the method name. There will be a maximum number of items returned based on the page_size field in the request |
| next_page_token | [string](#string) |  | Token to retrieve the next page of results, or empty if there are no more results in the list |

### RebuildIndexRequest




### RebuildIndexResponse




### RemoveMemberRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group_id | [string](#string) |  | The id of the group to remove a member from |
| account_id | [string](#string) |  | The account id to remove |

### UpdateAccountRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) |  | The account resource which replaces the resource on the server |
| update_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | The update mask applies to the resource. For the `FieldMask` definition, see https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#fieldmask |

### UpdateGroupRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| group | [ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) |  | The group resource which replaces the resource on the server |
| update_mask | [google.protobuf.FieldMask](#googleprotobuffieldmask) |  | The update mask applies to the resource. For the `FieldMask` definition, see https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#fieldmask |


### AccountsService

Follow recommended Methods for rpc APIs https://cloud.google.com/apis/design/resources?hl=de#methods
https://cloud.google.com/apis/design/standard_methods?hl=de#list
https://cloud.google.com/apis/design/naming_convention?hl=de

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListAccounts | [ListAccountsRequest](#listaccountsrequest) | [ListAccountsResponse](#listaccountsresponse) | Lists accounts |
| GetAccount | [GetAccountRequest](#getaccountrequest) | [.ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) | Gets an account rpc GetAccount(GetAccountRequest) returns (Account) { |
| CreateAccount | [CreateAccountRequest](#createaccountrequest) | [.ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) | Creates an account rpc CreateAccount(CreateAccountRequest) returns (Account) { |
| UpdateAccount | [UpdateAccountRequest](#updateaccountrequest) | [.ocis.messages.accounts.v1.Account](../../../messages/accounts/v1/grpc.md#account) | Updates an account rpc UpdateAccount(UpdateAccountRequest) returns (Account) { |
| DeleteAccount | [DeleteAccountRequest](#deleteaccountrequest) | [.google.protobuf.Empty](#googleprotobufempty) | Deletes an account |

### GroupsService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListGroups | [ListGroupsRequest](#listgroupsrequest) | [ListGroupsResponse](#listgroupsresponse) | Lists groups |
| GetGroup | [GetGroupRequest](#getgrouprequest) | [.ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) | Gets an groups |
| CreateGroup | [CreateGroupRequest](#creategrouprequest) | [.ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) | Creates a group |
| UpdateGroup | [UpdateGroupRequest](#updategrouprequest) | [.ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) | Updates a group |
| DeleteGroup | [DeleteGroupRequest](#deletegrouprequest) | [.google.protobuf.Empty](#googleprotobufempty) | Deletes a group |
| AddMember | [AddMemberRequest](#addmemberrequest) | [.ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) | group:addmember https://docs.microsoft.com/en-us/graph/api/group-post-members?view=graph-rest-1.0&tabs=http |
| RemoveMember | [RemoveMemberRequest](#removememberrequest) | [.ocis.messages.accounts.v1.Group](../../../messages/accounts/v1/grpc.md#group) | group:removemember https://docs.microsoft.com/en-us/graph/api/group-delete-members?view=graph-rest-1.0 |
| ListMembers | [ListMembersRequest](#listmembersrequest) | [ListMembersResponse](#listmembersresponse) | group:listmembers https://docs.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-1.0 |

### IndexService



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RebuildIndex | [RebuildIndexRequest](#rebuildindexrequest) | [RebuildIndexResponse](#rebuildindexresponse) |  |

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

