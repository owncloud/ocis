<template>
  <div class="oc-login-card oc-position-center">
    <img class="oc-login-logo" :src="logoImg" alt="" :aria-hidden="true" />
    <div class="oc-login-card-body">
      <h1 v-text="$gettext('Missing or invalid config')" class="oc-login-card-title" />
      <p v-text="$gettext('Please check if the file config.json exists and is correct.')" />
      <p v-text="$gettext('Also, make sure to check the browser console for more information.')" />
    </div>
    <div class="oc-login-card-footer">
      <p>
        <span v-text="$gettext('For help visit our')" />
        <a
          href="https://owncloud.dev/clients/web"
          target="_blank"
          v-text="$gettext('documentation')"
        />
        <span v-text="$gettext('or join our')" />
        <a href="https://matrix.to/#/#ocis:matrix.org" target="_blank" v-text="$gettext('chat')" />.
      </p>
      <p>
        {{ footerSlogan }}
      </p>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent } from 'vue'
import { useThemeStore } from '@ownclouders/web-pkg'
import { useHead } from '../composables/head'
import { storeToRefs } from 'pinia'

export default defineComponent({
  name: 'MissingConfigPage',
  setup() {
    const themeStore = useThemeStore()
    const { currentTheme } = storeToRefs(themeStore)

    const logoImg = computed(() => currentTheme.value.logo.login)
    const footerSlogan = computed(() => currentTheme.value.common.slogan)

    useHead()

    return {
      logoImg,
      footerSlogan
    }
  }
})
</script>
