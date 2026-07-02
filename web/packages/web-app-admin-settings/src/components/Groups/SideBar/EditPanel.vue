<template>
  <div id="group-edit-panel" class="oc-mt-xl">
    <group-info-box :group="group" />
    <form id="group-edit-form" class="oc-background-highlight oc-p-m" autocomplete="off">
      <oc-text-input
        id="displayName-input"
        v-model="editGroup.displayName"
        class="oc-mb-s"
        :label="$gettext('Group name')"
        :error-message="formData.displayName.errorMessage"
        :fix-message-line="true"
        @update:model-value="validateDisplayName"
      />
      <compare-save-dialog
        class="edit-compare-save-dialog oc-mb-l"
        :original-object="group"
        :compare-object="editGroup"
        :confirm-button-disabled="invalidFormData"
        @revert="revertChanges"
        @confirm="onEditGroup(editGroup)"
      ></compare-save-dialog>
    </form>
  </div>
</template>
<script lang="ts" setup>
import { ref, unref, computed, watch } from 'vue'
import { Group } from '@ownclouders/web-client/graph/generated'
import { CompareSaveDialog, eventBus, useMessages } from '@ownclouders/web-pkg'
import { useClientService } from '@ownclouders/web-pkg'
import GroupInfoBox from './GroupInfoBox.vue'
import { useGroupSettingsStore } from '../../../composables'
import { useGettext } from 'vue3-gettext'

interface Props {
  group?: Group
}

const { group = null } = defineProps<Props>()
const clientService = useClientService()
const { showErrorMessage } = useMessages()
const groupSettingsStore = useGroupSettingsStore()
const { $gettext } = useGettext()

const editGroup = ref<Group>({})
const formData = ref({
  displayName: {
    errorMessage: '',
    valid: true
  }
})

const onEditGroup = async (editGroup: Group) => {
  try {
    const client = clientService.graphAuthenticated
    await client.groups.editGroup(editGroup.id, editGroup)
    const updatedGroup = await client.groups.getGroup(editGroup.id)
    groupSettingsStore.upsertGroup(updatedGroup)

    eventBus.publish('sidebar.entity.saved')

    return updatedGroup
  } catch (error) {
    console.error(error)
    showErrorMessage({
      title: $gettext('Failed to edit group'),
      errors: [error]
    })
  }
}
const invalidFormData = computed(() =>
  Object.values(unref(formData))
    .map((v: any) => !!v.valid)
    .includes(false)
)
watch(
  () => group,
  () => (editGroup.value = { ...group }),
  {
    deep: true,
    immediate: true
  }
)
async function validateDisplayName() {
  unref(formData).displayName.valid = false

  if (unref(editGroup).displayName.trim() === '') {
    unref(formData).displayName.errorMessage = $gettext('Group name cannot be empty')
    return false
  }

  if (unref(editGroup).displayName.length > 255) {
    unref(formData).displayName.errorMessage = $gettext('Group name cannot exceed 255 characters')
    return false
  }

  if (group.displayName !== unref(editGroup).displayName) {
    try {
      const client = clientService.graphAuthenticated
      await client.groups.getGroup(unref(editGroup).displayName)
      unref(formData).displayName.errorMessage = $gettext('Group "%{groupName}" already exists', {
        groupName: unref(editGroup).displayName
      })
      return false
    } catch {}
  }

  unref(formData).displayName.errorMessage = ''
  unref(formData).displayName.valid = true
  return true
}

function revertChanges() {
  editGroup.value = { ...group }
  Object.values(unref(formData)).forEach((formDataValue: any) => {
    formDataValue.valid = true
    formDataValue.errorMessage = ''
  })
}
</script>
<style lang="scss">
#group-edit-panel {
  #group-edit-form {
    border-top-left-radius: 5px;
    border-top-right-radius: 5px;
  }

  .edit-compare-save-dialog {
    border-bottom-left-radius: 5px;
    border-bottom-right-radius: 5px;
  }

  .group-info {
    align-items: center;
    flex-direction: column;
  }

  .group-info-display-name {
    font-size: 1.5rem;
  }
}
</style>
