import { useUserActionsEdit } from '../../../../../src/composables/actions/users/useUserActionsEdit'
import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'

describe('useUserActionsEdit', () => {
  describe('method "isVisible"', () => {
    it.each([
      { resources: [mock<User>()], isVisible: true },
      { resources: [], isVisible: false },
      { resources: [mock<User>(), mock<User>()], isVisible: false }
    ])('should only return true for one user', ({ resources, isVisible }) => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources })).toEqual(isVisible)
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useUserActionsEdit>) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useUserActionsEdit()
      setup(instance)
    })
  }
}
