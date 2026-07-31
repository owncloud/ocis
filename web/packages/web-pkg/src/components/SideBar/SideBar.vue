<template>
  <div
    id="app-sidebar"
    ref="appSideBar"
    data-testid="app-sidebar"
    :class="{
      'has-active-sub-panel': hasActiveSubPanel,
      'oc-flex oc-flex-center oc-flex-middle': loading,
      'app-sidebar-full-width': fullWidthSideBar
    }"
  >
    <span ref="focusGuard" class="sidebar-focus-guard" tabindex="-1" />
    <oc-spinner v-if="loading" :aria-label="$gettext('Loading sidebar')" />
    <template v-else>
      <div
        v-for="panel in displayPanels"
        :id="`sidebar-panel-${panel.name}`"
        :key="`panel-${panel.name}`"
        :data-testid="`sidebar-panel-${panel.name}`"
        :tabindex="activePanelName === panel.name ? -1 : null"
        class="sidebar-panel"
        :inert="activePanelName !== panel.name"
        :class="{
          'is-root-panel': panel.isRoot?.(panelContext),
          'is-active-sub-panel': hasActiveSubPanel && activeSubPanelName === panel.name, // only one specific sub panel can be active
          'is-active-root-panel': hasActiveRootPanel && panel.isRoot?.(panelContext) // all root panels are active if no sub panel is active
        }"
      >
        <div
          v-if="[activePanelName, oldPanelName].includes(panel.name)"
          class="sidebar-panel__header header"
        >
          <oc-button
            v-if="!panel.isRoot?.(panelContext)"
            v-oc-tooltip="accessibleLabelBack"
            class="header__back"
            appearance="raw"
            :aria-label="accessibleLabelBack"
            @click="closePanel"
          >
            <oc-icon name="arrow-left-s" fill-type="line" />
          </oc-button>

          <h2 class="header__title oc-my-rm" :tabindex="activePanelName === panel.name ? -1 : null">
            {{ panel.title(panelContext) }}
          </h2>

          <oc-button
            appearance="raw"
            class="header__close"
            :aria-label="$gettext('Close file sidebar')"
            @click="closeSidebar"
          >
            <oc-icon name="close" />
          </oc-button>
        </div>

        <div>
          <slot v-if="panel.isRoot?.(panelContext)" name="rootHeader" />
          <slot v-else name="subHeader" />
        </div>
        <div class="sidebar-panel__body" :class="[`sidebar-panel__body-${panel.name}`]">
          <div
            class="sidebar-panel__body-content"
            :class="{ 'sidebar-panel__body-content-stretch': !panel.isRoot?.(panelContext) }"
          >
            <slot name="body">
              <component
                :is="p.component"
                v-for="(p, index) in panel.isRoot?.(panelContext) ? rootPanels : [panel]"
                :key="`sidebar-panel-${p.name}`"
                :class="{ 'multi-root-panel-separator oc-mt oc-pt-s': index > 0 }"
                v-bind="p.componentAttrs?.(panelContext) || {}"
              />
            </slot>
          </div>

          <div
            v-if="panel.isRoot?.(panelContext) && subPanels.length > 0"
            class="sidebar-panel__navigation oc-mt-m"
          >
            <oc-button
              v-for="panelSelect in subPanels"
              :id="`sidebar-panel-${panelSelect.name}-select`"
              :key="`panel-select-${panelSelect.name}`"
              :data-testid="`sidebar-panel-${panelSelect.name}-select`"
              appearance="raw"
              @click="openPanel(panelSelect.name)"
            >
              <oc-icon :name="panelSelect.icon" :fill-type="panelSelect.iconFillType" />
              {{ panelSelect.title(panelContext) }}
              <oc-icon name="arrow-right-s" fill-type="line" />
            </oc-button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import {
  computed,
  nextTick,
  onBeforeUnmount,
  onMounted,
  ref,
  unref,
  useTemplateRef,
  watch
} from 'vue'
import { SideBarPanel, SideBarPanelContext } from './types'
import { useGettext } from 'vue3-gettext'

interface Props {
  isOpen: boolean
  loading: boolean
  availablePanels: SideBarPanel<unknown, unknown, unknown>[]
  panelContext: SideBarPanelContext<unknown, unknown, unknown>
  activePanel?: string
}

interface Emits {
  (e: 'close'): void
  (e: 'selectPanel', value: string | null): void
}

const { isOpen, loading, availablePanels, panelContext, activePanel = '' } = defineProps<Props>()
const emit = defineEmits<Emits>()
const { $gettext } = useGettext()
const appSideBar = useTemplateRef<HTMLElement>('appSideBar')
const focusGuard = useTemplateRef<HTMLElement>('focusGuard')
const focusHandoffPending = ref(false)
const previouslyFocusedElement =
  document.activeElement instanceof HTMLElement ? document.activeElement : null

const rootPanels = computed(() => {
  return availablePanels.filter((p) => p.isVisible(panelContext) && p.isRoot?.(panelContext))
})
const subPanels = computed(() =>
  availablePanels.filter((p) => p.isVisible(panelContext) && !p.isRoot?.(panelContext))
)
const displayPanels = computed<SideBarPanel<unknown, unknown, unknown>[]>(() => {
  if (unref(rootPanels).length) {
    return [unref(rootPanels)[0], ...unref(subPanels)]
  }
  return unref(subPanels)
})

const activeSubPanelName = computed(() => {
  const panelName = activePanel?.split('#')[0]
  if (!panelName) {
    return null
  }
  if (
    !unref(subPanels)
      .map((p) => p.name)
      .includes(panelName)
  ) {
    return null
  }
  return panelName
})
const hasActiveSubPanel = computed(() => {
  return unref(activeSubPanelName) !== null
})
const hasActiveRootPanel = computed(() => {
  return unref(activeSubPanelName) === null
})

const oldPanelName = ref<string>(null)
const setOldPanelName = (name: string) => {
  oldPanelName.value = name
}
const activePanelName = computed<string>(() => {
  if (unref(hasActiveSubPanel)) {
    return unref(activeSubPanelName)
  }
  return unref(rootPanels)[0].name
})

const accessibleLabelBack = computed(() => {
  if (unref(rootPanels).length === 1) {
    return $gettext('Back to %{panel} panel', {
      panel: unref(rootPanels)[0].title(panelContext)
    })
  }
  return $gettext('Back to main panels')
})

const windowWidth = ref(window.innerWidth)

const fullWidthSideBar = computed(() => unref(windowWidth) <= 960)
const backgroundContentEl = computed(() => {
  return unref(appSideBar)?.parentElement?.querySelector('div') as HTMLElement
})

const onResize = () => {
  if (!isOpen) {
    return
  }

  windowWidth.value = window.innerWidth
}

const focusActivePanel = async () => {
  if (!unref(focusHandoffPending) || !isOpen || loading) {
    return
  }

  const focusOrigin = document.activeElement

  // Let the browser apply `inert` and settle native focus before notifying
  // assistive technology about the active panel's focus target.
  await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()))

  const activePanel = Array.from(
    unref(appSideBar)?.querySelectorAll<HTMLElement>('.sidebar-panel') || []
  ).find((panel) => panel.id === `sidebar-panel-${unref(activePanelName)}`)

  if (
    activePanel?.contains(document.activeElement) &&
    document.activeElement !== unref(focusGuard)
  ) {
    focusHandoffPending.value = false
    return
  }

  if (unref(hasActiveSubPanel) && activePanel) {
    const transitionDuration = Number.parseFloat(getComputedStyle(activePanel).transitionDuration)
    if (transitionDuration > 0) {
      await new Promise<void>((resolve) => {
        const finish = () => {
          window.clearTimeout(timeout)
          activePanel.removeEventListener('transitionend', finish)
          resolve()
        }
        const timeout = window.setTimeout(finish, transitionDuration * 1000 + 50)
        activePanel.addEventListener('transitionend', finish, { once: true })
      })
    }
  }

  const focusTarget = unref(hasActiveSubPanel)
    ? activePanel?.querySelector<HTMLElement>('.header__back')
    : activePanel?.querySelector<HTMLElement>('.header__title')
  if (document.activeElement === unref(focusGuard) || document.activeElement === focusOrigin) {
    focusTarget?.focus({ preventScroll: true })
  }
  focusHandoffPending.value = false
}

const requestPanelFocus = () => {
  focusHandoffPending.value = true
  void focusActivePanel()
}

watch(
  () => activePanel,
  () => requestPanelFocus(),
  { flush: 'post' }
)

watch(
  () => loading,
  (isLoading) => {
    if (!isLoading) {
      void focusActivePanel()
    }
  },
  { flush: 'post' }
)

watch(
  () => loading,
  (isLoading) => {
    if (isLoading && unref(appSideBar)?.contains(document.activeElement)) {
      focusHandoffPending.value = true
      unref(focusGuard)?.focus({ preventScroll: true })
    }
  },
  { flush: 'sync' }
)

watch(
  () => isOpen,
  async (isOpen) => {
    if (!isOpen) {
      return
    }
    await nextTick()
    if (unref(fullWidthSideBar) && unref(backgroundContentEl)) {
      // hide content behind sidebar when it has full width to avoid focusable elements
      unref(backgroundContentEl).style.visibility = 'hidden'
    }
  },
  { immediate: true }
)

function setSidebarPanel(panel: string) {
  emit('selectPanel', panel)
}

function resetSidebarPanel() {
  emit('selectPanel', null)
}

function closeSidebar() {
  emit('close')
}

function openPanel(panel: string) {
  focusHandoffPending.value = true
  unref(focusGuard)?.focus({ preventScroll: true })
  setOldPanelName(unref(activePanelName))
  setSidebarPanel(panel)
}

function closePanel() {
  focusHandoffPending.value = true
  setOldPanelName(unref(activePanelName))
  resetSidebarPanel()
  unref(focusGuard)?.focus({ preventScroll: true })
}

onMounted(() => {
  window.addEventListener('resize', onResize)
  requestPanelFocus()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)

  if (unref(backgroundContentEl)) {
    unref(backgroundContentEl).style.visibility = 'visible'
  }

  nextTick(() => {
    if (previouslyFocusedElement?.isConnected) {
      previouslyFocusedElement.focus()
    }
  })
})
</script>

<style lang="scss">
#app-sidebar {
  border-left: 1px solid var(--oc-color-border);
  position: relative;
  overflow: hidden;
  min-width: 440px;
  width: 440px;
}

.sidebar-focus-guard {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip-path: inset(50%);
}
.app-sidebar-full-width {
  min-width: 100% !important;
  width: 100% !important;
}

@media only screen and (max-width: $oc-breakpoint-medium-default) {
  .files-wrapper {
    flex-wrap: nowrap !important;
  }
}

.sidebar-panel {
  $root: &;
  overflow: hidden;
  width: 100%;
  max-width: 100%;
  height: 100%;
  max-height: 100%;
  display: grid;
  grid-template-rows: auto auto 1fr;
  background-color: var(--oc-color-background-sidebar, var(--oc-color-background-default));
  top: 0;
  position: absolute;
  transform: translateX(100%);
  transition:
    transform 0.4s ease,
    visibility 0.4s 0s;
  // visibility is here to prevent focusing panel child elements,
  // the transition delay keeps care that it will only apply if the element is visible or not.
  // hidden: if element is off screen
  // visible: if element is on screen
  visibility: hidden;
  border-top-right-radius: var(--oc-space-medium);
  border-bottom-right-radius: var(--oc-space-medium);

  @media screen and (prefers-reduced-motion: reduce), (update: slow) {
    transition-duration: 0.001ms !important;
  }

  &.is-active-root-panel,
  &.is-active-sub-panel {
    visibility: unset;
    transform: translateX(0);
  }

  &.is-active-root-panel {
    right: 0 !important;
    transition: right 0.4s 0s;
  }

  &.is-root-panel {
    transform: translateX(0);
    visibility: visible;
    transition: right 0.4s 0s;
    right: 100px;
  }

  .multi-root-panel-separator {
    border-top: 1px solid var(--oc-color-border);
  }

  &__header {
    padding: var(--oc-space-small) var(--oc-space-small) 0 var(--oc-space-small);

    &.header {
      display: grid;
      grid-template-columns: auto 1fr auto;
      align-items: center;
    }

    & .header {
      &__back {
        grid-column-start: 1;
      }

      &__title {
        text-align: center;
        color: var(--oc-color-text-default);
        font-size: var(--oc-font-size-large);
        grid-column-start: 2;
      }

      &__close {
        grid-column-start: 3;
      }
    }
  }

  &__body {
    overflow-y: auto;
    overflow-x: hidden;
    padding: var(--oc-space-small);
    display: flex;
    flex-direction: column;

    &-content-stretch {
      flex: 1;
    }
  }

  &__navigation {
    margin: var(--oc-space-small) - var(--oc-space-small) - var(--oc-space-small);

    > button {
      border-bottom: 1px solid var(--oc-color-border);
      width: 100%;
      border-radius: 0;
      color: var(--oc-color-text-default) !important;
      display: grid;
      grid-template-columns: auto 1fr auto;
      text-align: left;
      height: 50px;
      padding: 0 var(--oc-space-small);

      &:first-of-type {
        border-top: 1px solid var(--oc-color-border);
      }

      &:last-of-type {
        border-bottom: 0;
      }

      &:hover,
      &:focus {
        border-color: var(--oc-color-border) !important;
      }

      &:hover {
        background-color: var(--oc-color-background-muted) !important;
      }
    }
  }
}
</style>
