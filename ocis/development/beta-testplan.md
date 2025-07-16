---
title: "Beta testplan"
date: 2022-03-24T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: beta-testplan.md
---

# Beta Testing

This document is supposed to give you some ideas how and what to test on ocis. It's not meant to be an extensive list of all tests to be done, rather it should help you, as beta-tester, to get started and enable you to get creative and create your own test-cases. [Derive from these examples, be creative, do unusual and unconventional things, to try to break things](https://twitter.com/sempf/status/514473420277694465).

One option to create new test-cases and to stress the system is to examine what the [API acceptance-tests](https://owncloud.dev/ocis/development/testing/#testing-with-test-suite-natively-installed) or the [web-UI](#web) does, [examine the requests](#decode-https-traffic-with-wireshark) and do something a bit different with curl. This is also a good way to find out how APIs work that are not already fully documented.

Some cases have suggested setup steps, but feel free to use other setups. This can include:
- different deployment methods (e.g. running single binary, docker-container, docker-compose setup, [individual services in own docker containers](https://owncloud.dev/ocis/deployment/ocis_individual_services/))
- different identity managers (e.g. [different external LDAP](https://owncloud.dev/ocis/deployment/ocis_ldap/), internal IDM)
- different reverse proxies (e.g. [traefik](https://owncloud.dev/ocis/deployment/ocis_traefik/))
- different OpenIDConnect IDPs (e.g builtin IDP, [keycloak](https://owncloud.dev/ocis/deployment/ocis_keycloak/), AzureAD)
  - for some functionalities you will need an [LDAP server](https://owncloud.dev/ocis/deployment/ocis_ldap/) where the IDP and oCIS both get the users from
  - [keycloak example](https://owncloud.dev/ocis/deployment/ocis_keycloak/)
  - [service configuration](https://owncloud.dev/extensions/idp/configuration/)
- different storage systems (decomposedFS on local POSIX (default), [decomposedFS on NFS](https://owncloud.dev/ocis/storage-backends/dcfsnfs/), [S3](https://owncloud.dev/ocis/deployment/ocis_s3/) )

It's a good idea to test ocis in the same environment where you are planning to use it later (with the LDAP server, storage system, etc. of your organisation).

# run oCIS
For a quick start, please have a look into the [getting started documentation of oCIS]({{< ref "../getting-started/#run-ocis" >}})
If you would like
to access oCIS remotely please refer to the [Basic Remote Setup]({{< ref "../deployment/basic-remote-setup" >}}) section.

## additional tips
- to allow basic auth (e.g. to easily access oCIS using curl commands) you have to set `PROXY_ENABLE_BASIC_AUTH=true` environment variable
- if you cannot use real SSL Certificates set `OCIS_INSECURE=true` as environment variable


# Testplan

## user / groups from LDAP

Prerequisite:
- connect ocis to your preferred LDAP server
- create users and groups in LDAP
- start ocis with basic auth

documentation resources:
  - [configure ocis with LDAP](https://owncloud.dev/ocis/deployment/ocis_ldap/)
  - [sharing API is compatible to ownCloud 10](https://doc.owncloud.com/server/10.9/developer_manual/core/apis/ocs-share-api.html)
  - [webDav operations](#webdav)

| Test Case                                                                                             | Expected Result                                                                       | Comment |
|-------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|---------|
| share file / folder to a group                                                                        | member of the group can access shared item                                            |         |
| share file / folder to a group, remove member from group in LDAP                                      | removed member should not have access to the shared item                              |         |
| share file / folder to a group with different permissions, as receiver try to violate the permissions | receiver should not be able to violate the permissions                                |         |
| try to login with wrong credentials                                                                   | login should not be possible                                                          |         |
| set a quota in LDAP, upload files till the quota is exceeded                                          | upload should work till quota is full, uploads should not work when quota is full     |         |
| try to access files / folders of other users                                                          | access should not be possible                                                         |         |
| try to share with non-existing users and groups                                                       | sharing should not be possible                                                        |         |
| try to share with user/groups-names that contain special characters                                   | sharing should be possible, access shares with that user does not create any problems |         |

## other sharing

should be tried in various ways and in different environments

documentation resources:
- [sharing API is compatible to ownCloud 10](https://doc.owncloud.com/server/10.9/developer_manual/core/apis/ocs-share-api.html)
- [webDav operations](#webdav)

| Test Case                                                                             | Expected Result                                                       | Comment                                                   |
|---------------------------------------------------------------------------------------|-----------------------------------------------------------------------|-----------------------------------------------------------|
| share a file/folder with the same name from different users                           | receiver can accept and access both file/folders and distinguish them |  |
| share a file/folder with the same name but different permissions from different users | receiver can access both file/folders according to the permissions    |  |
| share a file/folder with the same name but different locations from one user          | receiver can accept and access both file/folders and distinguish them |  |
| share a file/folder back to the sharer                                                | sharing back should not be possible                                   |                                                           |
| re-share a file/folder with different permissions                                     | sharing with lower permissions is possible, but not with higher       |                                                           |
| decline received share                                                                | shared resource should not be shown to the receiver                   |                                                           |


## parallel deployment

- [configure ocis with LDAP](https://owncloud.dev/ocis/deployment/ocis_ldap/)
- [setup oC10 and ocis in parallel](https://owncloud.dev/ocis/deployment/oc10_ocis_parallel/)
- create users and groups in LDAP

documentation resources:
- [sharing API is compatible to ownCloud 10](https://doc.owncloud.com/server/10.9/developer_manual/core/apis/ocs-share-api.html)

| Test Case                                                                                                                                                        | Expected Result                                          | Comment |
|------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------|---------|
| share file / folder to a group in one implementation (use different permissions), access the items with the other implementation, try to violate the permissions | receiver should not be able to violate the permissions   |         |
| share file / folder to a group, remove member from group in LDAP, try to access items with the removed member from both implementations                          | removed member should not have access to the shared item |         |

## Spaces

Prerequisite:
- start ocis with basic auth
- variable declaration, for curl examples:
```shell
SERVER_URI=https://localhost:9200
GRAPH_API_PATH=graph/v1.0
SHARE_API_PATH=ocs/v2.php/apps/files_sharing/api/v1/shares
USER=admin
PASSWORD=admin
```
- create a new user `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$GRAPH_API_PATH/users" -X POST -d'{"displayName":"Example User","mail":"example@example.org","onPremisesSamAccountName":"example","passwordProfile":{"password":"ThePassword"}}'`
- give the user the "Space Admin" role
  1. get the id of the user: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/users/<user-name> | jq '.id'`
  2. assign role user to role: `curl -k -u $USER:$PASSWORD  $SERVER_URI/api/v0/settings/assignments-add -d'{"account_uuid":"<user-id>","role_id":"2aadd357-682c-406b-8874-293091995fdd"}`
- create a space: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$GRAPH_API_PATH/drives" -X POST -d'{"Name":"Space for finance","driveType":"project", "description":"we need some space"}'`
- list spaces: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/me/drives`
- disable a space: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X DELETE`
- delete a space: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X DELETE -H'Purge: T'`
- restore a space: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X PATCH -d"{}" -H"Restore: true"`
- rename the space: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X PATCH -d'{"name":"नेपालि नाम"}'`
- change description of the space: `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X PATCH -d'{"description":"this contains important data"}'`
- change quota of the space `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X PATCH -d'{"quota":{"total":<bytes>}}'`
- change image of space:
  1. upload an image file to the `.spaces` folder: `curl -k -u $USER:$PASSWORD https://localhost:9200/dav/spaces/<space-id>/.space/<file-name> -T <path-of-local-image> -v`
  2. note the id provided in the `Oc-Fileid` header
  3. set the image as "special-folder": `curl -k -u $USER:$PASSWORD  $SERVER_URI/$GRAPH_API_PATH/drives/<space-id> -X PATCH -d'{"special":[{"specialFolder":{"name":"image"},"id":"<oc-fileid>"}]}'`
- share a space: `curl -k -u $USER:$PASSWORD $SERVER_URI/$SHARE_API_PATH  -d'space_ref=<space-id>&shareType=7&shareWith=<receiver-username>&role=<role>'`
  - Roles:
    - viewer
    - editor
    - manager
- share a resource within a space `curl -k -u $USER:$PASSWORD $SERVER_URI/$SHARE_API_PATH -d'shareType=0&shareWith=<receiver-username>&space_ref=<space-id>%2F<resource>&permissions=<permissions>`
  - possible permissions
    - 1 = read
    - 2 = update
    - 4 = create
    - 8 = delete
    - 15 = read/write
    - 16 = share
    - 31 = All permissions
- WebDAV
  - root: `https://<server-uri>/dav/spaces/<space-id>`
  - [example commands for operations](#webdav)

| Test Case                                                                                                         | Expected Result                                                                   | Comment |
|-------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------|---------|
| create a space                                                                                                    | space should exist                                                                |         |
| create a space with special characters as a name & description                                                    | space should exist                                                                |         |
| create a space, delete the space                                                                                  | space should not exist                                                            |         |
| create a space, share the space with a user                                                                       | space should be accessible                                                        |         |
| create a space, share the space with a group                                                                      | space should be accessible, space content is shared among all users               |         |
| create a space, share the space with a group, disable the space                                                   | space should not be accessible                                                    |         |
| create a space, share the space with a user, disable the space, restore the space                                 | space should be accessible                                                        |         |
| create a space, disable the space, delete the space, restore the space                                            | it should not be possible to restore the space                                    |         |
| create a space, disable the space, try to share the space                                                         | sharing the space should not be possible                                          |         |
| create a space, try to delete the space                                                                           | it should not be possible to delete an enabled space                              |         |
| create & share a space with a group with viewer role, do CRUD file/folder operations with WebDAV                  | space content is readable but neither space not content should not be writable    |         |
| create & share a space with a group with editor role, do CRUD file/folder operations with WebDAV                  | space and content should be writable                                              |         |
| create a space, try CRUD file/folder operations with WebDAV on the space with a user that its not shared with     | space and content should not be accessible                                        |         |
| create a space with a quota, share the space, upload files till the quota is exceeded                             | upload should work till quota is full, uploads should not work when quota is full |         |
| share file/folders from inside a space (see other sharing section)                                                | sharing works and obeys the permissions                                           |         |
| create a space, rename the space                                                                                  | new name should be displayed in API calls and web                                 |         |
| create a space, change description of the space                                                                   | new description should be displayed in API calls and web                          |         |
| create a space, set quota, change quota of the space                                                              | new quota is obeyed                                                               |         |
| create a space, set quota, change quota of the space to a value that is lower than the sum of data already stored | new quota is obeyed, new files cannot be uploaded                                 |         |
| try the various space operations with invalid data                                                                | good error output, server does not crash                                          |         |
| try the various space operations without the correct permissions                                                  | operations are not possible without sufficient permissions                        |         |
| try the various space operations on personal and virtual spaces                                                   | server should not crash, good error responses                                     |         |
| try the various space operations sending invalid data (invalid numbers, wrong types, invalid JSON, etc.)          | server should not crash, good error responses                                     |         |


## Web

Prerequisite:
- [connect ocis to your preferred LDAP server](https://owncloud.dev/ocis/deployment/ocis_ldap/)
- create users and groups in LDAP
- Use your preferred browser (Firefox 94-96, Chrome 92-97, Opera 81-82, Edge 96-97, Safari 14-15) to access the built-in webUI (by default: https://localhost:9200)

| Test Case                                                                                                          | Expected Result                                                                       | Comment |
|--------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|---------|
| Login with the created user                                                                                        | User logs in.                                                                        |         |
| Create a text file.                                                                                                | Text editor can open, file is saved.                                                  |         |
| Create a text file with special characters as name                                                                 | file is created if the name is legal otherwise an error is displayed                  |         |
| Modify a text file.                                                                                                | File can be modified, no problems found.                                              |         |
| Rename a file.                                                                                                     | File is renamed.                                                                      |         |
| Upload a file.                                                                                                     | File is uploaded, no problems found.                                                  |         |
| Upload multiple files at once.                                                                                     | Files are uploaded, no problems found.                                                |         |
| delete all content of a folder at once.                                                                            | Folder is cleaned, items are visible in the trashbin                                  |         |
| Overwrite a file by uploading a new version.                                                                       | File is uploaded and overwritten, file versions are displayed                         |         |
| Overwrite a file by uploading a new version, restore the original version.                                         | File is restored correctly                                                            |         |
| upload a huge file                                                                                                 | File is uploaded, no problems found.                                                  |         |
| upload a huge file, cancel the upload, restart the upload                                                          | Upload continues at the position where it was cancelled, file is uploaded completely. |         |
| Remove a file.                                                                                                     | File is removed correctly, it appears in the trashbin.                                |         |
| Restore the deleted file from trashbin                                                                             | File is restored correctly                                                            |         |
| Remove multiple files that have the same name but are located in different folders                                 | Files are removed correctly, they appear in the trashbin.                             |         |
| Restore some of the deleted files from trashbin                                                                    | Files are restored correctly in the correct folders.                                  |         |
| Restore some of the deleted files from trashbin, but delete the original containing folder before                  | Files are restored correctly                                                          |         |
| Clean files from the trashbin                                                                                      | files are permanently deleted                                                         |         |
| Create a lot of files, delete a lot of files, empty the trashbin                                                   | trashbin is cleaned                                                                   |         |
| Move a file inside a folder.                                                                                       | There are not problems on the process.                                                |         |
| Move a file inside a folder that already contains a file with the same name                                        | File is not moved, content in the destination is not overwritten                      |         |
| Create a folder.                                                                                                   | Folder is created, no MKCOL problems appear.                                          |         |
| Create a folder with special characters as name                                                                    | Folder is created if the name is legal otherwise an error is displayed                |         |
| Create a folder with a name of an already existing file/folder                                                     | Folder is not created, an error is displayed                                          |         |
| Create a folder with a lot of subfolders, use special characters in the name                                       | Folder is created, no MKCOL problems appear.                                          |         |
| Delete a folder.                                                                                                   | Folder is removed.                                                                    |         |
| Move a folder inside another.                                                                                      | No problems while moving the folder.                                                  |         |
| open images in mediaviewer                                                                                         | files are displayed correctly.                                                        |         |
| open videos in mediaviewer                                                                                         | files are displayed correctly.                                                        |         |
| switch through videos and images in mediaviewer                                                                    | files are displayed correctly.                                                        |         |
| Share a file by public link.                                                                                       | Link is created and can be accessed.                                                  |         |
| Share a folder by public link.                                                                                     | Link is created and can be accessed.                                                  |         |
| Share a file with another user.                                                                                    | It is shared correctly.                                                               |         |
| Share a folder with another user.                                                                                  | It is shared correctly.                                                               |         |
| Share a file with a group.                                                                                         | It is shared correctly.                                                               |         |
| Share a folder with a group.                                                                                       | It is shared correctly.                                                               |         |
| Share a folder with userB giving edit permissions. As userB do CRUD operations on items inside the received folder | userB doesn't find any problem while interacting with files.                          |         |
| Use your mobile device to access the UI                                                                            | All elements reachable                                                                |         |
| do tests mentioned in the [spaces](#spaces) section using the web-UI                                               |                                                                                       |         |

## Desktop Client

Prerequisite:
- [install the desktop client on your preferred OS](https://owncloud.com/desktop-app/)
- start ocis
- connect a new account in the desktop client to ocis e.g. `https://localhost:9200` if you are running the server and the client on the same machine
- accept the self-signed certificate
- you will be redirected to the browser, accept the certificate there also
- login as any user
- allow the "ownCloud desktop app" to access ocis
- select a folder to sync

| Test Case                                                                                               | Expected Result                                                              | Comment |
|---------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|---------|
| Set up two clients with the same user. Change files, add some, delete some, move some, create folders.  | Changes sync properly in both clients without errors.                        |         |
| Share a file using contextual menu with userB.                                                          | Option to share appears in the contextual menu and file is correctly shared. |         |


## Mobile Clients (iOS || Android)

Prerequisite:
- [install the ownCloud app on your mobile](https://owncloud.com/mobile-apps/)
- start oCIS as described in [Basic Remote Setup]({{< ref "../deployment/basic-remote-setup" >}})
- connect a new account in the mobile client to ocis
- accept the self-signed certificate
- you will be redirected to the browser, accept the certificate there also
- login as any user
- allow the mobile app to access ocis


| Test Case                                     | Expected Result                          | Comment |
|-----------------------------------------------|------------------------------------------|---------|
| Connect to server, see files, download one.   | No problems while downloading.           |         |
| Upload a file using mobile client.            | No problems while uploading.             |         |
| Share a file with userB using mobile client.  | File is correctly shared.                |         |

## other WebDAV clients

Prerequisite:
- start ocis with basic auth

| Test Case                                                     | Expected Result                                             | Comment                                                                                                      |
|---------------------------------------------------------------|-------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| use ocis as webDAV external storage in ownCloud 10            | resource access works                                       |                                                                                                              |
| access webDAV with your file-manager                          | that will not give you a good UX, but ocis should not crash | Urls: https://\<ocis-server\>/remote.php/webdav  & https://\<ocis-server\>/remote.php/dav/files/\<username\> |
| access webDAV with the "remote-files" function of LibreOffice | files are accessible and can be written back                |                                                                                                              |

# Tips for testing

## WebDav
WebDav is accessible under different path
- old: https://\<ocis-server\>/remote.php/webdav
- new: https://\<ocis-server\>/remote.php/dav/files/\<username\>
- spaces: https://\<ocis-server\>/dav/spaces/\<space-id\>

WebDav specifications can be found on http://webdav.org/

here some general WebDav request examples:

variable declaration:
```shell
SERVER_URI=https://localhost:9200
API_PATH=remote.php/webdav
USER=admin
PASSWORD=admin
```
- list content of root folder: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$API_PATH/" -X PROPFIND`
- list content of sub-folder: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$API_PATH/f1" -X PROPFIND`
- create a folder: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$API_PATH/folder" -X MKCOL`
- delete a resource: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$API_PATH/folder" -X DELETE`
- rename / move a resource: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$API_PATH/folder" -X MOVE -H "Destination: $SERVER_URI/$API_PATH/renamed"`
- copy a resource: `curl -k -u $USER:$PASSWORD  "$SERVER_URI/$API_PATH/folder" -X COPY -H "Destination: $SERVER_URI/$API_PATH/folder-copy"`

## decode HTTPS traffic with wireshark
To decode the HTTPS traffic we need the keys that were used to encrypt the traffic. Those keys are kept secret by the clients, but we can request the clients to save them in a specific file, so that wireshark can use them to decrypt the traffic again.

1. create key file: `touch /tmp/sslkey.log`
2. start wireshark
3. set log filename
    - navigate to Edit=>Preferences=>Protocols=>TLS
    - in the field `(Pre)-Master-Secret log filename` enter `/tmp/sslkey.log`
4. decode as HTTP
    - navigate to Analyze=>Decode As...
    - click the + button
    - set Field: `TLS Port; Value=9200; Type: Integer, base 10; Default (none); Current HTTP` (adjust the port if you are using another one than 9200)
5. start recording
    - use `port 9200` as capture filter to only record ocis packages
    - use `http` as display filter to see only decoded traffic
6. run test-software with `SSLKEYLOGFILE=/tmp/sslkey.log` as env. variable e.g.
   - curl: `SSLKEYLOGFILE=/tmp/sslkey.log curl -k -u admin:admin https://localhost:9200/ocs/v1.php/cloud/users`
   - Browser: `SSLKEYLOGFILE=/tmp/sslkey.log firefox`
   - LibreOffice: `SSLKEYLOGFILE=/tmp/sslkey.log libreoffice`
   - acceptance tests: `SSLKEYLOGFILE=/tmp/sslkey.log make test-acceptance-api ...`

## format output
- piping **xml** results to `xmllint` gives you nice formats. E.g. `curl -k --user marie:radioactivity "https://localhost:9200/ocs/v1.php/apps/files_sharing/api/v1/shares" | xmllint --format -`
- piping **json** results to `jq` gives you nice formats. E.g. `curl -k --user marie:radioactivity "https://localhost:9200/ocs/v1.php/apps/files_sharing/api/v1/shares?format=json" | jq`

## create edge cases
- [Big List of Naughty Strings](https://github.com/minimaxir/big-list-of-naughty-strings)
