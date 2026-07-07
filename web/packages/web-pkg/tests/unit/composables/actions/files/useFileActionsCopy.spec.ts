import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
} from '@ownclouders/web-test-helpers'
import { useFileActionsCopy } from '../../../../../src/composables/actions/files'
import { useClipboardStore } from '../../../../../src/composables/piniaStores'
import { describe } from 'vitest'

describe('copy', () => {
  describe('search context', () => {
    describe('computed property "actions"', () => {
      describe('handler', () => {
        it.each([
          {
            resources: [{ id: '1' }, { id: '2' }] as Resource[],
            copyAbleResources: ['1', '2']
          },
          {
            resources: [
              { id: '1' },
              { id: '2' },
              { id: '3' },
              { id: '4', fileId: '5', canDownload: () => true, driveType: 'project' }
            ] as Resource[],
            copyAbleResources: ['1', '2', '3']
          }
        ])('should filter non copyable resources', ({ resources, copyAbleResources }) => {
          getWrapper({
            searchLocation: true,
            setup: ({ actions }) => {
              unref(actions)[0].handler({ space: null, resources })
              const clipboardStore = useClipboardStore()
              expect(clipboardStore.copyResources).toHaveBeenCalledWith(
                resources.filter((r) => copyAbleResources.includes(r.id as string))
              )
            }
          })
        })
      })
    })
  })
  describe('computed property "actions"', () => {
    describe('isVisible', () => {
      it('returns true if "canDownload" is true', () => {
        getWrapper({
          searchLocation: false,
          setup: ({ actions }) => {
            expect(
              unref(actions)[0].isVisible({
                space: null,
                resources: [mock<Resource>({ id: '1', canDownload: () => true })]
              })
            ).toBeTruthy()
          }
        })
      })
      it('returns false if "canDownload" is false', () => {
        getWrapper({
          searchLocation: false,
          setup: ({ actions }) => {
            expect(
              unref(actions)[0].isVisible({
                space: null,
                resources: [mock<Resource>({ id: '1', canDownload: () => false })]
              })
            ).toBeFalsy()
          }
        })
      })
    })
  })
})

function getWrapper({
  searchLocation = false,
  setup
}: {
  searchLocation: boolean
  setup: (instance: ReturnType<typeof useFileActionsCopy>) => void
}) {
  const routeName = searchLocation ? 'files-common-search' : 'files-spaces-generic'

  const mocks = {
    ...defaultComponentMocks({ currentRoute: mock<RouteLocation>({ name: routeName }) }),
    space: {
      driveType: 'personal'
    } as unknown as SpaceResource
  }
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsCopy()
        setup(instance)
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
