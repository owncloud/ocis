<template>
  <div class="oc-modal-background">
    <focus-trap :active="true" :initial-focus="initialFocusRef" :tabbable-options="tabbableOptions">
      <!-- eslint-disable-next-line vuejs-accessibility/no-static-element-interactions -->
      <div
        :id="elementId"
        ref="ocModal"
        :class="classes"
        tabindex="0"
        role="dialog"
        aria-modal="true"
        aria-labelledby="oc-modal-title"
        @keydown.esc="cancelModalAction"
      >
        <div class="oc-modal-title">
          <oc-icon v-if="iconName !== ''" :name="iconName" :variation="variation" />
          <h2 id="oc-modal-title" class="oc-text-truncate" v-text="title" />
        </div>
        <div class="oc-modal-body">
          <div v-if="$slots.content" key="modal-slot-content" class="oc-modal-body-message">
            <slot name="content" />
          </div>
          <template v-else>
            <p
              v-if="message"
              key="modal-message"
              class="oc-modal-body-message oc-mt-rm"
              :class="{ 'oc-mb-rm': !hasInput || contextualHelperData }"
              v-text="message"
            />
            <div
              v-if="contextualHelperData"
              class="oc-modal-body-contextual-helper"
              :class="{ 'oc-mb-rm': !hasInput }"
            >
              <span class="text" v-text="contextualHelperLabel" />
              <oc-contextual-helper class="oc-pl-xs" v-bind="contextualHelperData" />
            </div>
            <oc-text-input
              v-if="hasInput"
              key="modal-input"
              ref="ocModalInput"
              v-model="userInputValue"
              class="oc-modal-body-input"
              :error-message="inputError"
              :label="inputLabel"
              :type="inputType"
              :description-message="inputDescription"
              :fix-message-line="true"
              :selection-range="inputSelectionRange"
              @update:model-value="inputOnInput"
              @enter-key-down="confirm"
            />
          </template>
        </div>

        <div v-if="!hideActions" class="oc-modal-body-actions oc-flex oc-flex-right">
          <div class="oc-modal-body-actions-grid">
            <oc-button
              class="oc-modal-body-actions-cancel"
              variation="passive"
              appearance="outline"
              :disabled="isLoading"
              @click="cancelModalAction"
              >{{ $gettext(buttonCancelText) }}
            </oc-button>
            <oc-button
              v-if="!hideConfirmButton"
              class="oc-modal-body-actions-confirm oc-ml-s"
              variation="primary"
              :appearance="buttonConfirmAppearance"
              :disabled="isLoading || buttonConfirmDisabled || !!inputError"
              :show-spinner="showSpinner"
              @click="confirm"
              >{{ $gettext(buttonConfirmText) }}
            </oc-button>
          </div>
        </div>
      </div>
    </focus-trap>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, computed, unref } from 'vue'
import OcButton from '../OcButton/OcButton.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcTextInput from '../OcTextInput/OcTextInput.vue'
import { FocusTrap } from 'focus-trap-vue'
import { FocusTargetOrFalse, FocusTrapTabbableOptions } from 'focus-trap'
import { ContextualHelperData } from '../../helpers'

/**
 * @component OcModal
 *
 * @description
 * A reusable modal component designed to focus user attention on a single action or confirmation.
 *
 * @features
 * - Displays a modal with customizable title, message, and actions.
 * - Supports optional input fields with validation and error messages.
 * - Includes contextual helper data for additional guidance.
 * - Configurable cancel and confirm buttons with loading states.
 * - Focus management using `focus-trap` for accessibility.
 *
 * Modals are generally used to force the user to focus on confirming or completing a single action.
 *
 * ## Background and position
 * Every modal gets automatically added a background which spans the whole width and height.
 * The modal itself is aligned to center both vertically and horizontally.
 *
 * ## Variations
 * Only use the `danger` variation if the action cannot be undone.
 *
 * The overall variation defines the modal's top border, heading (including an optional item) text color and the
 * variation of the confirm button, while the cancel buttons defaults to the `passive` variation. Both button's
 * variations and appearances can be targeted individually (see examples and API docs below).
 *
 * @props
 * @prop {string} [elementId] - Optional modal ID.
 * @prop {string} [elementClass] - Optional modal class.
 * @prop {'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'info'} [variation='passive'] - Modal variation.
 * @prop {string} [icon] - Optional icon to display next to the title.
 * @prop {string} title - Modal title (required).
 * @prop {string} [message] - Modal message (can be replaced by content slot).
 * @prop {string} [contextualHelperLabel] - Label for contextual helper data.
 * @prop {ContextualHelperData} [contextualHelperData] - Data for contextual helper.
 * @prop {string} [buttonCancelText='Cancel'] - Text for the cancel button.
 * @prop {string} [buttonConfirmText='Confirm'] - Text for the confirm button.
 * @prop {boolean} [buttonConfirmDisabled=false] - Disables the confirm button.
 * @prop {boolean} [hideConfirmButton=false] - Hides the confirm button.
 * @prop {boolean} [hasInput=false] - Enables an input field in the modal.
 * @prop {string} [inputType='text'] - Type of the input field.
 * @prop {string} [inputValue] - Value of the input field.
 * @prop {[number, number]} [inputSelectionRange] - Selection range for the input field.
 * @prop {string} [inputLabel] - Label for the input field.
 * @prop {string} [inputDescription] - Description message for the input field.
 * @prop {string} [inputError] - Error message for the input field.
 * @prop {string | boolean} [focusTrapInitial] - Custom initial focus target.
 * @prop {boolean} [hideActions=false] - Hides the action buttons at the bottom.
 * @prop {boolean} [isLoading=false] - Enables loading state for the modal.
 *
 * @emits
 * @event cancel - Triggered when the cancel button is clicked or the escape key is pressed.
 * @event confirm - Triggered when the confirm button is clicked. Emits the input value if present.
 * @event input - Triggered when the user types into the input field. Emits the input value.
 *
 * @slots
 * @slot content - Custom content to replace the default message.
 *
 *
 */

interface Props {
  elementId?: string
  elementClass?: string
  variation?: 'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'info'
  icon?: string
  title: string
  message?: string
  contextualHelperLabel?: string
  contextualHelperData?: ContextualHelperData
  buttonCancelText?: string
  buttonConfirmText?: string
  buttonConfirmDisabled?: boolean
  hideConfirmButton?: boolean
  hasInput?: boolean
  inputType?: string
  inputValue?: string
  inputSelectionRange?: [number, number]
  inputLabel?: string
  inputDescription?: string
  inputError?: string
  focusTrapInitial?: string | boolean
  hideActions?: boolean
  isLoading?: boolean
}

interface Emits {
  (e: 'cancel'): void
  (e: 'confirm', value: string): void
  (e: 'input', value: string): void
}

defineOptions({
  name: 'OcModal',
  status: 'ready',
  release: '1.3.0'
})

const {
  elementId = null,
  elementClass = null,
  variation = 'passive',
  icon = null,
  title,
  message = null,
  contextualHelperLabel = '',
  contextualHelperData = null,
  buttonCancelText = 'Cancel',
  buttonConfirmText = 'Confirm',
  buttonConfirmDisabled = false,
  hideConfirmButton = false,
  hasInput = false,
  inputType = 'text',
  inputValue = null,
  inputSelectionRange = null,
  inputLabel = null,
  inputDescription = null,
  inputError = null,
  focusTrapInitial = null,
  hideActions = false,
  isLoading = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const userInputValue = ref<string>()
const showSpinner = ref(false)
const buttonConfirmAppearance = ref('filled')
const ocModal = ref<HTMLElement>()
const ocModalInput = ref<typeof OcTextInput>()

const tabbableOptions = computed((): FocusTrapTabbableOptions => {
  // Enable shadow DOM support for e.g. emoji-picker
  return {
    getShadowRoot: true
  }
})

const resetLoadingState = () => {
  showSpinner.value = false
  buttonConfirmAppearance.value = 'filled'
}

const setLoadingState = () => {
  showSpinner.value = true
  buttonConfirmAppearance.value = 'outline'
}
function cancelModalAction() {
  /**
   * The user clicked on the cancel button or hit the escape key
   */
  emit('cancel')
}
function confirm() {
  if (buttonConfirmDisabled || inputError) {
    return
  }
  /**
   * The user clicked on the confirm button. If input exists, emits its value
   *
   * @property {String} value Value of the input
   */
  emit('confirm', unref(userInputValue))
}
function inputOnInput(value: string) {
  /**
   * The user typed into the input
   *
   * @property {String} value Value of the input
   */
  emit('input', value)
}
watch(
  () => isLoading,
  () => {
    if (!isLoading) {
      return resetLoadingState()
    }
    setTimeout(() => {
      if (!isLoading) {
        return resetLoadingState()
      }
      setLoadingState()
    }, 700)
  },
  { immediate: true }
)

watch(
  () => inputValue,
  (value: string) => {
    userInputValue.value = value
  },
  { immediate: true }
)
const initialFocusRef = computed<FocusTargetOrFalse>(() => {
  if (focusTrapInitial || focusTrapInitial === false) {
    return focusTrapInitial as FocusTargetOrFalse
  }

  return () => unref(ocModalInput)?.$el || unref(ocModal)
})
const classes = computed(() => {
  return ['oc-modal', `oc-modal-${variation}`, elementClass]
})
const iconName = computed(() => {
  if (icon) {
    return icon
  }

  switch (variation) {
    case 'danger':
      return 'alert'
    case 'warning':
      return 'error-warning'
    case 'success':
      return 'checkbox-circle'
    case 'info':
      return 'information'
    default:
      return ''
  }
})
</script>

<style lang="scss">
@mixin oc-modal-variation($color) {
  span {
    color: $color;
  }
}

.oc-modal {
  background-color: var(--oc-color-background-default);
  border: 1px solid var(--oc-color-input-border);
  border-radius: 15px;
  box-shadow: 5px 0 25px rgba(0, 0, 0, 0.3);
  max-height: 90dvh;
  max-width: 500px;
  overflow: auto;
  width: 100%;

  &:focus {
    outline: none;
  }

  &-background {
    align-items: center;
    background-color: rgba(0, 0, 0, 0.4);
    display: flex;
    flex-flow: row wrap;
    height: 100%;
    justify-content: center;
    left: 0;
    position: fixed;
    top: 0;
    width: 100%;
    z-index: var(--oc-z-index-modal);
  }

  &-primary {
    @include oc-modal-variation(var(--oc-color-swatch-primary-default));
  }

  &-success {
    @include oc-modal-variation(var(--oc-color-swatch-success-default));
  }

  &-warning {
    @include oc-modal-variation(var(--oc-color-swatch-warning-default));
  }

  &-danger {
    @include oc-modal-variation(var(--oc-color-swatch-danger-default));
  }

  &-title {
    align-items: center;
    border-bottom: 1px solid var(--oc-color-input-border);
    border-top-left-radius: 15px;
    border-top-right-radius: 15px;
    display: flex;
    flex-flow: row wrap;
    line-height: 1.625;
    padding: var(--oc-space-medium) var(--oc-space-medium);

    > .oc-icon {
      margin-right: var(--oc-space-small);
    }

    > h2 {
      font-size: 1rem;
      font-weight: bold;
      margin: 0;
    }
  }

  &-body {
    color: var(--oc-color-text-default);
    line-height: 1.625;
    padding: var(--oc-space-medium) var(--oc-space-medium) 0;

    &-message {
      margin-bottom: var(--oc-space-medium);
      margin-top: var(--oc-space-small);
    }

    &-contextual-helper {
      margin-bottom: var(--oc-space-medium);
    }

    .oc-input {
      line-height: normal;
    }

    &-input {
      /* FIXME: this is ugly, but required so that the bottom padding doesn't look off when reserving vertical space for error messages below the input. */
      margin-bottom: -20px;
      padding-bottom: var(--oc-space-medium);

      .oc-text-input-message {
        margin-bottom: var(--oc-space-xsmall);
      }
    }

    &-actions {
      text-align: right;
      background: var(--oc-color-background-default);
      border-bottom-right-radius: 15px;
      border-bottom-left-radius: 15px;
      padding: var(--oc-space-medium);

      .oc-button {
        border-radius: 4px;
      }

      &-grid {
        display: inline-grid;
        grid-auto-flow: column;
        grid-auto-columns: 1fr;
      }
    }
  }

  .oc-text-input-password-wrapper {
    button {
      background-color: var(--oc-color-background-highlight) !important;
    }
  }
}
</style>
