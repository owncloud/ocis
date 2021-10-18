---
title: "GRPC API"
date: 2018-05-02T00:00:00+00:00
weight: 50
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: grpc.md
---

{{< toc >}}

## proto/v0/thumbnails.proto

### CS3Source



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| authorization | [string](#string) |  |  |

### GetThumbnailRequest

A request to retrieve a thumbnail

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filepath | [string](#string) |  | The path to the source image |
| thumbnail_type | [GetThumbnailRequest.ThumbnailType](#getthumbnailrequestthumbnailtype) |  | The type to which the thumbnail should get encoded to. |
| width | [int32](#int32) |  | The width of the thumbnail |
| height | [int32](#int32) |  | The height of the thumbnail |
| webdav_source | [WebdavSource](#webdavsource) |  |  |
| cs3_source | [CS3Source](#cs3source) |  |  |

### GetThumbnailResponse

The service response

| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thumbnail | [bytes](#bytes) |  | The thumbnail as a binary |
| mimetype | [string](#string) |  | The mimetype of the thumbnail |

### WebdavSource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url | [string](#string) |  | REQUIRED. |
| is_public_link | [bool](#bool) |  | REQUIRED. |
| webdav_authorization | [string](#string) |  | OPTIONAL. |
| reva_authorization | [string](#string) |  | OPTIONAL. |
| public_link_token | [string](#string) |  | OPTIONAL. |

### GetThumbnailRequest.ThumbnailType

The file types to which the thumbnail can get encoded to.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PNG | 0 | Represents PNG type |
| JPG | 1 | Represents JPG type |

### ThumbnailService

A Service for handling thumbnail generation

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetThumbnail | [GetThumbnailRequest](#getthumbnailrequest) | [GetThumbnailResponse](#getthumbnailresponse) | Generates the thumbnail and returns it. |

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
