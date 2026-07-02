<template>
  <span :key="`oc-switch-${checked.toString()}`" class="oc-switch">
    <span :id="labelId" v-text="label" />
    <button
      data-testid="oc-switch-btn"
      class="oc-switch-btn"
      role="switch"
      :aria-checked="checked"
      :aria-labelledby="labelId"
      @click="toggle"
    />
  </span>
</template>

<script lang="ts" setup>
import { uniqueId } from '../../helpers'

/**
 * OcSwitch Component
 *
 * A toggle switch component that allows users to switch between two states.
 *
 * @component
 * @name OcSwitch
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {boolean} [checked=false] - The current state of the switch (true for on, false for off).
 * @prop {string} [label=null] - The accessible name of the switch.
 * @prop {string} [labelId=uniqueId('oc-switch-label-')] - The ID of the label element. If not provided, a unique ID is generated.
 *
 * @emits
 * @event update:checked - Emitted when the switch state changes.
 * @type {boolean}
 *
 * @example
 *   <OcSwitch
 *     :checked="isOn"
 *     label="Enable notifications"
 *     @update:checked="handleToggle"
 *   />
 *
 */

interface Props {
  checked?: boolean
  label?: string
  labelId?: string
}

interface Emits {
  (e: 'update:checked', value: boolean): void
}

defineOptions({
  name: 'OcSwitch',
  status: 'ready',
  release: '1.0.0'
})
const {
  checked = false,
  label = null,
  labelId = uniqueId('oc-switch-label-')
} = defineProps<Props>()

const emit = defineEmits<Emits>()

function toggle() {
  emit('update:checked', !checked)
}
</script>

<style lang="scss">
.oc-switch {
  align-items: center;
  display: inline-flex;
  gap: var(--oc-space-small);

  &-btn {
    border: 1px solid var(--oc-color-input-bg);
    border-radius: 20px;
    cursor: pointer;
    display: block;
    height: 18px;
    margin: 0;
    padding: 0;
    position: relative;
    transition: background-color 0.25s;
    width: 31px;

    &::before {
      background-color: var(--oc-color-swatch-inverse-hover);
      box-shadow: rgb(0 0 0 / 25%) 0px 0px 2px 1px;
      border-radius: 50%;
      content: '';
      height: 12px;
      left: 1px;
      position: absolute;
      top: 2px;
      transition: transform 0.25s;
      width: 12px;
    }

    &[aria-checked='false'] {
      background-color: var(--oc-color-swatch-inverse-muted);

      &::before {
        transform: translateX(0);
        left: 2px;
      }
    }

    &[aria-checked='true'] {
      background-color: var(--oc-color-swatch-primary-default);

      &::before {
        transform: translateX(calc(100% + 2px));
        left: 1px;
      }
    }
  }
}
</style>
