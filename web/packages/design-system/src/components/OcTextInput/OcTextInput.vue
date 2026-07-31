<template>
  <div :class="attrs.class">
    <slot name="label">
      <label class="oc-label" :for="id" v-text="label" />
    </slot>
    <div class="oc-position-relative">
      <oc-icon
        v-if="readOnly"
        name="lock"
        size="small"
        class="oc-mt-s oc-ml-s oc-position-absolute"
      />
      <component
        :is="inputComponent"
        :id="id"
        v-bind="additionalAttributes"
        ref="input"
        :aria-invalid="ariaInvalid"
        class="oc-text-input oc-input oc-rounded"
        :class="{
          'oc-text-input-warning': !!warningMessage,
          'oc-text-input-danger': !!errorMessage,
          'oc-pl-l': !!readOnly,
          'clear-action-visible': showClearButton
        }"
        :type="type"
        :value="displayValue"
        :disabled="disabled || readOnly"
        v-on="additionalListeners"
        @change="onChange(($event.target as HTMLInputElement).value)"
        @input="onInput(($event.target as HTMLInputElement).value)"
        @password-challenge-completed="$emit('passwordChallengeCompleted')"
        @password-challenge-failed="$emit('passwordChallengeFailed')"
        @focus="onFocus($event.target)"
        @keydown.enter="$emit('enterKeyDown')"
      />
      <oc-button
        v-if="showClearButton"
        :aria-label="clearButtonAccessibleLabelValue"
        class="oc-pr-s oc-position-center-right oc-text-input-btn-clear"
        :class="{ 'oc-mr-rm': type === 'date' }"
        appearance="raw"
        @click="onClear"
      >
        <oc-icon name="close" size="small" variation="passive" />
      </oc-button>
    </div>
    <div
      v-if="showMessageLine"
      class="oc-text-input-message"
      :class="{
        'oc-text-input-description': !!descriptionMessage,
        'oc-text-input-warning': !!warningMessage,
        'oc-text-input-danger': !!errorMessage
      }"
      :aria-live="errorMessage || warningMessage ? 'assertive' : 'off'"
      aria-atomic="true"
      role="alert"
    >
      <oc-icon
        v-if="messageText !== null && !!descriptionMessage"
        name="information"
        size="small"
        fill-type="line"
        accessible-label="info"
        aria-hidden="true"
      />

      <span
        :id="messageId"
        :class="{
          'oc-text-input-description': !!descriptionMessage,
          'oc-text-input-warning': !!warningMessage,
          'oc-text-input-danger': !!errorMessage
        }"
        v-text="messageText"
      />
    </div>
    <portal-target v-if="type === 'password'" name="app.design-system.password-policy" />
  </div>
</template>

<script lang="ts" setup>
import { HTMLAttributes, nextTick, computed, useAttrs, useTemplateRef } from 'vue'

import { uniqueId } from '../../helpers'
import OcButton from '../OcButton/OcButton.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcTextInputPassword from '../OcTextInputPassword/OcTextInputPassword.vue'
import { PasswordPolicy } from '../../helpers'
import { useGettext } from 'vue3-gettext'

/**
 * OcTextInput component
 *
 * Form Inputs are used to allow users to provide text input when the expected
 * input is short. Form Input has a range of options and supports several text
 * formats including numbers. For longer input, use the form `Textarea` element.
 *
 * ## Accessibility
 * The label is required and represents the name of the input.
 *
 * The description-message can be used additionally to give further information about the input field. When a
 * description is given, it will be automatically referenced via the `aria-describedby` property.
 * An error or warning will replace the description as well as the `aria-describedby` property until the error
 * or warning is fixed.
 *
 * @component
 * @example
 * <oc-text-input label="Text" v-model="inputValue"/>
 * <oc-text-input type="password" label="Password" v-model="passwordValue"/>
 *
 * @prop {string} [id] - The ID of the element.
 * @prop {'text'|'number'|'email'|'password'|'date'} [type='text'] - The type of the form input field.
 * @prop {string} [modelValue] - Text value of the form input field.
 * @prop {[number, number]} [selectionRange] - Selection range to accomplish partial selection.
 * @prop {boolean} [clearButtonEnabled=false] - Whether or not the input element should have a dedicated button for clearing the input content.
 * @prop {string} [clearButtonAccessibleLabel] - The aria label for the clear button. Only used if it's enabled at all.
 * @prop {string} [defaultValue] - Value to show when no value is provided. This does not set `value` automatically. The user needs to explicitly enter a text to set it as `value`.
 * @prop {boolean} [disabled=false] - Disables the input field.
 * @prop {string} label - Accessible label of the form input field.
 * @prop {string} [warningMessage] - A warning message which is shown below the input.
 * @prop {string} [errorMessage] - An error message which is shown below the input.
 * @prop {boolean} [fixMessageLine=false] - Whether or not vertical space below the input should be reserved for a one line message, so that content actually appearing there doesn't shift the layout.
 * @prop {string} [descriptionMessage] - A description text which is shown below the input field.
 * @prop {boolean} [readOnly=false] - Determines if the input field is read only. Read only field will be visualized by a lock item and additionally behaves like a disabled field.
 * @prop {PasswordPolicy} [passwordPolicy] - Array of password policy rules, if type is password and password policy is given, the entered value will be checked against these rules.
 * @prop {function} [generatePasswordMethod] - Method to generate random password.
 *
 * @emits change - Emitted when the input value changes.
 * @emits update:modelValue - Emitted when the input value is updated.
 * @emits focus - Emitted when the input field is focused.
 * @emits passwordChallengeCompleted - Emitted when the password challenge is completed.
 * @emits passwordChallengeFailed - Emitted when the password challenge fails.
 * @emits enterKeyDown - Emitted when enter key is pressed.
 */

interface Props {
  id?: string
  type?: string | 'text' | 'number' | 'email' | 'password' | 'date'
  modelValue?: string
  selectionRange?: [number, number]
  clearButtonEnabled?: boolean
  clearButtonAccessibleLabel?: string
  defaultValue?: string
  disabled?: boolean
  label: string
  warningMessage?: string
  errorMessage?: string
  fixMessageLine?: boolean
  descriptionMessage?: string
  readOnly?: boolean
  passwordPolicy?: PasswordPolicy
  generatePasswordMethod?: (...args: unknown[]) => string
}
interface Emits {
  (e: 'change', value: string): void
  (e: 'update:modelValue', value: string): void
  (e: 'focus', value: string): void
  (e: 'passwordChallengeCompleted'): void
  (e: 'passwordChallengeFailed'): void
  (e: 'enterKeyDown'): void
}

defineOptions({
  name: 'OcTextInput',
  status: 'ready',
  release: '1.0.0'
})

defineExpose({
  focus
})
const emit = defineEmits<Emits>()
const attrs = useAttrs()
const {
  id = uniqueId('oc-textinput-'),
  type = 'text',
  modelValue = null,
  selectionRange = null,
  clearButtonEnabled = false,
  clearButtonAccessibleLabel = '',
  defaultValue = null,
  disabled = false,
  label,
  warningMessage = null,
  errorMessage = null,
  fixMessageLine = false,
  descriptionMessage = null,
  readOnly = false,
  passwordPolicy = {},
  generatePasswordMethod = null
} = defineProps<Props>()

const { $gettext } = useGettext()
const input = useTemplateRef('input')
const showMessageLine = computed(() => {
  return fixMessageLine || !!warningMessage || !!errorMessage || !!descriptionMessage
})
const messageId = computed(() => {
  return `${id}-message`
})
const additionalListeners = computed(() => {
  if (type === 'password') {
    return { passwordGenerated: onInput }
  }

  return {}
})

const additionalAttributes = computed(() => {
  const additionalAttrs: Record<string, unknown> = {}
  const describedByIds: string[] = []

  if (!!warningMessage || !!errorMessage || !!descriptionMessage) {
    describedByIds.push(messageId.value)
  }

  if (type === 'password' && Object.keys(passwordPolicy).length) {
    describedByIds.push(`${id}-password-policy`)
  }

  if (describedByIds.length > 0) {
    additionalAttrs['aria-describedby'] = describedByIds.join(' ')
  }

  // FIXME: placeholder usage is discouraged, we need to find a better UX concept
  if (defaultValue) {
    additionalAttrs['placeholder'] = defaultValue
  }
  if (type === 'password') {
    additionalAttrs['password-policy'] = passwordPolicy
    additionalAttrs['generate-password-method'] = generatePasswordMethod
    additionalAttrs['has-warning'] = !!warningMessage
    additionalAttrs['has-error'] = !!errorMessage
  }
  // Exclude listeners for events which are handled via methods in this component

  const { change, input, focus, class: classes, ...attrs } = useAttrs()

  return { ...attrs, ...additionalAttrs }
})

const ariaInvalid = computed(() => {
  return (!!errorMessage).toString() as HTMLAttributes['aria-invalid']
})
const messageText = computed(() => {
  if (errorMessage) {
    return errorMessage
  }

  if (warningMessage) {
    return warningMessage
  }

  return descriptionMessage
})
const showClearButton = computed(() => {
  return !disabled && clearButtonEnabled && !!modelValue
})
const clearButtonAccessibleLabelValue = computed(() => {
  return clearButtonAccessibleLabel || $gettext('Clear input')
})
const displayValue = computed(() => {
  return modelValue || ''
})
const inputComponent = computed(() => {
  return type === 'password' ? OcTextInputPassword : 'input'
})

/**
 * Puts focus on this input element
 * @public
 */
function focus() {
  ;(input.value as HTMLInputElement).focus()
}
function onClear() {
  focus()

  onInput(null)
  onChange(null)
}
function onChange(value: string) {
  /**
   * Change event
   * @type {event}
   **/
  emit('change', value)
}
function onInput(value: string) {
  /**
   * Input event
   * @type {event}
   **/
  emit('update:modelValue', value)
}
async function onFocus(target: HTMLInputElement) {
  await nextTick()
  target.select()
  if (selectionRange && selectionRange.length > 1) {
    target.setSelectionRange(selectionRange[0], selectionRange[1])
  }
  /**
   * Focus event - emitted as soon as the input field is focused
   * @type {event}
   **/
  emit('focus', target.value)
}
</script>

<style lang="scss">
.oc-text-input-message.oc-text-input-description {
  display: flex;
  align-items: center;
  position: relative;
  padding-left: var(--oc-space-large);
  padding-top: calc(var(--oc-space-xsmall) - 2px);

  .oc-icon {
    position: absolute;
    left: var(--oc-space-xsmall);
    top: var(--oc-space-xsmall);
  }
}

.oc-text-input {
  &-description {
    color: var(--oc-color-text-muted);
  }

  &-success,
  &-success:focus {
    border-color: var(--oc-color-swatch-success-default) !important;
    color: var(--oc-color-swatch-success-default) !important;
  }

  &-warning,
  &-warning:focus {
    border-color: var(--oc-color-swatch-warning-default) !important;
    color: var(--oc-color-swatch-warning-default) !important;
  }

  &-danger,
  &-danger:focus {
    border-color: var(--oc-color-swatch-danger-default) !important;
    color: var(--oc-color-swatch-danger-default) !important;
  }

  &-message {
    display: flex;
    align-items: center;
    margin-top: var(--oc-space-xsmall);
    min-height: $oc-font-size-default * 1.5;
  }

  &.clear-action-visible {
    padding-right: ($oc-size-icon-default * 0.7) + 7px;
  }
}
</style>

<docs>
```js
<template>
  <section>
    <h3 class="oc-heading-divider">
      Input Types
    </h3>
    <oc-text-input class="oc-mb-s" label="Text"/>
    <oc-text-input class="oc-mb-s" disabled label="Disabled" value="I am disabled"/>
    <oc-text-input class="oc-mb-s" read-only="true" label="Read only" value="I am read only"/>
    <oc-text-input class="oc-mb-s" type="number" label="Number"/>
    <oc-text-input class="oc-mb-s" type="email" label="Email"/>
    <oc-text-input class="oc-mb-s" type="password" label="Password"/>
    <h3 class="oc-heading-divider">
      Binding
    </h3>
    <oc-text-input label="Text" v-model="inputValue"/>
    <oc-text-input disabled label="Text" v-model="inputValue"/>
    <h3 class="oc-heading-divider">
      Interactions
    </h3>
    <oc-button @click="_focus" class="oc-my-m">Focus input below</oc-button>
    <oc-text-input label="Focus field" ref="inputForFocus"/>
    <oc-button @click="_focusAndSelect" class="oc-my-m">Focus and select input below</oc-button>
    <oc-text-input label="Select field" value="Will you select this existing text?" ref="inputForFocusSelect"/>
    <oc-text-input label="Clear input" v-model="inputValueForClearing" :clear-button-enabled="true"/>
    <oc-text-input label="Input with default" v-model="inputValueWithDefault" :clear-button-enabled="true"
                   default-value="Some default"/>
    <p>
      Value: {{ inputValueWithDefault !== null ? inputValueWithDefault : "null" }}
    </p>
    <h3 class="oc-heading-divider">
      Messages
    </h3>
    <oc-text-input
      label="Input with description message below"
      class="oc-mb-s"
      description-message="This is a description message."
      :fix-message-line="true"
    />
    <oc-text-input
      label="Input with error and warning messages with reserved space below"
      class="oc-mb-s"
      v-model="valueForMessages"
      :error-message="errorMessage"
      :warning-message="warningMessage"
      :fix-message-line="true"
    />
    <oc-text-input
      label="Input with error and warning messages without reserved space below"
      class="oc-mb-s"
      v-model="valueForMessages"
      :error-message="errorMessage"
      :warning-message="warningMessage"
    />
  </section>
</template>
<script>
  export default {
    data: () => {
      return {
        inputValue: 'initial',
        valueForMessages: '',
        inputValueForClearing: 'clear me',
        inputValueWithDefault: null,
      }
    },
    computed: {
      errorMessage() {
        return this.valueForMessages.length === 0 ? 'Value is required.' : ''
      },
      warningMessage() {
        return this.valueForMessages.endsWith(' ') ? 'Trailing whitespace should be avoided.' : ''
      }
    },
    methods: {
      _focus() {
        this.$refs.inputForFocus.focus()
      },
      _focusAndSelect() {
        this.$refs.inputForFocusSelect.focus()
      }
    }
  }
</script>
```
</docs>
