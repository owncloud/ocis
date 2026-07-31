<template>
  <oc-button
    :id="`context-menu-trigger-${resourceDomSelector(item)}`"
    v-oc-tooltip="contextMenuLabel"
    :aria-label="contextMenuLabel"
    appearance="raw"
    @click.stop.prevent="
      $emit('quickActionClicked', {
        event: $event,
        dropdown: ocDropRef
      })
    "
  >
    <oc-icon name="more-2" />
    <oc-drop
      :ref="`context-menu-drop-ref-${resourceDomSelector(item)}`"
      :drop-id="`context-menu-drop-${resourceDomSelector(item)}`"
      :toggle="`#context-menu-trigger-${resourceDomSelector(item)}`"
      class="oc-overflow-hidden"
      position="bottom-end"
      mode="click"
      close-on-click
      focus-on-open
      padding-size="small"
    >
      <slot name="contextMenu" :item="item" />
    </oc-drop>
  </oc-button>
</template>

<script lang="ts" setup>
import { computed, useTemplateRef, ComponentPublicInstance } from 'vue'
import { useGettext } from 'vue3-gettext'
import { Item, extractDomSelector } from '@ownclouders/web-client'
import { ContextMenuBtnClickEventData, type OcDropType } from '../../helpers'

interface Props {
  item: Item
  resourceDomSelector?: (resource: Item) => string
}
interface Emits {
  (event: 'quickActionClicked', payload: ContextMenuBtnClickEventData): void
}
const { item, resourceDomSelector = (resource: Item) => extractDomSelector(resource.id) } =
  defineProps<Props>()

defineEmits<Emits>()

const ocDropRef = useTemplateRef<ComponentPublicInstance<OcDropType>>(
  `context-menu-drop-ref-${resourceDomSelector(item)}`
)
const { $gettext } = useGettext()
const contextMenuLabel = computed(() => $gettext('Show context menu'))
</script>
