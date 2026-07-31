<template>
  <nav :aria-label="$gettext('Account menu')">
    <oc-button
      id="_userMenuButton"
      ref="menuButton"
      v-oc-tooltip="$gettext('My Account')"
      class="oc-topbar-personal"
      appearance="raw"
      :aria-label="$gettext('My Account')"
    >
      <avatar-image
        v-if="onPremisesSamAccountName"
        class="oc-topbar-avatar oc-topbar-personal-avatar oc-flex-inline oc-flex-center oc-flex-middle"
        :width="32"
        :userid="onPremisesSamAccountName"
        :user-name="user.displayName"
      />
      <oc-avatar-item
        v-else
        class="oc-topbar-avatar oc-topbar-unauthenticated-avatar oc-flex-inline oc-flex-center oc-flex-middle"
        :name="$gettext('User Menu login')"
        :width="32"
        icon="user"
        icon-fill-type="line"
        icon-color="var(--oc-color-swatch-brand-default)"
        background="var(--oc-color-swatch-brand-contrast)"
      />
    </oc-button>
    <oc-drop
      ref="menu"
      drop-id="account-info-container"
      toggle="#_userMenuButton"
      mode="click"
      close-on-click
      padding-size="small"
      class="oc-overflow-hidden"
    >
      <oc-list class="user-menu-list">
        <template v-if="!onPremisesSamAccountName">
          <li>
            <oc-button
              id="oc-topbar-account-manage"
              type="router-link"
              :to="accountPageRoute"
              appearance="raw"
            >
              <oc-icon name="settings-4" fill-type="line" class="oc-p-xs" />
              <span v-text="$gettext('Preferences')" />
            </oc-button>
          </li>
          <li>
            <oc-button
              id="oc-topbar-account-login"
              appearance="raw"
              type="router-link"
              :to="loginLink"
            >
              <oc-icon name="login-box" fill-type="line" class="oc-p-xs" />
              <span v-text="$gettext('Log in')" />
            </oc-button>
          </li>
        </template>
        <template v-else>
          <li class="profile-info-wrapper oc-pl-s">
            <avatar-image
              :width="32"
              :userid="onPremisesSamAccountName"
              :user-name="user.displayName"
            />
            <span class="profile-info-wrapper" :class="{ 'oc-py-xs': !user.mail }">
              <span class="oc-display-block" v-text="user.displayName" />
              <span v-if="user.mail" class="oc-text-small" v-text="user.mail" />
              <quota-information
                v-if="quotaEnabled"
                :quota="quota"
                class="oc-text-small oc-mt-xs"
              />
            </span>
          </li>
          <li>
            <oc-button
              id="oc-topbar-account-manage"
              type="router-link"
              :to="accountPageRoute"
              appearance="raw"
            >
              <oc-icon name="settings-4" fill-type="line" class="oc-p-xs" />
              <span v-text="$gettext('Preferences')" />
            </oc-button>
          </li>
          <li v-if="inlineInstances.length > 0">
            <div class="oc-width-1-1" data-testid="instance-switcher">
              <p class="oc-text-xs oc-text-muted">
                {{
                  $pgettext(
                    'The instance switcher section title in the user menu available when multiple instances are enabled in oCIS',
                    'Instances'
                  )
                }}
              </p>
              <oc-list>
                <li
                  v-for="instance in inlineInstances"
                  :key="instance.url"
                  data-testid="instance-switcher-item"
                >
                  <a :href="instance.url" class="instance-link oc-rounded">
                    <span
                      :class="[
                        'instance-indicator',
                        instance.primary && 'primary',
                        instance.active && 'current'
                      ]"
                    />
                    {{ instance.url }}
                  </a>
                </li>
                <li v-if="canOpenInstancesModal">
                  <oc-button
                    appearance="raw"
                    data-testid="instance-switcher-show-all-button"
                    @click="showInstancesModal"
                  >
                    <span class="instance-indicator placeholder" />
                    {{
                      $pgettext(
                        'The open instances modal action label in the user menu available when multiple instances are enabled in oCIS',
                        'Show all instances'
                      )
                    }}
                  </oc-button>
                </li>
              </oc-list>
            </div>
          </li>
          <li>
            <oc-button id="oc-topbar-account-logout" appearance="raw" @click="logout">
              <oc-icon name="logout-box-r" fill-type="line" class="oc-p-xs" />
              <span v-text="$gettext('Log out')" />
            </oc-button>
          </li>
        </template>
      </oc-list>
      <div v-if="footerLinks.length > 0" class="imprint-footer oc-py-s oc-mt-m oc-text-center">
        <template v-for="(link, index) in footerLinks" :key="link.label">
          <span v-if="index > 0"> · </span>
          <oc-button v-if="link.url" type="a" appearance="raw" :href="link.url" target="_blank">
            <span v-text="link.label" />
          </oc-button>
        </template>
      </div>
    </oc-drop>
  </nav>
</template>

<script lang="ts">
import { storeToRefs } from 'pinia'
import { defineComponent, ComponentPublicInstance, computed, unref } from 'vue'
import {
  useRoute,
  useSpacesStore,
  useThemeStore,
  useUserStore,
  routeToContextQuery,
  useAuthService
} from '@ownclouders/web-pkg'
import { OcDrop } from '@ownclouders/design-system/components'
import QuotaInformation from '../Account/QuotaInformation.vue'
import { useGettext } from 'vue3-gettext'
import { useInstances } from '../../composables/instances'

export default defineComponent({
  components: { QuotaInformation },
  setup() {
    const route = useRoute()
    const userStore = useUserStore()
    const themeStore = useThemeStore()
    const spacesStore = useSpacesStore()
    const authService = useAuthService()
    const { $pgettext } = useGettext()
    const { inlineInstances, canOpenInstancesModal, showInstancesModal } = useInstances()

    const { user } = storeToRefs(userStore)

    const accountPageRoute = computed(() => ({
      name: 'account',
      query: routeToContextQuery(unref(route)),
      params: {
        scope: unref(route).params.scope
      }
    }))

    const loginLink = computed(() => {
      return {
        name: 'login',
        query: { redirectUrl: unref(route).fullPath }
      }
    })
    const logout = () => {
      authService.logoutUser()
    }

    const imprintUrl = computed(() => themeStore.currentTheme.common.urls.imprint)
    const privacyUrl = computed(() => themeStore.currentTheme.common.urls.privacy)
    const accessibilityStatementUrl = computed(
      () => themeStore.currentTheme.common.urls.accessibilityStatement
    )

    const quota = computed(() => {
      return spacesStore.personalSpace?.spaceQuota
    })

    const footerLinks = computed(() =>
      [
        unref(imprintUrl) && {
          label: $pgettext(
            'User profile menu: link label; opens the service’s imprint.',
            'Imprint'
          ),
          url: unref(imprintUrl)
        },
        unref(accessibilityStatementUrl) && {
          label: $pgettext(
            'User profile menu: link label; opens the service’s accessibility statement.',
            'Accessibility'
          ),
          url: unref(accessibilityStatementUrl)
        },
        unref(privacyUrl) && {
          label: $pgettext(
            'User profile menu: link label; opens the service’s privacy policy.',
            'Privacy'
          ),
          url: unref(privacyUrl)
        }
      ].filter(Boolean)
    )

    return {
      user,
      accountPageRoute,
      loginLink,
      quota,
      footerLinks,
      inlineInstances,
      logout,
      showInstancesModal,
      canOpenInstancesModal
    }
  },
  computed: {
    onPremisesSamAccountName() {
      return this.user?.onPremisesSamAccountName
    },
    quotaEnabled() {
      return !!this.quota
    }
  },
  mounted() {
    ;(this.$refs.menu as InstanceType<typeof OcDrop>)?.tippy?.setProps({
      onHidden: () => (this.$refs.menuButton as ComponentPublicInstance).$el.focus(),
      onShown: () =>
        (this.$refs.menu as ComponentPublicInstance).$el.querySelector('a:first-of-type').focus()
    })
  }
})
</script>

<style lang="scss" scoped>
.user-menu-list li {
  align-items: center;
  display: flex;
  margin: var(--oc-space-xsmall) 0;

  &:first-child {
    margin-top: 0;
  }

  &:last-child {
    margin-bottom: 0;
  }

  a,
  button {
    gap: var(--oc-space-medium);
    justify-content: flex-start;
    min-height: 3rem;
    padding-left: var(--oc-space-small);
    width: 100%;

    &:focus,
    &:hover {
      background-color: var(--oc-color-background-hover);
      color: var(--oc-color-swatch-passive-default);
      text-decoration: none;
    }
  }

  &.profile-info-wrapper {
    gap: var(--oc-space-medium);
    min-height: 3rem;
  }
}

.imprint-footer {
  background-color: var(--oc-color-background-hover);
  margin-left: calc(var(--oc-space-small) * -1);
  width: calc(100% + var(--oc-space-small) * 2);
  margin-bottom: calc(var(--oc-space-small) * -1) !important;

  a {
    font-size: var(--oc-font-size-medium) !important;
    color: var(--oc-color-text-default);
  }
}

.instance-link {
  align-items: center;
  display: flex;
  gap: var(--oc-space-medium);
  justify-content: flex-start;
  width: 100%;
}

.instance-indicator {
  background-color: var(--oc-color-swatch-passive-default);
  border-radius: 50%;
  height: 0.5rem;
  width: 0.5rem;

  &.primary {
    background-color: var(--oc-color-swatch-primary-default);
  }

  &.current {
    background-color: var(--oc-color-swatch-success-default);
  }

  &.placeholder {
    background-color: transparent;
  }
}
</style>
