import DriveResolver from '../../../../src/views/spaces/DriveResolver.vue'
import { queryItemAsString, useDriveResolver, useRouteParam } from '@ownclouders/web-pkg'
import { computed, ref } from 'vue'
import { mock, mockDeep } from 'vitest-mock-extended'
import { ClientService } from '@ownclouders/web-pkg'
import { useGetMatchingSpace } from '@ownclouders/web-pkg'
import { locationPublicUpload } from '@ownclouders/web-pkg'
import { PublicSpaceResource, Resource, SpaceResource } from '@ownclouders/web-client'
import { SharePermissionBit } from '@ownclouders/web-client'
import {
  defaultPlugins,
  mount,
  defaultComponentMocks,
  defaultStubs,
  RouteLocation,
  useGetMatchingSpaceMock
} from '@ownclouders/web-test-helpers'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useGetMatchingSpace: vi.fn(),
  useDriveResolver: vi.fn(),
  useRouteParam: vi.fn(),
  queryItemAsString: vi.fn()
}))

describe('DriveResolver view', () => {
  it('renders the "drive-redirect"-component when no space is given', async () => {
    const { wrapper } = getMountedWrapper()
    await wrapper.vm.$nextTick()
    expect(wrapper.find('drive-redirect-stub').exists()).toBeTruthy()
  })
  it('renders the "generic-trash"-component when on a trash route', async () => {
    const { wrapper } = getMountedWrapper({
      space: mock<SpaceResource>({ driveType: 'project' }),
      currentRouteName: 'files-trash-generic'
    })
    await wrapper.vm.$nextTick()

    expect(wrapper.find('generic-trash-stub').exists()).toBeTruthy()
  })
  it('renders the "generic-space"-component when a space is given', async () => {
    const { wrapper } = getMountedWrapper({
      space: mock<SpaceResource>({ driveType: 'project' })
    })
    await wrapper.vm.$nextTick()

    expect(wrapper.find('generic-space-stub').exists()).toBeTruthy()
  })
  it('redirects to the public drop page in a public context with "upload-only"-permissions', async () => {
    const space = mock<SpaceResource>({
      id: '1',
      fileId: '1',
      getDriveAliasAndItem: vi.fn(),
      driveType: 'public'
    })
    const clientService = mockDeep<ClientService>()
    clientService.webdav.getFileInfo.mockResolvedValue(
      mockDeep<PublicSpaceResource>({ publicLinkPermission: SharePermissionBit.Create })
    )
    const { wrapper, mocks } = getMountedWrapper({
      space,
      mocks: { $clientService: clientService }
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(mocks.$router.push).toHaveBeenCalledWith({
      name: locationPublicUpload.name,
      params: { token: space.id }
    })
  })
  it('redirects to personal space if user has access to the resource via their personal space', async () => {
    const space = mock<SpaceResource>({
      id: '1',
      getDriveAliasAndItem: vi.fn(),
      driveType: 'public'
    })
    const internalSpace = mock<SpaceResource>({
      id: '1',
      getDriveAliasAndItem: vi.fn(),
      driveType: 'personal'
    })
    const clientService = mockDeep<ClientService>()
    clientService.webdav.getPathForFileId.mockResolvedValue('/path')
    clientService.webdav.getFileInfo.mockResolvedValue(mock<Resource>())
    const { wrapper, mocks } = getMountedWrapper({
      space,
      internalSpace,
      mocks: { $clientService: clientService },
      userContextReady: true
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(mocks.$router.push).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'files-spaces-generic' })
    )
  })
  it('redirects to private link if user has access to the resource via a share', async () => {
    const space = mock<SpaceResource>({
      id: '1',
      getDriveAliasAndItem: vi.fn(),
      driveType: 'public'
    })
    const clientService = mockDeep<ClientService>()
    clientService.webdav.getPathForFileId.mockResolvedValue('/path')
    const { wrapper, mocks } = getMountedWrapper({
      space,
      mocks: { $clientService: clientService },
      userContextReady: true
    })

    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(mocks.$router.push).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'resolvePrivateLink'
      })
    )
  })
  it('redirects to private link if no drive alias but a fileId is given', async () => {
    const { wrapper, mocks } = getMountedWrapper({ driveAliasAndItem: '' })
    await wrapper.vm.$nextTick()

    expect(mocks.$router.push).toHaveBeenCalledWith(
      expect.objectContaining({
        name: 'resolvePrivateLink'
      })
    )
  })
})

function getMountedWrapper({
  mocks = {},
  space = undefined,
  internalSpace = undefined,
  currentRouteName = 'files-spaces-generic',
  userContextReady = false,
  driveAliasAndItem = 'personal/einstein/file',
  fileId = '1'
}: {
  mocks?: Record<string, unknown>
  space?: SpaceResource
  internalSpace?: SpaceResource
  currentRouteName?: string
  userContextReady?: boolean
  driveAliasAndItem?: string
  fileId?: string
} = {}) {
  vi.mocked(useRouteParam).mockReturnValue(ref(driveAliasAndItem))
  vi.mocked(queryItemAsString).mockReturnValue(fileId)
  vi.mocked(useDriveResolver).mockImplementation(() => ({
    space: ref(space),
    item: ref('/'),
    itemId: computed(() => 'id'),
    loading: ref(false)
  }))
  vi.mocked(useGetMatchingSpace).mockImplementation(() =>
    useGetMatchingSpaceMock({
      getInternalSpace: () => internalSpace
    })
  )

  const defaultMocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({
        name: currentRouteName,
        params: { driveAliasAndItem: '/' }
      })
    }),
    ...(mocks && mocks)
  }

  return {
    mocks: defaultMocks,
    wrapper: mount(DriveResolver, {
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: { authState: { userContextReady: userContextReady } }
          })
        ],
        mocks: defaultMocks,
        provide: defaultMocks,
        stubs: { ...defaultStubs, 'app-banner': true }
      }
    })
  }
}
