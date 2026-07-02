<template>
  <oc-text-input
    v-model="currentPassword"
    :label="$gettext('Current password')"
    type="password"
    :fix-message-line="true"
  />
  <oc-text-input
    v-model="newPassword"
    :label="$gettext('New password')"
    type="password"
    :fix-message-line="true"
    :error-message="newPasswordErrorMessage"
    @change="validatePasswordConfirm"
  />
  <oc-text-input
    v-model="newPasswordConfirm"
    :label="$gettext('Repeat new password')"
    type="password"
    :fix-message-line="true"
    :error-message="passwordConfirmErrorMessage"
    @change="validatePasswordConfirm"
  />
</template>

<script lang="ts">
import { computed, defineComponent, ref, PropType, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import { Modal, useClientService, useMessages } from '@ownclouders/web-pkg'

export default defineComponent({
  name: 'EditPasswordModal',
  props: { modal: { type: Object as PropType<Modal>, required: true } },
  emits: ['update:confirmDisabled'],
  setup(props, { emit, expose }) {
    const { showMessage, showErrorMessage } = useMessages()
    const clientService = useClientService()
    const { $gettext } = useGettext()

    const currentPassword = ref('')
    const newPassword = ref('')
    const newPasswordConfirm = ref('')
    const passwordConfirmErrorMessage = ref('')
    const newPasswordErrorMessage = ref('')

    const confirmButtonDisabled = computed(() => {
      return (
        !unref(currentPassword).trim().length ||
        !unref(newPassword).trim().length ||
        unref(newPassword).trim() !== unref(newPasswordConfirm).trim() ||
        unref(currentPassword).trim() === unref(newPassword).trim()
      )
    })

    watch(
      confirmButtonDisabled,
      () => {
        emit('update:confirmDisabled', unref(confirmButtonDisabled))
      },
      { immediate: true }
    )

    const onConfirm = () => {
      return clientService.graphAuthenticated.users
        .changeOwnPassword({
          currentPassword: unref(currentPassword).trim(),
          newPassword: unref(newPassword).trim()
        })
        .then(() => {
          showMessage({ title: $gettext('Password was changed successfully') })
        })
        .catch((error) => {
          console.error(error)
          showErrorMessage({
            title: $gettext('Failed to change password'),
            errors: [error]
          })
        })
    }

    expose({ onConfirm })

    watch([currentPassword, newPassword], () => {
      newPasswordErrorMessage.value = ''
      if (!unref(currentPassword).trim().length || !unref(newPassword).trim().length) {
        return
      }
      if (unref(currentPassword).trim() != unref(newPassword).trim()) {
        return
      }
      newPasswordErrorMessage.value = $gettext(
        'New password must be different from current password'
      )
    })

    return {
      currentPassword,
      newPassword,
      newPasswordConfirm,
      passwordConfirmErrorMessage,
      newPasswordErrorMessage,

      // unit tests
      confirmButtonDisabled
    }
  },

  methods: {
    validatePasswordConfirm() {
      this.passwordConfirmErrorMessage = ''

      if (
        this.newPassword.trim().length &&
        this.newPasswordConfirm.trim().length &&
        this.newPassword !== this.newPasswordConfirm
      ) {
        this.passwordConfirmErrorMessage = this.$gettext(
          'Password and password confirmation must be identical'
        )
        return false
      }

      return true
    }
  }
})
</script>
