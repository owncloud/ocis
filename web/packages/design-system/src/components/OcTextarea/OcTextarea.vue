<template>
  <div>
    <label class="oc-label" :for="id" v-text="label" />
    <textarea
      :id="id"
      v-bind="additionalAttributes"
      ref="textAreaRef"
      class="oc-textarea oc-rounded"
      :class="{
        'oc-textarea-warning': !!warningMessage,
        'oc-textarea-danger': !!errorMessage
      }"
      :value="modelValue"
      :aria-invalid="ariaInvalid"
      @input="onInput(($event.target as HTMLInputElement).value)"
      @focus="onFocus(true)"
      @keydown="onKeyDown($event)"
    />
    <div v-if="showMessageLine" class="oc-textarea-message">
      <span
        :id="messageId"
        :class="{
          'oc-textarea-description': !!descriptionMessage,
          'oc-textarea-warning': !!warningMessage,
          'oc-textarea-danger': !!errorMessage
        }"
        v-text="messageText"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, HTMLAttributes, useAttrs, useTemplateRef } from 'vue'
import { uniqueId } from '../../helpers'

/**
 * @component OcTextarea
 * @description A customizable textarea component that supports validation states,
 * accessibility features, and various user interactions.
 *
 * @example
 *   <OcTextarea
 *     v-model="textContent"
 *     label="Description"
 *     description-message="description"
 *   />
 *
 *   <OcTextarea
 *     v-model="comment"
 *     label="Comment"
 *     error-message="This field is required"
 *     @change="submitComment"
 *   />
 * </template>
 *
 * @prop {string} [id] - Unique identifier for the textarea. Auto-generated if not provided.
 * @prop {string} [modelValue] - The value of the textarea (for v-model binding).
 * @prop {string} [label] - Label text displayed above the textarea.
 * @prop {string} [warningMessage] - Warning message to display below the textarea.
 * @prop {string} [errorMessage] - Error message to display below the textarea.
 * @prop {string} [descriptionMessage] - Description message to display below the textarea.
 * @prop {boolean} [fixMessageLine=false] - Whether to always show the message line, even without a message.
 * @prop {boolean} [submitOnEnter=true] - Whether to emit change event when Enter is pressed.
 *
 * @event {string} update:modelValue - Emitted when the input value changes.
 * @event {boolean} focus - Emitted when the textarea is focused.
 * @event {string} change - Emitted when Enter is pressed and submitOnEnter is true.
 * @event {KeyboardEvent} keydown - Emitted when any key is pressed.
 *
 * @method focus() - Programmatically focus the textarea.
 */

interface Props {
  id?: string
  modelValue?: string
  label?: string
  warningMessage?: string
  errorMessage?: string
  descriptionMessage?: string
  fixMessageLine?: boolean
  submitOnEnter?: boolean
}
interface Emits {
  (event: 'update:modelValue', value: string): void
  (event: 'focus', value: boolean): void
  (event: 'change', value: string): void
  (event: 'keydown', value: KeyboardEvent): void
}
defineOptions({
  name: 'OcTextarea',
  status: 'ready',
  release: '1.0.0'
})
const {
  id = uniqueId('oc-textarea-'),
  modelValue = null,
  label = null,
  warningMessage = null,
  errorMessage = null,
  descriptionMessage = null,
  fixMessageLine = false,
  submitOnEnter = true
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const attrs = useAttrs()
const textAreaRef = useTemplateRef('textAreaRef')

function focus() {
  textAreaRef.value.focus()
}
function onInput(value: string) {
  /**
   * Input event
   * @type {event}
   **/
  emit('update:modelValue', value)
}
function onFocus(value: boolean) {
  /**
   * Focus event - emitted as soon as the input field is focused
   * @type {event}
   **/
  emit('focus', value)
}
function onKeyDown(e: KeyboardEvent) {
  const enterKey = e.key?.toLowerCase() === 'enter'
  if (submitOnEnter && enterKey && !e.ctrlKey && !e.shiftKey) {
    /**
     * Change event - emitted as soon as the user hits enter (without ctrl or shift)
     * Only applies if submitOnEnter is set to true
     * @type {string}
     */
    emit('change', (e.target as HTMLInputElement).value)
  }

  /**
   * KeyDown event - emitted as soon as the user hits a key
   * @type {event}
   */
  emit('keydown', e)
}
const showMessageLine = computed(() => {
  return fixMessageLine || !!warningMessage || !!errorMessage || !!descriptionMessage
})
const messageId = computed(() => {
  return `${id}-message`
})
const additionalAttributes = computed(() => {
  const additionalAttrs: Record<string, unknown> = {}
  if (!!warningMessage || !!errorMessage || !!descriptionMessage) {
    additionalAttrs['aria-describedby'] = messageId
  }
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

defineExpose({
  focus
})
</script>

<style lang="scss">
.oc-textarea {
  padding-bottom: var(--oc-space-xsmall);
  padding-top: var(--oc-space-xsmall);
  box-sizing: border-box;
  background: var(--oc-color-background-muted);
  border: 0 none;
  margin: 0;
  color: var(--oc-color-text-default);
  max-width: 100%;
  width: 100%;
  overflow: auto;

  &:disabled {
    color: var(--oc-color-input-text-muted);
  }

  &:focus {
    border-color: var(--oc-color-input-text-default);
    background-color: var(--oc-color-background-muted);
    color: var(--oc-color-text-default);
  }

  &-warning,
  &-warning:focus {
    border-color: var(--oc-color-swatch-warning-default);
    color: var(--oc-color-swatch-warning-default);
  }

  &-danger,
  &-danger:focus {
    border-color: var(--oc-color-swatch-danger-default);
    color: var(--oc-color-swatch-danger-default);
  }

  &-description {
    color: var(--oc-color-text-muted);
  }

  &-message {
    display: flex;
    align-items: center;
    margin-top: var(--oc-space-xsmall);

    min-height: $oc-font-size-default * 1.5;
  }
}
</style>
