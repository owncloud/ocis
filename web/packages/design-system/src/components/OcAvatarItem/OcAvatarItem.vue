<template>
  <div
    :data-test-item-name="name"
    :aria-label="accessibleLabel === '' ? null : accessibleLabel"
    :aria-hidden="accessibleLabel === '' ? 'true' : null"
    :focusable="accessibleLabel === '' ? 'false' : null"
    :role="accessibleLabel === '' ? null : 'img'"
  >
    <span
      class="oc-avatar-item"
      :style="{
        backgroundColor,
        '--icon-color': iconColor,
        '--width': avatarWidth
      }"
    >
      <oc-icon v-if="hasIcon" :name="icon" :size="iconSize" :fill-type="iconFillType" />
    </span>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import OcIcon from '../OcIcon/OcIcon.vue'

/**
 * OcAvatarItem - A base component for displaying customizable avatar items with icons.
 *
 * @prop {string} name - Name of the public link used as an accessible label
 * @prop {string} [icon=null] - Icon that should be used for the avatar
 * @prop {string} [iconColor='var(--oc-color-text-inverse)'] - Color that should be used for the icon
 * @prop {string} [iconFillType='fill'] - Fill-type that should be used for the icon
 * @prop {string} [iconSize='small'] - Describes the size of the avatar icon e.g.(small)
 * @prop {string} [background='var(--oc-color-swatch-passive-default)'] - Background color that should be used for the avatar.
 *   If empty a random color will be picked
 * @prop {string} [accessibleLabel=''] - Accessibility label used as alt. Use only in case the avatar is used alone.
 *   In case the avatar is used next to username or display name leave empty.
 *   If not specified, avatar will get `aria-hidden="true"`.
 * @prop {number} [width=30] - Describes the width of the avatar
 *
 * @example
 * ```vue
 * <!-- Basic usage -->
 * <oc-avatar-item name="name" accessible-label="name" />
 *
 * <!-- With icon and default background -->
 * <oc-avatar-item name="name" icon="close" accessible-label="name" />
 *
 * <!-- With custom styling -->
 * <oc-avatar-item
 *   name="name"
 *   icon="close"
 *   background="#465a64"
 *   :width="100"
 *   icon-size="large"
 *   accessible-label="name"
 * />
 * ```
 */

interface Props {
  name: string
  icon?: string
  iconColor?: string
  iconFillType?: string
  iconSize?: string
  background?: string
  accessibleLabel?: string
  width?: number
}
defineOptions({
  name: 'OcAvatarItem',
  status: 'ready',
  release: '10.0.0'
})
const {
  name,
  icon = null,
  iconColor = 'var(--oc-color-text-inverse)',
  iconFillType = 'fill',
  iconSize = 'small',
  background = 'var(--oc-color-swatch-passive-default)',
  accessibleLabel = '',
  width = 30
} = defineProps<Props>()

const avatarWidth = computed(() => {
  return `${unref(width)}px`
})
const hasIcon = computed(() => {
  return unref(icon) !== null
})
const backgroundColor = computed(() => {
  return unref(background) || unref(pickBackgroundColor)
})
const pickBackgroundColor = computed(() => {
  const backgroundColors = [
    '#b82015',
    '#c21c53',
    '#9C27B0',
    '#673AB7',
    '#3F51B5',
    '#106892',
    '#055c68',
    '#208377',
    '#1a761d',
    '#476e1a',
    '#636d0b',
    '#8e5c11',
    '#795548',
    '#465a64'
  ]
  return backgroundColors[Math.floor(Math.random() * backgroundColors.length)]
})
</script>

<style lang="scss">
.oc-avatar-item {
  align-items: center;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 18px;
  border-radius: 50%;
  display: inline-flex;
  height: var(--width);
  justify-content: center;
  width: var(--width);

  .oc-icon > svg {
    fill: var(--icon-color) !important;
  }
}
</style>

<docs>
```js
<h3>Empty OcAvatarItem</h3>
<oc-avatar-item name="Public link" accessible-label="Public link" />
<h3>OcAvatarItem with icon and default background</h3>
<oc-avatar-item name="Public link" icon="close" accessible-label="Public link" />
<h3>OcAvatarItem with icon and custom background</h3>
<oc-avatar-item name="Public link" icon="close" background="#465a64" accessible-label="Public link" />
<h3>OcAvatarItem with icon and custom background and custom width and iconsize</h3>
<oc-avatar-item :width="100" iconSize="large" name="Public link" icon="close" background="#465a64" accessible-label="Public link" />
```
</docs>
