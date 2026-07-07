<template>
  <div>
    <oc-select
      :model-value="selectedOption"
      :label="$gettext('Login')"
      :options="options"
      :placeholder="$gettext('Select...')"
      :warning-message="
        currentUserSelected ? $gettext('Your own login status will remain unchanged.') : ''
      "
      :position-fixed="true"
      @update:model-value="changeSelectedOption"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import { User } from '@ownclouders/web-client/graph/generated'
import { useClientService, useUserStore, Modal, useMessages } from '@ownclouders/web-pkg'
import { useUserSettingsStore } from '../../composables/stores/userSettings'

type LoginOption = {
  label: string
  value: boolean
}

interface Props {
  modal: Modal
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
const userStore = useUserStore()
const userSettingsStore = useUserSettingsStore()

const selectedOption = ref<LoginOption>()
const options = ref([
  { label: $gettext('Allowed'), value: true },
  { label: $gettext('Forbidden'), value: false }
])

watch(
  selectedOption,
  () => {
    emit('update:confirmDisabled', !unref(selectedOption))
  },
  { immediate: true }
)

const changeSelectedOption = (option: LoginOption) => {
  selectedOption.value = option
}

const currentUserSelected = computed(() => {
  return props.users.some((u) => u.id === userStore.user.id)
})

onMounted(() => {
  if (props.users.every((u) => u.accountEnabled !== false)) {
    selectedOption.value = unref(options).find(({ value }) => !!value)
  } else if (props.users.every((u) => u.accountEnabled === false)) {
    selectedOption.value = unref(options).find(({ value }) => !value)
  }
})

const onConfirm = async () => {
  const affectedUsers = props.users.filter(({ id }) => userStore.user.id !== id)
  const client = clientService.graphAuthenticated
  const promises = affectedUsers.map(({ id }) =>
    client.users.editUser(id, { accountEnabled: unref(selectedOption).value } as User)
  )
  const results = await Promise.allSettled(promises)

  function isFulfilled<T>(result: PromiseSettledResult<T>): result is PromiseFulfilledResult<T> {
    return result.status === 'fulfilled'
  }

  const succeeded = results.filter(isFulfilled)
  if (succeeded.length) {
    const title =
      succeeded.length === 1 && affectedUsers.length === 1
        ? $gettext('Login for user "%{user}" was edited successfully', {
            user: affectedUsers[0].displayName
          })
        : $ngettext(
            '%{userCount} user login was edited successfully',
            '%{userCount} users logins edited successfully',
            succeeded.length,
            { userCount: succeeded.length.toString() },
            true
          )
    showMessage({ title })
  }

  function isRejected<T>(result: PromiseSettledResult<T>): result is PromiseRejectedResult {
    return result.status === 'rejected'
  }

  const failed = results.filter(isRejected)
  if (failed.length) {
    failed.forEach(console.error)

    const title =
      failed.length === 1 && affectedUsers.length === 1
        ? $gettext('Failed edit login for user "%{user}"', {
            user: affectedUsers[0].displayName
          })
        : $ngettext(
            'Failed to edit %{userCount} user login',
            'Failed to edit %{userCount} user logins',
            failed.length,
            { userCount: failed.length.toString() },
            true
          )
    showErrorMessage({
      title,
      errors: (failed as PromiseRejectedResult[]).map((f) => f.reason)
    })
  }

  try {
    const usersResponse = await Promise.all(
      succeeded.map(({ value }) => {
        return client.users.getUser(value.id)
      })
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
