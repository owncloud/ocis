<template>
  <context-action-menu :menu-sections="menuSections" :action-options="actionOptions" />
</template>

<script lang="ts" setup>
import ContextActionMenu from '../ContextActions/ContextActionMenu.vue'
import { computed, Ref, toRef, unref } from 'vue'
import {
  ActionExtension,
  FileActionOptions,
  useExtensionRegistry,
  useFileActionsToggleHideShare,
  useFileActionsCopyPermanentLink,
  useFileActionsPaste,
  useFileActionsShowDetails,
  useFileActionsShowShares,
  useFileActionsEnableSync,
  useFileActionsCopy,
  useFileActionsDisableSync,
  useFileActionsDelete,
  useFileActionsDownloadArchive,
  useFileActionsEmptyTrashBin,
  useFileActionsMove,
  useFileActionsRestore,
  useFileActionsDownloadFile,
  useFileActionsRename,
  useFileActionsSetImage,
  useFileActionsNavigate,
  useFileActionsFavorite,
  useFileActionsCreateSpaceFromResource,
  useFileActions,
  useFileActionsDuplicate
} from '../../composables'
import { isNil } from 'lodash-es'

interface Props {
  actionOptions: FileActionOptions
}

const props = defineProps<Props>()
const { editorActions } = useFileActions()

const { actions: enableSyncActions } = useFileActionsEnableSync()
const { actions: hideShareActions } = useFileActionsToggleHideShare()
const { actions: copyActions } = useFileActionsCopy()
const { actions: copyPermanentLinkActions } = useFileActionsCopyPermanentLink()
const { actions: disableSyncActions } = useFileActionsDisableSync()
const { actions: deleteActions } = useFileActionsDelete()
const { actions: downloadArchiveActions } = useFileActionsDownloadArchive()
const { actions: downloadFileActions } = useFileActionsDownloadFile()
const { actions: favoriteActions } = useFileActionsFavorite()
const { actions: emptyTrashBinActions } = useFileActionsEmptyTrashBin()
const { actions: moveActions } = useFileActionsMove()
const { actions: navigateActions } = useFileActionsNavigate()
const { actions: pasteActions } = useFileActionsPaste()
const { actions: renameActions } = useFileActionsRename()
const { actions: restoreActions } = useFileActionsRestore()
const { actions: setSpaceImageActions } = useFileActionsSetImage()
const { actions: showDetailsActions } = useFileActionsShowDetails()
const { actions: createSpaceFromResourceActions } = useFileActionsCreateSpaceFromResource()
const { actions: showSharesActions } = useFileActionsShowShares()
const { actions: duplicateActions } = useFileActionsDuplicate()

const extensionRegistry = useExtensionRegistry()
const extensionsContextActions = computed(() => {
  return extensionRegistry
    .requestExtensions<ActionExtension>({
      id: 'global.files.context-actions',
      extensionType: 'action'
    })
    .map((e) => e.action)
})
const extensionsBatchActions = computed(() => {
  return extensionRegistry
    .requestExtensions<ActionExtension>({
      id: 'global.files.batch-actions',
      extensionType: 'action'
    })
    .map((e) => e.action)
})

// type cast to make vue-tsc aware of the type
const actionOptions = toRef(props, 'actionOptions') as Ref<FileActionOptions>

const menuItemsBatchActions = computed(() =>
  [
    ...unref(enableSyncActions),
    ...unref(disableSyncActions),
    ...unref(downloadArchiveActions),
    ...unref(moveActions),
    ...unref(copyActions),
    ...unref(duplicateActions),
    ...unref(emptyTrashBinActions),
    ...unref(deleteActions),
    ...unref(restoreActions),
    ...unref(createSpaceFromResourceActions),
    ...unref(extensionsBatchActions).filter((a) => a.category === 'actions' || isNil(a.category))
  ].filter((item) => item.isVisible(unref(actionOptions)))
)
const menuItemsBatchSideBar = computed(() =>
  [
    ...unref(showDetailsActions),
    ...unref(extensionsBatchActions).filter((a) => a.category === 'sidebar')
  ].filter((item) => item.isVisible(unref(actionOptions)))
)

const menuItemsContext = computed(() => {
  return [
    ...unref(editorActions),
    ...unref(extensionsContextActions).filter((a) => a.category === 'context')
  ]
    .filter((item) => item.isVisible(unref(actionOptions)))
    .sort((x, y) => Number(y.hasPriority) - Number(x.hasPriority))
})

const menuItemsShare = computed(() => {
  return [
    ...unref(showSharesActions),
    ...unref(copyPermanentLinkActions),
    ...unref(extensionsContextActions).filter((a) => a.category === 'share')
  ].filter((item) => item.isVisible(unref(actionOptions)))
})

const menuItemsActions = computed(() => {
  return [
    ...unref(downloadArchiveActions),
    ...unref(downloadFileActions),
    ...unref(deleteActions),
    ...unref(moveActions),
    ...unref(copyActions),
    ...unref(pasteActions),
    ...unref(renameActions),
    ...unref(duplicateActions),
    ...unref(createSpaceFromResourceActions),
    ...unref(restoreActions),
    ...unref(enableSyncActions),
    ...unref(disableSyncActions),
    ...unref(hideShareActions),
    ...unref(setSpaceImageActions),
    ...unref(extensionsContextActions).filter((a) => a.category === 'actions' || isNil(a.category))
  ].filter((item) => item.isVisible(unref(actionOptions)))
})

const menuItemsSidebar = computed(() => {
  return [
    ...unref(favoriteActions).map((action) => {
      action.keepOpen = true
      return action
    }),
    ...unref(navigateActions),
    ...unref(showDetailsActions),
    ...unref(extensionsContextActions).filter((a) => a.category === 'sidebar')
  ].filter((item) => item.isVisible(unref(actionOptions)))
})

const menuSections = computed(() => {
  const sections = []
  if (unref(actionOptions).resources.length > 1) {
    if (unref(menuItemsBatchActions).length) {
      sections.push({
        name: 'batch-actions',
        items: [...unref(menuItemsBatchActions)]
      })
    }
    sections.push({
      name: 'batch-details',
      items: [...unref(menuItemsBatchSideBar)]
    })
    return sections
  }

  if (unref(menuItemsContext).length) {
    sections.push({
      name: 'context',
      items: unref(menuItemsContext)
    })
  }
  if (unref(menuItemsShare).length) {
    sections.push({
      name: 'share',
      items: unref(menuItemsShare)
    })
  }
  if (unref(menuItemsActions).length) {
    sections.push({
      name: 'actions',
      items: unref(menuItemsActions)
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
