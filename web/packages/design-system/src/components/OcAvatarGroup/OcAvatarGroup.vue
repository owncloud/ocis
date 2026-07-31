<template>
  <oc-avatar-item
    :width="width"
    :icon-size="iconSize"
    icon="group"
    :name="name"
    :accessible-label="accessibleLabel"
  />
</template>

<script lang="ts" setup>
import OcAvatarItem from '../OcAvatarItem/OcAvatarItem.vue'

/**
 * OcAvatarGroup - A component for displaying a group avatar with a group icon.
 *
 * @prop {string} name - Name of the group used as an accessible label
 * @prop {string} [accessibleLabel=''] - Accessibility label used as alt. Use only in case the avatar is used alone. In case the avatar is used next to username or display name leave empty. If not specified, avatar will get `aria-hidden="true"`.
 * @prop {number} [width=30] - The width of the avatar in pixels. Defaults to 30px.
 * @prop {string} [iconSize='small'] - The size of the group icon. Defaults to "small".
 *
 * @example
 * ```vue
 * <!-- Default usage -->
 * <oc-avatar-group name="group" accessible-label="group" />
 *
 * <!-- Custom size -->
 * <oc-avatar-group name="team" accessible-label="team" :width="50" icon-size="medium" />
 * ```
 */

interface Props {
  name: string
  accessibleLabel?: string
  width?: number
  iconSize?: string
}
defineOptions({
  name: 'OcAvatarGroup',
  status: 'ready',
  release: '10.0.0'
})

const { name, accessibleLabel = '', width = 30, iconSize = 'small' } = defineProps<Props>()
</script>

<docs>
```js
<oc-avatar-group name="Group" accessible-label="Group" />
```
</docs>
