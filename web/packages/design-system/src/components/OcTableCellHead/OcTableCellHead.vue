<template>
  <oc-table-cell
    type="th"
    :align-h="alignH"
    :align-v="alignV"
    :width="width"
    :wrap="wrap"
    class="oc-th"
    @click="emit('click', $event)"
  >
    <slot />
  </oc-table-cell>
</template>
<script lang="ts" setup>
import OcTableCell from '../OcTableCell/OcTableCell.vue'

interface Props {
  alignH?: 'left' | 'center' | 'right'
  alignV?: 'top' | 'middle' | 'bottom'
  width?: 'auto' | 'shrink' | 'expand'
  wrap?: 'break' | 'nowrap' | 'truncate'
}
interface Emits {
  (e: 'click', event: MouseEvent): void
}
defineOptions({
  name: 'OcTh',
  status: 'ready',
  release: '2.1.0'
})
const { alignH = 'left', alignV = 'middle', width = 'auto', wrap = 'nowrap' } = defineProps<Props>()
const emit = defineEmits<Emits>()
</script>
<style lang="scss">
.oc-th {
  color: var(--oc-color-text-muted);
  font-weight: inherit;

  &[aria-sort] {
    cursor: pointer;
  }
}
</style>
