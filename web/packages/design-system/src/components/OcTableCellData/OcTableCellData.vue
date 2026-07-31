<template>
  <oc-table-cell
    type="td"
    :align-h="alignH"
    :align-v="alignV"
    :width="width"
    :wrap="wrap"
    class="oc-td"
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
defineOptions({
  name: 'OcTd',
  status: 'ready',
  release: '2.1.0'
})

const { alignH = 'left', alignV = 'middle', width = 'auto', wrap = null } = defineProps<Props>()
</script>
