---
title: "Deploy as an app in ownCloud Classic"
date: 2018-05-02T00:00:00+00:00
weight: 1
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/deployments
geekdocFilePath: oc10-app.md
---

{{< toc >}}


## Compatibility

Please note that the usage of Web UI and ownCloud Classic as backend is not recommended starting with version 7.1.0 of the Web UI. Therefore, this section only applies to versions < 7.1.0.

## Introduction

ownCloud Web is being deployed as an app to [ownCloud marketplace](https://marketplace.owncloud.com/) to enable easy integration into existing ownCloud Classic instances.
After completing this setup, ownCloud Web will be available on `https://<your-owncloud-server>/index.php/apps/web`.

## Prerequisites
- Running [ownCloud Classic server](https://owncloud.com/download-server/) with version 10.8
- Installed [oauth2 app](https://marketplace.owncloud.com/apps/oauth2)
- Command line access to your server

## Deploying ownCloud Web
Download the [ownCloud Web app](https://marketplace.owncloud.com/apps/web) from the marketplace and enable it:
```console
occ market:install web
```

## Configure oauth2
Within the `Admin` page of ownCloud Classic, head into `User Authentication` and add a new client with arbitrary name (e.g. `ownCloud Web`) and redirection URL `https://<your-owncloud-server>/index.php/apps/web/oidc-callback.html`.

{{< figure src="/clients/web/static/oauth2.png" alt="Example OAuth2 entry" >}}

{{< hint >}}
You can mark the ownCloud web client as `trusted` by clicking the respective checkbox so authorization after authentication gets omitted.
{{< /hint >}}

{{< hint >}}
If you use OpenID Connect you need to add a new client for ownCloud Web to your identity provider instead.
{{< /hint >}}

## Configure ownCloud Classic
### Set ownCloud Web address
To set the ownCloud Web address and to display ownCloud Web in the app switcher, add the following line into `config/config.php`:

```php
'web.baseUrl' => 'https://<your-owncloud-server>/index.php/apps/web',
```

### Configure link routing
Administrators can optionally decide whether ownCloud Links (public and private links) should be provided by the Classic web interface or by ownCloud Web using the `web.rewriteLinks` option in `config/config.php`. The option defaults to `false` so that the links open in the Classic web interface. Setting it to `true` will redirect all links to ownCloud Web. To redirect all private and public links to ownCloud Web, add the following line into `config/config.php`:

```php
'web.rewriteLinks' => true,
```

### Make ownCloud Web the default web interface
Administrators can optionally decide to make ownCloud Web the default web interface that users see after they log in to ownCloud. By default, the Classic web interface will be presented to users. To present ownCloud Web to users by default, add the following line into `config/config.php`:

```php
'defaultapp' => 'web',
```

{{< hint info >}}
While it is possible to make ownCloud Web the default web interface, the decision should be carefully evaluated. Features are still being added to ownCloud Web and users might need to use the Classic web interface to do certain actions.
{{< /hint >}}

## Configure ownCloud Web
There are a few config values which need to be set in order for ownCloud Web to work correctly. Please copy the example config below into `config/config.json` and adjust it for your environment:

```json
{
  "server" : "https://<your-owncloud-server>",
  "theme": "https://<your-owncloud-server>/index.php/apps/web/themes/owncloud/theme.json",
  "auth": {
    "clientId": "<client-id-from-oauth2>",
    "url": "https://<your-owncloud-server>/index.php/apps/oauth2/api/v1/token",
    "authUrl": "https://<your-owncloud-server>/index.php/apps/oauth2/authorize",
    "logoutUrl": "https://<your-owncloud-server>/index.php/logout"
  },
  "apps" : [
    "files",
    "preview",
    "search"
  ],
  "applications" : [
    {
      "title": {
        "en": "Classic Design",
        "de": "Klassisches Design",
        "fr": "Design classique",
        "zh": "文件"
      },
      "icon": "swap-box",
      "url": "https://<your-owncloud-server>/index.php/apps/files"
    },
    {
      "icon": "settings-4",
      "menu": "user",
      "target": "_self",
      "title": {
        "de": "Einstellungen",
        "en": "Settings"
      },
      "url": "https://<your-owncloud-server>/index.php/settings/personal"
    }
  ]
}
```

{{< hint info >}}
If any issues arise when trying to access the new design, a good start for debugging it is to run your `config.json` file through a json validator of your choice.
{{< /hint >}}

|config parameter|explanation|
|---|---|
|server|ownCloud Classic server address|
|theme|Theme to be used in ownCloud Web pointing to a json file inside of `themes` folder|
|auth.clientId|Client ID received when adding ownCloud Web in the `User Authentication` section in `Admin`|
|apps|List of internal extensions to be loaded|
|applications|Additional apps and links to be displayed in the application switcher or in the user menu|
|applications[0].title|Visible title in the application switcher or user menu, localizable|
|applications[1].menu|Use `user` to move the menu item into the user menu. Defaults to app switcher|

{{< hint info >}}
It is important that you don't edit or place the `config.json` within the app folder. If you do, the integrity check of the app will fail and raise warnings.
{{< /hint >}}

{{< hint >}}
If you use OpenID Connect you need to replace the `"auth"` part with following configuration:

```json
  "openIdConnect": {
    "metadata_url": "<fqdn-of-the-identity-provider>/.well-known/openid-configuration",
    "authority": "<fqdn-of-the-identity-provider>",
    "client_id": "<client-id-from-the-identity-provider>",
    "response_type": "code",
    "scope": "openid profile email"
  }
```
{{< /hint >}}

## Integrate ownCloud Classic features in ownCloud Web
### Add links to the app switcher
ownCloud Classic features that are not deeply integrated with the Classic UI (e.g., full screen apps) can be added to the ownCloud Web app switcher so that users can easily access them from ownCloud Web. You can use the following example and customize it according to your needs.

{{< hint info >}}
All apps that are listed in the ownCloud Classic app switcher will be added as links to the app switcher of the new ownCloud Web automatically. All of those links will open in a new browser tab on click.
{{< /hint >}}

To add new elements in the app switcher, paste the following into the `applications` section of `config.json`:

```json
    {
      "title": {
        "en": "Custom Groups",
        "de": "Benutzerdefinierte Gruppen" 
      },
      "icon": "settings-4",
      "url": "https://<your-owncloud-server>/settings/personal?sectionid=customgroups"
    }
```

{{< hint info >}}
The URL in the example might need adaptations depending on the configuration of your ownCloud Classic. App switcher elements added this way will open the respective page in a new tab. This method can also be used to link external sites like Help pages or similar.
{{< /hint >}}

### Add links to the user menu
Just like adding links to the app switcher, you can also add links to the user menu.

```json
    {
      "icon": "settings-4",
      "menu": "user",
      "target": "_self",
      "title": {
        "de": "Hilfe",
        "en": "Help"
      },
      "url": "https://help-link.example"
    }
```

This will add a link to the specified URL in the user menu. This way, the link will open in the same tab. If you instead want to open it in a new tab, just remove the line `"target": "_self",`.

### ONLYOFFICE
For ONLYOFFICE there is a [native integration](https://github.com/ONLYOFFICE/onlyoffice-owncloud-web) available for ownCloud Web when it is used with ownCloud Classic. It fully integrates the ONLYOFFICE Document Editors and allows users to create and open documents right from ownCloud Web.

To be able to use ONLYOFFICE in ownCloud Web, it is required to run
- ownCloud Classic >= 10.8
- ownCloud Web >= 4.0.0
- [ONLYOFFICE Connector for ownCloud Classic](https://marketplace.owncloud.com/apps/onlyoffice) >= 7.1.1

Make sure that ONLYOFFICE works as expected in the Classic UI and add the following to `config.json` to make it available in ownCloud Web:

```json
"external_apps": [
    {
        "id": "onlyoffice",
        "path": "https://<your-owncloud-server>/apps-external/onlyoffice/js/web/onlyoffice.js"
    }
]
```

{{< hint info >}}
The URL in the example might need adaptations depending on the configuration of your ownCloud Classic.
{{< /hint >}}

### Collabora Online
For Collabora Online there is a native integration available for ownCloud Web when it is used with ownCloud Classic. It fully integrates the Collabora Online Document Editors and allows users to create and open documents right from ownCloud Web.

To be able to use Collabora Online in ownCloud Web, it is required to run
- ownCloud Classic >= 10.8
- ownCloud Web >= 4.0.0
- [Collabora Online Connector for ownCloud Classic](https://marketplace.owncloud.com/apps/richdocuments) >= 2.7.0

Make sure that Collabora Online works as expected in the Classic UI and add the following to `config.json` to make it available in ownCloud Web:

```json
"external_apps": [
    {
        "id": "richdocuments",
        "path": "https://<your-owncloud-server>/apps/richdocuments/js/richdocuments.js"
    }
]
```

{{< hint info >}}
The URL in the example might need adaptations depending on the configuration of your ownCloud Classic.
{{< /hint >}}

## Additional configuration for certain core apps
There is additional configuration available for certain core apps. You can find them listed below.

### Preview app
In case the backend has additional preview providers configured there is no mechanism, yet, to announce those to the `Preview` app in ownCloud Web. As an intermediate solution you can add the additional supported mimeTypes to the `Preview` app by following these steps:
1. Remove the `"preview"` string from the `"apps"` section in your `config.json` file
2. Add the following config to your `config.json` file:
```json
"external_apps": [
    {
      "id": "preview",
      "path": "web-app-preview",
      "config": {
        "mimeTypes": ["image/tiff", "image/webp"]
      }
    }
  ]
```

If you already have an `"external_apps"` section, just add the preview app to the list. Please adjust the `"mimeTypes"` list according to your additional preview providers. See https://github.com/owncloud/files_mediaviewer#supporting-more-media-types for advise on how to add preview providers to the backend.

### Text-Editor app
The `text-editor` app provides a list of file extensions that the app is associated with, both for opening files and for creating new files. 
By default, only `.txt` and `.md` files appear in the file creation menu and offer the text-editor as default app on a left mouse click 
in the file list. For other file types the text-editor app only appears in the right mouse click context menu. In case you want to change this 
default set of primary file extensions for the text-editor you can overwrite it as follows:
1. Remove the `"text-editor"` string from the `"apps"` section in your `config.json` file
2. Add the following config to your `config.json` file:
```json
"external_apps": [
    {
      "id": "text-editor",
      "path": "web-app-text-editor",
      "config": {
        "primaryExtensions": ["txt", "yaml"]
      }
    }
  ]
```
With the above example config the text editor will offer creation of new files for `.txt` and `.yaml` files instead of `.txt` and `.md` files. 
Also, a left mouse click on any `.txt` or `.yaml` file will open the respective file in the text-editor app. In this example, `.md` files would 
not be opened in the text-editor by default anymore, but the text-editor will would appear in the context menu for the file as alternative app. 

If you already have an `"external_apps"` section, just add the preview app to the list. Please adjust the `"mimeTypes"` list according to your additional preview providers. See https://github.com/owncloud/files_mediaviewer#supporting-more-media-types for advise on how to add preview providers to the backend.
 
{{< hint info >}}
The reason why the app needs to be ported from the `apps` section to the `external_apps` section is that only the `external_apps` support additional configuration. There are plans to change the configuration of apps to give you a coherent admin experience in that regard. 
{{< /hint >}}

## Accessing ownCloud Web
After following all the steps, you should see a new entry in the application switcher called `New Design` which points to the ownCloud web.

{{< figure src="/clients/web/static/application-switcher-oc10.jpg" alt="ownCloud Classic application switcher" >}}
