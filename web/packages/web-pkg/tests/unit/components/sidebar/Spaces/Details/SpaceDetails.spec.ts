import SpaceDetails from '../../../../../../src/components/SideBar/Spaces/Details/SpaceDetails.vue'
import {
  CollaboratorShare,
  ShareRole,
  SpaceResource,
  Resource,
  SpaceMember,
  GraphSharePermission
} from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { RouteLocation } from 'vue-router'
import { User } from '@ownclouders/web-client/graph/generated'

const spaceMock = {
  type: 'space',
  name: ' space',
  id: '1',
  mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
  members: [
    mock<SpaceMember>({
      permissions: [GraphSharePermission.deletePermissions],
      grantedTo: { user: { id: '1', displayName: 'alice' }, group: undefined }
    })
  ],
  spaceQuota: {
    used: 100,
    total: 1000
  }
} as unknown as SpaceResource

const spaceShare = {
  id: '1',
  sharedWith: {
    id: 'Alice',
    displayName: 'alice'
  },
  role: mock<ShareRole>()
} as CollaboratorShare

const selectors = {
  spaceDefaultImage: '.space-default-image',
  spaceMembers: '.oc-space-details-sidebar-members'
}

describe('Details SideBar Panel', () => {
  it('displays the details side panel', () => {
    const { wrapper } = createWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('does render the space default image if "showSpaceImage" is false', () => {
    const { wrapper } = createWrapper({ props: { showSpaceImage: false } })
    expect(wrapper.find(selectors.spaceDefaultImage).exists()).toBeTruthy()
  })
  it('does not render share indicators if "showShareIndicators" is false', () => {
    const { wrapper } = createWrapper({
      spaceResource: spaceMock,
      props: { showShareIndicators: false }
    })
    expect(wrapper.find(selectors.spaceMembers).exists()).toBeFalsy()
  })
  it('does not render share indicators if space is disabled', () => {
    const { wrapper } = createWrapper({
      spaceResource: { ...spaceMock, disabled: true },
      props: { showShareIndicators: true }
    })
    expect(wrapper.find(selectors.spaceMembers).exists()).toBeFalsy()
  })
})

function createWrapper({ spaceResource = spaceMock, props = {} } = {}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-spaces-generic' })
    })
  }

  return {
    wrapper: shallowMount(SpaceDetails, {
      props: { ...props },
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              userState: { user: { id: '1', onPremisesSamAccountName: 'marie' } as User },
              sharesState: { collaboratorShares: [spaceShare] },
              resourcesStore: { resources: [mock<Resource>({ name: 'file1', type: 'file' })] }
            }
          })
        ],
        mocks,
        provide: { resource: spaceResource, ...mocks }
      }
    })
  }
}
