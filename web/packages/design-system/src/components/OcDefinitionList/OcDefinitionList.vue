<template>
  <dl class="details-list">
    <!-- eslint-disable-next-line vue/no-v-for-template-key -->
    <template v-for="item in items" :key="item.term">
      <dt>{{ item.term }}</dt>
      <dd>{{ item.definition }}</dd>
    </template>
  </dl>
</template>

<script lang="ts" setup>
/**
 * OcDefinitionList Component
 *
 * A Vue component for rendering a semantic definition list (`<dl>`), which displays key-value pairs.
 * This is useful for presenting metadata, definitions, or other structured data in a readable format.
 *
 * @component
 * @name OcDefinitionList
 * @status ready
 * @release 16.0.0
 *
 * @props {Array<Object>} items - An array of objects representing the key-value pairs to display.
 * @props {string} items[].term - The term or key in the key-value pair.
 * @props {string} items[].definition - The definition or value in the key-value pair.
 *
 * @example
 * <template>
 *   <OcDefinitionList :items="[
 *     { term: 'Name', definition: 'John Doe' },
 *     { term: 'Age', definition: '30' }
 *   ]" />
 * </template>
 *
 */

/**
 * Definition lists are used to visually communicate key-value pairs.
 * They are useful for displaying metadata, definitions, and other data that is best represented in a key-value format.
 * @items - An array of objects with `term` and `definition` properties.
 *
 */

interface DefinitionItem {
  term: string
  definition: string
}
interface Props {
  items: DefinitionItem[]
}

defineOptions({
  name: 'OcDefinitionList',
  status: 'ready',
  release: '16.0.0'
})
const { items } = defineProps<Props>()
</script>

<style lang="scss">
.details-list {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  dt,
  dd {
    margin-bottom: var(--oc-space-small);
  }
  dt {
    font-weight: bold;
    white-space: nowrap;
  }
  dd {
    margin-inline-start: var(--oc-space-medium);
  }
}
</style>
