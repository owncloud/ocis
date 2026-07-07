<template>
  <div v-if="!isEmbedModeEnabled" class="oc-flex">
    <oc-button
      v-for="action in filteredActions"
      :key="action.label()"
      v-oc-tooltip="action.label()"
      :aria-label="action.label()"
      appearance="raw"
      class="oc-mr-xs quick-action-button oc-p-xs"
      :class="`files-quick-action-${action.name}`"
      @click="action.handler({ space, resources: [item] })"
    >
      <oc-icon :name="action.icon" fill-type="line" />
    </oc-button>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useEmbedMode, useExtensionRegistry } from '@ownclouders/web-pkg'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { unref } from 'vue'
import { quickActionsExtensionPoint } from '../../extensionPoints'

interface Props {
  item: Resource
  space?: SpaceResource
}
const { item, space = undefined } = defineProps<Props>()
const extensionRegistry = useExtensionRegistry()
const { isEnabled: isEmbedModeEnabled } = useEmbedMode()

const filteredActions = computed(() => {
  return unref(extensionRegistry)
    .requestExtensions(quickActionsExtensionPoint)
    .map((e) => e.action)
    .filter(({ isVisible }) => isVisible({ space: space, resources: [item] }))
})
</script>

<style lang="scss">
.quick-action-button {
  &:hover {
    background-color: var(--oc-color-background-secondary) !important;
  }
}
</style>
