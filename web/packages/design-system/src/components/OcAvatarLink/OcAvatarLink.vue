<template>
  <oc-avatar-item
    :width="width"
    :icon-size="iconSize"
    icon="link"
    :name="name"
    :accessible-label="accessibleLabel"
  />
</template>

<script lang="ts" setup>
import OcAvatarItem from '../OcAvatarItem/OcAvatarItem.vue'

/**
 * OcAvatarLink - A component for displaying a link avatar with a link icon.
 *
 * @prop {string} name - Name of the public link used as an accessible label
 * @prop {string} [accessibleLabel=''] - Accessibility label used as alt.
 *   Use only in case the avatar is used alone.
 *   In case the avatar is used next to username or display name leave empty.
 *   If not specified, avatar will get `aria-hidden="true"`.
 * @prop {number} [width=30] - Describes the width of the avatar
 * @prop {string} [iconSize='small'] - Describes the size of the avatar icon e.g.(small)
 *
 * @example
 * ```vue
 * <!-- Default usage -->
 * <oc-avatar-link name="name" accessible-label="name" />
 *
 * <!-- Custom size -->
 * <oc-avatar-link
 *   name="name"
 *   accessible-label="name"
 *   :width="50"
 *   icon-size="medium"
 * />
 * ```
 */

interface Props {
  name: string
  accessibleLabel?: string
  width?: number
  iconSize?: string
}
defineOptions({
  name: 'OcAvatarLink',
  status: 'ready',
  release: '2.1.0'
})
const { name, accessibleLabel = '', width = 30, iconSize = 'small' } = defineProps<Props>()
</script>

<docs>
```js
<oc-avatar-link name="Public link" accessible-label="Public link" />
```
</docs>
