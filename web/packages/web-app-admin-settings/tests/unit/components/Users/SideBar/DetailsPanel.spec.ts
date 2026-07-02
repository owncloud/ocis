import { User } from '@ownclouders/web-client/graph/generated'
import DetailsPanel from '../../../../../src/components/Users/SideBar/DetailsPanel.vue'
import UserInfoBox from '../../../../../src/components/Users/SideBar/UserInfoBox.vue'
import { PartialComponentProps, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

const defaultUser = { displayName: 'user', memberOf: [] } as User

describe('DetailsPanel', () => {
  describe('computed method "user"', () => {
    it('should be set if only one user is given', () => {
      const { wrapper } = getWrapper({ props: { user: defaultUser, users: [defaultUser] } })
      expect(wrapper.vm.user).toEqual(defaultUser)
    })
    it('should not be set if no users are given', () => {
      const { wrapper } = getWrapper({
        props: { user: null, users: [] }
      })
      expect(wrapper.vm.user).toEqual(null)
    })
    it('should not be set if multiple users are given', () => {
      const { wrapper } = getWrapper({
        props: { user: null, users: [defaultUser, { displayName: 'user2' } as User] }
      })
      expect(wrapper.vm.user).toEqual(null)
    })
  })

  describe('computed method "noUsers"', () => {
    it('should be true if no users are given', () => {
      const { wrapper } = getWrapper({
        props: { user: null, users: [] }
      })
      expect((wrapper.vm as any).noUsers).toBeTruthy()
    })
    it('should be false if users are given', () => {
      const { wrapper } = getWrapper({ props: { user: defaultUser, users: [defaultUser] } })
      expect((wrapper.vm as any).noUsers).toBeFalsy()
    })
  })

  describe("should render either user's info when a user is selected or a hint to select a user when no user is selected, but not both at the same time", () => {
    it('should not show user info when no users are selected', () => {
      const { wrapper } = getWrapper({
        props: { user: null, users: [] }
      })

      expect(wrapper.find('[data-testid="no-user-selected"]').text()).toEqual(
        'Select a user to view details'
      )

      expect(wrapper.findComponent(UserInfoBox).exists()).toBeFalsy()
    })

    it('should show user info when a user is selected', () => {
      const { wrapper } = getWrapper({ props: { user: defaultUser, users: [defaultUser] } })

      expect(wrapper.findComponent(UserInfoBox).exists()).toBeTruthy()
      expect(wrapper.find('[data-testid="no-user-selected"]').exists()).toBeFalsy()
    })
  })
  describe('computed method "multipleUsers"', () => {
    it('should be false if no users are given', () => {
      const { wrapper } = getWrapper({ props: { user: null, users: [] } })
      expect((wrapper.vm as any).multipleUsers).toBeFalsy()
    })
    it('should be false if one user is given', () => {
      const { wrapper } = getWrapper({ props: { user: defaultUser, users: [defaultUser] } })
      expect((wrapper.vm as any).multipleUsers).toBeFalsy()
    })
    it('should be true if multiple users are given', () => {
      const { wrapper } = getWrapper({
        props: { user: null, users: [defaultUser, { displayName: 'user2' } as User] }
      })
      expect((wrapper.vm as any).multipleUsers).toBeTruthy()
    })
  })
})

function getWrapper({ props }: { props: PartialComponentProps<typeof DetailsPanel> }) {
  return {
    wrapper: shallowMount(DetailsPanel, {
      props: {
        users: [],
        roles: [],
        ...props
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
