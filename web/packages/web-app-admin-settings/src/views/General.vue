<template>
  <div>
    <app-template
      ref="template"
      :breadcrumbs="breadcrumbs"
      :show-app-bar="false"
      :is-side-bar-open="isSideBarOpen"
      :side-bar-active-panel="sideBarActivePanel"
      :side-bar-available-panels="sideBarAvailablePanels"
    >
      <template #mainContent>
        <div class="oc-px-m">
          <InfoSection />
          <AppearanceSection />
        </div>
      </template>
    </app-template>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import AppTemplate from '../components/AppTemplate.vue'
import InfoSection from '../components/General/InfoSection.vue'
import AppearanceSection from '../components/General/AppearanceSection.vue'
import DetailsPanel from '../components/General/SideBar/DetailsPanel.vue'
import { useGettext } from 'vue3-gettext'
import { useSideBar } from '@ownclouders/web-pkg'

const template = ref()
const { $gettext } = useGettext()

const sideBarAvailablePanels = [
  {
    name: 'DetailsPanel',
    icon: 'settings-4',
    title: () => $gettext('Details'),
    component: DetailsPanel,
    isRoot: () => true,
    isVisible: () => true
  }
]
const { isSideBarOpen, sideBarActivePanel } = useSideBar()
const breadcrumbs = computed(() => {
  return [
    { text: $gettext('Administration Settings'), to: { path: '/admin-settings' } },
    {
      text: $gettext('General')
    }
  ]
})
</script>
