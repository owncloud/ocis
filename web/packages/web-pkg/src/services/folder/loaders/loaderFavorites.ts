import { FolderLoader, FolderLoaderTask, TaskContext } from '../folderService'
import { Router } from 'vue-router'
import { useTask } from 'vue-concurrency'
import { buildResource } from '@ownclouders/web-client'
import { isLocationCommonActive } from '../../../router'
import { unref } from 'vue'

export class FolderLoaderFavorites implements FolderLoader {
  public isEnabled(): boolean {
    return true
  }

  public isActive(router: Router): boolean {
    const currentRoute = unref(router.currentRoute)
    return (
      isLocationCommonActive(router, 'files-common-favorites') ||
      currentRoute?.query?.contextRouteName === 'files-common-favorites'
    )
  }

  public getTask(context: TaskContext): FolderLoaderTask {
    const { resourcesStore, clientService, userStore } = context

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    return useTask(function* (signal1, signal2) {
      resourcesStore.clearResourceList()
      resourcesStore.setAncestorMetaData({})

      const { results } = yield clientService.webdav.listFavoriteFiles({
        username: userStore.user?.onPremisesSamAccountName,
        signal: signal1
      })

      const resources = results.map(buildResource)
      resourcesStore.initResourceList({ currentFolder: null, resources })
    })
  }
}
