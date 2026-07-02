<template>
  <ul class="oc-list oc-my-rm oc-mx-rm" :class="{ 'oc-list-raw': raw }">
    <slot />
  </ul>
</template>

<script lang="ts" setup>
/**
 * OcList Component
 *
 * Renders a list with optional raw styling.
 *
 * @component
 * @name OcList
 * @status ready
 * @release 12.0.0
 *
 * @prop {boolean} [raw=false] - If true, applies raw styling to the list.
 *
 * @slot default - Slot for list items.
 *
 * @example
 * <OcList :raw="true">
 *   <li>Item 1</li>
 *   <li>Item 2</li>
 * </OcList>
 */

interface Props {
  raw?: boolean
}

defineOptions({
  name: 'OcList',
  status: 'ready',
  release: '12.0.0'
})

const { raw = false } = defineProps<Props>()
</script>

<style lang="scss">
.oc-list {
  list-style: none;
  margin: 0;
  padding: 0;

  &-divider > :nth-child(n + 2) {
    border-top: 1px solid var(--oc-color-border);
    margin-top: var(--oc-space-small);
    padding-top: var(--oc-space-small);
  }
  &-raw {
    list-style-type: none;
    a {
      text-decoration: none !important;
      &:hover {
        color: inherit;
      }
    }
  }
}
</style>
