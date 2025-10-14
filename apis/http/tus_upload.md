---
title: "Resumable Upload"
date: 2023-10-10T00:00:00+00:00
weight: 21
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/apis/http
geekdocFilePath: tus_upload.md
geekdocCollapseSection: true
---

Infinite Scale supports the tus resumable-upload protocol, which is a robust, modular, and open protocol designed to resume large file uploads reliably over HTTP.
In situations where file uploads might be interrupted due to network issues, browser crashes, or other unforeseen interruptions,
tus ensures that uploads can be resumed from the point of failure without losing data.
This documentation shows some basic examples, refer [tus official site](https://tus.io/protocols/resumable-upload) for more details.

## Supported tus Features

The backend announces certain tus features to clients. WebDAV responses come with tus HTTP headers for the offical tus features, and additional, ownCloud specific features are announced via the capabilities endpoint (e.g. `https://localhost:9200/ocs/v1.php/cloud/capabilities?format=json`).

The following snippet shows the relevant part of the server capabilities of Infinite Scale that concerns the tus upload:
```json
{
  "ocs": {
    "data": {
      "capabilities": {
        "files": {
          "tus_support": {
              "version": "1.0.0",
              "resumable": "1.0.0",
              "extension": "creation,creation-with-upload",
              "max_chunk_size": 10000000,
              "http_method_override": ""
            }
          }
        }
      }
    }
  }
}
```

| Parameter      | Environment Variable           | Default Value | Description                                                         |
| -------------- | ------------------------------ | ------------- | ------------------------------------------------------------------- |
| max_chunk_size | FRONTEND_UPLOAD_MAX_CHUNK_SIZE | 10000000      | Announces the max chunk sizes in bytes for uploads via the clients. |

## Upload in Chunks

### Create an Upload URL

The client must send a POST request against a known upload creation URL to request a new upload resource.
The filename has to be provided in base64-encoded format.

Example:
```shell
# base64 encoded filename 'tustest.txt' is 'dHVzdGVzdC50eHQ='
echo -n 'tustest.txt' | base64
```

{{< tabs "create-upload-url" >}}
{{< tab "Request" >}}
```shell
curl -ks -XPOST https://ocis.test/remote.php/dav/spaces/8d72036d-14a5-490f-889e-414064156402$196ac304-7b88-44ce-a4db-c4becef0d2e0 \
-H "Authorization: Bearer eyJhbGciOiJQUzI..."\
-H "Tus-Resumable: 1.0.0" \
-H "Upload-Length: 10" \
-H "Upload-Metadata: filename dHVzdGVzdC50eHQ="
```
{{< /tab >}}

{{< tab "Response - 201 Created" >}}
```
< HTTP/1.1 201 Created
< Access-Control-Allow-Headers: Tus-Resumable, Upload-Length, Upload-Metadata, If-Match
< Access-Control-Allow-Origin: *
< Access-Control-Expose-Headers: Tus-Resumable, Upload-Offset, Location
< Content-Length: 0
< Content-Security-Policy: default-src 'none';
< Date: Mon, 16 Oct 2023 08:49:39 GMT
< Location: https://ocis.test/data/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJyZXZhIiwiZXhwIjoxNjk3NTMyNTc5LCJpYXQiOjE2OTc0NDYxNzksInRhcmdldCI6Imh0dHA6Ly9sb2NhbGhvc3Q6OTE1OC9kYXRhL3R1cy8zYTU3ZWZlMS04MzE0LTQ4MGEtOWY5Ny04N2Q1YzBjYTJhMTgifQ.FbrlY7mdOfsbFgMrP8OtcHlCEq72a2ZVnPD2iBo9MfM
< Tus-Extension: creation,creation-with-upload,checksum,expiration
< Tus-Resumable: 1.0.0
< Vary: Origin
< X-Content-Type-Options: nosniff
< X-Download-Options: noopen
< X-Frame-Options: SAMEORIGIN
< X-Permitted-Cross-Domain-Policies: none
< X-Request-Id: xxxxxxxxxxxxxxxxxxxxxx
< X-Robots-Tag: none
<
* Connection #0 to host localhost left intact
```
{{< /tab >}}
{{< /tabs >}}

The server will return a temporary upload URL in the location header of the response:
```
< Location: https://ocis.test/data/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJyZXZhIiwiZXhwIjoxNjk3NTMyNTc5LCJpYXQiOjE2OTc0NDYxNzksInRhcmdldCI6Imh0dHA6Ly9sb2NhbGhvc3Q6OTE1OC9kYXRhL3R1cy8zYTU3ZWZlMS04MzE0LTQ4MGEtOWY5Ny04N2Q1YzBjYTJhMTgifQ.FbrlY7mdOfsbFgMrP8OtcHlCEq72a2ZVnPD2iBo9MfM
```

### Upload the First Chunk

Once a temporary upload URL has been created, a client can send a PATCH request to upload a file. The file content should be sent in the body of the request:
{{< tabs "upload-the-first-chunk" >}}
{{< tab "Request" >}}
```shell
curl -ks -XPATCH https://temporary-upload-url \
-H "Authorization: Bearer eyJhbGciOiJQUzI..." \
-H "Tus-Resumable: 1.0.0" \
-H "Upload-Offset: 0" \
-H "Content-Type: application/offset+octet-stream" -d "01234"
```
{{< /tab >}}

{{< tab "Response - 204 No Content" >}}
```
< HTTP/1.1 204 No Content
< Date: Tue, 17 Oct 2023 04:10:52 GMT
< Oc-Fileid: 8d72036d-14a5-490f-889e-414064156402$73bb5450-816b-4cae-90aa-1f96adc95bd4!84e319e4-de1d-4dd8-bbd0-e51d933cdbcd
< Tus-Resumable: 1.0.0
< Upload-Expires: 1697602157
< Upload-Offset: 5
< Vary: Origin
< X-Content-Type-Options: nosniff
< X-Request-Id: xxxxxxxxxxxxxxxxxxxxxx
<
* Connection #0 to host localhost left intact
```
{{< /tab >}}
{{< /tabs >}}

### Upload Further Chunks

After the first chunk is uploaded, the second chunk can be uploaded by pointing `Upload-Offset` to exact position that was returned in the first response.
Upload process will not be marked as complete until the total uploaded content size matches the `Upload-Length` specified during the creation of the temporary URL.

{{< tabs "upload-the-second-chunk" >}}
{{< tab "Request" >}}
```shell
curl -ks -XPATCH https://temporary-upload-url \
-H "Authorization: Bearer eyJhbGciOiJQUzI..." \
-H "Tus-Resumable: 1.0.0" \
-H "Upload-Offset: 5" \
-H "Content-Type: application/offset+octet-stream" -d "56789"
```
{{< /tab >}}

{{< tab "Response - 204 No Content" >}}
```
< HTTP/1.1 204 No Content
< Date: Tue, 17 Oct 2023 04:11:00 GMT
< Oc-Fileid: 8d72036d-14a5-490f-889e-414064156402$73bb5450-816b-4cae-90aa-1f96adc95bd4!84e319e4-de1d-4dd8-bbd0-e51d933cdbcd
< Tus-Resumable: 1.0.0
< Upload-Expires: 1697602157
< Upload-Offset: 10
< Vary: Origin
< X-Content-Type-Options: nosniff
< X-Request-Id: xxxxxxxxxxxxxxxxxxxxxx
<
* Connection #0 to host localhost left intact
```
{{< /tab >}}
{{< /tabs >}}
{{< hint type=warning title="Important Warning" >}}
`Upload-Offset` header indicates the byte position in the target file where the server should start writing the upload content.
It ensures data integrity and order during the upload process.
{{< /hint >}}

## Creation with Upload

{{< tabs "creation-with-upload" >}}
{{< tab "Request" >}}
```shell
curl -ks -XPOST https://ocis.test/remote.php/dav/spaces/{space-id} \
-H "Authorization: Bearer eyJhbGciOiJQUzI..." \
-H "Tus-Resumable: 1.0.0" \
-H "Upload-Length: 14" \
-H "Content-Type: application/offset+octet-stream" \
-H "Upload-Metadata: filename dGVzdC50eHQ=" \
-H "Tus-Extension: creation-with-upload" \
-d "upload content"
```
{{< /tab >}}

{{< tab "Response - 201 Created" >}}
```shell
< HTTP/1.1 201 Created
< Access-Control-Allow-Headers: Tus-Resumable, Upload-Length, Upload-Metadata, If-Match
< Access-Control-Allow-Origin: *
< Access-Control-Expose-Headers: Tus-Resumable, Upload-Offset, Location
< Content-Length: 0
< Content-Security-Policy: default-src 'none';
< Content-Type: text/plain
< Date: Mon, 16 Oct 2023 04:18:25 GMT
< Etag: "372c96743f68bc40e789124d30567371"
< Last-Modified: Mon, 16 Oct 2023 04:18:25 +0000
< Location: https://ocis.test/data/eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJyZXZhIiwiZXhwIjoxNjk3NTE2MzA1LCJpYXQiOjE2OTc0Mjk5MDUsInRhcmdldCI6Imh0dHA6Ly9sb2NhbGhvc3Q6OTE1OC9kYXRhL3R1cy82NjlhODBlZi1hN2VjLTQwYTAtOGNmOS05MTgwNTVhYzlkZjAifQ.yq-ofJYnJ9FLML7Z_jki1FJQ7Ulbt9O_cmLe6V411A4
< Oc-Etag: "372c96743f68bc40e789124d30567371"
< Oc-Fileid: 44d3e1e0-6c01-4b94-9145-9d0068239fcd$446bdad4-4b27-41f1-afce-0881f202a214!d7c292a6-c395-4e92-bf07-2c1663aec8dd
< Oc-Perm: RDNVWZP
< Tus-Extension: creation,creation-with-upload,checksum,expiration
< Tus-Resumable: 1.0.0
< Upload-Expires: 1697516305
< Upload-Offset: 14
< Vary: Origin
< X-Content-Type-Options: nosniff
< X-Download-Options: noopen
< X-Frame-Options: SAMEORIGIN
* TLSv1.2 (IN), TLS header, Supplemental data (23):
{ [5 bytes data]
< X-Permitted-Cross-Domain-Policies: none
< X-Request-Id: xxxxxxxxxxxxxxxxxxxxxx
< X-Robots-Tag: none
<
* Connection #0 to host localhost left intact
```
{{< /tab >}}
{{< /tabs >}}

{{< hint type=warning title="Important Warning" >}}
The `Upload-Length` header of the request has to contain the exact size of the upload content in byte.
{{< /hint >}}

## Supported Upload-Metadata

Upload-metadata key-value pairs aren't specified in the general tus docs. The following ones are supported in the ownCloud ecosystem:

| Parameter (key)                  | Example (value, MUST be Base64 encoded)                                                                                                         | Description                                                                                         |
| -------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| `name` OR `filename` (mandatory) | example.pdf                                                                                                                                     | Filename                                                                                            |
| `mtime` (recommended)            | 1701708712                                                                                                                                      | Modification time (Unix time format)                                                                |
| `checksum` (recommended)         | SHA1 a330de5886e5a92d78fb3f8d59fe469857759e72                                                                                                   | Checksum, computed from the client                                                                  |
| `type` OR `filetype`             | application/pdf                                                                                                                                 | MIME Type, sent by the web UI                                                                       |
| `relativePath`                   | undefined                                                                                                                                       | File path relative to the folder that is being uploaded, including the filename. Sent by the web UI |
| `spaceId`                        | 8748cddf-66b7-4b85-91a7-e6d08d8e1639$a9778d63-21e7-4d92-9b47-1b81144b9993                                                                       | Sent by the web UI                                                                                  |
| `spaceName`                      | Personal                                                                                                                                        | Sent by the web UI                                                                                  |
| `driveAlias`                     | personal/admin                                                                                                                                  | Sent by the web UI                                                                                  |
| `driveType`                      | personal                                                                                                                                        | Sent by the web UI                                                                                  |
| `currentFolder`                  | /                                                                                                                                               | Sent by the web UI                                                                                  |
| `currentFolderId`                | 8748cddf-66b7-4b85-91a7-e6d08d8e1639$a9778d63-21e7-4d92-9b47-1b81144b9993!a9778d63-21e7-4d92-9b47-1b81144b9993                                  | Sent by the web UI                                                                                  |
| `uppyId`                         | uppy-example/pdf-1e-application/pdf-238300                                                                                                      | Sent by the web UI                                                                                  |
| `relativeFolder`                 |                                                                                                                                                 | File path relative to the folder that is being uploaded, without filename. Sent by the web UI.      |
| `tusEndpoint`                    | https://ocis.ocis-traefik.latest.owncloud.works/remote.php/dav/spaces/8748cddf-66b7-4b85-91a7-e6d08d8e1639$a9778d63-21e7-4d92-9b47-1b81144b9993 | Sent by the web UI                                                                                  |
| `uploadId`                       | 71d5f878-a96c-4d7b-9627-658d782c93d7                                                                                                            | Sent by the web UI                                                                                  |
| `topLevelFolderId`               | undefined                                                                                                                                       | Sent by the web UI                                                                                  |
| `routeName`                      | files-spaces-generic                                                                                                                            | Sent by the web UI                                                                                  |
| `routeDriveAliasAndItem`         | cGVyc29uYWwvYWRtaW4=                                                                                                                            | Sent by the web UI                                                                                  |
| `routeShareId`                   |                                                                                                                                                 | Share ID when uploading into a received folder share. Sent by the web UI                            |
