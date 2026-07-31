import { Ref, ref, unref, MaybeRef } from 'vue'
import { dirname } from 'path'
import { ClientService, folderService } from '../../services'
import { useAppFileHandling } from './useAppFileHandling'
import { isSearchResource, Resource } from '@ownclouders/web-client'
import { FileContext } from './types'
import { RouteLocationNormalizedLoaded } from 'vue-router'
import { useFileRouteReplace } from '../router/useFileRouteReplace'
import { DavProperty } from '@ownclouders/web-client/webdav'
import { useAuthService } from '../authContext/useAuthService'
import { isMountPointSpaceResource } from '@ownclouders/web-client'
import { useResourcesStore, useSpacesStore } from '../piniaStores'
import { storeToRefs } from 'pinia'
import { useRouteQuery } from '../router'
import { useSearch } from '../search'

interface AppFolderHandlingOptions {
  currentRoute: Ref<RouteLocationNormalizedLoaded>
  clientService?: ClientService
}

export interface AppFolderHandlingResult {
  isFolderLoading: Ref<boolean>
  activeFiles: Ref<Array<Resource>>

  loadFolderForFileContext(context: MaybeRef<FileContext>): Promise<any>
}

export function useAppFolderHandling({
  currentRoute,
  clientService
}: AppFolderHandlingOptions): AppFolderHandlingResult {
  const isFolderLoading = ref(false)
  const { webdav } = clientService
  const { replaceInvalidFileRoute } = useFileRouteReplace()
  const { getFileInfo } = useAppFileHandling({ clientService })
  const authService = useAuthService()
  const spacesStore = useSpacesStore()
  const { buildSearchTerm, search } = useSearch()
  const currentRouteQuery = useRouteQuery('contextRouteQuery')

  const resourcesStore = useResourcesStore()
  const { activeResources } = storeToRefs(resourcesStore)

  const loadFolderForFileContext = async (context: MaybeRef<FileContext>) => {
    isFolderLoading.value = true

    try {
      context = unref(context)

      if ((unref(currentRouteQuery) as any)?.term) {
        // run search query to load all results
        // TODO: add filters from query params
        const searchTerm = buildSearchTerm({ term: (unref(currentRouteQuery) as any)?.term })
        const { values } = await search(searchTerm, 200)
        const resources = values
          .filter(({ data }) => isSearchResource(data as Resource))
          .map<Resource>((v) => v.data as Resource)

        resourcesStore.initResourceList({ currentFolder: null, resources })
        isFolderLoading.value = false
        return
      }

      const flatFileLists = [
        'files-shares-with-me',
        'files-shares-with-others',
        'files-shares-via-link',
        'files-common-favorites'
      ]

      if (flatFileLists.includes(unref(context.routeName))) {
        // use the folder loader to load the resources for flat file lists
        const loaderTask = folderService.getTask()
        await loaderTask.perform()
        isFolderLoading.value = false
        return
      }

      resourcesStore.clearResourceList()
      const space = unref(context.space)
      const pathResource = await getFileInfo(context, {
        davProperties: [DavProperty.FileId]
      })
      replaceInvalidFileRoute({
        space,
        resource: pathResource,
        path: unref(context.item),
        fileId: unref(context.itemId)
      })

      const isSpaceRoot = spacesStore.spaces.some(
        (s) => isMountPointSpaceResource(s) && s.root.remoteItem?.id === pathResource.id
      )

      if (isSpaceRoot) {
        const resource = await getFileInfo(context)
        resourcesStore.initResourceList({ currentFolder: resource, resources: [resource] })
        isFolderLoading.value = false
        return
      }

      const path = dirname(pathResource.path)
      const { resource, children } = await webdav.listFiles(space, {
        path
      })

      if (resource.type === 'file') {
        resourcesStore.initResourceList({
          // FIXME: currentFolder should be null?!
          currentFolder: resource,
          resources: [resource]
        })
      } else {
        resourcesStore.initResourceList({ currentFolder: resource, resources: children })
      }
    } catch (error) {
      if (error.statusCode === 401) {
        return authService.handleAuthError(unref(currentRoute))
      }
      resourcesStore.setCurrentFolder(null)
      console.error(error)

      throw error
    } finally {
      isFolderLoading.value = false
    }
  }

  return {
    isFolderLoading,
    loadFolderForFileContext,
    activeFiles: activeResources
  }
}
