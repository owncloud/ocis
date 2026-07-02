<template>
  <oc-avatar-item
    :width="width"
    :icon-size="iconSize"
    icon="global"
    icon-fill-type="line"
    icon-color="#D78841"
    :name="name"
    :accessible-label="accessibleLabel"
  />
</template>

<script lang="ts" setup>
import OcAvatarItem from '../OcAvatarItem/OcAvatarItem.vue'

/**
 * OcAvatarGuest - A component for displaying a guest avatar with a distinctive icon.
 *
 * @prop {string} name - Name of the guest used as an accessible label
 * @prop {string} [accessibleLabel=''] - Accessibility label used as alt. Use only in case the avatar is used alone. In case the avatar is used next to username or display name leave empty. If not specified, avatar will get `aria-hidden="true"`.
 * @prop {number} [width=30] - The width of the avatar in pixels. Defaults to 30px.
 * @prop {string} [iconSize='small'] - The size of the guest icon. Defaults to "small".
 *
 * @example
 * ```vue
 * <!-- Default usage -->
 * <oc-avatar-guest name="guest" accessible-label="guest" />
 *
 * <!-- Custom size -->
 * <oc-avatar-guest name="external" accessible-label="external" :width="50" icon-size="medium" />
 * ```
 */

interface Props {
  name: string
  accessibleLabel?: string
  width?: number
  iconSize?: string
}
defineOptions({
  name: 'OcAvatarGuest',
  status: 'ready',
  release: '10.0.0'
})
const { name, accessibleLabel = '', width = 30, iconSize = 'small' } = defineProps<Props>()
</script>

<docs>
```js
<oc-avatar-guest name="Guest" accessible-label="Guest" />
```
</docs>
