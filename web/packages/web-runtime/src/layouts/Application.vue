<template>
  <div id="web-content">
    <div id="global-progress-bar">
      <custom-component-target :extension-point="progressBarExtensionPoint" />
    </div>
    <div id="web-content-header">
      <div v-if="isIE11" class="oc-background-muted oc-text-center oc-py-m">
        <p class="oc-m-rm" v-text="ieDeprecationWarning" />
      </div>
      <top-bar :applications-list="Object.values(apps)" />
    </div>
    <div id="web-content-main" class="oc-px-s oc-pb-s">
      <div class="app-container oc-flex">
        <app-loading-spinner v-if="isLoading" />
        <template v-else>
          <sidebar-nav
            v-if="isSidebarVisible"
            class="app-navigation"
            :nav-items="navItems"
            :closed="navBarClosed"
            @update:nav-bar-closed="setNavBarClosed"
          />
          <portal to="app.runtime.mobile.nav">
            <mobile-nav
              v-if="isMobileWidth && navItems.length && !hideNavigation"
              :nav-items="navItems"
            />
          </portal>
          <router-view
            v-for="name in ['default', 'app', 'fullscreen']"
            :key="`router-view-${name}`"
            class="app-content oc-width-1-1"
            :name="name"
          />
        </template>
      </div>

      <portal-target name="app.runtime.footer" />
    </div>
    <div class="snackbars">
      <div class="oc-invisible-sr" aria-live="polite" aria-atomic="true">
        <div v-for="message in allMessages" :key="message.id">
          {{ message.desc ? `${message.title}. ${message.desc}` : message.title }}
        </div>
      </div>
      <message-bar />
      <upload-bar v-if="!isUploadSnackbarHidden" id="upload-info-snackbar" />
    </div>
  </div>
</template>

<script lang="ts">
import orderBy from 'lodash-es/orderBy'
import {
  AppLoadingSpinner,
  CustomComponentExtension,
  CustomComponentTarget,
  Extension,
  ExtensionPoint,
  useAppsStore,
  useAuthStore,
  useConfigStore,
  useExtensionRegistry,
  useLocalStorage
} from '@ownclouders/web-pkg'
import TopBar from '../components/Topbar/TopBar.vue'
import MessageBar from '../components/MessageBar.vue'
import SidebarNav from '../components/SidebarNav/SidebarNav.vue'
import UploadBar from '../components/UploadBar.vue'
import MobileNav from '../components/MobileNav.vue'
import { NavItem, getExtensionNavItems } from '../helpers/navItems'
import { LoadingIndicator } from '@ownclouders/web-pkg'
import { useActiveApp, useRoute, useRouteMeta, useSpacesLoading } from '@ownclouders/web-pkg'
import {
  computed,
  defineComponent,
  nextTick,
  onBeforeUnmount,
  onMounted,
  provide,
  ref,
  unref,
  watch
} from 'vue'
import { RouteLocationAsRelativeTyped, useRouter } from 'vue-router'
import { useGettext } from 'vue3-gettext'
import { useMessages } from '@ownclouders/web-pkg'

import '@uppy/core/css/style.min.css'
import { storeToRefs } from 'pinia'

const MOBILE_BREAKPOINT = 640

export default defineComponent({
  name: 'ApplicationLayout',
  components: {
    AppLoadingSpinner,
    CustomComponentTarget,
    MessageBar,
    MobileNav,
    TopBar,
    SidebarNav,
    UploadBar
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const { $gettext } = useGettext()
    const authStore = useAuthStore()
    const activeApp = useActiveApp()
    const extensionRegistry = useExtensionRegistry()
    const messageStore = useMessages()

    const allMessages = ref<{ id: string; title: string; desc?: string }[]>([])

    watch(
      () => route.value.params?.scope,
      () => {
        extensionRegistry.rebuild({ route })
      },
      {
        immediate: true
      }
    )

    watch(
      () => messageStore.messages,
      (messages) => {
        if (messages && messages.length > 0) {
          allMessages.value = messages.map((msg) => ({
            id: msg.id,
            title: msg.title,
            desc: msg.desc
          }))
        }
      },
      { deep: true, immediate: true }
    )

    const appsStore = useAppsStore()
    const { apps } = storeToRefs(appsStore)

    const configStore = useConfigStore()
    const { options: configOptions } = storeToRefs(configStore)

    const extensionNavItems = computed(() =>
      getExtensionNavItems({ extensionRegistry, appId: unref(activeApp) })
    )

    // FIXME: we can convert to a single router-view without name (thus without the loop) and without this watcher when we release v6.0.0
    watch(
      useRoute(),
      (route) => {
        if (unref(route).matched.length) {
          unref(route).matched.forEach((match) => {
            const keys = Object.keys(match.components).filter((key) => key !== 'default')
            if (keys.length) {
              console.warn(
                `named components are deprecated, use "default" instead of "${keys.join(
                  ', '
                )}" on route ${String(route.name)}`
              )
            }
          })
        }
      },
      { immediate: true }
    )

    const uploadSnackbarRouteMeta = useRouteMeta('isUploadSnackbarHidden', 'false')
    const isUploadSnackbarHidden = computed<boolean>(() => {
      return JSON.parse(unref(uploadSnackbarRouteMeta))
    })

    const requiredAuthContext = useRouteMeta('authContext')
    const { areSpacesLoading } = useSpacesLoading()
    const isLoading = computed(() => {
      if (['anonymous', 'idp'].includes(unref(requiredAuthContext))) {
        return false
      }
      return unref(areSpacesLoading)
    })

    const isMobileWidth = ref<boolean>(window.innerWidth < MOBILE_BREAKPOINT)
    provide('isMobileWidth', isMobileWidth)
    const onResize = () => {
      isMobileWidth.value = window.innerWidth < MOBILE_BREAKPOINT
    }

    const navItems = computed<NavItem[]>(() => {
      if (!authStore.userContextReady) {
        return []
      }

      const { href: currentHref } = router.resolve(unref(route))
      const newRef = currentHref.replace(/^\/vault/, '')

      return orderBy(
        unref(extensionNavItems).map((item) => {
          let active = typeof item.isActive !== 'function' || item.isActive()

          if (active) {
            active = [item.route, ...(item.activeFor || [])].filter(Boolean).some((currentItem) => {
              try {
                const comparativeHref = router
                  .resolve(currentItem as RouteLocationAsRelativeTyped)
                  .href.replace(/^\/vault/, '')

                return newRef.startsWith(comparativeHref)
              } catch (e) {
                console.error(e)
                return false
              }
            })
          }

          const name = typeof item.name === 'function' ? item.name() : item.name

          return {
            ...item,
            name: $gettext(name),
            active
          }
        }),
        ['priority', 'name']
      )
    })

    const hideNavigation = computed(() => unref(configOptions).hideNavigation)
    const isSidebarVisible = computed(() => {
      return unref(navItems).length && !unref(isMobileWidth) && !unref(hideNavigation)
    })

    const navBarClosed = useLocalStorage(`oc_navBarClosed`, false)
    const setNavBarClosed = (value: boolean) => {
      navBarClosed.value = value
    }

    const progressBarExtensionId = 'com.github.owncloud.web.runtime.default-progress-bar'
    const progressBarExtensionPointId = 'app.runtime.global-progress-bar'
    const defaultProgressBarExtension = computed<CustomComponentExtension>(() => ({
      id: progressBarExtensionId,
      type: 'customComponent',
      extensionPointIds: [progressBarExtensionPointId],
      content: LoadingIndicator,
      userPreference: {
        optionLabel: $gettext('Default progress bar')
      }
    }))

    extensionRegistry.registerExtensions(
      computed(() => [unref(defaultProgressBarExtension)] satisfies Extension[])
    )
    const progressBarExtensionPoint = computed<ExtensionPoint<CustomComponentExtension>>(() => ({
      id: progressBarExtensionPointId,
      extensionType: 'customComponent',
      multiple: false,
      defaultExtensionId: unref(defaultProgressBarExtension).id,
      userPreference: {
        label: $gettext('Global progress bar'),
        description: $gettext('Customize your progress bar')
      }
    }))
    const extensionPoints = computed<ExtensionPoint<Extension>[]>(() => [
      unref(progressBarExtensionPoint)
    ])
    extensionRegistry.registerExtensionPoints(extensionPoints)

    onMounted(async () => {
      await nextTick()
      window.addEventListener('resize', onResize)
      onResize()
    })

    onBeforeUnmount(() => {
      window.removeEventListener('resize', onResize)

      extensionRegistry.unregisterExtensions([progressBarExtensionId])
      extensionRegistry.unregisterExtensionPoints(unref(extensionPoints).flatMap((e) => e.id))
    })

    return {
      apps,
      progressBarExtensionPoint,
      isSidebarVisible,
      isLoading,
      navItems,
      isMobileWidth,
      navBarClosed,
      hideNavigation,
      isUploadSnackbarHidden,
      setNavBarClosed,
      allMessages
    }
  },
  computed: {
    isIE11() {
      return !!(window as any).MSInputMethodContext && !!(document as any).documentMode
    },
    ieDeprecationWarning() {
      return this.$gettext(
        'Internet Explorer (your current browser) is not officially supported. For security reasons, please switch to another browser.'
      )
    }
  }
})
</script>
<style lang="scss">
#web-content {
  display: flex;
  flex-flow: column;
  flex-wrap: nowrap;
  height: calc(100dvh - var(--web-runtime-maintenance-banner-height));

  #global-progress-bar {
    z-index: 10;
    position: absolute;
    top: 0;
    width: 100%;
  }

  #web-content-header,
  #web-content-main {
    flex-shrink: 1;
    flex-basis: auto;
  }

  #web-content-header {
    flex-grow: 0;
  }

  #web-content-main {
    align-items: flex-start;
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    justify-content: flex-start;
    overflow-y: hidden;

    .app-container {
      height: 100%;
      background-color: var(--oc-color-background-default);
      border-radius: 15px;
      overflow: hidden;
      width: 100%;

      .app-content {
        transition: all 0.35s cubic-bezier(0.34, 0.11, 0, 1.12);
      }
    }
  }

  .snackbars {
    position: absolute;
    right: 20px;
    bottom: 20px;
    z-index: calc(var(--oc-z-index-modal) + 1);

    @media (max-width: 640px) {
      left: 0;
      right: 0;
      margin: 0 auto;
      width: 100%;
      max-width: 500px;
    }

    #upload-info-snackbar {
      width: 400px;
      @media (max-width: 640px) {
        width: 100%;
        max-width: 500px;
      }
    }
  }
}
</style>
