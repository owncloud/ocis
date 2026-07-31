import { isSpaceResource, Resource, SpaceResource } from '@ownclouders/web-client'
import { useAbility } from '../ability'
import { useCapabilityStore, useUserStore } from '../piniaStores'
import { isProjectSpaceResource, isShareSpaceResource } from '@ownclouders/web-client'

export const useCanShare = () => {
  const capabilityStore = useCapabilityStore()
  const ability = useAbility()
  const userStore = useUserStore()

  const canShare = ({ space, resource }: { space: SpaceResource; resource: Resource }) => {
    if (!capabilityStore.sharingApiEnabled) {
      return false
    }

    if (isShareSpaceResource(space)) {
      return false
    }

    if (isProjectSpaceResource(space) && !space.canShare({ user: userStore.user })) {
      return false
    }

    if (isSpaceResource(resource) && capabilityStore.capabilities.spaces.server_managed) {
      return false
    }

    return resource.canShare({ ability, user: userStore.user })
  }

  return { canShare }
}
