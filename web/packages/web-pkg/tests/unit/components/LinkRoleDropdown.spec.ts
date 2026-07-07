import LinkRoleDropdown from '../../../src/components/LinkRoleDropdown.vue'
import { defaultComponentMocks, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { ShareRole } from '@ownclouders/web-client'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { useLinkTypes } from '../../../src/composables/links/useLinkTypes'

vi.mock('../../../src/composables/links/useLinkTypes', () => ({
  useLinkTypes: vi.fn()
}))

const selectors = {
  currentRole: '.link-current-role',
  roleOption: '.role-dropdown-list li',
  roleOptionLabel: '.role-dropdown-list-option-label',
  roleDropdownBtn: '.link-role-dropdown-toggle'
}

describe('LinkRoleDropdown', () => {
  it('renders the label of the corresponding role to the given link type', () => {
    const modelValue = SharingLinkType.Internal
    const { wrapper } = getWrapper({ modelValue })

    expect(wrapper.find(selectors.currentRole).text()).toEqual(modelValue)
  })
  it('renders all available role options based on the link types', () => {
    const modelValue = SharingLinkType.Internal
    const availableLinkTypeOptions = [SharingLinkType.Internal, SharingLinkType.View]
    const { wrapper } = getWrapper({ modelValue, availableLinkTypeOptions })

    expect(wrapper.findAll(selectors.roleOption).length).toEqual(availableLinkTypeOptions.length)
    availableLinkTypeOptions.forEach((role, index) => {
      expect(wrapper.findAll(selectors.roleOptionLabel).at(index).text()).toEqual(role)
    })
  })
  it('does not render a button but a span if only one link type is available', () => {
    const availableLinkTypeOptions = [SharingLinkType.View]
    const { wrapper } = getWrapper({ availableLinkTypeOptions })

    expect(wrapper.find(selectors.roleDropdownBtn).exists()).toBeFalsy()
    expect(wrapper.find(selectors.currentRole).exists()).toBeTruthy()
  })
})

function getWrapper({
  modelValue = SharingLinkType.View,
  availableLinkTypeOptions = []
}: { modelValue?: SharingLinkType; availableLinkTypeOptions?: SharingLinkType[] } = {}) {
  vi.mocked(useLinkTypes).mockReturnValue(
    mock<ReturnType<typeof useLinkTypes>>({
      getLinkRoleByType: (value) => mock<ShareRole>({ displayName: value, description: value })
    })
  )

  const mocks = { ...defaultComponentMocks() }

  return {
    mocks,
    wrapper: mount(LinkRoleDropdown, {
      props: {
        modelValue,
        availableLinkTypeOptions
      },
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: mocks
      }
    })
  }
}
