import { ref } from 'vue'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useKeyboardActions, useResourcesStore, eventBus } from '@ownclouders/web-pkg'
import { useKeyboardTableMouseActions } from '../../../../src/composables/keyboardActions/useKeyboardTableMouseActions'

describe('useKeyboardTableMouseActions', () => {
  describe('table view (shift click)', () => {
    it('replaces the selection with the clicked range, regardless of direction', () => {
      const { resourcesStore, wrapper } = setUpTableRows(['a', 'b', 'c', 'd'])

      resourcesStore.setSelection(['a'])
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'c' },
        skipTargetSelection: false,
        extend: false
      })
      expect(resourcesStore.selectedIds).toEqual(['a', 'b', 'c'])

      // shift click "backwards" from the new anchor
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'a' },
        skipTargetSelection: false,
        extend: false
      })
      expect(resourcesStore.selectedIds).toEqual(['a', 'b', 'c'])

      wrapper.unmount()
    })

    it('clears any selection outside of the newly clicked range', () => {
      const { resourcesStore, wrapper } = setUpTableRows(['a', 'b', 'c', 'd'])

      resourcesStore.setSelection(['a', 'd'])
      resourcesStore.setLastSelectedId('a')
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'b' },
        skipTargetSelection: false,
        extend: false
      })

      expect(resourcesStore.selectedIds).toEqual(['a', 'b'])

      wrapper.unmount()
    })

    it('extends instead of replacing when ctrl/cmd is held', () => {
      const { resourcesStore, wrapper } = setUpTableRows(['a', 'b', 'c', 'd'])

      resourcesStore.setSelection(['d'])
      resourcesStore.setLastSelectedId('a')
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'b' },
        skipTargetSelection: false,
        extend: true
      })

      expect(resourcesStore.selectedIds).toEqual(['d', 'a', 'b'])

      wrapper.unmount()
    })

    it('skips disabled rows', () => {
      const { resourcesStore, wrapper } = setUpTableRows(['a', 'b', 'c', 'd'], { disabled: ['b'] })

      resourcesStore.setSelection(['a'])
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'c' },
        skipTargetSelection: false,
        extend: false
      })

      expect(resourcesStore.selectedIds).toEqual(['a', 'c'])

      wrapper.unmount()
    })
  })

  describe('tiles view (shift click)', () => {
    it('replaces the selection with the clicked range', () => {
      const { resourcesStore, wrapper } = setUpTiles(['a', 'b', 'c', 'd'])

      resourcesStore.setSelection(['a'])
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'c' },
        skipTargetSelection: false,
        extend: false
      })

      expect(resourcesStore.selectedIds).toEqual(['a', 'b', 'c'])

      wrapper.unmount()
    })

    it('extends instead of replacing when ctrl/cmd is held', () => {
      const { resourcesStore, wrapper } = setUpTiles(['a', 'b', 'c', 'd'])

      resourcesStore.setSelection(['d'])
      resourcesStore.setLastSelectedId('a')
      eventBus.publish('app.files.list.clicked.shift', {
        resource: { id: 'b' },
        skipTargetSelection: false,
        extend: true
      })

      expect(resourcesStore.selectedIds).toEqual(['d', 'a', 'b'])

      wrapper.unmount()
    })
  })
})

function setUpTableRows(ids: string[], { disabled = [] as string[] } = {}) {
  document.body.innerHTML = `
    <table><tbody>
      ${ids
        .map(
          (id) =>
            `<tr data-item-id="${id}" class="${disabled.includes(id) ? 'oc-table-disabled' : ''}"></tr>`
        )
        .join('')}
    </tbody></table>
  `
  return setUpComposable(ref('resource-table'))
}

function setUpTiles(ids: string[]) {
  document.body.innerHTML = `
    <div id="tiles-view"><ul>
      ${ids.map((id) => `<li><div data-item-id="${id}"></div></li>`).join('')}
    </ul></div>
  `
  return setUpComposable(ref('resource-tiles'))
}

function setUpComposable(viewMode: ReturnType<typeof ref<string>>) {
  let resourcesStore: ReturnType<typeof useResourcesStore>
  const wrapper = getComposableWrapper(
    () => {
      const keyActions = useKeyboardActions()
      resourcesStore = useResourcesStore()
      useKeyboardTableMouseActions(keyActions, viewMode)
      return {}
    },
    { pluginOptions: { piniaOptions: { stubActions: false } } }
  )
  return { resourcesStore, wrapper }
}
