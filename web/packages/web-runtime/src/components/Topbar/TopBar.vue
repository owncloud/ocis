<template>
  <header
    id="oc-topbar"
    :class="{ 'open-app': contentOnLeftPortal }"
    :aria-label="$gettext('Top bar')"
  >
    <div class="oc-topbar-left oc-flex oc-flex-middle oc-flex-start">
      <applications-menu
        v-if="appMenuExtensions.length && !isEmbedModeEnabled && !hideAppSwitcher"
        :menu-items="appMenuExtensions"
      />
      <!-- eslint-disable-next-line vuejs-accessibility/anchor-has-content -->
      <a v-if="!hideLogo && logoHref" :href="logoHref" class="oc-logo-href">
        <oc-responsive-image
          :src="{
            xs: currentTheme.logo.topbarSm,
            md: currentTheme.logo.topbar
          }"
          :alt="sidebarLogoAlt"
          class="oc-logo-image"
        />
      </a>
      <router-link v-else-if="!hideLogo" :to="homeLink" class="oc-logo-href">
        <oc-responsive-image
          :src="{
            xs: currentTheme.logo.topbarSm,
            md: currentTheme.logo.topbar
          }"
          :alt="sidebarLogoAlt"
          class="oc-logo-image"
        />
      </router-link>
      <div v-if="!isEmbedModeEnabled && canAccessVault" class="oc-flex oc-flex-middle">
        <oc-button
          id="oc-topbar-mode-switch-btn"
          class="oc-topbar-mode-switch"
          appearance="raw"
          gap-size="none"
        >
          <span class="oc-mr-xs" v-text="selectedMode.label" />
          <oc-icon name="arrow-down-s" />
        </oc-button>
        <oc-drop
          toggle="#oc-topbar-mode-switch-btn"
          mode="click"
          padding-size="small"
          close-on-click
        >
          <oc-list class="oc-topbar-mode-switch-list">
            <li v-for="option in modeOptions" :key="option.id">
              <oc-button
                appearance="raw"
                justify-content="space-between"
                class="oc-topbar-mode-switch-option oc-p-s oc-width-1-1"
                @click="selectedMode = option"
              >
                <span>{{ option.label }}</span>
                <oc-icon v-if="selectedMode.id === option.id" name="check" />
              </oc-button>
            </li>
          </oc-list>
        </oc-drop>
      </div>
    </div>
    <div v-if="!contentOnLeftPortal" class="oc-topbar-center oc-width-1-1">
      <custom-component-target :extension-point="topBarCenterExtensionPoint" />
    </div>
    <div class="oc-topbar-right oc-flex oc-flex-middle">
      <portal-target name="app.runtime.header.right" multiple />
    </div>
    <template v-if="!isEmbedModeEnabled">
      <portal to="app.runtime.header.right" :order="50">
        <feedback-link
          v-if="isFeedbackLinkEnabled && !hideAccountMenu"
          v-bind="feedbackLinkOptions"
        />
        <universal-access-dropdown v-if="isUniversalAccessEnabled" />
        <help-menu v-if="!hideAccountMenu" />
      </portal>
      <portal to="app.runtime.header.right" :order="100">
        <notifications v-if="isNotificationBellEnabled && !hideAccountMenu" />
        <side-bar-toggle v-if="isSideBarToggleVisible" :disabled="isSideBarToggleDisabled" />
        <user-menu v-if="!hideAccountMenu" />
      </portal>
    </template>
    <portal-target name="app.runtime.header.left" @change="updateLeftPortal" />
  </header>
</template>

<script lang="ts">
import { storeToRefs } from 'pinia'
import { computed, unref, PropType, ref } from 'vue'
import ApplicationsMenu from './ApplicationsMenu.vue'
import UserMenu from './UserMenu.vue'
import Notifications from './Notifications.vue'
import FeedbackLink from './FeedbackLink.vue'
import SideBarToggle from './SideBarToggle.vue'
import UniversalAccessDropdown from './UniversalAccessDropdown.vue'
import HelpMenu from './HelpMenu.vue'
import {
  ApplicationInformation,
  CustomComponentTarget,
  useAuthStore,
  useCapabilityStore,
  useConfigStore,
  useEmbedMode,
  useExtensionRegistry,
  useOpenEmptyEditor,
  useRouter,
  useThemeStore,
  useClipboardStore,
  useAbility
} from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { isRuntimeRoute } from '../../router'
import { appMenuExtensionPoint, topBarCenterExtensionPoint } from '../../extensionPoints'

export default {
  components: {
    ApplicationsMenu,
    CustomComponentTarget,
    FeedbackLink,
    HelpMenu,
    Notifications,
    SideBarToggle,
    UserMenu,
    UniversalAccessDropdown
  },
  props: {
    applicationsList: {
      type: Array as PropType<ApplicationInformation[]>,
      required: false,
      default: (): ApplicationInformation[] => []
    }
  },
  setup() {
    const { $gettext } = useGettext()
    const capabilityStore = useCapabilityStore()
    const themeStore = useThemeStore()
    const { currentTheme } = storeToRefs(themeStore)
    const configStore = useConfigStore()
    const { options: configOptions } = storeToRefs(configStore)
    const extensionRegistry = useExtensionRegistry()
    const { openEmptyEditor } = useOpenEmptyEditor()
    const { clearClipboard } = useClipboardStore()
    const ability = useAbility()

    const authStore = useAuthStore()
    const router = useRouter()
    const { isEnabled: isEmbedModeEnabled } = useEmbedMode()

    const appMenuExtensions = computed(() => {
      return extensionRegistry.requestExtensions(appMenuExtensionPoint)
    })

    const hideLogo = computed(() => unref(configOptions).hideLogo)
    const hideAppSwitcher = computed(() => unref(configOptions).hideAppSwitcher)
    const hideAccountMenu = computed(() => unref(configOptions).hideAccountMenu)

    const isNotificationBellEnabled = computed(() => {
      return (
        authStore.userContextReady && capabilityStore.notificationsOcsEndpoints.includes('list')
      )
    })

    const logoHref = computed(() => unref(currentTheme).logo.href)

    const homeLink = computed(() => {
      if (authStore.publicLinkContextReady && !authStore.userContextReady) {
        return {
          name: 'resolvePublicLink',
          params: { token: authStore.publicLinkToken }
        }
      }

      const isVaultScope = unref(router.currentRoute).params?.scope === 'vault'
      return isVaultScope ? '/vault/files' : '/'
    })

    const isFilesPublicUpoad = computed(() => {
      return unref(router.currentRoute).name === 'files-public-upload'
    })

    const isSideBarToggleVisible = computed(() => {
      return authStore.userContextReady || authStore.publicLinkContextReady
    })
    const isSideBarToggleDisabled = computed(() => {
      return isRuntimeRoute(unref(router.currentRoute)) || unref(isFilesPublicUpoad)
    })

    const contentOnLeftPortal = ref(false)
    const updateLeftPortal = (newContent: { hasContent: boolean; sources: string[] }) => {
      contentOnLeftPortal.value = newContent.hasContent
    }

    const isUniversalAccessEnabled = computed(() => {
      return (
        unref(currentTheme).common?.urls?.universalAccessEasyLanguage ||
        unref(currentTheme).common?.urls?.universalAccessSignLanguage
      )
    })

    const modeOptions = computed(() => {
      return [
        {
          id: 'default-mode',
          label: $gettext('Drive'),
          route: '/'
        },
        {
          id: 'vault-mode',
          label: $gettext('Vault'),
          route: '/vault/files'
        }
      ]
    })

    const selectedMode = computed({
      get() {
        const currentPath = window.location.pathname
        return currentPath.startsWith('/vault') ? unref(modeOptions)[1] : unref(modeOptions)[0]
      },
      set(mode) {
        if (mode.id === 'default-mode') {
          clearClipboard()
        }
        if (mode?.route && mode.route !== window.location.pathname) {
          window.location.href = mode.route
        }
      }
    })

    const canAccessVault = computed(
      () => capabilityStore.vaultEnabled && ability.can('read-all', 'Vault')
    )

    return {
      configOptions,
      contentOnLeftPortal,
      currentTheme,
      updateLeftPortal,
      isNotificationBellEnabled,
      hideLogo,
      isEmbedModeEnabled,
      isSideBarToggleVisible,
      isSideBarToggleDisabled,
      logoHref,
      homeLink,
      topBarCenterExtensionPoint,
      appMenuExtensions,
      hideAppSwitcher,
      hideAccountMenu,
      isUniversalAccessEnabled,
      modeOptions,
      selectedMode,
      canAccessVault
    }
  },
  computed: {
    sidebarLogoAlt() {
      return this.$gettext('Navigate to personal files page')
    },

    isFeedbackLinkEnabled() {
      return !this.configOptions.disableFeedbackLink
    },

    feedbackLinkOptions() {
      const feedback = this.configOptions.feedbackLink
      if (!this.isFeedbackLinkEnabled || !feedback) {
        return {}
      }

      return {
        ...(feedback.href && { href: feedback.href }),
        ...(feedback.ariaLabel && { ariaLabel: feedback.ariaLabel }),
        ...(feedback.description && { description: feedback.description })
      }
    }
  }
}
</script>

<style lang="scss">
#oc-topbar {
  align-items: center;
  display: grid;
  grid-template-areas: 'logo center right' 'secondRow secondRow secondRow';
  grid-template-columns: 30% 30% 40%;
  grid-template-rows: 52px auto;
  padding: 0 1rem;
  position: sticky;
  z-index: 5;

  @media (min-width: $oc-breakpoint-small-default) {
    column-gap: 10px;
    grid-template-columns: max-content 9fr 1fr;
    grid-template-rows: 1;
    height: 52px;
    justify-content: center;
    padding: 0 1.1rem;
  }

  &.open-app {
    grid-template-columns: 30% 30% 40%;

    @media (min-width: $oc-breakpoint-small-default) {
      grid-template-columns: max-content 1fr 1fr;
    }
  }

  img {
    max-height: 38px;
    image-rendering: auto;
    image-rendering: crisp-edges;
    image-rendering: pixelated;
    image-rendering: -webkit-optimize-contrast;
    user-select: none;
    width: 100%;
  }

  .oc-topbar-left {
    gap: 10px;
    grid-area: logo;
    @media (min-width: $oc-breakpoint-small-default) {
      gap: 20px;
    }

    .oc-logo-href {
      flex: 1;
    }
  }

  .oc-topbar-center {
    display: flex;
    grid-area: center;
    justify-content: flex-end;

    @media (min-width: $oc-breakpoint-small-default) {
      justify-content: center;
    }
  }

  .oc-topbar-right {
    gap: 10px;
    grid-area: right;
    justify-content: space-between;

    @media (min-width: $oc-breakpoint-small-default) {
      gap: 20px;
      justify-content: flex-end;
    }
  }

  .oc-topbar-mode-switch {
    color: var(--oc-color-swatch-brand-contrast);
    flex-shrink: 0;
    text-transform: uppercase;
    font-weight: bold;

    .oc-icon > svg {
      fill: var(--oc-color-swatch-brand-contrast);
    }
  }

  .oc-topbar-mode-switch-list li {
    margin: var(--oc-space-xsmall) 0;

    &:first-child {
      margin-top: 0;
    }

    &:last-child {
      margin-bottom: 0;
    }
  }

  .oc-topbar-mode-switch-option {
    &:hover,
    &:focus {
      background-color: var(--oc-color-background-hover);
      text-decoration: none;
    }
  }
}
</style>
