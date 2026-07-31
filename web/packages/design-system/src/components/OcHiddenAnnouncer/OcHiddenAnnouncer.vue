<template>
  <span
    :id="id"
    class="oc-invisible-sr oc-hidden-announcer"
    :aria-live="level"
    aria-atomic="true"
    v-text="announcement"
  />
</template>

<script lang="ts" setup>
import { computed, HTMLAttributes } from 'vue'

/**
 * OcHiddenAnnouncer Component
 *
 * Provides a live region for screen reader announcements.
 * It uses ARIA live region attributes to dynamically announce changes to assistive technologies.
 *
 * @component
 * @name OcHiddenAnnouncer
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {('polite' | 'assertive' | 'off')} [level='polite'] - The ARIA live region level.
 *   - `polite`: Adds the announcement to the screen reader speech queue at the end.
 *   - `assertive`: Forces the announcement to output directly.
 *   - `off`: Disables the live region.
 * @prop {string} announcement - The announcement text to be read by the screen reader.
 *
 * @example
 * <OcHiddenAnnouncer
 *   level="assertive"
 *   announcement="Form submitted successfully."
 * />
 *
 * @accessibility
 * Screen reader software detects dynamic changes in live regions (elements with attributes like `aria-live="polite"`, `aria-live="assertive"`).
 * Ensure the live region exists in the DOM before making changes to it, so assistive technology can subscribe to its updates.
 * For more details, refer to the [MDN ARIA Live Regions documentation](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/ARIA_Live_Regions).
 *
 * @debugging
 * Debug live regions without starting a screen reader using [NerdeRegion](https://chrome.google.com/webstore/detail/nerderegion/lkcampbojgmgobcfinlkgkodlnlpjieb).
 */

interface Props {
  level?: HTMLAttributes['aria-live']
  announcement: string
}
defineOptions({
  name: 'OcHiddenAnnouncer',
  status: 'ready',
  release: '1.0.0'
})

const { level = 'polite', announcement } = defineProps<Props>()

const id = computed(() => {
  return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
})
</script>
