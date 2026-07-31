import { useRouter } from './useRouter'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { createFileRouteOptions } from '../../helpers/router'
import { Router } from 'vue-router'
import { ConfigStore, useConfigStore } from '../piniaStores'

export interface FileRouteReplaceOptions {
  router?: Router
  configStore?: ConfigStore
}

export const useFileRouteReplace = (options: FileRouteReplaceOptions = {}) => {
  const router = options.router || useRouter()
  const configStore = options.configStore || useConfigStore()

  const replaceInvalidFileRoute = ({
    space,
    resource,
    path,
    fileId
  }: {
    space: SpaceResource
    resource: Resource
    path: string
    fileId?: string | number
  }): boolean => {
    if (!configStore.options.routing?.idBased) {
      return false
    }
    if (path === resource.path && fileId === resource.fileId) {
      return false
    }

    const routeOptions = createFileRouteOptions(space, resource, { configStore })
    router.replace(routeOptions)
    return true
  }

  return {
    replaceInvalidFileRoute
  }
}
