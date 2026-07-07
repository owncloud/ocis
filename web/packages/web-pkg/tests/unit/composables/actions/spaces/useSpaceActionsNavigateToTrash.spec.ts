import { useSpaceActionsNavigateToTrash } from '../../../../../src'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'

describe('navigateToSpace', () => {
  describe('isVisible property', () => {
    it('should be false when no resource given', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ resources: [] })).toBe(false)
        }
      })
    })
    it('should be false when the space is disabled', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [
                mock<SpaceResource>({
                  disabled: true,
                  driveType: 'project'
                })
              ]
            })
          ).toBe(false)
        }
      })
    })
    it('should be false when the space is no project space', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({
              resources: [
                mock<SpaceResource>({
                  disabled: false,
                  driveType: 'personal'
                })
              ]
            })
          ).toBe(false)
        }
      })
    })
  })
  describe('handler', () => {
    it('should redirect to respective trash', () => {
      const { mocks } = getWrapper({
        setup: async ({ actions }) => {
          await unref(actions)[0].handler({
            resources: [mock<SpaceResource>()]
          })
          expect(mocks.$router.push).toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useSpaceActionsNavigateToTrash>) => void
}) {
  const mocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'files-spaces-projects' })
  })
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useSpaceActionsNavigateToTrash()
        setup(instance)
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
