<template>
  <div
    class="oc-fade-in oc-flex oc-flex-wrap oc-notification-message oc-box-shadow-medium oc-rounded oc-p-m"
    :class="classes"
  >
    <div class="oc-flex oc-flex-wrap oc-flex-middle oc-flex-1" :role="role" :aria-live="ariaLive">
      <div class="oc-flex oc-flex-middle oc-flex-between oc-width-1-1">
        <div class="oc-flex oc-flex-middle">
          <oc-icon :variation="iconVariation" name="information" fill-type="line" class="oc-mr-s" />
          <div class="oc-notification-message-title">
            {{ title }}
          </div>
        </div>
        <oc-button appearance="raw" :aria-label="$gettext('Close')" @click="close"
          ><oc-icon name="close"
        /></oc-button>
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

interface Props {
  status?: 'passive' | 'primary' | 'success' | 'warning' | 'danger'
  title: string
  message?: string
  errorLogContent?: string
  timeout?: number
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
  timeout = 5
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

onMounted(() => {
  /**
   * Notification will be destroyed if timeout is set
   */
  if (timeout !== 0) {
    setTimeout(() => {
      close()
    }, timeout * 1000)
  }
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
  }

  &-error-log-toggle-button {
    word-break: keep-all;
  }
}
</style>
