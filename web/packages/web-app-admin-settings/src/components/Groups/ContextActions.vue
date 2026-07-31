<template>
  <div>
    <context-action-menu :menu-sections="menuSections" :action-options="actionOptions" />
  </div>
</template>

<script lang="ts" setup>
import { useActionsShowDetails } from '@ownclouders/web-pkg'
import { computed, unref } from 'vue'
import { ContextActionMenu } from '@ownclouders/web-pkg'
import { GroupActionOptions } from '@ownclouders/web-pkg'
import { useGroupActionsEdit, useGroupActionsDelete } from '../../composables/actions/groups'

interface Props {
  actionOptions: GroupActionOptions
}

const props = defineProps<Props>()
const { actions: showDetailsActions } = useActionsShowDetails()
const { actions: deleteActions } = useGroupActionsDelete()
const { actions: editActions } = useGroupActionsEdit()

const menuItemsPrimaryActions = computed(() =>
  [...unref(editActions), ...unref(deleteActions)].filter((item) =>
    item.isVisible(props.actionOptions)
  )
)

const menuItemsSidebar = computed(() =>
  [...unref(showDetailsActions)].filter((item) => item.isVisible(props.actionOptions))
)

const menuSections = computed(() => {
  const sections = []

  if (unref(menuItemsPrimaryActions).length) {
    sections.push({
      name: 'primaryActions',
      items: unref(menuItemsPrimaryActions)
    })
  }
  if (unref(menuItemsSidebar).length) {
    sections.push({
      name: 'sidebar',
      items: unref(menuItemsSidebar)
    })
  }
  return sections
})
</script>
