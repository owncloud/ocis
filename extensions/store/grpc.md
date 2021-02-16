---
title: "GRPC API"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: grpc.md
---

{{< toc >}}

## store.proto

### DatabasesRequest




### DatabasesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| databases | [string](#string) | repeated |  |

### DeleteOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| database | [string](#string) |  |  |
| table | [string](#string) |  |  |

### DeleteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| options | [DeleteOptions](#deleteoptions) |  |  |

### DeleteResponse




### Field



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | type of value e.g string, int, int64, bool, float64 |
| value | [string](#string) |  | the actual value |

### ListOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| database | [string](#string) |  |  |
| table | [string](#string) |  |  |
| prefix | [string](#string) |  |  |
| suffix | [string](#string) |  |  |
| limit | [uint64](#uint64) |  |  |
| offset | [uint64](#uint64) |  |  |

### ListRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| options | [ListOptions](#listoptions) |  |  |

### ListResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| keys | [string](#string) | repeated |  |

### ReadOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| database | [string](#string) |  |  |
| table | [string](#string) |  |  |
| prefix | [bool](#bool) |  |  |
| suffix | [bool](#bool) |  |  |
| limit | [uint64](#uint64) |  |  |
| offset | [uint64](#uint64) |  |  |
| where | [ReadOptions.WhereEntry](#readoptionswhereentry) | repeated |  |

### ReadOptions.WhereEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Field](#field) |  |  |

### ReadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| options | [ReadOptions](#readoptions) |  |  |

### ReadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| records | [Record](#record) | repeated |  |

### Record



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | key of the recorda |
| value | [bytes](#bytes) |  | value in the record |
| expiry | [int64](#int64) |  | time.Duration (signed int64 nanoseconds) |
| metadata | [Record.MetadataEntry](#recordmetadataentry) | repeated | the associated metadata |

### Record.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Field](#field) |  |  |

### TablesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| database | [string](#string) |  |  |

### TablesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tables | [string](#string) | repeated |  |

### WriteOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| database | [string](#string) |  |  |
| table | [string](#string) |  |  |
| expiry | [int64](#int64) |  | time.Time |
| ttl | [int64](#int64) |  | time.Duration |

### WriteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| record | [Record](#record) |  |  |
| options | [WriteOptions](#writeoptions) |  |  |

### WriteResponse





### Store



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Read | [ReadRequest](#readrequest) | [ReadResponse](#readresponse) |  |
| Write | [WriteRequest](#writerequest) | [WriteResponse](#writeresponse) |  |
| Delete | [DeleteRequest](#deleterequest) | [DeleteResponse](#deleteresponse) |  |
| List | [ListRequest](#listrequest) | [ListResponse](#listresponse) stream |  |
| Databases | [DatabasesRequest](#databasesrequest) | [DatabasesResponse](#databasesresponse) |  |
| Tables | [TablesRequest](#tablesrequest) | [TablesResponse](#tablesresponse) |  |

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
