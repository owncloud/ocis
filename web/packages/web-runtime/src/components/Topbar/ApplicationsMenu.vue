<template>
  <nav
    id="applications-menu"
    :aria-label="$gettext('Main navigation')"
    class="oc-flex oc-flex-middle"
  >
    <oc-button
      id="_appSwitcherButton"
      v-oc-tooltip="applicationSwitcherLabel"
      appearance="raw-inverse"
      variation="brand"
      class="oc-topbar-menu-burger"
      :aria-label="applicationSwitcherLabel"
    >
      <oc-icon name="grid" size="large" class="oc-flex" />
    </oc-button>
    <oc-drop
      ref="menu"
      drop-id="app-switcher-dropdown"
      toggle="#_appSwitcherButton"
      mode="click"
      padding-size="small"
      close-on-click
      @show-drop="updateAppIcons"
    >
      <div class="oc-display-block oc-position-relative">
        <oc-list class="applications-list">
          <li v-for="(n, nid) in sortedMenuItems" :key="`apps-menu-${nid}`">
            <oc-button
              :key="n.url ? 'apps-menu-external-link' : 'apps-menu-internal-link'"
              :appearance="isMenuItemActive(n) ? 'raw-inverse' : 'raw'"
              :variation="isMenuItemActive(n) ? 'primary' : 'passive'"
              :class="{ 'oc-background-primary-gradient router-link-active': isMenuItemActive(n) }"
              :data-test-id="n.id"
              v-bind="getAdditionalAttributes(n)"
              v-on="getAdditionalEventBindings(n)"
            >
              <oc-application-icon
                :key="`apps-menu-icon-${nid}-${appIconKey}`"
                :icon="n.icon || 'link'"
                :color-primary="n.color"
              />
              <span v-text="$gettext(n.label())" />
            </oc-button>
          </li>
        </oc-list>
      </div>
    </oc-drop>
  </nav>
</template>

<script lang="ts">
import { defineComponent, PropType, ref, computed, unref, useTemplateRef } from 'vue'
import { OcDrop, OcApplicationIcon } from '@ownclouders/design-system/components'
import { useGettext } from 'vue3-gettext'
import { v4 as uuidV4 } from 'uuid'
import { AppMenuItemExtension, useRouter } from '@ownclouders/web-pkg'

export default defineComponent({
  components: {
    OcApplicationIcon
  },
  props: {
    menuItems: {
      type: Array as PropType<AppMenuItemExtension[]>,
      required: false,
      default: (): AppMenuItemExtension[] => []
    }
  },
  setup(props) {
    const router = useRouter()
    const { $gettext } = useGettext()
    const appIconKey = ref('')

    const menu = useTemplateRef<typeof OcDrop>('menu')

    const activeRoutePath = computed(() => unref(router.currentRoute).path)
    const sortedMenuItems = computed(() => {
      return [...props.menuItems].sort(
        (a, b) => (a.priority || Number.MAX_SAFE_INTEGER) - (b.priority || Number.MAX_SAFE_INTEGER)
      )
    })

    const applicationSwitcherLabel = computed(() => {
      return $gettext('Application Switcher')
    })
    const updateAppIcons = () => {
      appIconKey.value = uuidV4().replaceAll('-', '')
    }

    const getAdditionalEventBindings = (item: AppMenuItemExtension) => {
      return item.handler ? { click: item.handler } : {}
    }

    const getScopeAwarePath = (rawPath?: string) => {
      if (!rawPath) {
        return rawPath
      }

      const normalizedPath = rawPath.startsWith('/') ? rawPath : `/${rawPath}`
      const pathWithoutVault = normalizedPath.replace(/^\/vault(?=\/|$)/, '')
      const isVaultScope = unref(router.currentRoute).params?.scope === 'vault'

      return isVaultScope ? `/vault${pathWithoutVault}` : pathWithoutVault
    }

    const getAdditionalAttributes = (item: AppMenuItemExtension) => {
      let type: string
      if (item.handler) {
        type = 'button'
      } else if (item.path) {
        type = 'router-link'
      } else {
        type = 'a'
      }

      return {
        type,
        ...(type === 'router-link' && { to: getScopeAwarePath(item.path) }),
        ...(type === 'a' && { href: item.url, target: '_blank' })
      }
    }
    const isMenuItemActive = (item: AppMenuItemExtension) => {
      const routePath = getScopeAwarePath(item.path)
      return !!routePath && unref(activeRoutePath)?.startsWith(routePath)
    }

    return {
      menu,
      sortedMenuItems,
      appIconKey,
      updateAppIcons,
      applicationSwitcherLabel,
      getAdditionalAttributes,
      getAdditionalEventBindings,
      isMenuItemActive
    }
  },
  mounted() {
    this.menu?.tippy?.setProps({
      onShown: () => this.menu.$el.querySelector('a:first-of-type').focus()
    })
  }
})
</script>

<style lang="scss" scoped>
.oc-drop {
  width: 280px;
}

.applications-list li {
  margin: var(--oc-space-xsmall) 0;

  &:first-child {
    margin-top: 0;
  }

  &:last-child {
    margin-bottom: 0;
  }

  a,
  button {
    padding: 5px;
    border-radius: 8px;
    gap: var(--oc-space-medium);
    justify-content: flex-start;
    width: 100%;

    &.oc-button-primary-raw-inverse {
      &:focus,
      &:hover {
        color: var(--oc-color-swatch-primary-contrast) !important;
      }
    }

    &.oc-button-passive-raw {
      &:focus,
      &:hover {
        color: var(--oc-color-swatch-passive-default) !important;
      }
    }

    &:focus {
      text-decoration: none;
    }

    &:hover {
      background-color: var(--oc-color-background-hover);
      text-decoration: none;
    }

    .icon-box {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      width: 40px;
      height: 40px;
    }

    .active-check {
      position: absolute;
      right: 1rem;
    }
  }
}
</style>
