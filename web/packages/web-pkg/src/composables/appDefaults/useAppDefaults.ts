import { computed, unref, Ref } from 'vue'
import { useRouter, useRoute, useRouteParam } from '../router'
import { ClientService } from '../../services'
import { basename } from 'path'

import { FileContext } from './types'
import {
  useAppNavigation,
  AppNavigationResult,
  contextQueryToFileContextProps,
  contextRouteNameKey,
  queryItemAsString
} from './useAppNavigation'
import { useAppConfig, AppConfigResult } from './useAppConfig'
import { useAppFileHandling, AppFileHandlingResult } from './useAppFileHandling'
import { useAppFolderHandling, AppFolderHandlingResult } from './useAppFolderHandling'
import { useAppDocumentTitle } from './useAppDocumentTitle'
import { RequestResult, useRequest } from '../authContext'
import { useClientService } from '../clientService'
import { MaybeRef } from '../../utils'
import { useDriveResolver } from '../driveResolver'
import { urlJoin } from '@ownclouders/web-client'
import { useAppsStore, useAuthStore } from '../piniaStores'
import { storeToRefs } from 'pinia'

// TODO: this file/folder contains file/folder loading logic extracted from preview and drawio extensions
// Discussion how to progress from here can be found in this issue:
// https://github.com/owncloud/web/issues/3301

interface AppDefaultsOptions {
  applicationId: string
  applicationName?: MaybeRef<string>
  clientService?: ClientService
}

export type AppDefaultsResult = AppConfigResult &
  AppNavigationResult &
  AppFileHandlingResult &
  RequestResult &
  AppFolderHandlingResult & {
    isPublicLinkContext: Ref<boolean>
    currentFileContext: Ref<FileContext>
  }

export function useAppDefaults(options: AppDefaultsOptions): AppDefaultsResult {
  const router = useRouter()
  const appsStore = useAppsStore()
  const currentRoute = useRoute()
  const clientService = options.clientService ?? useClientService()
  const applicationId = options.applicationId

  const authStore = useAuthStore()
  const { publicLinkContextReady } = storeToRefs(authStore)

  const driveAliasAndItem = useRouteParam('driveAliasAndItem')
  const { space, item, itemId, loading } = useDriveResolver({ driveAliasAndItem })
  const currentFileContext = computed((): FileContext => {
    if (unref(loading)) {
      return null
    }
    let path: string
    if (unref(space)) {
      path = urlJoin(unref(space).webDavPath, unref(item))
    } else {
      // deprecated.
      path = urlJoin(queryItemAsString(unref(currentRoute)?.params?.filePath))
    }

    return {
      path,
      driveAliasAndItem: unref(driveAliasAndItem),
      space: unref(space),
      item: unref(item),
      itemId: unref(itemId),
      fileName: basename(path),
      routeName: queryItemAsString(unref(currentRoute).query[contextRouteNameKey]),
      ...contextQueryToFileContextProps(unref(currentRoute).query)
    }
  })

  useAppDocumentTitle({
    appsStore,
    applicationId,
    applicationName: options.applicationName,
    currentFileContext,
    currentRoute
  })

  return {
    isPublicLinkContext: publicLinkContextReady,
    currentFileContext,
    ...useAppConfig({ appsStore, ...options }),
    ...useAppNavigation({ router, currentFileContext }),
    ...useAppFileHandling({
      clientService
    }),
    ...useAppFolderHandling({
      clientService,
      currentRoute
    }),
    ...useRequest({ clientService, currentRoute: unref(currentRoute) })
  }
}
