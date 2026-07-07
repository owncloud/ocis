<template>
  <span
    class="oc-avatar-count"
    :style="{ width: size + 'px', height: size + 'px', fontSize: fontSize }"
    v-text="`+${count}`"
  />
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'

/**
 * OcAvatarCount - A component for displaying a count badge, typically used to indicate the number of additional users or items.
 *
 * @prop {number} count - The count to display. This is required and will be prefixed with a "+" (e.g., "+21").
 * @prop {number} [size=30] - The size of the badge in pixels. This determines both the width/height and the font size.
 *
 * @example
 * ```vue
 * <!-- Default size with a count -->
 * <oc-avatar-count :count="21" />
 *
 * <!-- Custom size -->
 * <oc-avatar-count :count="99" :size="50" />
 * ```
 */

interface Props {
  count: number
  size?: number
}

defineOptions({
  name: 'OcAvatarCount',
  status: 'ready',
  release: '2.1.0'
})

const { count, size = 30 } = defineProps<Props>()

const fontSize = computed(() => {
  return Math.floor(unref(size) / 2.5) + 'px'
})
</script>

<style lang="scss">
.oc-avatar-count {
  align-items: center;
  background-color: var(--oc-color-swatch-passive-default);
  border-radius: 50%;
  color: var(--oc-color-text-inverse);
  display: flex;
  justify-content: center;
}
</style>

<docs>
```js
<oc-avatar-count :count="21" />
```
</docs>
