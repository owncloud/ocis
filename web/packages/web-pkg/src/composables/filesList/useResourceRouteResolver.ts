import { unref, Ref } from 'vue'

import { useGetMatchingSpace } from '../spaces'
import { createFileRouteOptions } from '../../helpers/router'
import { createLocationSpaces } from '../../router'
import { CreateTargetRouteOptions } from '../../helpers/folderLink/types'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { ConfigStore } from '../piniaStores'

export type ResourceRouteResolverOptions = {
  configStore?: ConfigStore
  targetRouteCallback?: Ref<(arg: CreateTargetRouteOptions) => unknown>
  space?: Ref<SpaceResource>
}

export const useResourceRouteResolver = (
  options: ResourceRouteResolverOptions = {},
  emit?: any
) => {
  const targetRouteCallback = options.targetRouteCallback
  const { getMatchingSpace } = useGetMatchingSpace(options)

  const createFolderLink = (createTargetRouteOptions: CreateTargetRouteOptions) => {
    if (unref(targetRouteCallback)) {
      return unref(targetRouteCallback)(createTargetRouteOptions)
    }

    const { path, fileId, resource } = createTargetRouteOptions

    const space = unref(options.space) || getMatchingSpace(resource)
    if (!space) {
      return {}
    }
    return createLocationSpaces(
      'files-spaces-generic',
      createFileRouteOptions(space, { path, fileId })
    )
  }

  const createFileAction = (resource: Resource) => {
    const space = unref(options.space) || getMatchingSpace(resource)
    /**
     * Triggered when a default action is triggered on a file
     * @property {object} resource resource for which the event is triggered
     */
    emit('fileClick', { space, resources: [resource] })
  }

  return {
    createFileAction,
    createFolderLink
  }
}
