<template>
  <span class="oc-recipient">
    <template v-if="recipient.hasAvatar">
      <oc-spinner
        v-if="recipient.isLoadingAvatar"
        key="recipient-avatar-spinner"
        size="small"
        :aria-label="$gettext('Loading avatar')"
        data-testid="recipient-avatar-spinner"
      />
      <oc-avatar
        v-else
        :key="recipient.avatar || recipient.name"
        data-testid="recipient-avatar"
        class="oc-recipient-avatar"
        :src="recipient.avatar"
        :user-name="recipient.name"
        :width="16.8"
      />
    </template>
    <oc-icon
      v-else-if="recipient.icon && recipient.icon.name"
      class="oc-recipient-icon"
      size="small"
      :name="recipient.icon.name"
      :accessible-label="recipient.icon.label"
      data-testid="recipient-icon"
    />
    <p class="oc-recipient-name" data-testid="recipient-name" v-text="recipient.name" />
    <!-- @slot Append content (actions, additional info, etc.)  -->
    <slot name="append" />
  </span>
</template>

<script lang="ts" setup>
import OcAvatar from '../OcAvatar/OcAvatar.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcSpinner from '../OcSpinner/OcSpinner.vue'
import { Recipient } from '../../helpers'

/**
 * @component OcRecipient
 * @description A component that displays a recipient's avatar, icon, or name with optional appendable content.
 *
 * @props {Object} recipient - The recipient object containing details to render.
 * @props {string} recipient.name - The name of the recipient.
 * @props {string} [recipient.avatar] - The URL of the recipient's avatar image.
 * @props {boolean} [recipient.hasAvatar] - Indicates if the recipient has an avatar.
 * @props {boolean} [recipient.isLoadingAvatar] - Indicates if the avatar is currently loading.
 * @props {Object} [recipient.icon] - The icon object for the recipient.
 * @props {string} recipient.icon.name - The name of the icon to display.
 * @props {string} [recipient.icon.label] - The accessible label for the icon.
 *
 * @slot append - Slot for appending additional content (e.g., actions or extra information).
 *
 * @example
 *   <OcRecipient
 *     :recipient="{
 *       name: 'John Doe',
 *       avatar: 'https://example.com/avatar.jpg',
 *       hasAvatar: true,
 *       isLoadingAvatar: false,
 *       icon: { name: 'user', label: 'User Icon' }
 *     }"
 *   />
 */

interface Props {
  recipient: Recipient
}
defineOptions({
  name: 'OcRecipient',
  status: 'ready',
  release: '8.2.0'
})

const { recipient } = defineProps<Props>()
</script>

<style lang="scss">
.oc-recipient {
  align-items: center;
  background-color: var(--oc-color-background-default);
  border: 1px solid var(--oc-color-input-border);
  border-radius: 6px;
  display: flex;
  gap: var(--oc-space-xsmall);
  justify-content: flex-start;
  padding: var(--oc-space-xsmall);
  width: auto;

  &-icon > svg {
    fill: var(--oc-color-text-default);
  }

  &-name {
    color: var(--oc-color-text-default);
    margin: 0;
    padding: 0;
  }
}
</style>
