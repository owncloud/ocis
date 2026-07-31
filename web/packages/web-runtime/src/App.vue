<template>
  <maintenance-banner v-show="maintenanceMode" ref="maintenance-banner" />
  <portal-target name="app.app-banner" multiple />
  <div
    id="web"
    :style="{ '--web-runtime-maintenance-banner-height': maintenanceBannerHeight + 'px' }"
  >
    <oc-hidden-announcer :announcement="announcement" level="polite" />
    <skip-to target="web-content-main">
      <span v-text="$gettext('Skip to main')" />
    </skip-to>
    <component :is="layout"></component>
    <modal-wrapper />
  </div>
</template>
<script lang="ts" setup>
import { setCurrentLanguage } from './helpers/language'
import SkipTo from './components/SkipTo.vue'
import ModalWrapper from './components/ModalWrapper.vue'
import { useLayout } from './composables/layout'
import { additionalTranslations } from './helpers/additionalTranslations' // eslint-disable-line
import {
  eventBus,
  useConfigStore,
  useResourcesStore,
  useRouter,
  useThemeStore
} from '@ownclouders/web-pkg'
import { useHead } from './composables/head'
import { RouteLocation, useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { isEqual } from 'lodash-es'
import { MaintenanceBanner } from './components/MaintenanceBanner'
import { useTemplateRef, ref, onMounted, unref, computed, onWatcherCleanup, watch } from 'vue'
import { MaybeElement, useElementSize } from '@vueuse/core'
import { useMaintenanceMode } from './composables/maintenanceMode'
import { useGettext } from 'vue3-gettext'

const resourcesStore = useResourcesStore()
const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const configStore = useConfigStore()
const { maintenanceMode } = storeToRefs(configStore)

const { startCheckingMaintenanceMode, stopCheckingMaintenanceMode } = useMaintenanceMode()

const router = useRouter()
const route = useRoute()
useHead()

const activeRoute = computed(() => router.resolve(unref(router.currentRoute)))

const { layout, layoutType } = useLayout({ router })

const maintenanceBannerElement = useTemplateRef('maintenance-banner')
const { height: maintenanceBannerHeight } = useElementSize(
  maintenanceBannerElement as unknown as MaybeElement
)
const announcement = ref('')

const getText = useGettext()
const { $gettext, current: currentLanguage } = getText

function announceRouteChange(pageTitle: string) {
  announcement.value = $gettext('Navigated to %{ pageTitle }', { pageTitle })
}

function extractPageTitleFromRoute(route: RouteLocation) {
  const routeTitle = route.meta.title ? $gettext(route.meta.title.toString()) : undefined
  if (!routeTitle) {
    return
  }
  const glue = ' - '
  const titleSegments = [routeTitle]
  return {
    shortDocumentTitle: titleSegments.join(glue),
    fullDocumentTitle: [...titleSegments, unref(currentTheme).common.name].join(glue)
  }
}
watch(
  () => unref(activeRoute),
  (newRoute, oldRoute) => {
    /**
     * Hide global loading spinner. It usually gets hidden after all apps
     * have been loaded, but in some scenarios (plain layouts) we never load them.
     */
    if (unref(layoutType) !== 'application') {
      const loader = document.getElementById('splash-loading')
      if (!loader?.classList.contains('splash-hide')) {
        loader.classList.add('splash-hide')
      }
    }

    const getAppContextFromRoute = (route: RouteLocation): string[] => {
      return route?.path?.split('/').slice(1, 4)
    }

    const oldAppContext = getAppContextFromRoute(oldRoute)
    const newAppContext = getAppContextFromRoute(newRoute)

    if (isEqual(oldAppContext, newAppContext)) {
      return
    }

    if ('driveAliasAndItem' in newRoute.params) {
      return
    }

    /*
     * If app context has been changed and no file context is set, we will reset current folder.
     */
    resourcesStore.setCurrentFolder(null)
  }
)

watch(
  () => route,
  (to) => {
    const extracted = extractPageTitleFromRoute(to)
    if (!extracted) {
      return
    }
    const { shortDocumentTitle, fullDocumentTitle } = extracted
    announceRouteChange(shortDocumentTitle)
    document.title = fullDocumentTitle
  },
  { immediate: true }
)

watch(
  maintenanceMode,
  (maintenanceMode) => {
    if (maintenanceMode) {
      startCheckingMaintenanceMode()
    } else {
      stopCheckingMaintenanceMode()
    }

    onWatcherCleanup(() => {
      stopCheckingMaintenanceMode()
    })
  },
  { immediate: true }
)

onMounted(() => {
  eventBus.subscribe(
    'runtime.documentTitle.changed',
    ({
      shortDocumentTitle,
      fullDocumentTitle
    }: {
      shortDocumentTitle: string
      fullDocumentTitle: string
    }) => {
      document.title = fullDocumentTitle
      announceRouteChange(shortDocumentTitle)
    }
  )

  setCurrentLanguage({
    language: getText,
    languageSetting: currentLanguage
  })
})
</script>
<style lang="scss">
body {
  margin: 0;
}

#web {
  background-color: var(--oc-color-swatch-brand-default);
  height: calc(100dvh - var(--web-runtime-maintenance-banner-height));
  max-height: 100dvh;
  overflow-y: hidden;

  .mark-highlight {
    font-weight: 600;
  }
}

iframe {
  border: 0;
}
</style>
