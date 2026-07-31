<template>
  <div
    class="oc-fade-in oc-flex oc-flex-wrap oc-notification-message oc-box-shadow-medium oc-rounded oc-p-m"
    :class="classes"
  >
    <div class="oc-flex oc-flex-wrap oc-flex-middle oc-flex-1" :role="role" :aria-live="ariaLive">
      <div class="oc-flex oc-flex-middle oc-width-1-1">
        <oc-icon :variation="iconVariation" name="information" fill-type="line" class="oc-mr-s" />
        <div class="oc-notification-message-title oc-flex-1">
          {{ title }}
        </div>
        <!-- eslint-disable-next-line vuejs-accessibility/no-static-element-interactions -->
        <div
          class="oc-notification-message-interactive oc-flex oc-flex-middle oc-ml-m"
          @mouseenter="onMouseEnter"
          @mouseleave="onMouseLeave"
          @focusin="onFocusIn"
          @focusout="onFocusOut"
        >
          <div v-if="actions.length" class="oc-notification-message-actions oc-flex oc-flex-middle">
            <oc-button
              v-for="(action, index) in actions"
              :key="index"
              appearance="raw"
              variation="primary"
              class="oc-notification-message-action-button"
              :aria-label="action.ariaLabel || action.label"
              @click="onActionClick(action)"
            >
              <span class="oc-notification-message-action-button-label">{{ action.label }}</span>
            </oc-button>
          </div>
          <oc-button
            class="oc-notification-message-close-button oc-ml-s"
            appearance="raw"
            :aria-label="$gettext('Close')"
            @click="close"
            ><oc-icon name="close"
          /></oc-button>
        </div>
      </div>
      <div v-if="message || errorLogContent" class="oc-flex oc-flex-between oc-width-1-1 oc-mt-s">
        <span
          v-if="message"
          class="oc-notification-message-content oc-text-muted oc-mr-s"
          v-text="message"
        />
        <oc-button
          v-if="errorLogContent"
          class="oc-notification-message-error-log-toggle-button"
          gap-size="none"
          appearance="raw"
          :aria-expanded="showErrorLog"
          @click="showErrorLog = !showErrorLog"
        >
          <span v-text="$gettext('Details')"></span>
          <oc-icon :name="showErrorLog ? 'arrow-up-s' : 'arrow-down-s'" />
        </oc-button>
      </div>
      <oc-error-log v-if="showErrorLog" class="oc-mt-m" :content="errorLogContent" />
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, ref, onMounted } from 'vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcButton from '../OcButton/OcButton.vue'
import OcErrorLog from '../OcErrorLog/OcErrorLog.vue'

/**
 * OcNotificationMessage Component
 *
 * This component is used to display notification messages to users. Notifications can have different statuses
 * (e.g., passive, primary, success, warning, danger) and can include a title, message, and optional error log content.
 * The component also supports an auto-dismiss feature based on a timeout.
 *
 * @component
 * @name OcNotificationMessage
 * @status ready
 * @release 1.0.0
 *
 * @props {string} [status='passive'] - The status of the notification. Defines the color and icon variation.
 *                                      Possible values: 'passive', 'primary', 'success', 'warning', 'danger'.
 * @props {string} title - The title of the notification. This is a required property.
 * @props {string} [message=null] - The message content of the notification.
 * @props {string} [errorLogContent=null] - The error log content to display when the "Details" button is clicked.
 * @props {number} [timeout=5] - The number of seconds the notification is displayed before auto-dismiss.
 *                                If set to 0, the notification will not auto-dismiss.
 * @props {array} [actions=[]] - Optional action buttons rendered in the notification, e.g. an "Undo" action.
 *
 * @emits {void} close - Emitted when the user clicks the close button or when the notification auto-dismisses.
 *
 * @example
 * <OcNotificationMessage
 *   status="success"
 *   title="Operation Successful"
 *   message="Your changes have been saved."
 *   :timeout="10"
 *   @close="handleClose"
 * />
 *
 */

export interface OcNotificationMessageAction {
  label: string
  ariaLabel?: string
  onClick: () => void
}

interface Props {
  status?: 'passive' | 'primary' | 'success' | 'warning' | 'danger'
  title: string
  message?: string
  errorLogContent?: string
  timeout?: number
  actions?: OcNotificationMessageAction[]
}
interface Emits {
  (e: 'close'): void
}
defineOptions({
  name: 'OcNotificationMessage',
  status: 'ready',
  release: '1.0.0'
})

const {
  status = 'passive',
  title,
  message = null,
  errorLogContent = null,
  timeout = 5,
  actions = []
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const showErrorLog = ref(false)

function close() {
  /**
   * The close event is emitted when the user clicks the close icon.
   * @type {void}
   */
  emit('close')
}

function onActionClick(action: OcNotificationMessageAction) {
  action.onClick()
  close()
}

const classes = computed(() => {
  return `oc-notification-message-${status}`
})
const iconVariation = computed(() => {
  return status
})
const isStatusDanger = computed(() => {
  return status === 'danger'
})
const role = computed(() => {
  return isStatusDanger.value ? 'alert' : 'status'
})
const ariaLive = computed(() => {
  return isStatusDanger.value ? 'assertive' : 'polite'
})

let remainingTime = timeout * 1000
let timeoutStartedAt = 0
let timeoutHandle: ReturnType<typeof setTimeout> = null
let isHovering = false
let isFocused = false

function startTimeout() {
  if (timeout === 0 || remainingTime <= 0 || timeoutHandle || isHovering || isFocused) {
    return
  }
  timeoutStartedAt = Date.now()
  timeoutHandle = setTimeout(() => {
    timeoutHandle = null
    close()
  }, remainingTime)
}

function pauseTimeout() {
  if (timeout === 0 || !timeoutHandle) {
    return
  }
  clearTimeout(timeoutHandle)
  timeoutHandle = null
  remainingTime -= Date.now() - timeoutStartedAt
}

function onMouseEnter() {
  isHovering = true
  pauseTimeout()
}

function onMouseLeave() {
  isHovering = false
  startTimeout()
}

function onFocusIn() {
  isFocused = true
  pauseTimeout()
}

function onFocusOut() {
  isFocused = false
  startTimeout()
}

onMounted(() => {
  /**
   * Notification will be destroyed if timeout is set
   */
  startTimeout()
})
</script>

<style lang="scss">
.oc-notification-message {
  background-color: var(--oc-color-background-default) !important;
  margin-top: var(--oc-space-small);
  position: relative;
  word-break: break-word;

  &-title {
    font-size: 1.15rem;
    min-width: 0;
  }

  &-error-log-toggle-button {
    word-break: keep-all;
  }

  &-interactive {
    flex-shrink: 0;
  }

  &-action-button-label {
    white-space: nowrap;
  }
}
</style>
