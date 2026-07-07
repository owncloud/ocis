import { computed, unref } from 'vue'
import {
  queryItemAsString,
  useCapabilityStore,
  useMessages,
  useModals,
  useRouteQuery
} from '@ownclouders/web-pkg'
import { useClientService } from '@ownclouders/web-pkg'
import { UserAction, UserActionOptions } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { User } from '@ownclouders/web-client/graph/generated'
import { useUserSettingsStore } from '../../stores/userSettings'

export const useUserActionsDelete = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const capabilityStore = useCapabilityStore()
  const { $gettext, $ngettext } = useGettext()
  const clientService = useClientService()
  const { dispatchModal } = useModals()
  const userSettingsStore = useUserSettingsStore()

  const currentPageQuery = useRouteQuery('page', '1')
  const currentPage = computed(() => {
    return parseInt(queryItemAsString(unref(currentPageQuery)))
  })

  const itemsPerPageQuery = useRouteQuery('items-per-page', '1')
  const itemsPerPage = computed(() => {
    return parseInt(queryItemAsString(unref(itemsPerPageQuery)))
  })

  const deleteUsers = async (users: User[]) => {
    const graphClient = clientService.graphAuthenticated
    const promises = users.map((user) => graphClient.users.deleteUser(user.id))
    const results = await Promise.allSettled(promises)

    const succeeded = results.filter((r) => r.status === 'fulfilled')
    if (succeeded.length) {
      const title =
        succeeded.length === 1 && users.length === 1
          ? $gettext('User "%{user}" was deleted successfully', { user: users[0].displayName })
          : $ngettext(
              '%{userCount} user was deleted successfully',
              '%{userCount} users were deleted successfully',
              succeeded.length,
              { userCount: succeeded.length.toString() },
              true
            )
      showMessage({ title })
    }

    const failed = results.filter((r) => r.status === 'rejected')
    if (failed.length) {
      failed.forEach(console.error)

      const title =
        failed.length === 1 && users.length === 1
          ? $gettext('Failed to delete user "%{user}"', { user: users[0].displayName })
          : $ngettext(
              'Failed to delete %{userCount} user',
              'Failed to delete %{userCount} users',
              failed.length,
              { userCount: failed.length.toString() },
              true
            )
      showErrorMessage({
        title,
        errors: (failed as PromiseRejectedResult[]).map((f) => f.reason)
      })
    }

    userSettingsStore.removeUsers(users)
    userSettingsStore.setSelectedUsers([])

    const pageCount = Math.ceil(userSettingsStore.users.length / unref(itemsPerPage))
    if (unref(currentPage) > 1 && unref(currentPage) > pageCount) {
      // reset pagination to avoid empty lists (happens when deleting all items on the last page)
      currentPageQuery.value = pageCount.toString()
    }
  }

  const handler = ({ resources }: UserActionOptions) => {
    if (!resources.length) {
      return
    }

    dispatchModal({
      variation: 'danger',
      title: $ngettext('Delete user "%{user}"?', 'Delete %{userCount} users?', resources.length, {
        user: resources[0].displayName,
        userCount: resources.length.toString()
      }),
      confirmText: $gettext('Delete'),
      message: $ngettext(
        'Are you sure you want to delete this user?',
        'Are you sure you want to delete the %{userCount} selected users?',
        resources.length,
        {
          userCount: resources.length.toString()
        }
      ),
      hasInput: false,
      onConfirm: () => deleteUsers(resources)
    })
  }

  const actions = computed((): UserAction[] => [
    {
      name: 'delete',
      icon: 'delete-bin',
      label: () => {
        return $gettext('Delete')
      },
      handler,
      isVisible: ({ resources }) => {
        return !!resources.length && !capabilityStore.graphUsersDeleteDisabled
      },
      class: 'oc-users-actions-delete-trigger'
    }
  ])

  return {
    actions,
    // HACK: exported for unit tests:
    deleteUsers
  }
}
