<template>
  <div>
    <context-action-menu :menu-sections="menuSections" :action-options="actionOptions" />
    <label class="oc-invisible-sr" for="space-image-upload-input">
      {{
        $pgettext(
          'Accesibility label to upload a space image in the space context actions',
          'Upload space image'
        )
      }}
    </label>
    <input
      id="space-image-upload-input"
      ref="spaceImageInput"
      type="file"
      name="file"
      multiple
      tabindex="-1"
      :accept="supportedSpaceImageMimeTypes"
      @change="uploadImageSpace"
    />
  </div>
</template>

<script lang="ts" setup>
import { ContextActionMenu, useSpaceActionsNavigateToTrash } from '@ownclouders/web-pkg'
import { useFileActionsShowDetails } from '@ownclouders/web-pkg'
import { useSpaceActionsUploadImage } from '../../composables'
import {
  useSpaceActionsDelete,
  useSpaceActionsDisable,
  useSpaceActionsDuplicate,
  useSpaceActionsEditDescription,
  useSpaceActionsEditQuota,
  useSpaceActionsEditReadmeContent,
  useSpaceActionsRename,
  useSpaceActionsRestore,
  useSpaceActionsShowMembers,
  useSpaceActionsSetIcon
} from '@ownclouders/web-pkg'
import { isLocationSpacesActive } from '@ownclouders/web-pkg'
import { computed, Ref, ref, toRef, unref, VNodeRef } from 'vue'
import { useRouter, usePreviewService } from '@ownclouders/web-pkg'
import { FileActionOptions, SpaceActionOptions } from '@ownclouders/web-pkg'
import { useFileActionsDownloadArchive } from '@ownclouders/web-pkg'

const props = defineProps<{
  actionOptions: SpaceActionOptions
}>()

const router = useRouter()
const previewService = usePreviewService()

const actionOptions = toRef(props, 'actionOptions') as Ref<SpaceActionOptions>

const supportedSpaceImageMimeTypes = computed(() => {
  return previewService.getSupportedMimeTypes('image/').join(',')
})

const { actions: deleteActions } = useSpaceActionsDelete()
const { actions: disableActions } = useSpaceActionsDisable()
const { actions: duplicateActions } = useSpaceActionsDuplicate()
const { actions: editQuotaActions } = useSpaceActionsEditQuota()
const { actions: editReadmeContentActions } = useSpaceActionsEditReadmeContent()
const { actions: editDescriptionActions } = useSpaceActionsEditDescription()
const { actions: setSpaceIconActions } = useSpaceActionsSetIcon()
const { actions: renameActions } = useSpaceActionsRename()
const { actions: restoreActions } = useSpaceActionsRestore()
const { actions: showDetailsActions } = useFileActionsShowDetails()
const { actions: showMembersActions } = useSpaceActionsShowMembers()
const { actions: downloadArchiveActions } = useFileActionsDownloadArchive()
const { actions: navigateToTrashActions } = useSpaceActionsNavigateToTrash()

const spaceImageInput: VNodeRef = ref(null)
const { actions: uploadImageActions, uploadImageSpace } = useSpaceActionsUploadImage({
  spaceImageInput
})

const menuItemsMembers = computed(() => {
  const fileHandlers = [...unref(showMembersActions), ...unref(downloadArchiveActions)]
  // HACK: downloadArchiveActions requires FileActionOptions but we have SpaceActionOptions
  return [...fileHandlers].filter((item) => item.isVisible(unref(actionOptions) as any))
})

const menuItemsPrimaryActions = computed(() => {
  const fileHandlers = [
    ...unref(renameActions),
    ...unref(duplicateActions),
    ...unref(editDescriptionActions),
    ...unref(uploadImageActions),
    ...unref(setSpaceIconActions)
  ]

  if (isLocationSpacesActive(router, 'files-spaces-generic')) {
    fileHandlers.splice(2, 0, ...unref(editReadmeContentActions))
  }
  return [...fileHandlers].filter((item) => item.isVisible(unref(actionOptions)))
})

const menuItemsSecondaryActions = computed(() => {
  const fileHandlers = [
    ...unref(editQuotaActions),
    ...unref(disableActions),
    ...unref(restoreActions),
    ...unref(navigateToTrashActions),
    ...unref(deleteActions)
  ]

  return [...fileHandlers].filter((item) => item.isVisible(unref(actionOptions)))
})

const menuItemsSidebar = computed(() => {
  const fileHandlers = [...unref(showDetailsActions)]
  return [...fileHandlers].filter((item) =>
    // HACK: showDetails provides FileAction[] but we have SpaceActionOptions, so we need to cast them to FileActionOptions
    item.isVisible(unref(actionOptions) as unknown as FileActionOptions)
  )
})

const menuSections = computed(() => {
  const sections = []
  if (unref(menuItemsMembers).length) {
    sections.push({
      name: 'members',
      items: unref(menuItemsMembers)
    })
  }
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

<style lang="scss">
#space-image-upload-input {
  position: absolute;
  left: -99999px;
}
</style>
