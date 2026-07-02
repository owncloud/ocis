<template>
  <oc-button
    id="files-toggle-sidebar"
    v-oc-tooltip="buttonLabel"
    :aria-label="buttonLabel"
    appearance="raw-inverse"
    variation="brand"
    class="oc-my-s"
    @click.stop="toggleSideBar"
  >
    <oc-icon name="side-bar-right" :fill-type="buttonIconFillType" />
  </oc-button>
</template>

<script lang="ts">
import { computed, defineComponent, unref } from 'vue'
import { SideBarEventTopics, useEventBus, useSideBar } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

export default defineComponent({
  name: 'SideBarToggle',
  setup() {
    const eventBus = useEventBus()
    const { $gettext } = useGettext()
    const { isSideBarOpen } = useSideBar({ bus: eventBus })
    const buttonIconFillType = computed(() => {
      return unref(isSideBarOpen) ? 'fill' : 'line'
    })
    const buttonLabel = computed(() => {
      if (unref(isSideBarOpen)) {
        return $gettext('Close sidebar to hide details')
      }
      return $gettext('Open sidebar to view details')
    })
    const toggleSideBar = () => {
      eventBus.publish(SideBarEventTopics.toggle)
    }

    return {
      buttonIconFillType,
      buttonLabel,
      toggleSideBar
    }
  }
})
</script>
