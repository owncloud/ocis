import { useModals } from '@ownclouders/web-pkg'
import { useGroupActionsCreateGroup } from '../../../../../src/composables/actions/groups/useGroupActionsCreateGroup'
import { unref } from 'vue'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'

describe('useGroupActionsCreateGroup', () => {
  describe('method "handler"', () => {
    it('creates a modal', () => {
      getWrapper({
        setup: async ({ actions }) => {
          const { dispatchModal } = useModals()
          await unref(actions)[0].handler()
          expect(dispatchModal).toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useGroupActionsCreateGroup>) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useGroupActionsCreateGroup()
      setup(instance)
    })
  }
}
