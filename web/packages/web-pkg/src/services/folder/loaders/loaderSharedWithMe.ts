import { FolderLoader, FolderLoaderTask, TaskContext } from '../folderService'
import { Router } from 'vue-router'
import { useTask } from 'vue-concurrency'
import { buildIncomingShareResource, call } from '@ownclouders/web-client'
import { isLocationSharesActive } from '../../../router'
import { unref } from 'vue'

export class FolderLoaderSharedWithMe implements FolderLoader {
  public isEnabled(): boolean {
    return true
  }

  public isActive(router: Router): boolean {
    const currentRoute = unref(router.currentRoute)
    return (
      isLocationSharesActive(router, 'files-shares-with-me') ||
      currentRoute?.query?.contextRouteName === 'files-shares-with-me'
    )
  }

  public getTask(context: TaskContext): FolderLoaderTask {
    const { spacesStore, clientService, configStore, resourcesStore, sharesStore } = context

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    return useTask(function* (signal1, signal2) {
      resourcesStore.clearResourceList()
      resourcesStore.setAncestorMetaData({})

      if (configStore.options.routing.fullShareOwnerPaths) {
        yield spacesStore.loadMountPoints({
          graphClient: clientService.graphAuthenticated,
          signal: signal1
        })
      }

      const value = yield* call(
        clientService.graphAuthenticated.driveItems.listSharedWithMe({ signal: signal1 })
      )

      const resources = value.map((driveItem) =>
        buildIncomingShareResource({
          driveItem,
          graphRoles: sharesStore.graphRoles,
          serverUrl: configStore.serverUrl
        })
      )

      resourcesStore.initResourceList({ currentFolder: null, resources })
    })
  }
}
