<template>
  <div class="oc-notification oc-mb-s" :class="classes">
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
/**
 * OcNotifications Component
 *
 * This component is used to display a container for notification messages. It supports different positions
 * for displaying notifications on the screen.
 *
 * Notifications for screen reader users
 * This component uses so called live regions in order to announce its content to screen readers once the notification appeared (this is not the normal modus operandi for screen readers, since their reading order is usually the DOM order – when the user does not take shortcuts). There are two types of live regions: `aria-live="polite"` (equivalent to `role="status"`) and `aria-live="assertive"` (equivalent to `role="alert"`). The latter directly interrupts the current output of the screen reader, the former waits until the current output is finished and reads the announcement afterwards. Since 'assertive' should be used sparingly, only `<oc-notfication-message>`'s "danger" status prop value uses `aria-live="assertive"` (and `role="alert"`). Using `aria-live` and `role="assertive|status"` simultaneously is for compatibility reasons regarding different browser and assistive technology pairings.
 *
 * @component
 * @name OcNotifications
 * @status ready
 * @release 1.0.0
 *
 * @props {string} [position='default'] - The position of the notification container.
 *                                        Possible values: 'default', 'top-left', 'top-center', 'top-right'.
 *
 * @slots default - Slot to include notification messages, such as OcNotificationMessage components.
 *
 * @computed
 * @computed {string} classes - Dynamically computed CSS class based on the `position` prop.
 *
 * @example
 * <OcNotifications position="top-right">
 *   <OcNotificationMessage
 *     status="success"
 *     title="Success"
 *     message="Your operation was successful."
 *   />
 * </OcNotifications>
 *
 */

interface Props {
  position?: 'default' | 'top-left' | 'top-center' | 'top-right'
}

defineOptions({
  name: 'OcNotifications',
  status: 'ready',
  release: '1.0.0'
})
const { position = 'default' } = defineProps<Props>()

const classes = computed(() => `oc-notification-${position}`)
</script>

<style lang="scss">
.oc-notification {
  box-sizing: border-box;
  max-width: 100%;
  width: 400px;
  z-index: 1040;

  &-top-left {
    position: fixed;
    top: var(--oc-space-small);
    left: var(--oc-space-small);
  }
  &-top-center {
    position: fixed;
    top: var(--oc-space-small);
    left: 0;
    right: 0;
    margin-left: auto;
    margin-right: auto;
  }
  &-top-right {
    position: fixed;
    top: var(--oc-space-small);
    right: var(--oc-space-small);
  }
}
</style>
