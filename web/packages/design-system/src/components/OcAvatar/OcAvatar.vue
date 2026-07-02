<template>
  <span
    class="vue-avatar--wrapper oc-avatar"
    :style="style"
    :width="width"
    :aria-label="accessibleLabel === '' ? null : accessibleLabel"
    :aria-hidden="accessibleLabel === '' ? 'true' : null"
    :focusable="accessibleLabel === '' ? 'false' : null"
    :role="accessibleLabel === '' ? null : 'img'"
    :data-test-user-name="userName"
  >
    <oc-img v-if="isImage" loading-type="lazy" class="avatarImg" :src="src" @error="onImgError" />
    <span v-else class="avatarInitials">{{ userInitial }}</span>
  </span>
</template>

<script lang="ts" setup>
import OcImg from '../OcImage/OcImage.vue'
import { extractInitials } from './extractInitials'
import { ref, computed, unref } from 'vue'

/**
 * OcAvatar - A component for displaying user avatars with support for images, initials, and accessible labels.
 *
 * @prop {string} [src] - The source URL of the avatar image. If not provided or if the image fails to load, initials will be displayed instead.
 * @prop {string} [userName=''] - The name of the user. Used to generate initials if no image is provided or the image fails to load.
 * @prop {string} [accessibleLabel=''] - An accessible label for the avatar. If empty, the avatar will be hidden from assistive technologies.
 * @prop {number} [width=50] - The width and height of the avatar in pixels. Defaults to 50px.
 *
 * @example
 * ```vue
 * <!-- Avatar with an image -->
 * <oc-avatar src="https://example.com" accessible-label="accessibleLabel" />
 *
 * <!-- Avatar with user initials -->
 * <oc-avatar user-name="lorem" accessible-label="lorem" />
 *
 * <!-- Avatar with custom size -->
 * <oc-avatar user-name="lorem" width="100" accessible-label="lorem" />
 * ```
 */

interface Props {
  src?: string
  userName?: string
  accessibleLabel?: string
  width?: number
}

defineOptions({
  name: 'OcAvatar',
  status: 'ready',
  release: '1.0.0'
})
const { src = '', userName = '', accessibleLabel = '', width = 50 } = defineProps<Props>()
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
const imgError = ref(false)
function onImgError() {
  imgError.value = true
}

function randomBackgroundColor(seed: number, colors: string[]) {
  return colors[seed % colors.length]
}
const background = computed(() => {
  if (!unref(isImage)) {
    return unref(randomBackgroundColor(unref(userName).length, unref(backgroundColors)))
  }
  return ''
})

const isImage = computed(() => {
  return !unref(imgError) && Boolean(unref(src))
})

const style = computed(() => {
  const style = {
    width: `${unref(width)}px`,
    height: `${unref(width)}px`,
    lineHeight: `${unref(width)}px`
  }

  const initialBackgroundAndFontStyle = {
    backgroundColor: unref(background),
    fontSize: `${Math.floor(unref(width) / 2.5)}px`,
    fontFamily: 'Helvetica, Arial, sans-serif',
    color: 'white'
  }

  Object.assign(style, initialBackgroundAndFontStyle)

  return style
})

const userInitial = computed(() => {
  if (!unref(isImage)) {
    return extractInitials(unref(userName))
  }
  return ''
})
</script>

<style lang="scss">
.oc-avatar {
  font-weight: normal;
  align-items: center;
  justify-content: center;
  text-align: center;
  user-select: none;
  display: flex;
  border-radius: 50%;

  .avatarImg {
    width: 100%;
    height: auto;
    border-radius: 50%;
  }

  .avatarInitials {
    color: white !important;
  }
}
</style>

<docs>
```js
  <oc-avatar class="oc-mb-s" src="https://picsum.photos/50/50?image=1074" accessible-label="Lion" />
  <oc-avatar class="oc-mb-s" user-name="Bruce Lee" accessible-label="Lion" />
```
</docs>
