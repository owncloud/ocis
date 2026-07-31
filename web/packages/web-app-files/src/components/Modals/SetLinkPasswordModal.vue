<template>
  <oc-text-input
    :model-value="password"
    :label="$gettext('Password')"
    type="password"
    :password-policy="inputPasswordPolicy"
    :generate-password-method="inputGeneratePasswordMethod"
    :placeholder="link.hasPassword ? '●●●●●●●●' : null"
    :error-message="errorMessage"
    class="oc-modal-body-input"
    @password-challenge-completed="$emit('update:confirmDisabled', false)"
    @password-challenge-failed="$emit('update:confirmDisabled', true)"
    @keydown.enter.prevent="$emit('confirm')"
    @update:model-value="onInput"
  />
</template>

<script lang="ts" setup>
import { ref, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { upperFirst } from 'lodash-es'
import {
  useClientService,
  useMessages,
  usePasswordPolicyService,
  useSharesStore
} from '@ownclouders/web-pkg'
import { LinkShare, Resource, SpaceResource } from '@ownclouders/web-client'

interface Props {
  space: SpaceResource
  resource: Resource
  link: LinkShare
  callbackFn?: () => void
}

interface Emits {
  (e: 'confirm'): void
  (e: 'update:confirmDisabled', value: boolean): void
}
const { space, resource, link, callbackFn = undefined } = defineProps<Props>()

defineEmits<Emits>()

const { showMessage, showErrorMessage } = useMessages()
const clientService = useClientService()
const passwordPolicyService = usePasswordPolicyService()
const { $gettext } = useGettext()
const { updateLink } = useSharesStore()

const password = ref('')
const errorMessage = ref<string>()

const onInput = (value: string) => {
  password.value = value
  errorMessage.value = undefined
}

const onConfirm = async () => {
  try {
    await updateLink({
      clientService,
      space: space,
      resource: resource,
      linkShare: link,
      options: { password: unref(password) }
    })
    if (callbackFn) {
      callbackFn()
      return
    }
    showMessage({ title: $gettext('Link was updated successfully') })
  } catch (e) {
    // Human-readable error message is provided, for example when password is on banned list
    if (e.response?.status === 400) {
      const errorMsg = e.response.data.error.message
      errorMessage.value = $gettext(upperFirst(errorMsg))
      return Promise.reject()
    }

    showErrorMessage({
      title: $gettext('Failed to update link'),
      errors: [e]
    })
  }
}

const inputPasswordPolicy = passwordPolicyService.getPolicy({ enforcePassword: true })
const inputGeneratePasswordMethod = () => passwordPolicyService.generatePassword()

/*
 * onConfirm is called by modalWrapper component and should be exposed.
 */
defineExpose({ onConfirm })
</script>
