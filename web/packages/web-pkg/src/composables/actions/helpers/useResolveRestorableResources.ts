import { User } from '@ownclouders/web-client/graph/generated'
import {
  isProjectSpaceResource,
  Resource,
  SpaceResource,
  TrashResource
} from '@ownclouders/web-client'
import { DavProperties } from '@ownclouders/web-client/webdav'
import { useClientService } from '../../clientService'
import { useUserStore } from '../../piniaStores'

export const useResolveRestorableResources = () => {
  const clientService = useClientService()
  const userStore = useUserStore()

  const canRestoreDeleteForSpace = (space: SpaceResource, user: User) => {
    if (!isProjectSpaceResource(space)) {
      return true
    }
    return space.canRestoreFromTrashbin({ user })
  }

  /**
   * Resolves the trash-bin entries for a batch of just-deleted resources within a single space.
   * Returns null if the space doesn't support trash-restore, or if the trash listing fails.
   */
  const resolveRestorableResources = async (
    space: SpaceResource,
    deletedResources: Resource[]
  ): Promise<TrashResource[] | null> => {
    if (!canRestoreDeleteForSpace(space, userStore.user)) {
      return null
    }

    let trashResources: TrashResource[]
    try {
      const { children } = await clientService.webdav.listFiles(
        space,
        {},
        { davProperties: DavProperties.Trashbin, isTrash: true }
      )
      trashResources = children as TrashResource[]
    } catch {
      return null
    }

    const resolved: TrashResource[] = []
    for (const deleted of deletedResources) {
      const candidates = trashResources.filter(
        (t) => t.name === deleted.name && t.path === deleted.path
      )
      if (!candidates.length) {
        continue
      }
      // pick the most recently deleted match
      candidates.sort((a, b) => new Date(b.ddate).getTime() - new Date(a.ddate).getTime())
      resolved.push(candidates[0])
    }

    if (resolved.length !== deletedResources.length) {
      return null
    }

    return resolved
  }

  return { resolveRestorableResources }
}
