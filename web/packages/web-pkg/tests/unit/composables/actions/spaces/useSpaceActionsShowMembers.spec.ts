import { unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'
import { useSpaceActionsShowMembers } from '../../../../../src/composables/actions'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'

describe('showMembers', () => {
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be true when a resource is given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [{ id: '1' } as SpaceResource] })).toBe(
            true
          )
        }
      })
    })
    it('should be false when resource is disabled', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [{ id: '1', disabled: true } as SpaceResource]
            })
          ).toBe(false)
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useSpaceActionsShowMembers>) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useSpaceActionsShowMembers()
      setup(instance)
    })
  }
}
