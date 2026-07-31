<template>
  <form autocomplete="off" @submit.prevent="emit('confirm')">
    <oc-text-input
      id="input-folder-name"
      v-model="formData.folderName"
      :label="$gettext('Folder name')"
      :error-message="folderNameError"
    />
    <oc-text-input
      id="input-folder-password"
      v-model="formData.password"
      type="password"
      :label="$gettext('Password')"
      class="oc-mt-s"
      :password-policy="passwordPolicy"
      :generate-password-method="() => passwordPolicyService.generatePassword()"
    />

    <div class="oc-flex oc-flex-middle oc-mt-m">
      <oc-icon class="oc-mr-s" :name="selectedTypeIcon" fill-type="line" />
      <link-role-dropdown
        v-model="formData.selectedType"
        :available-link-type-options="availableLinkTypes"
      />
    </div>

    <input type="submit" class="oc-hidden" />
  </form>
</template>

<script lang="ts" setup>
import {
  LinkRoleDropdown,
  useLinkTypes,
  useMessages,
  usePasswordPolicyService,
  useResourcesStore,
  useSpacesStore
} from '@ownclouders/web-pkg'
import { computed, reactive, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import { DavHttpError } from '../../../web-client/src'
import { useCreateFileHandler } from '../composables/useCreateFileHandler'

const emit = defineEmits<{
  confirm: []
  'update:confirmDisabled': [isDisabled: boolean]
}>()

const { $gettext } = useGettext()
const { showErrorMessage } = useMessages()
const { createFileHandler } = useCreateFileHandler()
const { currentFolder } = useResourcesStore()
const { spaces, currentSpace } = useSpacesStore()
const { defaultLinkType, getAvailableLinkTypes, getLinkRoleByType } = useLinkTypes()
const passwordPolicyService = usePasswordPolicyService()

const formData = reactive({
  folderName: '',
  password: '',
  selectedType: unref(defaultLinkType)
})

const folderNameError = ref('')

const isFormValid = computed(() => {
  return formData.folderName !== '' && passwordPolicy.check(formData.password)
})
const availableLinkTypes = computed(() => getAvailableLinkTypes({ isFolder: true }))
const selectedTypeIcon = computed(() => getLinkRoleByType(formData.selectedType).icon)

const passwordPolicy = passwordPolicyService.getPolicy({
  enforcePassword: true
})

const onConfirm = async () => {
  if (!unref(isFormValid)) {
    return Promise.reject()
  }

  try {
    folderNameError.value = ''

    const personalSpace = unref(spaces).find((space) => space.driveType === 'personal')

    if (!personalSpace) {
      throw new Error('Could not find personal space')
    }

    await createFileHandler({
      fileName: formData.folderName,
      currentFolder: unref(currentFolder),
      personalSpace: personalSpace,
      currentSpace: unref(currentSpace),
      password: formData.password,
      type: formData.selectedType
    })
  } catch (error) {
    if (error instanceof DavHttpError && error.statusCode === 405) {
      folderNameError.value = $gettext('Folder already exists')
      return Promise.reject()
    }

    console.error(error)
    showErrorMessage({ title: $gettext('Failed to create folder'), errors: [error] })
    return Promise.reject()
  }
}

watch(
  isFormValid,
  (isValid) => {
    emit('update:confirmDisabled', !isValid)
  },
  { immediate: true }
)

defineExpose({ onConfirm })
</script>
