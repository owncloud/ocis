import { useCapabilityStore, useModals } from '@ownclouders/web-pkg'
import { useUserActionsCreateUser } from '../../../../../src/composables/actions/users/useUserActionsCreateUser'
import { unref } from 'vue'
import { getComposableWrapper, writable } from '@ownclouders/web-test-helpers'

describe('useUserActionsCreateUser', () => {
  describe('method "isVisible"', () => {
    it.each([true, false])(
      'is enabled based on the capability',
      (capabilityCreateUsersDisabled) => {
        getWrapper({
          setup: ({ actions }) => {
            const capabilityStore = useCapabilityStore()
            writable(capabilityStore).graphUsersCreateDisabled = capabilityCreateUsersDisabled
            expect(unref(actions)[0].isVisible()).toEqual(!capabilityCreateUsersDisabled)
          }
        })
      }
    )
  })
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
  setup: (instance: ReturnType<typeof useUserActionsCreateUser>) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useUserActionsCreateUser()
      setup(instance)
    })
  }
}
