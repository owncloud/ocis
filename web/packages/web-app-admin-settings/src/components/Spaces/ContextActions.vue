<template>
  <div>
    <context-action-menu :menu-sections="menuSections" :action-options="{ resources: items }" />
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'
import { ContextActionMenu } from '@ownclouders/web-pkg'

import {
  useSpaceActionsDelete,
  useSpaceActionsDisable,
  useSpaceActionsEditDescription,
  useSpaceActionsEditQuota,
  useSpaceActionsRename,
  useSpaceActionsRestore,
  useActionsShowDetails
} from '@ownclouders/web-pkg'

interface Props {
  items: SpaceResource[]
}

const props = defineProps<Props>()
const filterParams = computed(() => ({ resources: props.items }))

const { actions: deleteActions } = useSpaceActionsDelete()
const { actions: disableActions } = useSpaceActionsDisable()
const { actions: editQuotaActions } = useSpaceActionsEditQuota()
const { actions: editDescriptionActions } = useSpaceActionsEditDescription()
const { actions: renameActions } = useSpaceActionsRename()
const { actions: restoreActions } = useSpaceActionsRestore()
const { actions: showDetailsActions } = useActionsShowDetails()

const menuItemsPrimaryActions = computed(() =>
  [...unref(renameActions), ...unref(editDescriptionActions)].filter((item) =>
    item.isVisible(unref(filterParams))
  )
)
const menuItemsSecondaryActions = computed(() =>
  [
    ...unref(editQuotaActions),
    ...unref(disableActions),
    ...unref(restoreActions),
    ...unref(deleteActions)
  ].filter((item) => item.isVisible(unref(filterParams)))
)
const menuItemsSidebar = computed(() =>
  [...unref(showDetailsActions)].filter((item) => item.isVisible(unref(filterParams)))
)

const menuSections = computed(() => {
  const sections = []

  if (unref(menuItemsPrimaryActions).length) {
    sections.push({
      name: 'primaryActions',
      items: unref(menuItemsPrimaryActions)
    })
  }
  if (unref(menuItemsSecondaryActions).length) {
    sections.push({
      name: 'secondaryActions',
      items: unref(menuItemsSecondaryActions)
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
