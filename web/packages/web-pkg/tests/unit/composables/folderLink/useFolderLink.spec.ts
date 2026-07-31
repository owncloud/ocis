import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { CapabilityStore, useFolderLink } from '../../../../src/composables'
import { Resource, SpaceResource } from '@ownclouders/web-client'

describe('useFolderLink', () => {
  it('getFolderLink should return the correct folder link', () => {
    const resource = {
      path: '/my-folder',
      id: '2',
      fileId: '2',
      storageId: '1',
      spaceId: '1'
    }
    const wrapper = createWrapper()

    const folderLink = wrapper.vm.getFolderLink(resource)
    expect(folderLink).toEqual({
      name: 'files-spaces-generic',
      params: { driveAliasAndItem: 'personal/admin' },
      query: { fileId: '2' }
    })
  })

  it('getParentFolderLink should return the correct parent folder link', () => {
    const resource = {
      path: '/my-folder',
      id: '2',
      fileId: '2',
      storageId: '1',
      parentFolderId: '1',
      spaceId: '1'
    }

    const wrapper = createWrapper()
    const parentFolderLink = wrapper.vm.getParentFolderLink(resource)
    expect(parentFolderLink).toEqual({
      name: 'files-spaces-generic',
      params: { driveAliasAndItem: 'personal/admin' },
      query: { fileId: '1' }
    })
  })

  describe('getParentFolderName should return the correct parent folder name', () => {
    it('should equal "Personal" if share jail is enabled', () => {
      const resource = {
        path: '/my-folder',
        storageId: '1',
        spaceId: '1'
      } as Resource

      const wrapper = createWrapper()
      const parentFolderName = wrapper.vm.getParentFolderName(resource)
      expect(parentFolderName).toEqual('Personal')
    })
    it('should equal the space name if resource storage is representing a project space', () => {
      const resource = {
        path: '/my-folder',
        storageId: '2',
        spaceId: '2'
      } as Resource

      const wrapper = createWrapper()
      const parentFolderName = wrapper.vm.getParentFolderName(resource)
      expect(parentFolderName).toEqual('New space')
    })
    it('should equal the "Shared with me" if resource is representing the root share', () => {
      const resource = {
        path: '/My share',
        remoteItemPath: '/My share',
        remoteItemId: '1',
        spaceId: '1',
        isShareRoot: () => true
      } as Resource

      const wrapper = createWrapper()
      const parentFolderName = wrapper.vm.getParentFolderName(resource)
      expect(parentFolderName).toEqual('Shared with me')
    })
    it('should equal the share name if resource is representing a file or folder in the root of a share', () => {
      const resource = {
        path: '/My share/test.txt',
        remoteItemPath: '/My share',
        remoteItemId: '1',
        storageId: '1',
        spaceId: '1'
      } as Resource

      const wrapper = createWrapper()
      const parentFolderName = wrapper.vm.getParentFolderName(resource)
      expect(parentFolderName).toEqual('My share')
    })
  })
})

const createWrapper = () => {
  const spaces = [
    {
      id: '1',
      fileId: '1',
      driveType: 'personal',
      getDriveAliasAndItem: () => 'personal/admin'
    },
    {
      id: '2',
      driveType: 'project',
      name: 'New space',
      getDriveAliasAndItem: vi.fn()
    }
  ] as unknown as SpaceResource[]

  const mocks = defaultComponentMocks({})
  const capabilities = {
    spaces: { projects: true }
  } satisfies Partial<CapabilityStore['capabilities']>

  return getComposableWrapper(
    () => {
      const {
        getFolderLink,
        getParentFolderLink,
        getParentFolderName,
        getParentFolderLinkIconAdditionalAttributes
      } = useFolderLink()

      return {
        getFolderLink,
        getParentFolderLink,
        getParentFolderName,
        getParentFolderLinkIconAdditionalAttributes
      }
    },
    {
      mocks,
      provide: mocks,
      pluginOptions: {
        piniaOptions: { spacesState: { spaces }, capabilityState: { capabilities } }
      }
    }
  )
}
