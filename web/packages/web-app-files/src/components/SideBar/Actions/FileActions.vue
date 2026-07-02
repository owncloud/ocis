<template>
  <oc-list id="oc-files-actions-sidebar" class="oc-mt-s">
    <action-menu-item
      v-for="(action, index) in actions"
      :key="`action-${index}`"
      :action="action"
      :action-options="{ space, resources }"
      :has-limited-screen-space="true"
      class="oc-rounded"
    />
  </oc-list>
</template>

<script lang="ts" setup>
import { ActionMenuItem } from '@ownclouders/web-pkg'
import { useFileActions } from '@ownclouders/web-pkg'
import { computed, inject, Ref, unref } from 'vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'

const resource = inject<Ref<Resource>>('resource')
const space = inject<Ref<SpaceResource>>('space')
const resources = computed(() => {
  return [unref(resource)]
})
const { getAllAvailableActions } = useFileActions()
const actions = computed(() => {
  return getAllAvailableActions({
    space: unref(space),
    resources: unref(resources)
  })
})
</script>

<style lang="scss">
#oc-files-actions-sidebar {
  > li a,
  > li a:hover {
    color: var(--oc-color-swatch-passive-default);
    display: inline-flex;
    gap: 10px;
    vertical-align: top;
    text-decoration: none;
  }

  > li:hover {
    text-decoration: none !important;
    background-color: var(--oc-color-background-hover);
  }
}
</style>
