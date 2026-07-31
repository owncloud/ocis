<template>
  <div class="oc-page-size">
    <label
      class="oc-page-size-label"
      :for="selectId"
      data-testid="oc-page-size-label"
      :aria-hidden="true"
      v-text="label"
    />
    <oc-select
      :input-id="selectId"
      class="oc-page-size-select"
      data-testid="oc-page-size-select"
      :model-value="selected"
      :label="label"
      :label-hidden="true"
      :options="options"
      :clearable="false"
      :searchable="false"
      @update:model-value="emitChange"
    />
  </div>
</template>

<script lang="ts" setup>
import { uniqueId } from '../../helpers'
import OcSelect from '../OcSelect/OcSelect.vue'

/**
 * OcPageSize Component
 *
 * A reusable component for selecting a page size from a predefined set of options.
 *
 * @component
 * @name OcPageSize
 * @status ready
 * @release 8.0.0
 *
 * @props
 * @prop {Array<unknown>} options - The list of options to display in the dropdown.
 * @prop {string} label - The label for the dropdown.
 * @prop {string|number} selected - The currently selected value.
 * @prop {string} [selectId] - Optional unique ID for the dropdown. Defaults to a generated unique ID.
 *
 * @emits
 * @event change - Emitted when the selected value changes.
 * @param {string|boolean} value - The new selected value.
 *
 * @example
 * <OcPageSize
 *   :options="[{ value: 10, label: '10' }, { value: 20, label: '20' }]"
 *   label="Page Size"
 *   :selected="10"
 *   @change="handlePageSizeChange"
 * />
 */

interface Props {
  options: unknown[]
  label: string
  selected: string | number
  selectId?: string
}
interface Emits {
  (e: 'change', value: string | boolean): void
}

defineOptions({
  name: 'OcPageSize',
  status: 'ready',
  release: '8.0.0'
})
const { options, label, selected, selectId = uniqueId('oc-page-size-') } = defineProps<Props>()
const emit = defineEmits<Emits>()
function emitChange(value: boolean) {
  emit('change', value)
}
</script>

<style lang="scss">
.oc-page-size {
  align-items: center;
  display: flex;
  gap: var(--oc-space-xsmall);

  &-select,
  &-select .vs__dropdown-menu {
    min-width: var(--oc-size-width-xsmall);
  }
}
</style>
