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
 * Neither this component nor `<oc-notification-message>` is a live region: screen readers do not reliably announce live regions that are inserted together with their content. Consumers must announce messages through a live region that is already present in the DOM (see web-runtime's MessageBar).
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
