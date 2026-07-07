import { Resource, SpaceResource, extractNameWithoutExtension } from '@ownclouders/web-client'
import { computed, Ref, unref } from 'vue'
import { useClientService } from '../../clientService'
import { FileAction, FileActionOptions } from '../types'
import { useGettext } from 'vue3-gettext'
import { resolveFileNameDuplicate } from '../../../helpers/resource'
import { join } from 'path'
import { WebDAV } from '@ownclouders/web-client/webdav'
import { EDITOR_MODE_CREATE, useFileActions } from './useFileActions'
import {
  useMessages,
  useModals,
  useUserStore,
  useAppsStore,
  useResourcesStore
} from '../../piniaStores'
import { ApplicationFileExtension } from '../../../apps'
import { storeToRefs } from 'pinia'
import { useEmbedMode } from '../../embedMode'

export const useFileActionsCreateNewFile = ({ space }: { space?: Ref<SpaceResource> } = {}) => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const { $gettext } = useGettext()
  const { dispatchModal } = useModals()
  const appsStore = useAppsStore()
  const { isEnabled: isEmbedModeEnabled } = useEmbedMode()

  const { openEditor } = useFileActions()
  const clientService = useClientService()

  const resourcesStore = useResourcesStore()
  const { resources, currentFolder, areFileExtensionsShown } = storeToRefs(resourcesStore)

  const appNewFileMenuExtensions = computed(() =>
    appsStore.fileExtensions.filter(({ newFileMenu }) => !!newFileMenu)
  )

  const getNameErrorMsg = (fileName: string) => {
    if (fileName === '') {
      return $gettext('File name cannot be empty')
    }

    if (/[/]/.test(fileName)) {
      return $gettext('File name cannot contain "/"')
    }

    if (fileName === '.') {
      return $gettext('File name cannot be equal to "."')
    }

    if (fileName === '..') {
      return $gettext('File name cannot be equal to ".."')
    }

    if (/\s+$/.test(fileName)) {
      return $gettext('File name cannot end with whitespace')
    }

    const exists = unref(resources).find((file) => file.name === fileName)

    if (exists) {
      return $gettext('%{name} already exists', { name: fileName }, true)
    }

    return null
  }

  const openFile = (resource: Resource, appFileExtension: ApplicationFileExtension) => {
    resourcesStore.upsertResource(resource)

    return openEditor(appFileExtension, unref(space), resource, EDITOR_MODE_CREATE)
  }

  const handler = (
    fileActionOptions: FileActionOptions,
    extension: string,
    appFileExtension: ApplicationFileExtension
  ) => {
    let defaultName = $gettext('New file') + `.${extension}`

    if (unref(resources).some((f) => f.name === defaultName)) {
      defaultName = resolveFileNameDuplicate(defaultName, extension, unref(resources))
    }

    if (!areFileExtensionsShown.value) {
      defaultName = extractNameWithoutExtension({ name: defaultName, extension } as Resource)
    }

    const inputSelectionRange = !areFileExtensionsShown.value
      ? null
      : ([0, defaultName.length - (extension.length + 1)] as [number, number])

    dispatchModal({
      title: $gettext('Create a new file'),
      confirmText: $gettext('Create'),
      hasInput: true,
      inputValue: defaultName,
      inputLabel: $gettext('File name'),
      inputSelectionRange,
      onConfirm: async (fileName: string) => {
        if (!areFileExtensionsShown.value) {
          fileName = `${fileName}.${extension}`
        }

        try {
          let resource: Resource
          if (appFileExtension.createFileHandler) {
            resource = await appFileExtension.createFileHandler({
              fileName,
              space: unref(space),
              currentFolder: unref(currentFolder)
            })
          } else {
            const path = join(unref(currentFolder).path, fileName)
            resource = await (clientService.webdav as WebDAV).putFileContents(unref(space), {
              path
            })
          }

          resourcesStore.upsertResource(resource)

          showMessage({
            title: $gettext('"%{fileName}" was created successfully', { fileName: resource.name })
          })

          if (unref(isEmbedModeEnabled)) {
            return
          }

          return openFile(resource, appFileExtension)
        } catch (error) {
          console.error(error)
          showErrorMessage({
            title: $gettext('Failed to create file'),
            errors: [error]
          })
        }
      },
      onInput: (name, setError) =>
        setError(getNameErrorMsg(areFileExtensionsShown.value ? name : `${name}.${extension}`))
    })
  }

  const actions = computed((): FileAction[] => {
    const actions: FileAction[] = []
    // make sure there is only one action for a file extension/mime-type
    // if there are
    // - multiple ApplicationFileExtensions with priority
    // or
    // - multiple ApplicationFileExtensions without priority (and none with)
    // we do not guarantee which one is chosen
    const defaultMapping: Record<string, ApplicationFileExtension> = {}
    for (const appFileExtension of unref(appNewFileMenuExtensions) || []) {
      if (appFileExtension.hasPriority) {
        defaultMapping[appFileExtension.extension] = appFileExtension
      } else {
        defaultMapping[appFileExtension.extension] =
          defaultMapping[appFileExtension.extension] || appFileExtension
      }
    }

    for (const [, appFileExtension] of Object.entries(defaultMapping)) {
      actions.push({
        name: 'create-new-file',
        icon: 'add',
        handler: (args) =>
          appFileExtension.customHandler
            ? appFileExtension.customHandler(args, appFileExtension.extension, appFileExtension)
            : handler(args, appFileExtension.extension, appFileExtension),
        label: () => $gettext(appFileExtension.newFileMenu.menuTitle()),
        isVisible: () => {
          if (
            typeof appFileExtension.newFileMenu.isVisible === 'function' &&
            !appFileExtension.newFileMenu.isVisible({ currentFolder: unref(currentFolder) })
          ) {
            return false
          }

          return unref(currentFolder)?.canUpload({ user: userStore.user })
        },
        class: 'oc-files-actions-create-new-file',
        ext: appFileExtension.extension,
        isExternal: appFileExtension.app?.startsWith('external-')
      })
    }

    return actions
  })

  return {
    actions,
    getNameErrorMsg,
    openFile
  }
}
