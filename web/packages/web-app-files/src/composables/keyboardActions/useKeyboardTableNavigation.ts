import { QueryValue, useResourcesStore, FolderViewModeConstants } from '@ownclouders/web-pkg'
import { useScrollTo } from '@ownclouders/web-pkg'
import { Ref, ref, unref, nextTick, watchEffect } from 'vue'
import { Key, KeyboardActions, Modifier, focusCheckbox } from '@ownclouders/web-pkg'
import { Resource } from '@ownclouders/web-client'
import { findIndex } from 'lodash-es'
import { storeToRefs } from 'pinia'

const enum Direction {
  LEFT = 'left',
  RIGHT = 'right'
}
export const useKeyboardTableNavigation = (
  keyActions: KeyboardActions,
  paginatedResources: Ref<Resource[]>,
  viewMode: Ref<string | QueryValue>
) => {
  const { scrollToResource } = useScrollTo()
  const resourcesStore = useResourcesStore()
  const { latestSelectedId } = storeToRefs(resourcesStore)

  const bindKeyActionsIds: Ref<string[]> = ref([])
  const tileViewStart = ref<string>(null)
  const tileViewDirection = ref<Direction>(null)

  keyActions.bindKeyAction({ modifier: Modifier.Ctrl, primary: Key.A }, () =>
    handleSelectAllAction()
  )

  keyActions.bindKeyAction({ primary: Key.Space }, () => {
    resourcesStore.toggleSelection(unref(latestSelectedId))
  })

  keyActions.bindKeyAction({ primary: Key.Esc }, () => {
    keyActions.resetSelectionCursor()
    tileViewStart.value = null
    resourcesStore.resetSelection()
  })

  const bindTilesViewKeyActions = () => {
    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ primary: Key.ArrowLeft }, () => handleNavigateAction(true))
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ primary: Key.ArrowRight }, () => handleNavigateAction())
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ primary: Key.ArrowUp }, async () => {
        const elementsInRow = getElementsInRow()
        if (elementsInRow === -1) {
          return
        }
        await handleNavigateAction(true, elementsInRow)
      })
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ primary: Key.ArrowDown }, async () => {
        const elementsInRow = getElementsInRow()
        if (elementsInRow === -1) {
          return
        }
        await handleNavigateAction(false, elementsInRow)
      })
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowLeft }, () =>
        handleTilesShiftHorizontalAction(Direction.LEFT)
      )
    )
    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowRight }, () =>
        handleTilesShiftHorizontalAction(Direction.RIGHT)
      )
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowUp }, () => {
        handleTilesShiftUpAction()
      })
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowDown }, () => {
        handleTilesShiftDownAction()
      })
    )
  }

  const bindTableViewKeyActions = () => {
    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ primary: Key.ArrowUp }, () => handleNavigateAction(true))
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ primary: Key.ArrowDown }, () => handleNavigateAction())
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowUp }, () =>
        handleShiftUpAction()
      )
    )

    bindKeyActionsIds.value.push(
      keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowDown }, () =>
        handleShiftDownAction()
      )
    )
  }

  const handleNavigateAction = async (up = false, moveBy = 1) => {
    const nextId = !unref(latestSelectedId) ? getFirstResourceId() : getNextResourceId(up, moveBy)
    if (nextId === -1) {
      return
    }
    keyActions.resetSelectionCursor()
    tileViewStart.value = null
    resourcesStore.resetSelection()
    await nextTick()
    resourcesStore.addSelection(nextId)
    await nextTick()
    focusCheckbox(nextId)
    scrollToResource(nextId, { topbarElement: 'files-app-bar' })
  }

  const getNextResourceId = (previous = false, movedBy = 1) => {
    const latestSelectedResourceIndex = paginatedResources.value.findIndex(
      (resource) => resource.id === latestSelectedId.value
    )
    if (latestSelectedResourceIndex === -1) {
      return -1
    }

    const step = previous ? -movedBy : movedBy
    let nextResourceIndex = latestSelectedResourceIndex + step

    while (nextResourceIndex >= 0 && nextResourceIndex < paginatedResources.value.length) {
      if (paginatedResources.value[nextResourceIndex].processing !== true) {
        return paginatedResources.value[nextResourceIndex].id
      }
      nextResourceIndex += step
    }

    if (nextResourceIndex < 0 || nextResourceIndex >= paginatedResources.value.length) {
      return -1
    }
    return paginatedResources.value[nextResourceIndex].id
  }

  const getFirstResourceId = () => {
    return paginatedResources.value.length ? paginatedResources.value[0].id : -1
  }

  const getElementsInRow = () => {
    const tilesListCard = document.querySelectorAll('#tiles-view > ul > li > div')
    if (tilesListCard.length === 0) {
      return -1
    }
    const firstElement = Math.floor(tilesListCard[0].getBoundingClientRect().x)
    let elementsInRow = 1

    for (let i = 1; i < tilesListCard.length; i++) {
      if (Math.floor(tilesListCard[i].getBoundingClientRect().x) !== firstElement) {
        elementsInRow++
      } else {
        break
      }
    }
    return elementsInRow
  }

  const handleSelectAllAction = () => {
    keyActions.resetSelectionCursor()
    resourcesStore.setSelection(
      unref(paginatedResources)
        .filter((resource) => resource.processing !== true)
        .map(({ id }) => id)
    )
  }

  const getVerticalProperties = (viewDirection: Direction) => {
    const elementsInRow = getElementsInRow()
    if (elementsInRow === -1) {
      return {}
    }

    if (!unref(tileViewStart)) {
      tileViewStart.value = latestSelectedId.value
      tileViewDirection.value = viewDirection
    }

    const tilesListCard = document.querySelectorAll('#tiles-view > ul > li > div')

    const currentResourceIndex = findIndex(
      tilesListCard,
      (tile) => tile.getAttribute('data-item-id') === unref(latestSelectedId).toString()
    )

    const tileViewStartIndex = findIndex(tilesListCard, (tile) => {
      return tile.getAttribute('data-item-id') === unref(tileViewStart).toString()
    })

    const nextIndex =
      viewDirection === Direction.LEFT
        ? currentResourceIndex - elementsInRow
        : currentResourceIndex + elementsInRow

    if (!tilesListCard[nextIndex]) {
      return {}
    }

    const lastSelectedFileId = tilesListCard[nextIndex].getAttribute('data-item-id')

    return {
      currentResourceIndex,
      nextIndex,
      tilesListCard,
      tileViewStartIndex,
      lastSelectedFileId
    }
  }

  const handleTilesShiftUpAction = () => {
    const vp = getVerticalProperties(Direction.LEFT)
    if (Object.keys(vp).length === 0) {
      return
    }
    for (let i = vp.currentResourceIndex; i >= vp.nextIndex; i--) {
      if (i === vp.tileViewStartIndex) {
        continue
      }
      const id = vp.tilesListCard[i].getAttribute('data-item-id')
      i < vp.tileViewStartIndex
        ? resourcesStore.addSelection(id)
        : resourcesStore.removeSelection(id)
    }
    focusCheckbox(vp.lastSelectedFileId)
    resourcesStore.setLastSelectedId(vp.lastSelectedFileId)
  }
  const handleTilesShiftDownAction = () => {
    const vp = getVerticalProperties(Direction.RIGHT)
    if (Object.keys(vp).length === 0) {
      return
    }

    for (let i = vp.currentResourceIndex; i <= vp.nextIndex; i++) {
      if (i === vp.tileViewStartIndex) {
        continue
      }
      const id = vp.tilesListCard[i].getAttribute('data-item-id')
      i > vp.tileViewStartIndex
        ? resourcesStore.addSelection(id)
        : resourcesStore.removeSelection(id)
    }
    focusCheckbox(vp.lastSelectedFileId)
    resourcesStore.setLastSelectedId(vp.lastSelectedFileId)
  }

  const handleTilesShiftHorizontalAction = (direction: Direction) => {
    const nextResourceId = !unref(latestSelectedId)
      ? getFirstResourceId()
      : getNextResourceId(direction === Direction.LEFT, 1)
    if (nextResourceId === -1) {
      return
    }

    if (unref(latestSelectedId) === unref(tileViewStart)) {
      tileViewStart.value = unref(latestSelectedId)
      tileViewDirection.value = direction
    }

    if (!unref(tileViewStart)) {
      tileViewStart.value = latestSelectedId.value
      tileViewDirection.value = direction
    }
    if (tileViewDirection.value !== direction && tileViewDirection.value !== null) {
      resourcesStore.toggleSelection(unref(latestSelectedId))
      resourcesStore.setLastSelectedId(nextResourceId)
    }
    if (tileViewDirection.value === direction) {
      resourcesStore.addSelection(nextResourceId)
    }
    focusCheckbox(nextResourceId)
  }

  const handleShiftUpAction = (movedBy = 1) => {
    const nextResourceId = getNextResourceId(true, movedBy)
    if (nextResourceId === -1) {
      return
    }
    if (unref(keyActions.selectionCursor) > 0) {
      resourcesStore.toggleSelection(unref(latestSelectedId))
      resourcesStore.setLastSelectedId(nextResourceId)
    } else {
      resourcesStore.addSelection(nextResourceId)
    }
    focusCheckbox(nextResourceId)
    scrollToResource(nextResourceId, { topbarElement: 'files-app-bar' })
    keyActions.selectionCursor.value = unref(keyActions.selectionCursor) - 1
  }
  const handleShiftDownAction = (movedBy = 1) => {
    const nextResourceId = getNextResourceId(false, movedBy)
    if (nextResourceId === -1) {
      return
    }
    if (unref(keyActions.selectionCursor) < 0) {
      resourcesStore.toggleSelection(unref(latestSelectedId))
      resourcesStore.setLastSelectedId(nextResourceId)
    } else {
      resourcesStore.addSelection(nextResourceId)
    }
    focusCheckbox(nextResourceId)
    scrollToResource(nextResourceId, { topbarElement: 'files-app-bar' })
    keyActions.selectionCursor.value = unref(keyActions.selectionCursor) + 1
  }

  watchEffect(() => {
    bindKeyActionsIds.value.forEach((id) => keyActions.removeKeyAction(id))
    bindKeyActionsIds.value = []
    FolderViewModeConstants.name.tiles === viewMode.value
      ? bindTilesViewKeyActions()
      : bindTableViewKeyActions()
  })
}
