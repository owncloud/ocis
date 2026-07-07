<template>
  <div
    class="oc-login oc-height-viewport"
    :style="{ '--oc-login-background-image': 'url(' + backgroundImg + ')' }"
  >
    <h1 class="oc-invisible-sr" v-text="pageTitle" />
    <router-view />
  </div>
</template>

<script lang="ts">
import { storeToRefs } from 'pinia'
import { computed, defineComponent, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useRouteMeta, useThemeStore } from '@ownclouders/web-pkg'

export default defineComponent({
  name: 'PlainLayout',
  setup() {
    const { $gettext } = useGettext()
    const themeStore = useThemeStore()
    const { currentTheme } = storeToRefs(themeStore)

    const title = useRouteMeta('title')

    const pageTitle = computed(() => {
      return $gettext(unref(title) || '')
    })
    const backgroundImg = computed(() => currentTheme.value.loginPage.backgroundImg)

    return {
      pageTitle,
      backgroundImg
    }
  }
})
</script>
