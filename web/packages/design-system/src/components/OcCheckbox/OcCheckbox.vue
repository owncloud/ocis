<template>
  <span>
    <input
      :id="id"
      v-model="model"
      type="checkbox"
      name="checkbox"
      :class="classes"
      :value="option"
      :disabled="disabled"
      :aria-label="labelHidden ? label : null"
      @click="$emit('click', $event)"
      @keydown.enter="keydownEnter"
    />
    <label v-if="!labelHidden" :for="id" :class="computedLabelClasses" v-text="label" />
  </span>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { isEqual } from 'lodash-es'
import { getSizeClass, uniqueId } from '../../helpers'

/**
 * OcCheckbox - A customizable checkbox component supporting single checkboxes and checkbox groups.
 *
 * @prop {string} [id] - Unique identifier for the checkbox input. If not provided, an auto-generated ID will be used.
 * @prop {boolean|Array} [modelValue=false] - Value bound to the checkbox. Boolean for single checkboxes, array for checkbox groups.
 * @prop {boolean} [disabled=false] - Whether the checkbox is disabled.
 * @prop {any} [option=null] - Value associated with the checkbox when used in a checkbox group with array binding.
 * @prop {string} [label=null] - Text label displayed next to the checkbox.
 * @prop {boolean} [labelHidden=false] - Whether to hide the label visually but keep it accessible to screen readers.
 * @prop {'small'|'medium'|'large'} [size='medium'] - Size of the checkbox.
 *
 * @event {Event} click - Emitted when the checkbox is clicked.
 * @event {boolean|Array} update:modelValue - Emitted when the checkbox value changes.
 *
 * @example
 * ```vue
 * <!-- Basic checkbox -->
 * <oc-checkbox v-model="isChecked" label="label" />
 *
 * <!-- Disabled checkbox -->
 * <oc-checkbox v-model="isChecked" label="label" disabled />
 *
 * <!-- Different sizes -->
 * <oc-checkbox size="small" label="label" />
 * <oc-checkbox label="label" />
 * <oc-checkbox size="large" label="label" />
 *
 * <!-- Checkbox group with array model -->
 * <oc-checkbox
 *   v-for="option in ['a', 'b', 'c']"
 *   :key="option"
 *   v-model="selectedOptions"
 *   :option="option"
 *   :label="option"
 * />
 * ```
 */

interface Props {
  id?: string
  modelValue?: boolean | unknown[]
  disabled?: boolean
  option?: unknown
  label?: string
  labelHidden?: boolean
  labelClasses?: string[]
  size?: 'small' | 'medium' | 'large'
}
interface Emits {
  (e: 'click', event: Event): void
  (e: 'update:modelValue', value: boolean | unknown[]): void
}

defineOptions({
  name: 'OcCheckbox',
  status: 'ready',
  release: '1.0.0'
})

const {
  id = uniqueId('oc-checkbox-'),
  disabled = false,
  option = null,
  label = null,
  labelHidden = false,
  size = 'medium',
  modelValue = false,
  labelClasses = []
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const model = computed({
  get() {
    return modelValue
  },
  set(value: boolean) {
    emit('update:modelValue', value)
  }
})
function keydownEnter(event: KeyboardEvent) {
  model.value = !model.value
  emit('click', event)
}
const classes = computed(() => {
  return [
    'oc-checkbox',
    'oc-rounded',
    `oc-checkbox-${getSizeClass(size)}`,
    { 'oc-checkbox-checked': isChecked.value }
  ]
})

const computedLabelClasses = computed(() => [{ 'oc-cursor-pointer': !disabled }, ...labelClasses])

const isChecked = computed(() => {
  if (Array.isArray(model.value)) {
    return unref(model.value).some((m) => isEqual(m, option))
  }
  return unref(model)
})
</script>

<style lang="scss">
@mixin oc-form-check-size($factor) {
  height: $oc-size-form-check-default * $factor;
  width: $oc-size-form-check-default * $factor;
}

.oc-checkbox {
  @include oc-form-check-size(1);
  -webkit-appearance: none;
  -moz-appearance: none;

  background-position: 50% 50% !important;
  background-repeat: no-repeat !important;
  border: 2px solid var(--oc-color-input-border);
  display: inline-block;
  overflow: hidden;
  vertical-align: middle;
  background-color: transparent;
  outline: none;

  &-s {
    @include oc-form-check-size(0.7);
  }

  &-m {
    @include oc-form-check-size(1);
  }

  &-l {
    @include oc-form-check-size(1.5);
  }

  &:hover {
    cursor: pointer;
  }

  &:focus-visible {
    outline: var(--oc-color-swatch-primary-default) auto 1px;
  }

  &-checked,
  :checked,
  &:indeterminate {
    background-color: white;
  }

  &-checked,
  :checked {
    @include svg-fill($internal-form-checkbox-image, '#000', '#000');
  }

  &:indeterminate {
    @include svg-fill($internal-form-checkbox-indeterminate-image, '#000', '#000');
  }

  &:disabled {
    background-color: $form-radio-disabled-background;
    cursor: default;
    opacity: 0.4;
  }

  &:disabled:checked {
    @include svg-fill($internal-form-checkbox-image, '#000', $form-radio-disabled-icon-color);
  }

  &:disabled:indeterminate {
    @include svg-fill(
      $internal-form-checkbox-indeterminate-image,
      '#000',
      $form-radio-disabled-icon-color
    );
  }
}

label > .oc-checkbox + span {
  margin-left: var(--oc-space-xsmall);
}
</style>

<docs>
```js
<template>
  <section>
    <h3 class="oc-heading-divider oc-mt-s">
      Checkboxes Types
    </h3>
    <div class="oc-mb-s">
      <oc-checkbox size="small" label="Small checkbox" aria-label="Small checkbox"/>
    </div>
    <div class="oc-mb-s">
      <oc-checkbox :value="true" label="Medium checkbox (default)"/>
    </div>
    <div>
      <oc-checkbox size="large" label="Large checkbox"/>
    </div>
  </section>
</template>
```
```js
<template>
  <section>
    <h3 class="oc-heading-divider oc-mt-s">
      Checkbox group with array model
    </h3>
    <div class="oc-mb-s">
      <oc-checkbox
        v-for="o in availableOptions"
        :key="'option-' + o"
        v-model="selectedOptions"
        :option="o"
        :label="o"
        class="oc-mr-s"
      />
    </div>
    Selected option: {{ selectedOptions || "None" }}
  </section>
</template>
<script>
  export default {
    data: () => ({
      availableOptions: ["Water", "Wine", "Beer"],
      selectedOptions: []
    })
  }
</script>
```
</docs>
