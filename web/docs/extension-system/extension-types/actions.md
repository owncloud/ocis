---
title: 'Action Extensions'
date: 2024-01-23T00:00:00+00:00
weight: 60
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/extension-system/extension-types
geekdocFilePath: actions.md
geekdocCollapseSection: true
---

{{< toc >}}

## Action extension type

Actions are one of the possible extension types. Registered actions get rendered in various places across the UI, depending on their scope and targets.

### Configuration

This is what the `ActionExtension` interface looks like:

```typescript
interface ActionExtension {
  id: string
  type: 'action'
  extensionPointIds?: string[]
  action: Action // Please check the Action section below
}
```

For `id`, `type`, and `extensionPointIds`, please see [extension base section]({{< ref "../_index.md#extension-base-configuration" >}}) in top level docs.

#### Action

The most important configuration options are:

- `icon` - The icon to be displayed, can be picked from https://owncloud.design/#/Design%20Tokens/IconList
- `name` - The name of the action (not displayed in the UI)
- `label` - The text to be displayed
- `route` - The string/route to navigate to. The nav item will be a `<router-link>` tag.
- `href` - The URL to navigate to. The nav item will be a `<a>`tag.
- `handler` - The action to perform upon click. The nav item will be a `<button>` tag.
- `isVisible` - Determines whether the action is displayed to the user

Please check the [`Action` type](https://github.com/owncloud/web/blob/236c185540fc6758dc7bd84985c8834fa4145530/packages/web-pkg/src/composables/actions/types.ts#L6) for a full list of configuration options.

### Example

The following example shows how an action extension for downloading files could look like. Note that the extension is wrapped inside a Vue composable so it can easily be reused. All helper types and composables are being provided via the [web-pkg](https://github.com/owncloud/web/tree/master/packages/web-pkg) package.

```typescript
export const useDownloadFilesExtension = () => {
  const { $gettext } = useGettext()

  const extension = computed<ActionExtension>(() => ({
    id: 'com.github.owncloud.web.files.download-action',
    scopes: ['resource'],
    type: 'action',
    action: {
      name: 'download-files',
      icon: 'download',
      class: 'oc-files-actions-download-files',
      label: () => $gettext('Download'),
      isVisible: ({ space, resources }) => {
        if (resources.length === 0) {
          return false
        }

        return true
      },
      handler: ({ space, resources }) => {
        console.log('Triggering download...')
      }
    }
  }))

  return { extension }
}
```

The extension could then be registered in any app like so:

```typescript
export default defineWebApplication({
  setup() {
    const { extension } = useFileActionDownloadFiles()

    return {
      appInfo: {
        name: $gettext('Download app'),
        id: 'download-app'
      },
      extensions: computed(() => [unref(extension)])
    }
  }
})
```
