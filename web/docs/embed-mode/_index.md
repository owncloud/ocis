---
title: 'Embed Mode'
date: 2023-10-23T00:00:00+00:00
weight: 60
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/embed-mode
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc >}}

The ownCloud Web can be consumed by another application in a stripped down version called "Embed mode". This mode is supposed to be used in the context of selecting or sharing resources. If you're looking for even more minimalistic approach, you can take a look at the [File picker](https://owncloud.dev/integration/file_picker/).

## Getting started

To integrate ownCloud Web into your application, add an iframe element pointing to your ownCloud Web deployed instance with additional query parameter `embed=true`.

```html
<iframe src="<web-url>?embed=true"></iframe>
```

## Communication

To establish seamless cross-origin communication between the embedded instance and the parent application, our approach involves emitting events using the `postMessage` method. These events can be conveniently captured by utilizing the standard `window.addEventListener('message', listener)` pattern.

### Target origin

By default, the `postMessage` method does not specify the `targetOrigin` parameter. However, it is recommended best practice to explicitly pass in the URI of the iframe origin (not the parent application). To enhance security, you can specify this value by modifying the config option `options.embed.messagesOrigin`.

### Events

To maintain uniformity and ease of handling, each event encapsulates the same structure within its payload: `{ name: string, data: any }`.

| Name | Data | Description |
| --- | --- | --- |
| **owncloud-embed:select** | Resource[] | Gets emitted when user selects resources or location via the select action |
| **owncloud-embed:share** | string[] | **DEPRECATED**: Gets emitted when user selects resources and shares them via the "Share links" action. Use `owncloud-embed:share-links` instead. |
| **owncloud-embed:share-links** | Array<{ url: string; password?: string }> | Gets emitted when user selects resources and shares them via the "Share link(s)" or "Share link(s) and password(s)" action. Each object contains the link URL and optionally the password (when shared with password). |
| **owncloud-embed:cancel** | null | Gets emitted when user attempts to close the embedded instance via "Cancel" action |

### Example

#### Selecting resources

```html
<iframe src="https://my-owncloud-web-instance?embed=true"></iframe>

<script>
  function selectEventHandler(event) {
    if (event.data?.name !== 'owncloud-embed:select') {
      return
    }

    const resources = event.data.data

    doSomethingWithSelectedResources(resources)
  }

  window.addEventListener('message', selectEventHandler)
</script>
```

#### Sharing links with password

```html
<iframe src="https://my-owncloud-web-instance?embed=true"></iframe>

<script>
  function shareLinksEventHandler(event) {
    if (event.data?.name !== 'owncloud-embed:share-links') {
      return
    }

    const links = event.data.data // Array<{ url: string; password?: string }>

    links.forEach(link => console.log("Link", link.url, "Password", link.password))

    doSomethingWithSharedLinks(links)
  }

  window.addEventListener('message', shareLinksEventHandler)
</script>
```

## Location picker

By default, the Embed mode allows users to select resources. In certain cases (e.g. uploading a file), this needs to be changed to allow selecting a location. This can be achieved by running the embed mode with additional parameter `embed-target=location`. With this parameter, resource selection is disabled and the selected resources array always includes the current folder as the only item.
In special scenarios you also want the user to set a file name, this can be achieved by adding the `embed-choose-file-name=true` parameter, or if you also want to set a default file name, you can use `embed-choose-file-name-suggestion=my file.text`.


### Example

```html
<iframe src="https://my-owncloud-web-instance?embed=true&embed-target=location"></iframe>

<script>
  function selectEventHandler(event) {
    if (event.data?.name !== 'owncloud-embed:select') {
      return
    }

    const resources = event.data.data[0]

    doSomethingWithSelectedResources(resources)
  }

  window.addEventListener('message', selectEventHandler)
</script>
```

## File picker

The File Picker mode in ownCloud Web is designed for embedding an interface that allows users to pick a single file.
This mode can be configured to restrict the file types that users can select. To enable the File Picker mode, you need
to include the embed-target=file query parameter in the iframe URL. Furthermore, you can specify allowed file types
using the embed-file-types parameter. The file types can be specified using file extensions, MIME types, or a
combination of both. If the embed-file-types parameter is not provided, all file types will be selectable by default.

### Example

```html

<iframe src="https://my-owncloud-web-instance?embed=true&embed-target=file&embed-file-types=txt,image/png"></iframe>

<script>
    function selectEventHandler(event) {
        if (event.data?.name !== 'owncloud-embed:file-pick') {
            return
        }

        const file = event.data.data

        doSomethingWithPickedFile(file)
    }

    window.addEventListener('message', selectEventHandler)
</script>
```

## Delegate authentication

If you already have a valid `access_token` that can be used to call the API from within the Embed mode and do not want to force the user to authenticate again, you can delegate the authentication. Delegating authentication will disable internal login form in ownCloud Web and will instead use events to obtain the token and update it.

### Configuration

To allow authentication delegation, you need to set the config option `options.embed.delegateAuthentication` to `true`. This can be achieved via query parameter `embed-delegate-authentication=true`. Because we are using the `postMessage` method to communicate across different origins, it is best practice to verify that the event originated from a known origin and not from some malicious site. We highly recommend to allow this check in production environments. You can enable it by setting the config option `options.embed.delegateAuthenticationOrigin` via query parameter `embed-delegate-authentication-origin=my-origin`. The value of this parameter will be compared against the `MessageEvent.origin` value and if they do not match, the token will be rejected.

### Events

#### Opening Embed mode

As already mentioned, we're using the `postMessage` method to allow communication between the Embed mode and the parent application. When the Embed mode is opened for the first time, the user gets redirected to the `/web-oidc-callback` page where a message with payload `{ name: 'owncloud-embed:request-token', data: undefined }` is sent to request the `access_token` from the parent application. The parent application should set an event listener before opening the Embed mode and once received, it should send a message with payload `{ name: 'owncloud-embed:update-token', data: { access_token: '<bearer-token>' } }`. Once the Embed mode receives this message, it will save the token in the application state and will automatically authenticate the user.

{{< hint info >}}
When passing the token in the message payload, use only the token itself without `Bearer ` string as that will be added automatically in the Embed mode.
{{< /hint >}}

{{< hint info >}}
To save unnecessary duplication of messages with only different names, the name in the message payload above is exactly the same for both the initial authentication and subsequent token updates after renewal.
{{< /hint >}}

#### Updating the token

When authentication is delegated, the automatic renewal of the token inside of ownCloud Web is disabled. In order to update the token, a listener is created which awaits a message with payload `{ name: 'owncloud-embed:update-token', data: { access_token: '<bearer-token>' } }`. The token will then be replaced inside of the Embed mode automatically.
