<template>
  <oc-list class="app-actions">
    <action-menu-item
      v-for="action in actions"
      :key="`app-action-${action.name}`"
      size="small"
      :action="action"
      :action-options="{ app, version }"
    />
  </oc-list>
</template>
<script lang="ts" setup>
import { ActionMenuItem } from '@ownclouders/web-pkg'
import { useAppActionsDownload } from '../composables'
import { computed } from 'vue'
import { App, AppVersion } from '../types'

interface Props {
  app?: App
  version?: AppVersion | null
}
const { app = undefined, version = null } = defineProps<Props>()
const { downloadAppAction } = useAppActionsDownload()
const actions = computed(() => {
  return [downloadAppAction]
})
</script>

<style lang="scss">
.app-actions {
  display: flex;
  justify-content: flex-start;
  gap: 1rem;
}
</style>
