import { computed, unref } from 'vue'
import { useAbility } from '../ability'
import { useCapabilityStore } from '../piniaStores'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { useGettext } from 'vue3-gettext'
import { ShareRole } from '@ownclouders/web-client'

export const useLinkTypes = () => {
  const { $gettext } = useGettext()
  const capabilityStore = useCapabilityStore()
  const ability = useAbility()

  const canCreatePublicLinks = computed(() => ability.can('create-all', 'PublicLink'))

  const defaultLinkType = computed<SharingLinkType>(() => SharingLinkType.View)

  const isPasswordEnforcedForLinkType = (type: SharingLinkType) => {
    if (type === SharingLinkType.View) {
      return capabilityStore.sharingPublicPasswordEnforcedFor.read_only
    }
    if (type === SharingLinkType.Upload) {
      return capabilityStore.sharingPublicPasswordEnforcedFor.upload_only
    }
    if (type === SharingLinkType.CreateOnly) {
      return capabilityStore.sharingPublicPasswordEnforcedFor.read_write
    }
    if (type === SharingLinkType.Edit) {
      return capabilityStore.sharingPublicPasswordEnforcedFor.read_write_delete
    }
    return false
  }

  const getAvailableLinkTypes = ({ isFolder }: { isFolder: boolean }): SharingLinkType[] => {
    if (!unref(canCreatePublicLinks)) {
      return []
    }

    if (isFolder) {
      return [SharingLinkType.View, SharingLinkType.Edit, SharingLinkType.CreateOnly]
    }

    return [SharingLinkType.View, SharingLinkType.Edit]
  }

  // links don't have roles in graph API, hence we need to define them here
  const linkShareRoles = [
    {
      id: SharingLinkType.Internal,
      displayName: $gettext('Invited people'),
      description: $gettext('Link works only for invited people. Login is required.'),
      icon: 'user'
    },
    {
      id: SharingLinkType.View,
      displayName: $gettext('Can view'),
      description: $gettext('View, download'),
      icon: 'eye'
    },
    {
      id: SharingLinkType.Upload,
      displayName: $gettext('Can upload'),
      description: $gettext('View, upload, download'),
      icon: 'upload'
    },
    {
      id: SharingLinkType.Edit,
      displayName: $gettext('Can edit'),
      description: $gettext('View, upload, edit, download, delete'),
      icon: 'pencil'
    },
    {
      id: SharingLinkType.CreateOnly,
      displayName: $gettext('Secret File Drop'),
      description: $gettext('Upload only, existing content is not revealed'),
      icon: 'inbox-unarchive'
    }
  ] satisfies ShareRole[]

  const getLinkRoleByType = (type: SharingLinkType): ShareRole => {
    return linkShareRoles.find(({ id }) => id === type)
  }

  return {
    defaultLinkType,
    isPasswordEnforcedForLinkType,
    getAvailableLinkTypes,
    linkShareRoles,
    getLinkRoleByType
  }
}
