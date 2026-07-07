import EditDropdown from '../../../../../../src/components/SideBar/Shares/Collaborators/EditDropdown.vue'
import { defaultPlugins, PartialComponentProps, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import { OcButton } from '@ownclouders/design-system/components'

const selectors = {
  editBtn: '.collaborator-edit-dropdown-options-btn',
  removeShareSection: '.collaborator-edit-dropdown-options-list-remove',
  expireDateMenuAction: '.set-expiration-date',
  showAccessDetailsAction: '.show-access-details',
  navigateToParentAction: '.navigate-to-parent'
}

describe('EditDropdown', () => {
  describe('edit button', () => {
    it('is being rendered correctly', () => {
      const { wrapper } = getWrapper()
      const btn = wrapper.findComponent<typeof OcButton>(selectors.editBtn)

      expect(wrapper.find(selectors.editBtn).exists()).toBeTruthy()
      expect(btn.props('disabled')).toBeFalsy()
    })
    it('is being disabled when locked', () => {
      const { wrapper } = getWrapper({ isLocked: true })
      const btn = wrapper.findComponent<typeof OcButton>(selectors.editBtn)
      expect(btn.props('disabled')).toBeTruthy()
    })
  })
  describe('remove share action', () => {
    it('is being rendered when canEdit is true', () => {
      const { wrapper } = getWrapper({ canEdit: true })
      expect(wrapper.find(selectors.removeShareSection).exists()).toBeTruthy()
    })
    it('is not being rendered when canEdit is false', () => {
      const { wrapper } = getWrapper({ canEdit: false })
      expect(wrapper.find(selectors.removeShareSection).exists()).toBeFalsy()
    })
  })
  describe('expiration date', () => {
    it('is being rendered when canEdit is true', () => {
      const { wrapper } = getWrapper({ canEdit: true })
      expect(wrapper.find(selectors.expireDateMenuAction).exists()).toBeTruthy()
    })
    it('is not being rendered when canEdit is false', () => {
      const { wrapper } = getWrapper({ canEdit: false })
      expect(wrapper.find(selectors.expireDateMenuAction).exists()).toBeFalsy()
    })
  })
  describe('navigate to parent action', () => {
    it('is being rendered when sharedParentRoute is given', () => {
      const { wrapper } = getWrapper({
        sharedParentRoute: { params: { driveAliasAndItem: '/folder' } }
      })
      expect(wrapper.find(selectors.navigateToParentAction).exists()).toBeTruthy()
    })
    it('is not being rendered when sharedParentRoute is not given', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.navigateToParentAction).exists()).toBeFalsy()
    })
  })
  describe('show access details action', () => {
    it('is being rendered', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.showAccessDetailsAction).exists()).toBeTruthy()
    })
  })
})

function getWrapper(props: PartialComponentProps<typeof EditDropdown> = {}) {
  return {
    wrapper: shallowMount(EditDropdown, {
      props: {
        canEdit: true,
        shareCategory: 'user',
        accessDetails: [],
        ...props
      },
      global: {
        plugins: [...defaultPlugins()],
        provide: { resource: mock<Resource>() },
        stubs: { OcDrop: false, OcList: false, ContextMenuItem: false }
      }
    })
  }
}
