<template>
  <form autocomplete="off" @submit.prevent="$emit('confirm')">
    <oc-text-input
      id="create-user-input-user-name"
      v-model="user.onPremisesSamAccountName"
      class="oc-mb-s"
      :label="$gettext('User name') + '*'"
      :error-message="formData.userName.errorMessage"
      :fix-message-line="true"
      @update:model-value="validateUserName"
    />
    <oc-text-input
      id="create-user-input-display-name"
      v-model="user.displayName"
      class="oc-mb-s"
      :label="$gettext('First and last name') + '*'"
      :error-message="formData.displayName.errorMessage"
      :fix-message-line="true"
      @update:model-value="validateDisplayName"
    />
    <oc-text-input
      id="create-user-input-email"
      v-model="user.mail"
      class="oc-mb-s"
      :label="$gettext('Email') + '*'"
      :error-message="formData.email.errorMessage"
      type="email"
      :fix-message-line="true"
      @update:model-value="onInputEmail"
      @change="validateEmail"
    />
    <oc-text-input
      id="create-user-input-password"
      v-model="user.passwordProfile.password"
      autocomplete="new-password"
      class="oc-mb-s"
      :label="$gettext('Password') + '*'"
      :error-message="formData.password.errorMessage"
      type="password"
      :fix-message-line="true"
      @update:model-value="validatePassword"
    />
    <input type="submit" class="oc-hidden" />
  </form>
</template>

<script lang="ts" setup>
import { useGettext } from 'vue3-gettext'
import { computed, ref, unref, watch } from 'vue'
import * as EmailValidator from 'email-validator'
import { useClientService, useMessages } from '@ownclouders/web-pkg'
import { useUserSettingsStore } from '../../composables/stores/userSettings'

interface Emits {
  (e: 'confirm'): void
  (e: 'update:confirmDisabled', value: boolean): void
}

const emit = defineEmits<Emits>()
const { showMessage, showErrorMessage } = useMessages()
const clientService = useClientService()
const { $gettext } = useGettext()
const userSettingsStore = useUserSettingsStore()

const formData = ref<Record<string, { errorMessage: string; valid: boolean }>>({
  userName: {
    errorMessage: '',
    valid: false
  },
  displayName: {
    errorMessage: '',
    valid: false
  },
  email: {
    errorMessage: '',
    valid: false
  },
  password: {
    errorMessage: '',
    valid: false
  }
})

const user = ref({
  onPremisesSamAccountName: '',
  displayName: '',
  mail: '',
  passwordProfile: {
    password: ''
  }
})

const isFormInvalid = computed(() => {
  return Object.keys(unref(formData))
    .map((k) => !!unref(formData)[k].valid)
    .includes(false)
})

watch(
  isFormInvalid,
  () => {
    emit('update:confirmDisabled', unref(isFormInvalid))
  },
  { immediate: true }
)

const onConfirm = async () => {
  if (unref(isFormInvalid)) {
    return Promise.reject()
  }

  try {
    const client = clientService.graphAuthenticated
    const { id: createdUserId } = await client.users.createUser(unref(user))
    const createdUser = await client.users.getUser(createdUserId)
    showMessage({ title: $gettext('User was created successfully') })
    userSettingsStore.upsertUser(createdUser)
  } catch (error) {
    console.error(error)
    showErrorMessage({
      title: $gettext('Failed to create user'),
      errors: [error]
    })
  }
}

function onInputEmail() {
  if (!EmailValidator.validate(unref(user).mail)) {
    return
  }

  unref(formData).email.errorMessage = ''
  unref(formData).email.valid = true
}
async function validateUserName() {
  if (unref(user).onPremisesSamAccountName.trim() === '') {
    unref(formData).userName.errorMessage = $gettext('User name cannot be empty')
    unref(formData).userName.valid = false
    return false
  }

  if (unref(user).onPremisesSamAccountName.includes(' ')) {
    unref(formData).userName.errorMessage = $gettext('User name cannot contain white spaces')
    unref(formData).userName.valid = false
    return false
  }

  if (unref(user).onPremisesSamAccountName.length > 255) {
    unref(formData).userName.errorMessage = $gettext('User name cannot exceed 255 characters')
    unref(formData).userName.valid = false
    return false
  }

  // validate username against regex
  // shouldn't contain special characters except . and _
  // shouldn't start with a number
  // matching regex from server side
  const pattern =
    "^[a-zA-Z_][a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]*(@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)*$"
  if (!new RegExp(pattern).test(unref(user).onPremisesSamAccountName)) {
    unref(formData).userName.errorMessage = $gettext('User name cannot contain special characters')
    unref(formData).userName.valid = false
    return false
  }

  if (
    unref(user).onPremisesSamAccountName.length &&
    !isNaN(parseInt(unref(user).onPremisesSamAccountName[0]))
  ) {
    unref(formData).userName.errorMessage = $gettext('User name cannot start with a number')
    unref(formData).userName.valid = false
    return false
  }

  try {
    // Validate username by fetching the user. If the request succeeds, we throw a validation error
    const client = clientService.graphAuthenticated
    await client.users.getUser(unref(user).onPremisesSamAccountName)
    unref(formData).userName.errorMessage = $gettext('User "%{userName}" already exists', {
      userName: unref(user).onPremisesSamAccountName
    })
    unref(formData).userName.valid = false
    return false
  } catch {}

  unref(formData).userName.errorMessage = ''
  unref(formData).userName.valid = true
  return true
}
function validateDisplayName() {
  unref(formData).displayName.valid = false

  if (unref(user).displayName.trim() === '') {
    unref(formData).displayName.errorMessage = $gettext('First and last name cannot be empty')
    return false
  }

  if (unref(user).displayName.length > 255) {
    unref(formData).displayName.errorMessage = $gettext(
      'First and last name cannot exceed 255 characters'
    )
    return false
  }

  unref(formData).displayName.errorMessage = ''
  unref(formData).displayName.valid = true
  return true
}
function validateEmail() {
  unref(formData).email.valid = false

  if (!EmailValidator.validate(unref(user).mail)) {
    unref(formData).email.errorMessage = $gettext('Please enter a valid email')
    return false
  }

  unref(formData).email.errorMessage = ''
  unref(formData).email.valid = true
  return true
}
function validatePassword() {
  unref(formData).password.valid = false

  if (unref(user).passwordProfile.password.trim() === '') {
    unref(formData).password.errorMessage = $gettext('Password cannot be empty')
    return false
  }

  unref(formData).password.errorMessage = ''
  unref(formData).password.valid = true
  return true
}

defineExpose({ onConfirm })
</script>
