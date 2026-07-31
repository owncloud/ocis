import {
  getPermissionsForSpaceMember,
  GraphSharePermission,
  isIncomingShareResource,
  isPublicSpaceResource,
  isTrashResource,
  Resource,
  SpaceResource
} from '@ownclouders/web-client'
import { useCapabilityStore, useUserStore } from '../piniaStores'
import { isShareSpaceResource } from '@ownclouders/web-client'
import { useGetMatchingSpace } from '../spaces'

export const useCanListShares = () => {
  const capabilityStore = useCapabilityStore()
  const { isPersonalSpaceRoot } = useGetMatchingSpace()
  const userStore = useUserStore()

  const canListShares = ({ space, resource }: { space: SpaceResource; resource: Resource }) => {
    if (!capabilityStore.sharingApiEnabled) {
      return false
    }
    if (isPublicSpaceResource(space)) {
      return false
    }
    if (isPersonalSpaceRoot(resource)) {
      return false
    }
    if (isTrashResource(resource)) {
      return false
    }
    if (isIncomingShareResource(resource)) {
      return resource.sharePermissions.includes(GraphSharePermission.readPermissions)
    }
    if (isShareSpaceResource(space)) {
      const permissions = getPermissionsForSpaceMember(space, userStore.user)
      return permissions.includes(GraphSharePermission.readPermissions)
    }
    return true
  }

  return { canListShares }
}
