<template>
  <oc-button id="context-menu-trigger-whitespace" aria-hidden appearance="raw">
    <label class="oc-invisible-sr" for="context-menu-trigger-whitespace">
      {{ $pgettext('Accessibility label', 'Show context menu') }}
    </label>
    <oc-drop
      drop-id="context-menu-drop-whitespace"
      toggle="#context-menu-trigger-whitespace"
      position="bottom-end"
      mode="click"
      class="whitespace-context-actions-list"
      close-on-click
      padding-size="small"
    >
      <oc-list>
        <action-menu-item
          v-for="(action, actionIndex) in menuItemsActions"
          :key="`section-${action.name}-action-${actionIndex}`"
          :action="action"
          :action-options="actionOptions"
          class="oc-px-s oc-rounded oc-menu-item-hover"
          :data-testid="`whitespace-context-menu-item-${action.name}`"
        />
      </oc-list>
    </oc-drop>
  </oc-button>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import {
  useFileActionsPaste,
  useFileActionsShowDetails,
  useResourcesStore
} from '@ownclouders/web-pkg'
import { useFileActionsCreateNewFolder } from '@ownclouders/web-pkg'
import { SpaceResource } from '@ownclouders/web-client'
import { ActionMenuItem } from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'

interface Props {
  space?: SpaceResource
}
const { space = null } = defineProps<Props>()
const resourcesStore = useResourcesStore()
const { currentFolder } = storeToRefs(resourcesStore)

const currentSpace = computed(() => space)
// const contextMenuLabel = computed(() => $gettext('Show context menu'))
const actionOptions = computed(() => ({
  space: unref(currentSpace),
  resources: [currentFolder.value]
}))
const { actions: createNewFolderAction } = useFileActionsCreateNewFolder({ space: currentSpace })
const { actions: showDetailsAction } = useFileActionsShowDetails()
const { actions: pasteAction } = useFileActionsPaste()

const menuItemsActions = computed(() => {
  return [
    ...unref(createNewFolderAction),
    ...unref(pasteAction),
    ...unref(showDetailsAction)
  ].filter((item) => item.isVisible(unref(actionOptions)))
})
</script>

<style lang="scss">
#context-menu-trigger-whitespace {
  visibility: hidden;
  width: 0;
  height: 0;
}
.whitespace-context-actions-list {
  text-align: left;
  white-space: normal;

  .oc-card {
    padding-left: 0px !important;
    padding-right: 0px !important;
  }

  a,
  button,
  span {
    display: inline-flex;
    font-weight: normal !important;
    gap: 10px;
    justify-content: flex-start;
    vertical-align: top;
    width: 100%;
    text-align: left;
  }
}
</style>
