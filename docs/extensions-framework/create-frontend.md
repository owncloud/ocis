---
title: "Create frontend"
date: 2020-05-28T10:39:00+01:00
weight: 2
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: extensions-framework/create-frontend.md
---

{{< toc >}}

## Extension points
As of today, there are several extension points inside of Phoenix.

| Name          | Description |
| :------------ | :---------- |
| App switcher  | Enables navigation between different extensions. An entry for the extension gets registered automatically in the case that at least one nav item is defined. If you wish to register an entry manually, you can do so in the config.json |
| App container | Container for the UI of the extension which lives directly under the top bar. | 
| Routes        | Routes used by [Vue Router](https://router.vuejs.org/) to enable accessing views of the extension. |
| Nav items     | Nav items included in the navigation sidebar pointing to their assigned routes. |
| Store         | A global store which can be accessed by any other extension. |

In addition to all the Phoenix extension points, we have defined also the following extension points inside of our [files app](https://github.com/owncloud/phoenix/tree/master/apps/files).

| Name                   | Description |
| :------------          | :---------- |
| File action            | Item inside of file actions menu |
| Create new file action | Item inside of new file actions menu | 
| Sidebar                | Files app sidebar with highlighted file as a context |

## UI Bundle
### Format
The resulting UI bundle provided by the extension must be an [AMD module](https://en.wikipedia.org/wiki/Asynchronous_module_definition).

### Index file
To make sure that all the UI parts are loaded correctly, it is necessary to have an index file which follows our predefined structure.

#### AppInfo
All necessary information about the extension.

```js
const appInfo = {
  name: 'Example extension',
  id: 'example-extension',
  icon: 'document',
  // Following values are optional and part of the files app extension points
  // In case the extension is a file editor, you can register it by providing file extensions
  extensions: [
    {
      extension: 'txt',
      // Optionally you can register an action inside of the Create new menu
      newFileMenu: {
        menuTitle: () => 'New plain text file'
      }
    }
  ],
  // Right sidebar with the highlighted file as context
  fileSideBars: [
    {
      app: 'example-sidebar',
      component: ExampleSidebarComponent,
      enabled: () => true
    }
  ]
}
```

#### Routes and nav items
```js
const routes = [
  {
    name: 'example-route',
    path: '/',
    components: {
      app: ExampleExtensionComponent
    }
  }
]
const navItems = [
  {
    name: 'Example nav item',
    iconMaterial: 'home',
    route: {
      name: 'example-route'
    }
  }
]
```

#### Final export
In the final export should be also included store and translations, if exists.

```js
export default {
  appInfo,
  routes,
  navItems,
  store,
  translations
}
```
