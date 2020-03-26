---
title: "GRPC API"
date: 2018-05-02T00:00:00+00:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis-thumbnails
geekdocEditPath: edit/master/docs
geekdocFilePath: grpc.md
---

<a name="top"></a>

## Table of Contents

- [pkg/proto/v0/thumbnails.proto](#pkg/proto/v0/thumbnails.proto)
    - [GetRequest](#com.owncloud.ocis.thumbnails.v0.GetRequest)
    - [GetResponse](#com.owncloud.ocis.thumbnails.v0.GetResponse)
  
    - [GetRequest.FileType](#com.owncloud.ocis.thumbnails.v0.GetRequest.FileType)
  
    - [ThumbnailService](#com.owncloud.ocis.thumbnails.v0.ThumbnailService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="pkg/proto/v0/thumbnails.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pkg/proto/v0/thumbnails.proto



<a name="com.owncloud.ocis.thumbnails.v0.GetRequest"></a>

### GetRequest
A request to retrieve a thumbnail


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filepath | [string](#string) |  | The path to the source image |
| filetype | [GetRequest.FileType](#com.owncloud.ocis.thumbnails.v0.GetRequest.FileType) |  | The type to which the thumbnail should get encoded to. |
| etag | [string](#string) |  | The etag of the source image |
| width | [int32](#int32) |  | The width of the thumbnail |
| height | [int32](#int32) |  | The height of the thumbnail |
| authorization | [string](#string) |  | The authorization token |






<a name="com.owncloud.ocis.thumbnails.v0.GetResponse"></a>

### GetResponse
The service response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| thumbnail | [bytes](#bytes) |  | The thumbnail as a binary |
| mimetype | [string](#string) |  | The mimetype of the thumbnail |





 <!-- end messages -->


<a name="com.owncloud.ocis.thumbnails.v0.GetRequest.FileType"></a>

### GetRequest.FileType
The file types to which the thumbnail cna get encoded to.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PNG | 0 | Represents PNG type |
| JPG | 1 | Represents JPG type |


 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="com.owncloud.ocis.thumbnails.v0.ThumbnailService"></a>

### ThumbnailService
A Service for handling thumbnail generation

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetThumbnail | [GetRequest](#com.owncloud.ocis.thumbnails.v0.GetRequest) | [GetResponse](#com.owncloud.ocis.thumbnails.v0.GetResponse) | Generates the thumbnail and returns it. |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

