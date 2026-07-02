<template>
  <table :class="tableClasses">
    <slot />
  </table>
</template>
<script lang="ts" setup>
import { computed } from 'vue'

/**
 * @component OcTableSimple
 * @description A simple table component with manually built layout
 * @status ready
 * @release 2.1.0
 *
 * @example
 *   <oc-table-simple :hover="true">
 *     <oc-thead>
 *       <oc-tr>
 *         <oc-th>Resource</oc-th>
 *         <oc-th>Last modified</oc-th>
 *       </oc-tr>
 *     </oc-thead>
 *     <oc-tbody>
 *       <oc-tr v-for="item in items" :key="'item-' + item.id">
 *         <oc-td>{{ item.resource }}</oc-td>
 *         <oc-td>{{ item.last_modified }}</oc-td>
 *       </oc-tr>
 *     </oc-tbody>
 *   </oc-table-simple>
 */

interface Props {
  hover?: boolean
}
defineOptions({
  name: 'OcTableSimple',
  status: 'ready',
  release: '2.1.0'
})

const { hover = false } = defineProps<Props>()

const tableClasses = computed(() => {
  const result = ['oc-table-simple']
  if (hover) {
    result.push('oc-table-simple-hover')
  }
  return result
})
</script>
<style lang="scss">
.oc-table-simple {
  border-collapse: collapse;
  border-spacing: 0;
  width: 100%;

  &-hover tr {
    transition: background-color $transition-duration-short ease-in-out;
  }

  tr + tr {
    border-top: 1px solid var(--oc-color-border);
  }

  &-hover tr:hover {
    background-color: var(--oc-color-input-border);
  }
}
</style>
