---
title: Spaces
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/graph
geekdocFilePath: spaces.md
---

{{< toc >}}

## Graph Service

The Graph service is a reference implementation of the MS Graph API. There are no libraries doing any work only a set of routes and handlers.

## Spaces API

The Spaces API makes use of the [MS Graph API Drive resource](https://docs.microsoft.com/en-us/graph/api/resources/drive?view=graph-rest-1.0) to represent the concept of a Storage Space. Natively the MS Graph Specification [does not provide a way for creating Drives](https://docs.microsoft.com/en-us/graph/api/resources/drive?view=graph-rest-1.0#methods), as a Drive is a read only resource.

We circumvented this limitation by adding a `POST /drive/{drive-name}` to the Graph router. A major drawback of this solution is that this endpoint does not have support from the official MS Graph SDK, however it is reachable by any HTTP clients.

### Methods

```
POST /drive/{drive-name}
```

Calls to the following endpoint will create a Space with all the default parameters since we do not parse the request body just yet.

## Examples

We can now create a `Marketing` space and retrieve its WebDAV endpoint. Let's see how to do this.

### Starting conditions

This is the status of a DecomposedFS `users` tree. As we can see it is empty because we have not yet logged in with any users. It is a fresh new installation.

```
❯ tree -a /var/tmp/ocis/storage/users
/var/tmp/ocis/storage/users
├── blobs
├── nodes
│   └── root
├── spaces
│   ├── personal
│   └── share
├── trash
└── uploads
```

Let's start with creating a space:

`curl -k -X POST 'https://localhost:9200/graph/v1.0/drive/marketing' -u einstein:relativity -v`

```
❯ tree -a /var/tmp/ocis/storage/users
/var/tmp/ocis/storage/users
├── blobs
├── nodes
│   ├── 02dc1ec5-28b5-41c5-a48a-fabd4fa0562e
│   │   └── .space -> ../e85d185f-cdaa-4618-a312-e33ea435acfe
│   ├── 52efe3c2-c95a-47a1-8f3d-924aa473c711
│   ├── e85d185f-cdaa-4618-a312-e33ea435acfe
│   └── root
│       ├── 4c510ada-c86b-4815-8820-42cdf82c3d51 -> ../52efe3c2-c95a-47a1-8f3d-924aa473c711
│       └── c42debb8-926e-4a46-83b0-39dba56e59a4 -> ../02dc1ec5-28b5-41c5-a48a-fabd4fa0562e
├── spaces
│   ├── personal
│   │   └── 52efe3c2-c95a-47a1-8f3d-924aa473c711 -> ../../nodes/52efe3c2-c95a-47a1-8f3d-924aa473c711
│   ├── project
│   │   └── 02dc1ec5-28b5-41c5-a48a-fabd4fa0562e -> ../../nodes/02dc1ec5-28b5-41c5-a48a-fabd4fa0562e
│   └── share
├── trash
└── uploads
```

we can see that the `project` folder was added to the spaces as well as the `.space` folder to the space node `02dc1ec5-28b5-41c5-a48a-fabd4fa0562e`. For demonstration purposes, let's list the extended attributes of the new node:

```
xattr -l /var/tmp/ocis/storage/users/nodes/root/c42debb8-926e-4a46-83b0-39dba56e59a4
user.ocis.blobid:
user.ocis.blobsize: 0
user.ocis.name: c42debb8-926e-4a46-83b0-39dba56e59a4
user.ocis.owner.id: 4c510ada-c86b-4815-8820-42cdf82c3d51
user.ocis.owner.idp: https://localhost:9200
user.ocis.owner.type: primary
user.ocis.parentid: root
user.ocis.quota: 65536
user.ocis.space.name: marketing
```

As seen here it contains the metadata from the default list of requirements for this ticket.

Let's list the drive we just created using the graph API:

```
curl -k 'https://localhost:9200/graph/v1.0/me/drives' -u einstein:relativity -v | jq .value

[
  {
    "driveType": "personal",
    "id": "1284d238-aa92-42ce-bdc4-0b0000009157!52efe3c2-c95a-47a1-8f3d-924aa473c711",
    "lastModifiedDateTime": "2021-09-07T14:42:39.025050471+02:00",
    "name": "root",
    "owner": {
      "user": {
        "id": "4c510ada-c86b-4815-8820-42cdf82c3d51"
      }
    },
    "root": {
      "id": "1284d238-aa92-42ce-bdc4-0b0000009157!52efe3c2-c95a-47a1-8f3d-924aa473c711",
      "webDavUrl": "https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!52efe3c2-c95a-47a1-8f3d-924aa473c711"
    }
  },
  {
    "driveType": "project",
    "id": "1284d238-aa92-42ce-bdc4-0b0000009157!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e",
    "lastModifiedDateTime": "2021-09-07T14:42:39.030705579+02:00",
    "name": "root",
    "owner": {
      "user": {
        "id": "4c510ada-c86b-4815-8820-42cdf82c3d51"
      }
    },
    "quota": {
      "total": 65536
    },
    "root": {
      "id": "1284d238-aa92-42ce-bdc4-0b0000009157!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e",
      "webDavUrl": "https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e"
    }
  }
]

As we can see the response already contains a space-aware dav endpoint, which we can use to upload files to the space:

```
curl -k https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157\!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e/test.txt -X PUT -d "beep-sboop" -v -u einstein:relativity
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 9200 (#0)
* upload completely sent off: 10 out of 10 bytes
< HTTP/1.1 201 Created
< Access-Control-Allow-Origin: *
< Content-Length: 0
< Content-Security-Policy: default-src 'none';
< Content-Type: text/plain
< Date: Tue, 07 Sep 2021 12:45:54 GMT
< Etag: "e2942565a4eb52e8754c2806f215fe93"
< Last-Modified: Tue, 07 Sep 2021 12:45:54 +0000
< Oc-Etag: "e2942565a4eb52e8754c2806f215fe93"
< Oc-Fileid: MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OmU4ZDIwNDA5LTE1OWItNGI2Ny1iODZkLTlkM2U3ZjYyYmM0ZQ==
< Vary: Origin
< X-Content-Type-Options: nosniff
< X-Download-Options: noopen
< X-Frame-Options: SAMEORIGIN
< X-Permitted-Cross-Domain-Policies: none
< X-Robots-Tag: none
< X-Xss-Protection: 1; mode=block
<
* Connection #0 to host localhost left intact
* Closing connection 0
```

This is the state after every transformation:

```
tree -a /var/tmp/ocis/storage/users
/var/tmp/ocis/storage/users
├── blobs
│   └── 83842d56-91de-41d5-8800-b2fb7b2d31cf
├── nodes
│   ├── 02dc1ec5-28b5-41c5-a48a-fabd4fa0562e
│   │   ├── .space -> ../e85d185f-cdaa-4618-a312-e33ea435acfe
│   │   └── test.txt -> ../e8d20409-159b-4b67-b86d-9d3e7f62bc4e
│   ├── 52efe3c2-c95a-47a1-8f3d-924aa473c711
│   ├── e85d185f-cdaa-4618-a312-e33ea435acfe
│   ├── e8d20409-159b-4b67-b86d-9d3e7f62bc4e
│   └── root
│       ├── 4c510ada-c86b-4815-8820-42cdf82c3d51 -> ../52efe3c2-c95a-47a1-8f3d-924aa473c711
│       └── c42debb8-926e-4a46-83b0-39dba56e59a4 -> ../02dc1ec5-28b5-41c5-a48a-fabd4fa0562e
├── spaces
│   ├── personal
│   │   └── 52efe3c2-c95a-47a1-8f3d-924aa473c711 -> ../../nodes/52efe3c2-c95a-47a1-8f3d-924aa473c711
│   ├── project
│   │   └── 02dc1ec5-28b5-41c5-a48a-fabd4fa0562e -> ../../nodes/02dc1ec5-28b5-41c5-a48a-fabd4fa0562e
│   └── share
├── trash
└── uploads
```

Observe the `test.txt` in the `02dc1ec5-28b5-41c5-a48a-fabd4fa0562e` node.

To finalize, verify the new created file is webdav-listable:

```xml
curl -k https://localhost:9200/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157\!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e -X PROPFIND -v -u einstein:relativity | xmllint --format -

<?xml version="1.0" encoding="utf-8"?>
<d:multistatus xmlns:d="DAV:" xmlns:s="http://sabredav.org/ns" xmlns:oc="http://owncloud.org/ns">
  <d:response>
    <d:href>/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e/</d:href>
    <d:propstat>
      <d:prop>
        <oc:id>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OjAyZGMxZWM1LTI4YjUtNDFjNS1hNDhhLWZhYmQ0ZmEwNTYyZQ==</oc:id>
        <oc:fileid>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OjAyZGMxZWM1LTI4YjUtNDFjNS1hNDhhLWZhYmQ0ZmEwNTYyZQ==</oc:fileid>
        <d:getetag>"35a2ce5f56592d79d1b7233eff033347"</d:getetag>
        <oc:permissions>RDNVCK</oc:permissions>
        <d:resourcetype>
          <d:collection/>
        </d:resourcetype>
        <oc:size>0</oc:size>
        <d:getlastmodified>Tue, 07 Sep 2021 12:45:54 GMT</d:getlastmodified>
        <oc:favorite>0</oc:favorite>
      </d:prop>
      <d:status>HTTP/1.1 200 OK</d:status>
    </d:propstat>
  </d:response>
  <d:response>
    <d:href>/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e/.space/</d:href>
    <d:propstat>
      <d:prop>
        <oc:id>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OmU4NWQxODVmLWNkYWEtNDYxOC1hMzEyLWUzM2VhNDM1YWNmZQ==</oc:id>
        <oc:fileid>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OmU4NWQxODVmLWNkYWEtNDYxOC1hMzEyLWUzM2VhNDM1YWNmZQ==</oc:fileid>
        <d:getetag>"2e9a84bffce8b648ba626185800ee8fa"</d:getetag>
        <oc:permissions>SRDNVCK</oc:permissions>
        <d:resourcetype>
          <d:collection/>
        </d:resourcetype>
        <oc:size>0</oc:size>
        <d:getlastmodified>Tue, 07 Sep 2021 12:42:39 GMT</d:getlastmodified>
        <oc:favorite>0</oc:favorite>
      </d:prop>
      <d:status>HTTP/1.1 200 OK</d:status>
    </d:propstat>
  </d:response>
  <d:response>
    <d:href>/dav/spaces/1284d238-aa92-42ce-bdc4-0b0000009157!02dc1ec5-28b5-41c5-a48a-fabd4fa0562e/test.txt</d:href>
    <d:propstat>
      <d:prop>
        <oc:id>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OmU4ZDIwNDA5LTE1OWItNGI2Ny1iODZkLTlkM2U3ZjYyYmM0ZQ==</oc:id>
        <oc:fileid>MTI4NGQyMzgtYWE5Mi00MmNlLWJkYzQtMGIwMDAwMDA5MTU3OmU4ZDIwNDA5LTE1OWItNGI2Ny1iODZkLTlkM2U3ZjYyYmM0ZQ==</oc:fileid>
        <d:getetag>"e2942565a4eb52e8754c2806f215fe93"</d:getetag>
        <oc:permissions>RDNVW</oc:permissions>
        <d:resourcetype/>
        <d:getcontentlength>10</d:getcontentlength>
        <d:getcontenttype>text/plain</d:getcontenttype>
        <d:getlastmodified>Tue, 07 Sep 2021 12:45:54 GMT</d:getlastmodified>
        <oc:checksums>
          <oc:checksum>SHA1:8f4b4c83c565fc5ec54b78c30c94a6b65e411de5 MD5:6a3a4eca9a6726eef8f7be5b03ea9011 ADLER32:151303ed</oc:checksum>
        </oc:checksums>
        <oc:favorite>0</oc:favorite>
      </d:prop>
      <d:status>HTTP/1.1 200 OK</d:status>
    </d:propstat>
  </d:response>
</d:multistatus>
```
