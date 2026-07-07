import { computed, unref } from 'vue'
import { queryItemAsString, useMessages, useModals, useRouteQuery } from '@ownclouders/web-pkg'
import { useClientService } from '@ownclouders/web-pkg'
import { GroupAction, GroupActionOptions } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { Group } from '@ownclouders/web-client/graph/generated'
import { useGroupSettingsStore } from '../../stores'

export const useGroupActionsDelete = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const { $gettext, $ngettext } = useGettext()
  const clientService = useClientService()
  const { dispatchModal } = useModals()
  const groupSettingsStore = useGroupSettingsStore()

  const currentPageQuery = useRouteQuery('page', '1')
  const currentPage = computed(() => {
    return parseInt(queryItemAsString(unref(currentPageQuery)))
  })

  const itemsPerPageQuery = useRouteQuery('items-per-page', '1')
  const itemsPerPage = computed(() => {
    return parseInt(queryItemAsString(unref(itemsPerPageQuery)))
  })

  const deleteGroups = async (groups: Group[]) => {
    const graphClient = clientService.graphAuthenticated
    const promises = groups.map((group) => graphClient.groups.deleteGroup(group.id))
    const results = await Promise.allSettled(promises)

    const succeeded = results.filter((r) => r.status === 'fulfilled')
    if (succeeded.length) {
      const title =
        succeeded.length === 1 && groups.length === 1
          ? $gettext('Group "%{group}" was deleted successfully', { group: groups[0].displayName })
          : $ngettext(
              '%{groupCount} group was deleted successfully',
              '%{groupCount} groups were deleted successfully',
              succeeded.length,
              { groupCount: succeeded.length.toString() },
              true
            )
      showMessage({ title })
    }

    const failed = results.filter((r) => r.status === 'rejected')
    if (failed.length) {
      failed.forEach(console.error)

      const title =
        failed.length === 1 && groups.length === 1
          ? $gettext('Failed to delete group "%{group}"', { group: groups[0].displayName })
          : $ngettext(
              'Failed to delete %{groupCount} group',
              'Failed to delete %{groupCount} groups',
              failed.length,
              { groupCount: failed.length.toString() },
              true
            )
      showErrorMessage({
        title,
        errors: (failed as PromiseRejectedResult[]).map((f) => f.reason)
      })
    }

    groupSettingsStore.removeGroups(groups)
    groupSettingsStore.setSelectedGroups([])

    const pageCount = Math.ceil(groupSettingsStore.groups.length / unref(itemsPerPage))
    if (unref(currentPage) > 1 && unref(currentPage) > pageCount) {
      // reset pagination to avoid empty lists (happens when deleting all items on the last page)
      currentPageQuery.value = pageCount.toString()
    }
  }

  const handler = ({ resources }: GroupActionOptions) => {
    if (!resources.length) {
      return
    }

    dispatchModal({
      variation: 'danger',
      title: $ngettext(
        'Delete group "%{group}"?',
        'Delete %{groupCount} groups?',
        resources.length,
        {
          group: resources[0].displayName,
          groupCount: resources.length.toString()
        }
      ),
      confirmText: $gettext('Delete'),
      message: $ngettext(
        'Are you sure you want to delete this group?',
        'Are you sure you want to delete the %{groupCount} selected groups?',
        resources.length,
        {
          groupCount: resources.length.toString()
        }
      ),
      hasInput: false,
      onConfirm: () => deleteGroups(resources)
    })
  }

  const actions = computed((): GroupAction[] => [
    {
      name: 'delete',
      icon: 'delete-bin',
      label: () => {
        return $gettext('Delete')
      },
      handler,
      isVisible: ({ resources }) => {
        return !!resources.length && !resources.some((r) => r.groupTypes?.includes('ReadOnly'))
      },
      class: 'oc-groups-actions-delete-trigger'
    }
  ])

  return {
    actions,
    // HACK: exported for unit tests:
    deleteGroups
  }
}
