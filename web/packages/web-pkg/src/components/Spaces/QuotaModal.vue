<template>
  <quota-select
    :total-quota="selectedOption"
    :max-quota="capabilityRefs.spacesMaxQuota.value"
    :position-fixed="true"
    v-bind="attrs"
    @selected-option-change="changeSelectedQuotaOption"
  />
  <div v-if="warningMessage" class="oc-mt-s">
    <span class="oc-text-input-warning" v-text="warningMessage" />
    <oc-contextual-helper
      v-if="warningMessageContextualHelperData"
      class="oc-pl-xs"
      v-bind="warningMessageContextualHelperData"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, unref, ref, onMounted, watch, useAttrs } from 'vue'
import { useGettext } from 'vue3-gettext'
import QuotaSelect from '../QuotaSelect.vue'
import { SpaceResource } from '@ownclouders/web-client'
import {
  Modal,
  useClientService,
  useMessages,
  useSpacesStore,
  useCapabilityStore,
  useResourcesStore,
  useSharesStore
} from '../../composables'
import { useRouter } from '../../composables/router'
import { eventBus } from '../../services'
import { storeToRefs } from 'pinia'
import { ContextualHelperData } from '@ownclouders/design-system/helpers'

interface Props {
  modal: Modal
  spaces: SpaceResource[]
  warningMessage?: string
  warningMessageContextualHelperData?: ContextualHelperData
  resourceType?: 'space' | 'user'
}
interface Emits {
  (e: 'update:confirmDisabled', boolean): void
}
const {
  modal,
  spaces,
  warningMessage = '',
  warningMessageContextualHelperData = undefined,
  resourceType = 'space'
} = defineProps<Props>()

const emit = defineEmits<Emits>()

const { showMessage, showErrorMessage } = useMessages()
const capabilityStore = useCapabilityStore()
const capabilityRefs = storeToRefs(capabilityStore)
const { $gettext, $ngettext } = useGettext()
const clientService = useClientService()
const router = useRouter()
const spacesStore = useSpacesStore()
const sharesStore = useSharesStore()
const { updateResourceField } = useResourcesStore()
const attrs = useAttrs()

const selectedOption = ref(0)

const getSuccessMessage = (count: number) => {
  if (resourceType === 'space') {
    return $ngettext(
      'Space quota was changed successfully',
      'Quota of %{count} spaces was changed successfully',
      count,
      { count: count.toString() }
    )
  }
  if (resourceType === 'user') {
    return $ngettext(
      'User quota was changed successfully',
      'Quota of %{count} users was changed successfully',
      count,
      { count: count.toString() }
    )
  }
  return $gettext('Quota was changed successfully')
}
const getErrorMessage = (count: number) => {
  if (resourceType === 'space') {
    return $ngettext(
      'Failed to change space quota',
      'Failed to change quota for %{count} spaces',
      count,
      { count: count.toString() }
    )
  }
  if (resourceType === 'user') {
    return $ngettext(
      'Failed to change user quota',
      'Failed to change quota for %{count} users',
      count,
      { count: count.toString() }
    )
  }
  return $gettext('Failed to change quota')
}

const confirmButtonDisabled = computed(
  () => !spaces.some((space) => space.spaceQuota.total !== unref(selectedOption))
)

watch(
  confirmButtonDisabled,
  () => {
    emit('update:confirmDisabled', unref(confirmButtonDisabled))
  },
  { immediate: true }
)

const changeSelectedQuotaOption = (option: { value: number }) => {
  selectedOption.value = option.value
}

const onConfirm = async () => {
  const client = clientService.graphAuthenticated
  const requests = spaces.map(async (space): Promise<void> => {
    const updatedSpace = await client.drives.updateDrive(
      space.id,
      { name: space.name, quota: { total: unref(selectedOption) } },
      sharesStore.graphRoles
    )
    if (unref(router.currentRoute).name === 'admin-settings-spaces') {
      eventBus.publish('app.admin-settings.spaces.space.quota.updated', {
        spaceId: space.id,
        quota: updatedSpace.spaceQuota
      })
    }
    if (unref(router.currentRoute).name === 'admin-settings-users') {
      eventBus.publish('app.admin-settings.users.user.quota.updated', {
        spaceId: space.id,
        quota: updatedSpace.spaceQuota
      })
    }
    spacesStore.updateSpaceField({
      id: space.id,
      field: 'spaceQuota',
      value: updatedSpace.spaceQuota
    })
    updateResourceField<SpaceResource>({
      id: space.id,
      field: 'spaceQuota',
      value: updatedSpace.spaceQuota
    })
  })
  const results = await Promise.allSettled<Array<unknown>>(requests)
  const succeeded = results.filter((r) => r.status === 'fulfilled')
  if (succeeded.length) {
    showMessage({ title: getSuccessMessage(succeeded.length) })
  }
  const errors = results.filter((r) => r.status === 'rejected')
  if (errors.length) {
    console.error(errors)
    errors.forEach(console.error)
    showErrorMessage({
      title: getErrorMessage(errors.length),
      errors: (errors as PromiseRejectedResult[]).map((f) => f.reason)
    })
  }
}

defineExpose({ onConfirm })

onMounted(() => {
  selectedOption.value = spaces[0]?.spaceQuota?.total || 0
})
</script>
