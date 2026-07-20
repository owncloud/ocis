<template>
  <aside
    id="web-nav-sidebar"
    :aria-label="$gettext('Application navigation')"
    :class="{
      'oc-app-navigation-collapsed': closed,
      'oc-app-navigation-expanded': !closed
    }"
  >
    <oc-button
      appearance="raw"
      :class="toggleSidebarButtonClass"
      class="toggle-sidebar-button oc-pb-s oc-pt-m"
      :aria-label="$gettext('Toggle sidebar')"
      :aria-expanded="!closed"
      @click="$emit('update:nav-bar-closed', !closed)"
    >
      <oc-icon size="large" fill-type="line" :name="toggleSidebarButtonIcon" />
    </oc-button>
    <nav
      class="oc-sidebar-nav oc-mb-m oc-mt-s oc-px-xs"
      :aria-label="$gettext('Sidebar navigation menu')"
    >
      <div
        v-show="isAnyNavItemActive"
        id="nav-highlighter"
        class="oc-ml-s oc-background-primary-gradient"
        v-bind="highlighterAttrs"
        :aria-hidden="true"
      />
      <oc-list>
        <sidebar-nav-item
          v-for="(link, index) in navItems"
          :ref="(el) => (navItemRefs[index] = el as NavItemRef)"
          :key="index"
          :index="getUuid()"
          :target="link.route"
          :active="link.active"
          :icon="link.icon"
          :fill-type="link.fillType"
          :name="link.name"
          :collapsed="closed"
          :handler="link.handler"
        />
      </oc-list>
    </nav>
    <!-- @slot bottom content of the sidebar -->
    <slot name="bottom">
      <div v-if="!closed" class="versions oc-pb-m oc-pl-m oc-text-xsmall oc-text-muted">
        <span v-text="backendVersion" />
        <span v-text="webVersion" />
      </div>
    </slot>
  </aside>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  nextTick,
  onBeforeUnmount,
  onMounted,
  PropType,
  ref,
  unref,
  watch
} from 'vue'
import { v4 as uuidV4 } from 'uuid'
import SidebarNavItem from './SidebarNavItem.vue'
import { NavItem } from '../../helpers/navItems'
import { useCapabilityStore, getBackendVersion, getWebVersion } from '@ownclouders/web-pkg'

type NavItemRef = InstanceType<typeof SidebarNavItem>

export default defineComponent({
  components: {
    SidebarNavItem
  },
  props: {
    navItems: {
      type: Array as PropType<NavItem[]>,
      required: true
    },
    closed: { type: Boolean, default: false }
  },
  emits: ['update:nav-bar-closed'],
  setup(props) {
    let resizeObserver: ResizeObserver
    const navItemRefs = ref<Record<string, NavItemRef>>({})
    const highlighterAttrs = ref<Record<string, unknown>>({})
    const capabilityStore = useCapabilityStore()

    const webVersion = computed(() => getWebVersion())
    const backendVersion = computed(() => getBackendVersion({ capabilityStore }))

    onMounted(() => {
      const navBar = document.getElementById('web-nav-sidebar')
      const highlighter = document.getElementById('nav-highlighter')

      if (!highlighter || !navBar) {
        return
      }

      resizeObserver = new ResizeObserver(() => {
        const navItem = document.getElementsByClassName('oc-sidebar-nav-item-link')[0]
        if (!navItem) {
          return
        }
        highlighter.style.setProperty('transition-duration', `0.05s`)
        highlighter.style.setProperty('width', `${navItem.clientWidth}px`)
        highlighter.style.setProperty('height', `${navItem.clientHeight}px`)
      })
      resizeObserver.observe(navBar)
    })

    onBeforeUnmount(() => {
      resizeObserver.disconnect()
    })

    const updateHighlighterPosition = () => {
      const activeItemIndex = props.navItems.findIndex((n) => n.active)
      const activeEl = unref(navItemRefs)[activeItemIndex]
      if (activeEl) {
        highlighterAttrs.value = {
          style: {
            transform: `translateY(${activeEl.$el.offsetTop}px)`,
            'transition-duration': '0.2s'
          }
        }
      }
    }

    watch(
      () => props.navItems,
      async () => {
        await nextTick()
        updateHighlighterPosition()
      },
      { deep: true, immediate: true }
    )

    return { highlighterAttrs, navItemRefs, backendVersion, webVersion }
  },
  computed: {
    toggleSidebarButtonClass() {
      return this.closed
        ? 'toggle-sidebar-button-collapsed'
        : 'toggle-sidebar-button-expanded oc-pr-s'
    },

    toggleSidebarButtonIcon() {
      return this.closed ? 'arrow-drop-right' : 'arrow-drop-left'
    },

    isAnyNavItemActive() {
      return this.navItems.some((i) => i.active === true)
    }
  },
  methods: {
    getUuid() {
      return uuidV4().replaceAll('-', '')
    }
  }
})
</script>

<style lang="scss">
#nav-highlighter {
  position: absolute;
  border-radius: 5px;
  transition: transform 0.2s cubic-bezier(0.51, 0.06, 0.56, 1.37);
}

#web-nav-sidebar {
  background-color: var(--oc-color-background-sidebar, var(--oc-color-background-default));
  border-radius: 15px 0 0 15px;
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: all 0.35s cubic-bezier(0.34, 0.11, 0, 1.12);
  z-index: 4;

  .oc-list {
    position: relative;
  }

  .toggle-sidebar-button {
    min-height: 3rem;
    transition: all 0.2s ease-out;

    &:hover {
      overflow: hidden;
    }
    &:focus-visible {
      outline: none !important;
      box-shadow: inset 0 0 0 2px var(--oc-color-swatch-passive-default) !important;
    }
  }

  .toggle-sidebar-button-expanded {
    justify-content: flex-end !important;
  }

  .versions {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    justify-content: flex-end;
    flex-grow: 1;
  }

  .oc-sidebar-nav li a:not(.active),
  .oc-sidebar-nav li button:not(.active) {
    &:hover,
    &:focus {
      text-decoration: none !important;
      background-color: var(--oc-color-background-hover);
      color: var(--oc-color-swatch-passive-default);
    }
  }

  .oc-sidebar-nav li a.active,
  .oc-sidebar-nav li button.active {
    &:focus,
    &:hover {
      color: var(--oc-color-swatch-primary-contrast);
    }
  }
}

.oc-app-navigation-expanded {
  min-width: 230px !important;
  max-width: 230px !important;
}

.oc-app-navigation-collapsed {
  min-width: 62px !important;
  max-width: 62px !important;
}
</style>
