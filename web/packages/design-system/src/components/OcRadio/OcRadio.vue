<template>
  <span>
    <input
      :id="id"
      v-model="model"
      type="radio"
      name="radio"
      :class="classes"
      :aria-checked="option === modelValue"
      :value="option"
      :disabled="disabled"
    />
    <label :for="id" :class="labelClasses" v-text="label" />
  </span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { getSizeClass, uniqueId, AvailableSizeType } from '../../helpers'

/**
 * OcRadio Component
 *
 * A radio button component that can be grouped to allow users to select one option from a set.
 *
 * @component
 * @name OcRadio
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {string} [id] - Id for the radio button. If not provided, a unique id will be generated.
 * @prop {boolean} [disabled=false] - Disables the radio button.
 * @prop {unknown} [modelValue=false] - The model value of the radio button or group. Determines if the radio button is checked.
 * @prop {unknown} [option=null] - The value of this radio button. Used when part of a group.
 * @prop {string} [label=null] - Label for the radio button. Required for accessibility.
 * @prop {boolean} [hideLabel=false] - Hides the label visually but keeps it for screen readers.
 * @prop { 'xsmall' | 'small' | 'medium' | 'large' | 'xlarge' | 'xxlarge' | 'xxxlarge' } [size='medium'] - Size of the radio button.
 *
 * @emits
 * @event update:modelValue - Emitted when the model value changes.
 *
 * @example
 *  <oc-radio
 *    v-for="o in availableOptions"
 *    :key="'option-' + o"
 *    v-model="selectedOption"
 *    :option="o"
 *    :label="o"
 *  />
 */

interface Props {
  id?: string
  disabled?: boolean
  modelValue?: unknown
  option?: unknown
  label?: string
  hideLabel?: boolean
  size?: AvailableSizeType
}
interface Emits {
  (e: 'update:modelValue', value: unknown): void
}
defineOptions({
  name: 'OcRadio',
  status: 'ready',
  release: '1.0.0'
})
const {
  id = uniqueId('oc-radio-'),
  disabled = false,
  modelValue = false,
  option = null,
  label = null,
  hideLabel = false,
  size = 'medium'
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const model = computed({
  get() {
    return modelValue
  },
  set(value: unknown) {
    emit('update:modelValue', value)
  }
})
const classes = computed(() => {
  return ['oc-radio', 'oc-radio-' + getSizeClass(size)]
})
const labelClasses = computed(() => {
  return {
    'oc-invisible-sr': hideLabel,
    'oc-cursor-pointer': !disabled
  }
})
</script>

<style lang="scss">
@mixin oc-form-check-size($factor) {
  height: $oc-size-form-check-default * $factor;
  width: $oc-size-form-check-default * $factor;
}

.oc-radio {
  -webkit-appearance: none;
  -moz-appearance: none;

  border: 1px solid var(--oc-color-swatch-brand-default);
  border-radius: 50%;
  box-sizing: border-box;
  background-color: var(--oc-color-input-bg);
  background-position: 50% 50%;
  background-repeat: no-repeat;

  display: inline-block;
  margin: 0;
  overflow: hidden;

  transition: 0.2s ease-in-out;
  transition-property: background-color, border;
  vertical-align: middle;
  width: 1rem;

  &:not(:disabled) {
    cursor: pointer;
  }

  &:checked {
    background-color: var(--oc-color-background-highlight) !important;
  }

  &.oc-radio-s {
    @include oc-form-check-size(0.7);
  }

  &.oc-radio-m {
    @include oc-form-check-size(1);
  }

  &.oc-radio-l {
    @include oc-form-check-size(1.5);
  }
}

label > .oc-radio + span {
  margin-left: var(--oc-space-xsmall);
}
</style>
