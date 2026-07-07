<template>
  <div class="space-quota">
    <p class="oc-mb-s oc-mt-rm" v-text="spaceStorageDetailsLabel" />
    <oc-progress
      :value="quotaUsagePercent"
      :max="100"
      size="small"
      :variation="quotaProgressVariant"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { Quota } from '@ownclouders/web-client/graph/generated'
import { formatFileSize } from '../helpers'
import { useGettext } from 'vue3-gettext'

interface Props {
  spaceQuota: Quota
}

const props = defineProps<Props>()

const { current: currentLanguage, $gettext } = useGettext()
const spaceStorageDetailsLabel = computed(() => {
  if (props.spaceQuota.total) {
    return $gettext('%{used} of %{total} used (%{percentage}% used)', {
      used: quotaUsed.value,
      total: quotaTotal.value,
      percentage: quotaUsagePercent.value.toString()
    })
  }

  return $gettext('%{used} used (no restriction)', {
    used: quotaUsed.value
  })
})
const quotaTotal = computed(() => {
  return formatFileSize(props.spaceQuota.total, currentLanguage)
})
const quotaUsed = computed(() => {
  return formatFileSize(props.spaceQuota.used, currentLanguage)
})
const quotaUsagePercent = computed(() => {
  return parseFloat(((props.spaceQuota.used / props.spaceQuota.total) * 100).toFixed(2))
})
const quotaProgressVariant = computed(() => {
  switch (props.spaceQuota.state) {
    case 'normal':
      return 'primary'
    case 'nearing':
      return 'warning'
    case 'critical':
      return 'warning'
    default:
      return 'danger'
  }
})
</script>
