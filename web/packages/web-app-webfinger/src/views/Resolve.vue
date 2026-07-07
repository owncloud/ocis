<template>
  <main
    class="webfinger-resolve oc-height-viewport oc-flex oc-flex-column oc-flex-center oc-flex-middle"
  >
    <h1 class="oc-invisible-sr" v-text="pageTitle" />
    <div class="oc-card oc-card-body oc-text-center oc-width-large">
      <template v-if="hasError">
        <h2 key="webfinger-resolve-error">
          <span v-text="$gettext('Sorry!')" />
        </h2>
        <p v-text="$gettext('Something went wrong.')" />
        <p v-text="$gettext('We could not resolve the destination.')" />
      </template>
      <template v-else>
        <h2 key="webfinger-resolve-loading">
          <span v-text="$gettext('One moment please…')" />
        </h2>
        <p v-text="$gettext('You are being redirected.')" />
      </template>
    </div>
  </main>
</template>

<script lang="ts" setup>
import { computed, ref, unref, watch } from 'vue'
import {
  useClientService,
  useConfigStore,
  useLoadingService,
  useRoute,
  useRouteMeta
} from '@ownclouders/web-pkg'
import { OwnCloudServer, WebfingerDiscovery } from '../discovery'
import { useGettext } from 'vue3-gettext'
import { useAuthService } from '@ownclouders/web-pkg'

const configStore = useConfigStore()
const clientService = useClientService()
const loadingService = useLoadingService()
const authService = useAuthService()
const route = useRoute()
const { $gettext } = useGettext()

const title = useRouteMeta('title', '')
const pageTitle = computed(() => {
  return $gettext(unref(title))
})

const ownCloudServers = ref<OwnCloudServer[]>([])
const hasError = ref(false)
const webfingerDiscovery = new WebfingerDiscovery(configStore.serverUrl, clientService)
loadingService.addTask(async () => {
  try {
    const servers = await webfingerDiscovery.discoverOwnCloudServers()
    ownCloudServers.value = servers
    if (servers.length === 0) {
      hasError.value = true
    }
  } catch (e) {
    console.error(e)
    if (e.response?.status === 401) {
      return authService.handleAuthError(unref(route), { forceLogout: true })
    }
    hasError.value = true
  }
})

watch(ownCloudServers, (instances) => {
  if (instances.length === 0) {
    return
  }
  // we can't deal with multi-instance results. just pick the first one for now.
  window.location.href = ownCloudServers.value[0].href
})
</script>

<style lang="scss">
.webfinger-resolve {
  .oc-card {
    background: var(--oc-color-background-highlight);
    border-radius: 15px;

    &-body {
      h2 {
        margin-top: 0;
      }
      p {
        font-size: var(--oc-font-size-large);
      }
    }
  }
}
</style>
