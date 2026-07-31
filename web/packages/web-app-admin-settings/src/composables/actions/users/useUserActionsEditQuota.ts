import { computed, toRaw } from 'vue'
import { useGettext } from 'vue3-gettext'
import {
  QuotaModal,
  useAbility,
  useModals,
  UserAction,
  UserActionOptions,
  useCapabilityStore
} from '@ownclouders/web-pkg'
import { SpaceResource } from '@ownclouders/web-client'
import { isPersonalSpaceResource } from '@ownclouders/web-client'
import { User } from '@ownclouders/web-client/graph/generated'

export const useUserActionsEditQuota = () => {
  const { dispatchModal } = useModals()
  const capabilityStore = useCapabilityStore()
  const { $gettext } = useGettext()
  const ability = useAbility()

  const getModalTitle = ({ resources }: { resources: User[] }) => {
    if (resources.length === 1) {
      return $gettext('Change quota for user "%{name}"', {
        name: resources[0].displayName
      })
    }
    return $gettext('Change quota for %{count} users', {
      count: resources.length.toString()
    })
  }

  const getUserDrives = ({ resources }: { resources: User[] }) => {
    const selectedPersonalDrives: SpaceResource[] = []
    resources.forEach((user) => {
      const drive = toRaw(user.drive)
      if (drive === undefined || drive.id === undefined) {
        return
      }
      const spaceResource = {
        id: drive.id,
        name: user.displayName,
        spaceQuota: drive.quota
      } as SpaceResource
      selectedPersonalDrives.push(spaceResource)
    })
    return selectedPersonalDrives
  }

  const handler = ({ resources }: UserActionOptions) => {
    const usersWithoutDrive = resources.filter(
      ({ drive }) => !isPersonalSpaceResource(drive as SpaceResource)
    )

    dispatchModal({
      title: getModalTitle({ resources }),
      customComponent: QuotaModal,
      customComponentAttrs: () => ({
        spaces: getUserDrives({ resources }),
        resourceType: 'user',
        warningMessage: usersWithoutDrive.length
          ? $gettext('Quota will only be applied to users who logged in at least once.')
          : '',
        warningMessageContextualHelperData: usersWithoutDrive.length
          ? {
              title: $gettext('Unaffected users'),
              text: [...usersWithoutDrive]
                .sort((u1, u2) => u1.displayName.localeCompare(u2.displayName))
                .map((user) => user.displayName)
                .join(', ')
            }
          : {}
      })
    })
  }

  const actions = computed((): UserAction[] => [
    {
      name: 'editQuota',
      icon: 'cloud',
      label: () => {
        return $gettext('Edit quota')
      },
      handler,
      isVisible: ({ resources }) => {
        if (!resources || !resources.length) {
          return false
        }

        if (capabilityStore.graphUsersReadOnlyAttributes.includes('drive.quota')) {
          return false
        }

        if (!resources.some((r) => r.drive?.quota)) {
          return false
        }

        return ability.can('set-quota-all', 'Drive')
      },
      class: 'oc-users-actions-edit-quota-trigger'
    }
  ])

  return {
    actions
  }
}
