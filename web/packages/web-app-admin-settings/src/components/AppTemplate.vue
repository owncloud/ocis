<template>
  <main class="oc-flex oc-height-1-1 app-content oc-width-1-1">
    <app-loading-spinner v-if="loading" />
    <template v-else>
      <div id="admin-settings-wrapper" class="oc-width-expand oc-height-1-1 oc-position-relative">
        <div
          id="admin-settings-app-bar"
          ref="appBarRef"
          class="oc-app-bar oc-py-s"
          :class="{ 'admin-settings-app-bar-sticky': isSticky }"
        >
          <div class="admin-settings-app-bar-controls oc-flex oc-flex-between oc-flex-middle">
            <oc-breadcrumb
              v-if="!isMobileWidth"
              id="admin-settings-breadcrumb"
              class="oc-flex oc-flex-middle"
              :items="breadcrumbs"
            />
            <portal-target name="app.runtime.mobile.nav" />
            <div class="oc-flex">
              <view-options
                v-if="showViewOptions"
                :has-hidden-files="false"
                :has-file-extensions="false"
                :has-pagination="true"
                :should-show-flat-list-toggle="true"
                :pagination-options="paginationOptions"
                :per-page-default="perPageDefault"
                per-page-storage-prefix="admin-settings"
              />
            </div>
          </div>
          <div
            v-if="showAppBar"
            class="admin-settings-app-bar-actions oc-flex oc-flex-middle oc-mt-xs"
          >
            <slot
              name="topbarActions"
              :limited-screen-space="limitedScreenSpace"
              class="oc-flex-1 oc-flex oc-flex-start"
            />
            <batch-actions-component
              v-if="showBatchActions"
              :actions="batchActions"
              :action-options="{ resources: batchActionItems }"
              :limited-screen-space="limitedScreenSpace"
            />
          </div>
        </div>
        <slot name="mainContent" />
      </div>
      <side-bar
        v-if="isSideBarOpen"
        :active-panel="sideBarActivePanel"
        :available-panels="sideBarAvailablePanels"
        :panel-context="sideBarPanelContext"
        :loading="sideBarLoading"
        :is-open="isSideBarOpen"
        @select-panel="selectPanel"
        @close="closeSideBar"
      >
        <template #header>
          <slot name="sideBarHeader" />
        </template>
      </side-bar>
    </template>
  </main>
</template>

<script lang="ts" setup>
import { perPageDefault, paginationOptions } from '../defaults'
import {
  AppLoadingSpinner,
  SideBar,
  BatchActions as BatchActionsComponent,
  SideBarPanelContext,
  Action,
  useIsTopBarSticky
} from '@ownclouders/web-pkg'
import { inject, onBeforeUnmount, Ref, ref, unref, VNodeRef, watch } from 'vue'
import { eventBus } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { SideBarPanel } from '@ownclouders/web-pkg'
import { BreadcrumbItem } from '@ownclouders/design-system/helpers'
import { ViewOptions } from '@ownclouders/web-pkg'
import { Item } from '@ownclouders/web-client'

interface Props {
  breadcrumbs: BreadcrumbItem[]
  isSideBarOpen?: boolean
  sideBarAvailablePanels?: SideBarPanel<unknown, unknown, unknown>[]
  sideBarPanelContext?: SideBarPanelContext<unknown, unknown, unknown>
  sideBarActivePanel?: string | null
  loading?: boolean
  sideBarLoading?: boolean
  showViewOptions?: boolean
  showBatchActions?: boolean
  batchActionItems?: Item[]
  batchActions?: Action[]
  showAppBar?: boolean
}
const {
  breadcrumbs,
  isSideBarOpen = false,
  sideBarAvailablePanels = [],
  sideBarPanelContext = {},
  sideBarActivePanel = null,
  loading = false,
  sideBarLoading = false,
  showViewOptions = false,
  showBatchActions = false,
  batchActionItems = [],
  batchActions = [],
  showAppBar = true
} = defineProps<Props>()

const isMobileWidth = inject<Ref<boolean>>('isMobileWidth')
const appBarRef = ref<VNodeRef>()
const limitedScreenSpace = ref(false)
const { isSticky } = useIsTopBarSticky()

const onResize = () => {
  limitedScreenSpace.value = isSideBarOpen ? window.innerWidth <= 1600 : window.innerWidth <= 1200
}
const resizeObserver = new ResizeObserver(onResize)

const closeSideBar = () => {
  eventBus.publish(SideBarEventTopics.close)
}
const selectPanel = (panel: string) => {
  eventBus.publish(SideBarEventTopics.setActivePanel, panel)
}

watch(
  appBarRef,
  (ref) => {
    if (ref) {
      resizeObserver.observe(unref(appBarRef) as unknown as HTMLElement)
    }
  },
  { immediate: true }
)

onBeforeUnmount(() => {
  if (unref(appBarRef)) {
    resizeObserver.unobserve(unref(appBarRef) as unknown as HTMLElement)
  }
})
</script>

<style lang="scss">
#admin-settings-wrapper {
  overflow-y: auto;
}

#admin-settings-app-bar {
  background-color: var(--oc-color-background-default);
  border-top-right-radius: 15px;
  box-sizing: border-box;
  z-index: 2;
  position: inherit;
  padding: 0 var(--oc-space-medium);
  top: 0;

  &.admin-settings-app-bar-sticky {
    position: sticky;
  }
}

.admin-settings-app-bar-controls {
  height: 52px;

  @media (max-width: $oc-breakpoint-xsmall-max) {
    justify-content: space-between;
  }
}

.admin-settings-app-bar-actions {
  align-items: center;
  display: flex;
  min-height: 3rem;
}

#admin-settings-wrapper {
  overflow-y: auto;
}
</style>
