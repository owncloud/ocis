<template>
  <component :is="type" :class="cellClasses" @click="emit('click', $event)">
    <slot />
  </component>
</template>
<script lang="ts" setup>
import { computed } from 'vue'

interface Props {
  type?: 'td' | 'th'
  alignH?: 'left' | 'center' | 'right' | string
  alignV?: 'top' | 'middle' | 'bottom' | string
  width?: 'auto' | 'shrink' | 'expand' | string
  wrap?: 'break' | 'nowrap' | 'truncate' | string
}
interface Emits {
  (e: 'click', event: MouseEvent): void
}

const {
  type = 'td',
  alignH = 'left',
  alignV = 'middle',
  width = 'auto',
  wrap = null
} = defineProps<Props>()
defineOptions({
  name: 'OcTableCell',
  status: 'ready',
  release: '2.1.0'
})
const emit = defineEmits<Emits>()
const cellClasses = computed(() => {
  const classes = [
    'oc-table-cell',
    `oc-table-cell-align-${alignH}`,
    `oc-table-cell-align-${alignV}`,
    `oc-table-cell-width-${width}`
  ]
  if (wrap) {
    classes.push(`oc-text-${wrap}`)
  }
  return classes
})
</script>
<style lang="scss">
.oc-table-cell {
  /* padding is not configurable until we need it */
  padding: 0 var(--oc-space-small);
  position: relative;

  &-align {
    &-left {
      text-align: left;
    }

    &-center {
      text-align: center;
    }

    &-right {
      text-align: right;
    }

    &-top {
      vertical-align: top;
    }

    &-middle {
      vertical-align: middle;
    }

    &-bottom {
      vertical-align: bottom;
    }
  }

  &-width {
    &-shrink {
      width: 1px;
    }

    &-expand {
      min-width: 150px;
    }
  }
}
</style>
