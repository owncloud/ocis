import {
  defaultComponentMocks,
  defaultPlugins,
  RouteLocation,
  shallowMount
} from '@ownclouders/web-test-helpers'
import ResourceDetails from '../../../../src/components/FilesList/ResourceDetails.vue'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useOpenWithDefaultApp, useRouteQuery } from '@ownclouders/web-pkg'
import { ref } from 'vue'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  getIndicators: vi.fn(() => []),
  useRouteQuery: vi.fn(),
  useOpenWithDefaultApp: vi.fn()
}))

describe('ResourceDetails component', () => {
  describe('open with default actions', () => {
    it("doesn't open default action if query param 'openWithDefaultApp' isn't set true", () => {
      const { mocks } = getWrapper()
      expect(mocks.openWithDefaultAppMock).not.toHaveBeenCalled()
    })
    it("opens default action if query param 'openWithDefaultApp' is set true", () => {
      vi.mocked(useRouteQuery).mockImplementationOnce(() => ref('true'))
      const { mocks } = getWrapper()
      expect(mocks.openWithDefaultAppMock).toHaveBeenCalled()
    })
  })

  function getWrapper(isFolder = false) {
    const openWithDefaultAppMock = vi.fn()
    vi.mocked(useOpenWithDefaultApp).mockReturnValue({
      openWithDefaultApp: openWithDefaultAppMock
    })

    const mocks = {
      ...defaultComponentMocks({
        currentRoute: mock<RouteLocation>({
          name: 'files-public-link',
          query: { openWithDefaultAppQuery: 'true' }
        })
      }),
      openWithDefaultAppMock
    }

    const file = {
      id: '0',
      name: 'image.jpg',
      path: '/image.jpg',
      size: 24064,
      mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
      mimeType: 'image/jpg',
      isFolder,
      owner: {
        id: '1',
        displayName: 'admin'
      }
    } as Resource
    const space = mock<SpaceResource>()

    return {
      mocks,
      wrapper: shallowMount(ResourceDetails, {
        props: {
          space,
          singleResource: file
        },
        global: {
          mocks,
          plugins: [
            ...defaultPlugins({
              piniaOptions: { resourcesStore: { currentFolder: mock<Resource>() } }
            })
          ],
          provide: {
            ...mocks,
            space,
            resource: file,
            versions: []
          }
        }
      })
    }
  }
})
