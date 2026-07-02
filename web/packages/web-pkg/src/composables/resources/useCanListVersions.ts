import { useUserStore } from '../piniaStores'
import {
  IncomingShareResource,
  isSpaceResource,
  isTrashResource,
  Resource,
  SpaceResource
} from '@ownclouders/web-client'

export const useCanListVersions = () => {
  const userStore = useUserStore()

  const canListVersions = ({ space, resource }: { space: SpaceResource; resource: Resource }) => {
    if (resource.type === 'folder') {
      return false
    }
    if (isSpaceResource(resource)) {
      return false
    }
    if (isTrashResource(resource)) {
      return false
    }

    if (
      resource.isReceivedShare() &&
      typeof (resource as IncomingShareResource).canListVersions === 'function'
    ) {
      return (resource as IncomingShareResource).canListVersions()
    }

    return space?.canListVersions({ user: userStore.user })
  }

  return {
    canListVersions
  }
}
