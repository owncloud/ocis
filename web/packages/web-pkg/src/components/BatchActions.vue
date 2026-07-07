<template>
  <div>
    <oc-list
      id="oc-appbar-batch-actions"
      :class="{ 'oc-appbar-batch-actions-squashed': limitedScreenSpace }"
    >
      <action-menu-item
        v-for="(action, index) in actions"
        :key="`action-${index}`"
        :action="action"
        :action-options="actionOptions"
        appearance="raw"
        class="batch-actions oc-mr-s"
        :shortcut-hint="false"
        :show-tooltip="limitedScreenSpace"
        :has-limited-screen-space="limitedScreenSpace"
      />
    </oc-list>
  </div>
</template>

<script lang="ts" setup>
import ActionMenuItem from './ContextActions/ActionMenuItem.vue'
import { Action, ActionOptions } from '../composables'

interface Props {
  actions: Action[]
  actionOptions: ActionOptions
  limitedScreenSpace?: boolean
}

const { actions, actionOptions, limitedScreenSpace = false } = defineProps<Props>()
</script>

<style lang="scss">
#oc-appbar-batch-actions {
  display: block;

  .action-menu-item {
    padding-left: var(--oc-space-small) !important;
    padding-right: var(--oc-space-small) !important;
    gap: var(--oc-space-xsmall) !important;
  }
  .action-menu-item:hover:not([disabled]),
  .action-menu-item:focus:not([disabled]) {
    background-color: var(--oc-color-background-hover);
    border-color: var(--oc-color-background-hover);
  }

  li {
    float: left !important;
  }
  @media only screen and (min-width: 1200px) {
    align-items: center;
    display: flex;

    li {
      margin-top: 0;
      margin-bottom: 0;
    }
  }
}
.oc-appbar-batch-actions-squashed .oc-files-context-action-label {
  display: none;
}
</style>
