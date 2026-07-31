<template>
  <div>
    <p class="oc-mt-rm oc-mb-rm" v-text="message" />
    <p
      v-if="currentUserSelected"
      class="delete-user-own-account-hint oc-mt-s oc-mb-rm"
      v-text="$gettext('Your own account will not be deleted.')"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { User } from '@ownclouders/web-client/graph/generated'
import { Modal, useUserStore } from '@ownclouders/web-pkg'

interface Props {
  modal: Modal
  users: User[]
}
const props = defineProps<Props>()
const { $gettext, $ngettext } = useGettext()
const userStore = useUserStore()

const currentUserSelected = computed(() =>
  props.users.some((user) => user.id === userStore.user.id)
)

const message = computed(() =>
  $ngettext(
    'Are you sure you want to delete this user?',
    'Are you sure you want to delete the %{userCount} selected users?',
    props.users.length,
    { userCount: props.users.length.toString() }
  )
)
</script>

<style lang="scss" scoped>
.delete-user-own-account-hint {
  color: var(--oc-color-swatch-warning-default);
}
</style>
