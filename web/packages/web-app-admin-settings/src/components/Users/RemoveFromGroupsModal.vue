<template>
  <group-select
    :selected-groups="selectedOptions"
    :group-options="groups"
    :position-fixed="true"
    @selected-option-change="changeSelectedGroupOption"
  />
</template>

<script lang="ts" setup>
import { Ref, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import { Group, User } from '@ownclouders/web-client/graph/generated'
import GroupSelect from './GroupSelect.vue'
import { useClientService, Modal, useMessages } from '@ownclouders/web-pkg'
import { useUserSettingsStore } from '../../composables/stores/userSettings'

interface Props {
  modal: Modal
  groups: Group[]
  users: User[]
}

interface Emits {
  (e: 'update:confirmDisabled', value: boolean): void
}
const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { showMessage, showErrorMessage } = useMessages()
const clientService = useClientService()
const { $gettext, $ngettext } = useGettext()
const userSettingsStore = useUserSettingsStore()

const selectedOptions: Ref<Group[]> = ref([])
const changeSelectedGroupOption = (options: Group[]) => {
  selectedOptions.value = options
}

watch(
  selectedOptions,
  () => {
    emit('update:confirmDisabled', !unref(selectedOptions).length)
  },
  { immediate: true }
)

const onConfirm = async () => {
  const client = clientService.graphAuthenticated
  const usersToFetch: string[] = []
  const promises = unref(selectedOptions).reduce((acc, group) => {
    for (const user of props.users) {
      if (user.memberOf.find((userGroup) => userGroup.id === group.id)) {
        acc.push(client.groups.deleteMember(group.id, user.id))
        if (!usersToFetch.includes(user.id)) {
          usersToFetch.push(user.id)
        }
      }
    }
    return acc
  }, [])

  if (!promises.length) {
    const title = $ngettext(
      'Group assignment already removed',
      'Group assignments already removed',
      props.users.length * unref(selectedOptions).length
    )
    showMessage({ title })
    return
  }

  const results = await Promise.allSettled(promises)

  const succeeded = results.filter((r) => r.status === 'fulfilled')
  if (succeeded.length) {
    const title =
      succeeded.length === 1 && unref(selectedOptions).length === 1 && props.users.length === 1
        ? $gettext('Group assignment "%{group}" was deleted successfully', {
            group: unref(selectedOptions)[0].displayName
          })
        : $ngettext(
            '%{groupAssignmentCount} group assignment was deleted successfully',
            '%{groupAssignmentCount} group assignments were deleted successfully',
            succeeded.length,
            { groupAssignmentCount: succeeded.length.toString() },
            true
          )
    showMessage({ title })
  }

  const failed = results.filter((r) => r.status === 'rejected')
  if (failed.length) {
    failed.forEach(console.error)

    const title =
      failed.length === 1 && unref(selectedOptions).length === 1 && props.users.length === 1
        ? $gettext('Failed to delete group assignment "%{group}"', {
            group: unref(selectedOptions)[0].displayName
          })
        : $ngettext(
            'Failed to delete %{groupAssignmentCount} group assignment',
            'Failed to delete %{groupAssignmentCount} group assignments',
            failed.length,
            { groupAssignmentCount: failed.length.toString() },
            true
          )
    showErrorMessage({
      title,
      errors: (failed as PromiseRejectedResult[]).map((f) => f.reason)
    })
  }

  try {
    const usersResponse = await Promise.all(
      usersToFetch.map((userId) => client.users.getUser(userId))
    )

    usersResponse.forEach((user) => {
      userSettingsStore.upsertUser(user)
    })
  } catch (e) {
    console.error(e)
  }
}

defineExpose({ onConfirm })
</script>
