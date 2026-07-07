import kebabCase from 'lodash-es/kebabCase'
import isNil from 'lodash-es/isNil'
import { isShareSpaceResource } from '@ownclouders/web-client'
import { routeToContextQuery } from '../../appDefaults'
import { isLocationTrashActive } from '../../../router'
import { computed, unref } from 'vue'
import { useRouter, useRoute } from '../../router'
import { useGettext } from 'vue3-gettext'
import {
  Action,
  FileAction,
  FileActionOptions,
  useIsSearchActive,
  useWindowOpen
} from '../../actions'

import {
  useFileActionsEnableSync,
  useFileActionsToggleHideShare,
  useFileActionsCopy,
  useFileActionsDisableSync,
  useFileActionsDelete,
  useFileActionsDownloadArchive,
  useFileActionsDownloadFile,
  useFileActionsFavorite,
  useFileActionsMove,
  useFileActionsNavigate,
  useFileActionsRename,
  useFileActionsRestore,
  useFileActionsCreateSpaceFromResource,
  useFileActionsDuplicate
} from './index'
import {
  ActionExtension,
  useAppsStore,
  useConfigStore,
  useExtensionRegistry
} from '../../piniaStores'
import { ApplicationFileExtension } from '../../../apps'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { storeToRefs } from 'pinia'
import { useEmbedMode } from '../../embedMode'
import { RouteRecordName } from 'vue-router'
import { isLocationActive } from '../../../router/utils'

export const EDITOR_MODE_EDIT = 'edit'
export const EDITOR_MODE_CREATE = 'create'

export interface GetFileActionsOptions extends FileActionOptions {
  omitSystemActions?: boolean
}

export const useFileActions = () => {
  const appsStore = useAppsStore()
  const router = useRouter()
  const route = useRoute()
  const { $gettext } = useGettext()
  const isSearchActive = useIsSearchActive()
  const { isEnabled: isEmbedModeEnabled } = useEmbedMode()
  const { requestExtensions } = useExtensionRegistry()

  const { openUrl } = useWindowOpen()

  const configStore = useConfigStore()
  const { options } = storeToRefs(configStore)

  const { actions: enableSyncActions } = useFileActionsEnableSync()
  const { actions: hideShareActions } = useFileActionsToggleHideShare()
  const { actions: copyActions } = useFileActionsCopy()
  const { actions: deleteActions } = useFileActionsDelete()
  const { actions: disableSyncActions } = useFileActionsDisableSync()
  const { actions: downloadArchiveActions } = useFileActionsDownloadArchive()
  const { actions: downloadFileActions } = useFileActionsDownloadFile()
  const { actions: favoriteActions } = useFileActionsFavorite()
  const { actions: moveActions } = useFileActionsMove()
  const { actions: navigateActions } = useFileActionsNavigate()
  const { actions: renameActions } = useFileActionsRename()
  const { actions: restoreActions } = useFileActionsRestore()
  const { actions: createSpaceFromResource } = useFileActionsCreateSpaceFromResource()
  const { actions: duplicateActions } = useFileActionsDuplicate()

  const systemActions = computed((): Action[] => [
    ...unref(downloadArchiveActions),
    ...unref(downloadFileActions),
    ...unref(deleteActions),
    ...unref(moveActions),
    ...unref(copyActions),
    ...unref(renameActions),
    ...unref(duplicateActions),
    ...unref(createSpaceFromResource),
    ...unref(restoreActions),
    ...unref(enableSyncActions),
    ...unref(hideShareActions),
    ...unref(disableSyncActions),
    ...unref(favoriteActions),
    ...unref(navigateActions)
  ])

  const defaultActions = computed<FileAction[]>(() => {
    const contextActionExtensions = requestExtensions<ActionExtension>({
      id: 'global.files.default-actions',
      extensionType: 'action'
    })
    return contextActionExtensions.map((extension) => extension.action)
  })

  const extensionActions = computed(() => {
    return requestExtensions<ActionExtension>({
      id: 'global.files.context-actions',
      extensionType: 'action'
    }).map((e) => e.action)
  })

  const editorActions = computed(() => {
    if (unref(isEmbedModeEnabled)) {
      return []
    }

    return appsStore.fileExtensions
      .filter((fileExtension) => appsStore.apps[fileExtension.app]?.hasEditor)
      .map((fileExtension): FileAction => {
        const appInfo = appsStore.apps[fileExtension.app]

        return {
          name: `editor-${fileExtension.app}`,
          label: () => {
            if (fileExtension.label) {
              if (typeof fileExtension.label === 'function') {
                return fileExtension.label()
              }
              return fileExtension.label
            }
            return $gettext('Open in %{app}', { app: appInfo.name }, true)
          },
          showOpenInNewTabHint: true,
          icon: fileExtension.icon || appInfo.icon,
          ...(appInfo.iconFillType && {
            iconFillType: appInfo.iconFillType
          }),
          img: appInfo.img,
          route: ({ space, resources }) => {
            return getEditorRoute({
              appFileExtension: fileExtension,
              space,
              resource: resources[0],
              mode: EDITOR_MODE_EDIT
            })
          },
          handler: (options) =>
            openEditor(fileExtension, options.space, options.resources[0], EDITOR_MODE_EDIT),
          isVisible: ({ resources }) => {
            if (resources.length !== 1) {
              return false
            }

            if (!resources[0].canDownload() && !fileExtension.secureView) {
              return false
            }

            if (!unref(isSearchActive) && isLocationTrashActive(router, 'files-trash-generic')) {
              return false
            }

            if (isLocationActive(router, { name: fileExtension.routeName || fileExtension.app })) {
              return false
            }

            if (resources[0].extension && fileExtension.extension) {
              return resources[0].extension.toLowerCase() === fileExtension.extension.toLowerCase()
            }

            if (resources[0].mimeType && fileExtension.mimeType) {
              return (
                resources[0].mimeType.toLowerCase() === fileExtension.mimeType.toLowerCase() ||
                resources[0].mimeType.split('/')[0].toLowerCase() ===
                  fileExtension.mimeType.toLowerCase()
              )
            }

            return false
          },
          hasPriority: fileExtension.hasPriority,
          class: `oc-files-actions-${kebabCase(appInfo.name).toLowerCase()}-trigger`
        }
      })
      .sort((first, second) => {
        // Ensure default are listed first
        if (second.hasPriority !== first.hasPriority && second.hasPriority) {
          return 1
        }
        return 0
      })
  })

  const getEditorRoute = ({
    appFileExtension,
    space,
    resource,
    mode
  }: {
    appFileExtension: ApplicationFileExtension
    space: SpaceResource
    resource: Resource
    mode: string
  }) => {
    const remoteItemId = isShareSpaceResource(space) ? space.id : undefined
    const routeName = appFileExtension.routeName || appFileExtension.app
    const routeOpts = getEditorRouteOpts(routeName, space, resource, mode, remoteItemId)
    return router.resolve(routeOpts)
  }
  const getEditorRouteOpts = (
    routeName: RouteRecordName,
    space: SpaceResource,
    resource: Resource,
    mode: string,
    remoteItemId: string,
    templateId?: string
  ) => {
    return {
      name: routeName,
      params: {
        driveAliasAndItem: space?.getDriveAliasAndItem(resource),
        filePath: resource.path,
        fileId: resource.fileId,
        scope: route.value.params.scope,
        mode
      },
      query: {
        ...(remoteItemId && { shareId: remoteItemId }),
        ...(resource.fileId && unref(options).routing.idBased && { fileId: resource.fileId }),
        ...(templateId && { templateId }),
        ...routeToContextQuery(unref(router.currentRoute))
      }
    }
  }

  const openEditor = (
    appFileExtension: ApplicationFileExtension,
    space: SpaceResource,
    resource: Resource,
    mode: string
  ) => {
    const remoteItemId = isShareSpaceResource(space) ? space.id : undefined
    const routeName = appFileExtension.routeName || appFileExtension.app
    const routeOpts = getEditorRouteOpts(routeName, space, resource, mode, remoteItemId)

    if (unref(options).cernFeatures) {
      const path = router.resolve(routeOpts).href
      const target = `${appFileExtension.routeName}-${resource.path}`

      openUrl(path, target, true)
      return
    }

    router.push(routeOpts)
  }

  // TODO: Make user-configurable what is a defaultAction for a filetype/mimetype
  // returns the _first_ action from actions array which we now construct from
  // available mime-types coming from the app-provider and existing actions
  const triggerDefaultAction = (options: FileActionOptions) => {
    const action = getDefaultAction(options)
    action.handler({ ...options })
  }

  const getDefaultAction = (options: GetFileActionsOptions): Action | undefined => {
    const allActions = getAllAvailableActions(options)
    if (allActions.length) {
      return allActions[0]
    }
    return undefined
  }

  const getAllAvailableActions = (options: GetFileActionsOptions) => {
    const filterCallback = (action: FileAction) => action.isVisible(options)

    const primaryActions = [...unref(defaultActions), ...unref(editorActions)]
      .filter(filterCallback)
      .sort((a, b) => Number(b.hasPriority) - Number(a.hasPriority))

    const secondaryActions = options.omitSystemActions
      ? []
      : unref(systemActions).filter(filterCallback)

    return [
      ...primaryActions,
      ...secondaryActions,
      ...unref(extensionActions).filter(
        (a) =>
          a.isVisible(options as FileActionOptions) &&
          (a.category === 'actions' || isNil(a.category))
      )
    ]
  }

  return {
    editorActions,
    systemActions,
    defaultActions,
    getDefaultAction,
    getAllAvailableActions,
    getEditorRouteOpts,
    openEditor,
    triggerDefaultAction
  }
}
