import { FolderLoader, FolderLoaderTask, TaskContext } from '../folderService'
import { unref } from 'vue'
import { Router } from 'vue-router'
import { useTask } from 'vue-concurrency'
import { isLocationSharesActive } from '../../../router'
import { buildOutgoingShareResource, call } from '@ownclouders/web-client'

export class FolderLoaderSharedWithOthers implements FolderLoader {
  public isEnabled(): boolean {
    return true
  }

  public isActive(router: Router): boolean {
    const currentRoute = unref(router.currentRoute)
    return (
      isLocationSharesActive(router, 'files-shares-with-others') ||
      currentRoute?.query?.contextRouteName === 'files-shares-with-others'
    )
  }

  public getTask(context: TaskContext): FolderLoaderTask {
    const { userStore, spacesStore, clientService, configStore, resourcesStore } = context

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
        clientService.graphAuthenticated.driveItems.listSharedByMe({ signal: signal1 })
      )

      const resources = value
        .filter((s) => s.permissions.some(({ link }) => !link))
        .map((driveItem) =>
          buildOutgoingShareResource({
            driveItem,
            user: userStore.user,
            serverUrl: configStore.serverUrl
          })
        )

      resourcesStore.initResourceList({ currentFolder: null, resources })
    })
  }
}
