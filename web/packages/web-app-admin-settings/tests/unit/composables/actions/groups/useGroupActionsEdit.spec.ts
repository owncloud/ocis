import { useGroupActionsEdit } from '../../../../../src/composables/actions/groups/useGroupActionsEdit'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { Group } from '@ownclouders/web-client/graph/generated'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'

describe('useGroupActionsEdit', () => {
  describe('method "isVisible"', () => {
    it.each([
      { resources: [mock<Group>({ groupTypes: [] })], isVisible: true },
      { resources: [], isVisible: false },
      {
        resources: [mock<Group>({ groupTypes: [] }), mock<Group>({ groupTypes: [] })],
        isVisible: false
      }
    ])('should only return true for one group', ({ resources, isVisible }) => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources })).toEqual(isVisible)
        }
      })
    })
    it('should return false for read-only groups', () => {
      getWrapper({
        setup: ({ actions }) => {
          const resources = [mock<Group>({ groupTypes: ['ReadOnly'] })]
          expect(unref(actions)[0].isVisible({ resources })).toBeFalsy()
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useGroupActionsEdit>) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useGroupActionsEdit()
      setup(instance)
    })
  }
}
