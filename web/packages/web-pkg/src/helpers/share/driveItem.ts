import { SpaceResource } from '@ownclouders/web-client'
import { Graph } from '@ownclouders/web-client/graph'
import { SpacesStore } from '../../composables'

/**
 * Gets the drive item for a given shared space.
 */
export const getSharedDriveItem = async ({
  graphClient,
  spacesStore,
  space,
  signal
}: {
  graphClient: Graph
  spacesStore: SpacesStore
  space: SpaceResource
  signal?: AbortSignal
}) => {
  const matchingMountPoint = await spacesStore.getMountPointForSpace({
    graphClient,
    space,
    signal
  })
  if (!matchingMountPoint) {
    return null
  }
  const { id } = matchingMountPoint
  return graphClient.driveItems.getDriveItem(id.split('!')[0], id)
}
