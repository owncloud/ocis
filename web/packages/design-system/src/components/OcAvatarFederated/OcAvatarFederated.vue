<template>
  <oc-avatar-item
    :width="width"
    :icon-size="iconSize"
    icon="earth"
    icon-fill-type="line"
    icon-color="#85C1BA"
    :name="name"
    :accessible-label="accessibleLabel"
  />
</template>

<script lang="ts" setup>
import OcAvatarItem from '../OcAvatarItem/OcAvatarItem.vue'

/**
 * OcAvatarFederated - A component for displaying a federated avatar with a globe icon.
 *
 * @prop {string} name - The name associated with the avatar. This is required.
 * @prop {string} [accessibleLabel=''] - An accessible label for the avatar. If empty, the avatar will be hidden from assistive technologies.
 * @prop {number} [width=30] - The width and height of the avatar in pixels. Defaults to 30px.
 * @prop {string} [iconSize='small'] - The size of the globe icon. Defaults to "small".
 *
 * @example
 * ```vue
 * <!-- Default usage -->
 * <oc-avatar-federated name="federated" accessible-label="federated" />
 *
 * <!-- Custom size -->
 * <oc-avatar-federated name="federated" accessible-label="federated" :width="50" icon-size="medium" />
 * ```
 */

interface Props {
  name: string
  accessibleLabel?: string
  width?: number
  iconSize?: string
}
defineOptions({
  name: 'OcAvatarFederated',
  status: 'ready',
  release: '10.0.0'
})
const { name, accessibleLabel = '', width = 30, iconSize = 'small' } = defineProps<Props>()
</script>

<docs>
```js
<oc-avatar-federated name="Federated" accessible-label="Federated" />
```
</docs>
