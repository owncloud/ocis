import { onBeforeUnmount, onMounted, unref, Ref, watchEffect } from 'vue'
import { QueryValue, FolderViewModeConstants, useResourcesStore } from '@ownclouders/web-pkg'
import { eventBus } from '@ownclouders/web-pkg'
import { KeyboardActions } from '@ownclouders/web-pkg'
import { Resource } from '@ownclouders/web-client'
import { findIndex } from 'lodash-es'
import { storeToRefs } from 'pinia'

export const useKeyboardTableMouseActions = (
  keyActions: KeyboardActions,
  viewMode: Ref<string | QueryValue>
) => {
  const resourcesStore = useResourcesStore()
  const { latestSelectedId } = storeToRefs(resourcesStore)

  let fileListClickedEvent: string
  let fileListClickedMetaEvent: string
  let fileListClickedShiftEvent: string

  const handleCtrlClickAction = (resource: Resource) => {
    resourcesStore.toggleSelection(resource.id)
  }

  const handleShiftClickAction = ({
    resource,
    skipTargetSelection
  }: {
    resource: Resource
    skipTargetSelection: boolean
  }) => {
    const parent = document.querySelectorAll(`[data-item-id='${resource.id}']`)[0]
    const resourceNodes = Object.values(parent.parentNode.children)
    const latestNode = resourceNodes.find(
      (r) => r.getAttribute('data-item-id') === unref(latestSelectedId)
    )
    const clickedNode = resourceNodes.find((r) => r.getAttribute('data-item-id') === resource.id)

    let latestNodeIndex = resourceNodes.indexOf(latestNode)
    latestNodeIndex = latestNodeIndex === -1 ? 0 : latestNodeIndex

    const clickedNodeIndex = resourceNodes.indexOf(clickedNode)
    const minIndex = Math.min(latestNodeIndex, clickedNodeIndex)
    const maxIndex = Math.max(latestNodeIndex, clickedNodeIndex)

    for (let i = minIndex; i <= maxIndex; i++) {
      const nodeId = resourceNodes[i].getAttribute('data-item-id')
      const isDisabled = resourceNodes[i].classList.contains('oc-table-disabled')
      if ((skipTargetSelection && nodeId === resource.id) || isDisabled) {
        continue
      }
      resourcesStore.addSelection(nodeId)
    }
    resourcesStore.setLastSelectedId(resource.id)
  }

  const handleTilesShiftClickAction = ({
    resource,
    skipTargetSelection
  }: {
    resource: Resource
    skipTargetSelection: boolean
  }) => {
    const tilesListCard = document.querySelectorAll('#tiles-view > ul > li > div')
    const startIndex = findIndex(
      tilesListCard,
      (r) => r.getAttribute('data-item-id') === resource.id
    )
    const endIndex = findIndex(
      tilesListCard,
      (r) => r.getAttribute('data-item-id') === unref(latestSelectedId)
    )
    const minIndex = Math.min(endIndex, startIndex)
    const maxIndex = Math.max(endIndex, startIndex)

    for (let i = minIndex; i <= maxIndex; i++) {
      const nodeId = tilesListCard[i].getAttribute('data-item-id')
      const isDisabled = tilesListCard[i].classList.contains('oc-tile-card-disabled')

      if ((skipTargetSelection && nodeId === resource.id) || isDisabled) {
        continue
      }
      resourcesStore.addSelection(nodeId)
    }
    resourcesStore.setLastSelectedId(resource.id)
  }

  onMounted(() => {
    fileListClickedEvent = eventBus.subscribe(
      'app.files.list.clicked',
      keyActions.resetSelectionCursor
    )
    fileListClickedMetaEvent = eventBus.subscribe(
      'app.files.list.clicked.meta',
      handleCtrlClickAction
    )
  })

  onBeforeUnmount(() => {
    eventBus.unsubscribe('app.files.list.clicked', fileListClickedEvent)
    eventBus.unsubscribe('app.files.list.clicked.meta', fileListClickedMetaEvent)
    eventBus.unsubscribe('app.files.list.clicked.shift', fileListClickedShiftEvent)
  })
  watchEffect(() => {
    eventBus.unsubscribe('app.files.list.clicked.shift', fileListClickedShiftEvent)
    fileListClickedShiftEvent = eventBus.subscribe(
      'app.files.list.clicked.shift',
      FolderViewModeConstants.name.tiles === viewMode.value
        ? handleTilesShiftClickAction
        : handleShiftClickAction
    )
  })
}
