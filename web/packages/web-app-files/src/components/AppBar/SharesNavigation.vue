<template>
  <nav id="shares-navigation" class="oc-py-s" :aria-label="$gettext('Shares pages navigation')">
    <oc-list class="oc-flex oc-visible@s">
      <li v-for="navItem in navItems" :key="`shares-navigation-desktop-${navItem.id}`">
        <oc-button
          type="router-link"
          class="oc-mr-m oc-py-s shares-nav-desktop"
          appearance="raw"
          :to="navItem.to"
        >
          <oc-icon size="small" :name="navItem.icon" />
          <span v-text="navItem.text" />
        </oc-button>
      </li>
    </oc-list>
    <div class="oc-hidden@s">
      <oc-button id="shares_navigation_mobile" appearance="raw">
        <span v-text="currentNavItem.text" />
        <oc-icon name="arrow-down-s" fill-type="line" size="small" />
      </oc-button>
      <oc-drop toggle="#shares_navigation_mobile" mode="click" close-on-click padding-size="small">
        <oc-list>
          <li v-for="navItem in navItems" :key="`shares-navigation-mobile-${navItem.id}`">
            <oc-button
              type="router-link"
              class="oc-my-xs shares-nav-mobile"
              :to="navItem.to"
              :class="{ 'oc-background-primary-gradient': navItem.active }"
              :appearance="navItem.active ? 'raw-inverse' : 'raw'"
              :variation="navItem.active ? 'primary' : 'passive'"
            >
              <span class="icon-box" :class="{ 'icon-box-active': navItem.active }">
                <oc-icon :name="navItem.icon" />
              </span>
              <span v-text="navItem.text" />
            </oc-button>
          </li>
        </oc-list>
      </oc-drop>
    </div>
  </nav>
</template>

<script lang="ts" setup>
import { isLocationSharesActive, RouteShareTypes } from '@ownclouders/web-pkg'
import {
  locationSharesViaLink,
  locationSharesWithMe,
  locationSharesWithOthers
} from '@ownclouders/web-pkg'

import { computed, unref } from 'vue'
import { useRouter } from '@ownclouders/web-pkg'
import { useActiveLocation } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { RouteLocationRaw } from 'vue-router'

const { $gettext } = useGettext()
const router = useRouter()

const resolveScopeTemplatePath = (path: string, scopePrefix: string) =>
  path.replace(/^\/:scope\(vault\)\?/, scopePrefix)

const locationToPath = (location: RouteLocationRaw) => {
  const scope = unref(router.currentRoute).params?.scope
  const scopePrefix = scope === 'vault' ? '/vault' : ''

  if (typeof location === 'string') {
    return location
  }

  const locationWithScope: RouteLocationRaw = {
    ...location,
    ...(scope && {
      params: {
        ...((location as { params?: Record<string, unknown> }).params || {}),
        scope
      }
    })
  }

  const resolvedPath = router.resolve(locationWithScope).path
  if (resolvedPath) {
    return resolvedPath
  }

  const routeName = (location as { name?: string }).name
  const route = routeName ? router.getRoutes().find((r) => r.name === routeName) : undefined
  if (route?.path) {
    return resolveScopeTemplatePath(route.path, scopePrefix)
  }

  return ''
}
const sharesWithMeActive = useActiveLocation(
  isLocationSharesActive,
  locationSharesWithMe.name as RouteShareTypes
)
const sharesWithOthersActive = useActiveLocation(
  isLocationSharesActive,
  locationSharesWithOthers.name as RouteShareTypes
)
const sharesViaLinkActive = useActiveLocation(
  isLocationSharesActive,
  locationSharesViaLink.name as RouteShareTypes
)
const navItems = computed(() => [
  {
    id: locationSharesWithMe.name as string,
    icon: 'share-forward',
    to: locationToPath(locationSharesWithMe),
    text: $gettext('Shared with me'),
    active: unref(sharesWithMeActive)
  },
  {
    id: locationSharesWithOthers.name as string,
    icon: 'reply',
    to: locationToPath(locationSharesWithOthers),
    text: $gettext('Shared with others'),
    active: unref(sharesWithOthersActive)
  },
  {
    id: locationSharesViaLink.name as string,
    icon: 'link',
    to: locationToPath(locationSharesViaLink),
    text: $gettext('Shared via link'),
    active: unref(sharesViaLinkActive)
  }
])
const currentNavItem = computed(() => unref(navItems).find((navItem) => navItem.active))
</script>
<style lang="scss" scoped>
#shares-navigation {
  a {
    gap: var(--oc-space-medium);
    width: 100%;

    &:focus,
    &:hover {
      text-decoration: none;
    }

    &.shares-nav-mobile {
      justify-content: flex-start;
    }

    .icon-box {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      width: 40px;
      height: 40px;
    }
    .icon-box-active {
      box-shadow: 2px 0 6px rgba(0, 0, 0, 0.14);
    }
  }

  .shares-nav-desktop.router-link-active {
    border-bottom: 2px solid var(--oc-color-swatch-primary-default) !important;
    border-radius: 0;
  }
}
</style>
