# web-client

The `web-client` is a standalone package that allows you to interact with the [ownCloud Infinite Scale (oCIS)](https://github.com/owncloud/ocis/) APIs via TypeScript. It provides an abstraction layer between the server and a (web-) application that converts API data into objects with helpful types and utilities. This abstraction ensures that users of the APIs don't need in-depth knowledge about them, such as required methods or returned status codes.

The supported APIs are:

- Graph (drives, sharing, user & group management)
- OCS (capabilities & url signing)
- WebDAV (file operations)

## Installation

Depending on your package manager, run one of the following commands:

```
$ npm install @ownclouders/web-client

$ pnpm add @ownclouders/web-client

$ yarn add @ownclouders/web-client
```

## Usage

### Graph

The graph client needs to be instantiated with a base URI corresponding to your oCIS deployment and an axios instance. The axios instance is being used for all requests, which means it needs to include all relevant headers either statically or via interceptor.

```
import axios from axios
import { graph } from '@ownclouders/web-client'

const accessToken = 'some_access_token'
const baseURI = 'some_base_uri'

const axiosClient = axios.create({
	headers: { Authorization: accessToken }
})

const graphClient = graph(baseURI, axiosClient)
```

The following example demonstrates how to retrieve all spaces accessible to the user. A `SpaceResource` can then be used to e.g. fetch files and folders (see webdav example down below).

```
const mySpaces = await graphClient.drives.listMyDrives()
```

### OCS

The ocs client needs to be instantiated with a base URI corresponding to your oCIS deployment and an axios instance. The axios instance is being used for all requests, which means it needs to include all relevant headers either statically or via interceptor.

```
import axios from axios
import { ocs } from '@ownclouders/web-client'

const accessToken = 'some_access_token'
const baseURI = 'some_base_uri'

const axiosClient = axios.create({
	headers: { Authorization: accessToken }
})

const ocsClient = ocs(baseURI, axiosClient)
```

The following examples demonstrate how to fetch capabilities and sign URLs.

```
const capabilities = await ocsClient.getCapabilities()

const signedUrl = await ocsClient.signUrl('some_url_to_sign', 'your_username')
```

### WebDav

The webdav client needs to be instantiated with a base URI corresponding to your oCIS deployment. You can also pass a header callback which will be called with every dav request.

```
import { webdav } from '@ownclouders/web-client'

const accessToken = 'some_access_token'
const baseURI = 'some_base_uri'

const webDavClient = webdav(
	baseURI,
	() => ({ Authorization: accessToken })
)
```

The following example demonstrates how to list all resources of a given `SpaceResource` (see above how to fetch space resources).

```
const { resource, children } = await webDavClient.listFiles(spaceResource)
```
