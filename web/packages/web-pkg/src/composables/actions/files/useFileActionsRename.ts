import { isSameResource } from '../../../helpers/resource'
import { isLocationTrashActive, isLocationSharesActive } from '../../../router'
import { Resource } from '@ownclouders/web-client'
import { dirname, join } from 'path'
import { WebDAV } from '@ownclouders/web-client/webdav'
import {
  SpaceResource,
  isShareSpaceResource,
  extractNameWithoutExtension
} from '@ownclouders/web-client'
import { createFileRouteOptions } from '../../../helpers/router'
import { renameResource as _renameResource } from '../../../helpers/resource'
import { computed } from 'vue'
import { useClientService } from '../../clientService'
import { useRouter } from '../../router'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../types'
import {
  useMessages,
  useModals,
  useCapabilityStore,
  useConfigStore,
  useResourcesStore,
  useUserStore
} from '../../piniaStores'
import { useAbility } from '../../ability'

export const useFileActionsRename = () => {
  const { showErrorMessage } = useMessages()
  const capabilityStore = useCapabilityStore()
  const router = useRouter()
  const { $gettext } = useGettext()
  const clientService = useClientService()
  const configStore = useConfigStore()
  const { dispatchModal } = useModals()
  const userStore = useUserStore()
  const ability = useAbility()

  const resourcesStore = useResourcesStore()
  const { setCurrentFolder, upsertResource } = resourcesStore

  const getNameErrorMsg = (
    resource: Resource,
    newName: string,
    parentResources: Resource[] = undefined
  ) => {
    const newPath =
      resource.path.substring(0, resource.path.length - resource.name.length) + newName

    if (!newName) {
      return $gettext('The name cannot be empty')
    }

    if (/[/]/.test(newName)) {
      return $gettext('The name cannot contain "/"')
    }

    if (newName === '.') {
      return $gettext('The name cannot be equal to "."')
    }

    if (newName === '..') {
      return $gettext('The name cannot be equal to ".."')
    }

    if (/\s+$/.test(newName)) {
      return $gettext('The name cannot end with whitespace')
    }

    const exists = resourcesStore.resources.find(
      (file) => file.path === newPath && resource.name !== newName
    )
    if (exists) {
      const translated = $gettext('The name "%{name}" is already taken')
      return $gettext(translated, { name: newName }, true)
    }

    if (parentResources) {
      const exists = parentResources.find(
        (file) => file.path === newPath && resource.name !== newName
      )

      if (exists) {
        const translated = $gettext('The name "%{name}" is already taken')
        return $gettext(translated, { name: newName }, true)
      }
    }

    return null
  }

  const renameResource = async (space: SpaceResource, resource: Resource, newName: string) => {
    let currentFolder = resourcesStore.currentFolder

    try {
      const newPath = join(dirname(resource.path), newName)
      await (clientService.webdav as WebDAV).moveFiles(space, resource, space, {
        path: newPath
      })

      const isCurrentFolder = isSameResource(resource, currentFolder)

      if (isShareSpaceResource(space) && resource.isReceivedShare()) {
        space.rename(newName)

        if (isCurrentFolder) {
          currentFolder = { ...currentFolder } as Resource
          currentFolder.name = newName
          setCurrentFolder(currentFolder)
          return router.push(
            createFileRouteOptions(space, {
              path: '',
              fileId: resource.fileId
            })
          )
        }

        const sharedResource = { ...resource }
        sharedResource.name = newName
        upsertResource(sharedResource)
        return
      }

      if (isCurrentFolder) {
        currentFolder = { ...currentFolder } as Resource
        _renameResource(space, currentFolder, newPath)
        setCurrentFolder(currentFolder)
        return router.push(
          createFileRouteOptions(space, {
            path: newPath,
            fileId: resource.fileId
          })
        )
      }
      const fileResource = { ...resource } as Resource
      _renameResource(space, fileResource, newPath)
      upsertResource(fileResource)
    } catch (error) {
      console.error(error)
      let title = $gettext(
        'Failed to rename "%{file}" to "%{newName}"',
        { file: resource.name, newName },
        true
      )
      if (error.statusCode === 423) {
        title = $gettext(
          'Failed to rename "%{file}" to "%{newName}" - the file is locked',
          { file: resource.name, newName },
          true
        )
      }
      showErrorMessage({ title, errors: [error] })
    }
  }

  const handler = async ({ space, resources }: FileActionOptions) => {
    const currentFolder = resourcesStore.currentFolder
    let parentResources: Resource[]
    if (isSameResource(resources[0], currentFolder)) {
      const parentPath = dirname(currentFolder.path)
      parentResources = (await clientService.webdav.listFiles(space, { path: parentPath })).children
    }

    const areFileExtensionsShown = resourcesStore.areFileExtensionsShown
    const onConfirm = async (newName: string) => {
      if (!areFileExtensionsShown) {
        newName = `${newName}.${resources[0].extension}`
      }

      await renameResource(space, resources[0], newName)
    }
    const checkName = (newName: string, setError: (error: string) => void) => {
      if (!areFileExtensionsShown) {
        newName = `${newName}.${resources[0].extension}`
      }

      const error = getNameErrorMsg(resources[0], newName, parentResources)
      setError(error)
    }
    const nameWithoutExtension = extractNameWithoutExtension(resources[0])
    const modalTitle =
      !resources[0].isFolder && !areFileExtensionsShown ? nameWithoutExtension : resources[0].name

    const title = resources[0].isFolder
      ? $gettext('Rename folder %{name}', { name: modalTitle })
      : $gettext('Rename file %{name}', { name: modalTitle })

    const inputValue =
      !resources[0].isFolder && !areFileExtensionsShown ? nameWithoutExtension : resources[0].name

    const inputSelectionRange =
      resources[0].isFolder || !areFileExtensionsShown
        ? null
        : ([0, nameWithoutExtension.length] as [number, number])

    dispatchModal({
      variation: 'passive',
      title,
      confirmText: $gettext('Rename'),
      hasInput: true,
      inputValue,
      inputSelectionRange,
      inputLabel: resources[0].isFolder ? $gettext('Folder name') : $gettext('File name'),
      onConfirm,
      onInput: checkName
    })
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'rename',
      icon: 'pencil',
      label: () => {
        return $gettext('Rename')
      },
      handler,
      isVisible: ({ resources }) => {
        if (isLocationTrashActive(router, 'files-trash-generic')) {
          return false
        }
        if (
          isLocationSharesActive(router, 'files-shares-with-me') &&
          !capabilityStore.sharingCanRename
        ) {
          return false
        }
        if (resources.length !== 1) {
          return false
        }

        // FIXME: Remove this check as soon as renaming shares works as expected
        // see https://github.com/owncloud/ocis/issues/4866
        const rootShareIncluded = configStore.options.routing.fullShareOwnerPaths
          ? resources.some((r) => r.remoteItemPath && r.path)
          : resources.some((r) => r.remoteItemId && r.path === '/')
        if (rootShareIncluded) {
          return false
        }

        if (resources.length === 1 && resources[0].locked) {
          return false
        }

        const renameDisabled = resources.some((resource) => {
          return !resource.canRename({ user: userStore.user, ability }) || resource.processing
        })
        return !renameDisabled
      },
      class: 'oc-files-actions-rename-trigger'
    }
  ])

  return {
    actions,
    // HACK: exported for unit tests:
    getNameErrorMsg,
    renameResource
  }
}
