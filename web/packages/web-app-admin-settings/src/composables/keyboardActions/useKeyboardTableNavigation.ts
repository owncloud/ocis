import { useScrollTo } from '@ownclouders/web-pkg'
import { Ref, unref } from 'vue'
import { Key, KeyboardActions, Modifier, focusCheckbox } from '@ownclouders/web-pkg'
import { find, findIndex } from 'lodash-es'
import { Item } from '@ownclouders/web-client'

export const useKeyboardTableNavigation = (
  keyActions: KeyboardActions,
  paginatedResources: Ref<Item[]>,
  selectedRows: Ref<Item[]>,
  lastSelectedRowIndex: Ref<number>,
  lastSelectedRowId: Ref<string | null>
) => {
  const { scrollToResource } = useScrollTo()

  keyActions.bindKeyAction({ primary: Key.ArrowUp }, () => handleNavigateAction(true))

  keyActions.bindKeyAction({ primary: Key.ArrowDown }, () => handleNavigateAction())

  keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowUp }, () =>
    handleShiftUpAction()
  )

  keyActions.bindKeyAction({ modifier: Modifier.Shift, primary: Key.ArrowDown }, () =>
    handleShiftDownAction()
  )

  keyActions.bindKeyAction({ modifier: Modifier.Ctrl, primary: Key.A }, () =>
    handleSelectAllAction()
  )

  keyActions.bindKeyAction({ primary: Key.Space }, () => {
    const { lastSelectedRow, lastSelectedRowIndex } = getLastSelectedRow()
    if (lastSelectedRowIndex === -1) {
      selectedRows.value.push(lastSelectedRow)
    } else {
      selectedRows.value = unref(selectedRows).filter((item) => item.id !== lastSelectedRow.id)
    }
  })

  keyActions.bindKeyAction({ primary: Key.Esc }, () => {
    keyActions.resetSelectionCursor()
    selectedRows.value = []
  })

  const handleNavigateAction = (up = false) => {
    const nextResource = !unref(lastSelectedRowId) ? getFirstResource() : getNextResource(up)

    if (nextResource === -1) {
      return
    }

    const nextResourceIndex = findIndex(
      paginatedResources.value,
      (resource) => resource.id === nextResource.id
    )

    focusCheckbox(nextResource.id)
    keyActions.resetSelectionCursor()
    selectedRows.value = [nextResource]
    lastSelectedRowIndex.value = nextResourceIndex
    lastSelectedRowId.value = String(nextResource.id)

    scrollToResource(nextResource.id, { topbarElement: 'admin-settings-app-bar' })
  }

  const handleShiftUpAction = () => {
    const nextResource = getNextResource(true)
    if (nextResource === -1) {
      return
    }

    const nextResourceIndex = findIndex(
      paginatedResources.value,
      (resource) => resource.id === nextResource.id
    )

    if (unref(keyActions.selectionCursor) > 0) {
      const { lastSelectedRow, lastSelectedRowIndex } = getLastSelectedRow()

      lastSelectedRowIndex === -1
        ? selectedRows.value.push(lastSelectedRow)
        : (selectedRows.value = unref(selectedRows).filter(
            (item) => item.id !== lastSelectedRow.id
          ))
    } else {
      selectedRows.value.push(nextResource)
    }

    focusCheckbox(nextResource.id)
    lastSelectedRowIndex.value = nextResourceIndex
    lastSelectedRowId.value = String(nextResource.id)
    keyActions.selectionCursor.value = unref(keyActions.selectionCursor) - 1
    scrollToResource(nextResource.id, { topbarElement: 'admin-settings-app-bar' })
  }
  const handleShiftDownAction = () => {
    const nextResource = getNextResource(false)
    if (nextResource === -1) {
      return
    }

    const nextResourceIndex = findIndex(
      paginatedResources.value,
      (resource) => resource.id === nextResource.id
    )

    if (unref(keyActions.selectionCursor) < 0) {
      const lastSelectedRow = find(
        paginatedResources.value,
        (resource) => resource.id === lastSelectedRowId.value
      )
      const lastSelectedRowIndex = findIndex(
        unref(selectedRows),
        (resource: any) => resource.id === lastSelectedRowId.value
      )

      if (lastSelectedRowIndex === -1) {
        selectedRows.value.push(lastSelectedRow)
      } else {
        selectedRows.value = unref(selectedRows).filter((item) => item.id !== lastSelectedRow.id)
      }
    } else {
      selectedRows.value.push(nextResource)
    }

    focusCheckbox(nextResource.id)
    lastSelectedRowIndex.value = nextResourceIndex
    lastSelectedRowId.value = String(nextResource.id)
    keyActions.selectionCursor.value = unref(keyActions.selectionCursor) + 1
    scrollToResource(nextResource.id, { topbarElement: 'admin-settings-app-bar' })
  }

  const handleSelectAllAction = () => {
    keyActions.resetSelectionCursor()
    selectedRows.value = [...unref(paginatedResources)]
  }

  const getNextResource = (previous = false) => {
    const latestSelectedResourceIndex = paginatedResources.value.findIndex(
      (resource) => resource.id === lastSelectedRowId.value
    )
    if (latestSelectedResourceIndex === -1) {
      return -1
    }
    const nextResourceIndex = latestSelectedResourceIndex + (previous ? -1 : 1)
    if (nextResourceIndex < 0 || nextResourceIndex >= paginatedResources.value.length) {
      return -1
    }
    return paginatedResources.value[nextResourceIndex]
  }

  const getFirstResource = () => {
    return paginatedResources.value.length ? paginatedResources.value[0] : -1
  }

  const getLastSelectedRow = () => {
    const lastSelectedRow = find(
      paginatedResources.value,
      (resource) => resource.id === lastSelectedRowId.value
    )
    const lastSelectedRowIndex = findIndex(
      unref(selectedRows),
      (resource: any) => resource.id === lastSelectedRowId.value
    )
    return {
      lastSelectedRow,
      lastSelectedRowIndex
    }
  }
}
