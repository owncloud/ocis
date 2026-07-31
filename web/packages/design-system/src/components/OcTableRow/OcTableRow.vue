<template>
  <tr
    ref="observerTarget"
    @click="$emit('click', $event)"
    @contextmenu="$emit('contextmenu', $event)"
    @dragstart="$emit('dragstart', $event)"
    @drop="$emit('drop', $event)"
    @dragenter="$emit('dragenter', $event)"
    @dragleave="$emit('dragleave', $event)"
    @dragover="$emit('dragover', $event)"
    @mouseleave="$emit('mouseleave', $event)"
    @blur="$emit('blur', $event)"
  >
    <oc-td v-if="isHidden" :colspan="lazyColspan">
      <span class="shimmer" />
    </oc-td>
    <slot v-else />
  </tr>
</template>
<script lang="ts" setup>
import { customRef, computed, ref, unref } from 'vue'
import { useIsVisible } from '../../composables'
import OcTd from '../OcTableCellData/OcTableCellData.vue'

/**
 * @component OcTableRow
 * @description A table row component (`<tr>`). It supports lazy loading, visibility detection, and emits various events for interaction.
 *
 * @props
 * @prop {Object} [lazy] - Optional lazy loading configuration.
 * @prop {number} lazy.colspan - The number of columns the row should span when lazy loading is active.
 *
 * @emits
 * @event click - Emitted when the row is clicked.
 * @param {MouseEvent} event - The click event object.
 *
 * @event contextmenu - Emitted when the context menu is triggered on the row.
 * @param {MouseEvent} event - The context menu event object.
 *
 * @event dragstart - Emitted when a drag operation starts on the row.
 * @param {DragEvent} event - The dragstart event object.
 *
 * @event drop - Emitted when an item is dropped on the row.
 * @param {DragEvent} event - The drop event object.
 *
 * @event dragenter - Emitted when a dragged item enters the row.
 * @param {DragEvent} event - The dragenter event object.
 *
 * @event dragleave - Emitted when a dragged item leaves the row.
 * @param {DragEvent} event - The dragleave event object.
 *
 * @event dragover - Emitted when a dragged item is over the row.
 * @param {DragEvent} event - The dragover event object.
 *
 * @event mouseleave - Emitted when the mouse leaves the row.
 * @param {MouseEvent} event - The mouseleave event object.
 *
 * @event itemVisible - Emitted when the row becomes visible in the viewport.
 *
 * @example
 *   <oc-table-row
 *     :lazy="{ colspan: 3 }"
 *     @click="handleClick"
 *     @itemVisible="handleVisibility"
 *   >
 *     <td>Row Content</td>
 *   </oc-table-row>
 *
 */

interface Props {
  lazy?: { colspan: number }
}

interface Emits {
  (e: 'contextmenu', event: MouseEvent): void
  (e: 'click', event: MouseEvent): void
  (e: 'dragstart', event: DragEvent): void
  (e: 'drop', event: DragEvent): void
  (e: 'dragenter', event: DragEvent): void
  (e: 'dragleave', event: DragEvent): void
  (e: 'dragover', event: DragEvent): void
  (e: 'mouseleave', event: MouseEvent): void
  (e: 'itemVisible'): void
  (e: 'blur', event: FocusEvent): void
}
defineOptions({
  name: 'OcTr',
  status: 'ready',
  release: '1.0.0'
})
const { lazy = null } = defineProps<Props>()
const emit = defineEmits<Emits>()
const observerTarget = customRef((track, trigger) => {
  let $el: HTMLElement
  return {
    get() {
      track()
      return $el
    },
    set(value) {
      $el = value
      trigger()
    }
  }
})

const lazyColspan = computed(() => {
  return lazy ? lazy.colspan : 1
})

const { isVisible } = lazy
  ? useIsVisible({
      ...lazy,
      target: observerTarget,
      onVisibleCallback: () => emit('itemVisible')
    })
  : { isVisible: ref(true) }

const isHidden = computed(() => !unref(isVisible))

if (!lazy) {
  emit('itemVisible')
}
</script>
<style lang="scss">
.shimmer {
  background-color: var(--oc-color-input-text-muted);
  bottom: 12px;
  display: inline-block;
  left: var(--oc-space-small);
  opacity: 0.2;
  overflow: hidden;
  position: absolute;
  right: var(--oc-space-small);
  top: 12px;

  &::after {
    animation: shimmer 2s infinite;
    background-image: linear-gradient(
      90deg,
      rgba(#fff, 0) 0,
      rgba(#fff, 0.2) 20%,
      rgba(#fff, 0.5) 60%,
      rgba(#fff, 0)
    );
    bottom: 0;
    content: '';
    left: 0;
    position: absolute;
    right: 0;
    top: 0;
    transform: translateX(-100%);
  }

  @keyframes shimmer {
    100% {
      transform: translateX(100%);
    }
  }
}
</style>
