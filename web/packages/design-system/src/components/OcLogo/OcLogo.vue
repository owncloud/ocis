<template>
  <oc-img class="oc-logo" :src="src" :alt="alt" />
</template>

<script lang="ts" setup>
import OcImg from '../OcImage/OcImage.vue'

/**
 * @component OcLogo
 * @description A logo component that displays an image with specified source and alternative text.
 * @status ready
 * @release 7.3.0
 *
 * @props
 * @prop {string} src - The source URL of the logo image.
 * @prop {string} alt - The alternative text for the logo image.
 *
 * @example
 * <oc-logo src="src" alt="alt" />
 */

interface Props {
  src: string
  alt: string
}
defineOptions({
  name: 'OcLogo',
  status: 'ready',
  release: '7.3.0'
})

const { src, alt } = defineProps<Props>()
</script>

<style lang="scss">
.oc-logo {
  max-width: var(--oc-size-max-width-logo);
  max-height: var(--oc-size-max-height-logo);
  margin: var(--oc-space-small);
}
</style>

<docs>
```js
<oc-logo src="examples/placeholder_brand_logo.svg" alt="ownCloud" />
```
</docs>
