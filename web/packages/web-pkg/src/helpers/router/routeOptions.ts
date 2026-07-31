import { isShareSpaceResource, Resource, SpaceResource } from '@ownclouders/web-client'
import { ConfigStore, LocationQuery, useConfigStore } from '../../composables'
import { RouteParams } from 'vue-router'
import { isUndefined } from 'lodash-es'

/**
 * Creates route options for routing into a file location:
 * - params.driveAliasAndItem
 * - query.shareId
 * - query.fileId
 *
 * Both query options are optional.
 *
 * @param space {SpaceResource}
 * @param target {path: string, fileId: string | number}
 * @param options {configStore: ConfigStore}
 */
export const createFileRouteOptions = (
  space: SpaceResource,
  target: { path?: string; fileId?: string | number } = {},
  options?: { configStore: ConfigStore }
): { params: RouteParams; query: LocationQuery } => {
  const config = options?.configStore || useConfigStore()
  return {
    params: {
      driveAliasAndItem: space.getDriveAliasAndItem({ path: target.path || '' } as Resource)
    },
    query: {
      ...(isShareSpaceResource(space) && { shareId: space.id }),
      ...(config?.options?.routing?.idBased &&
        !isUndefined(target.fileId) && { fileId: `${target.fileId}` })
    }
  }
}
