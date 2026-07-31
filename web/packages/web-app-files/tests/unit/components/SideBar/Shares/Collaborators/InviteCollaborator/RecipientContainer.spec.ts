import RecipientContainer from '../../../../../../../src/components/SideBar/Shares/Collaborators/InviteCollaborator/RecipientContainer.vue'
import { CollaboratorAutoCompleteItem, ShareTypes } from '@ownclouders/web-client'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { CapabilityStore } from '@ownclouders/web-pkg'

vi.mock('../../../../../../../src/helpers/user/avatarUrl', () => ({
  avatarUrl: vi.fn().mockReturnValue('avatarUrl')
}))

const getRecipient = (shareType: number = ShareTypes.user.value): CollaboratorAutoCompleteItem => ({
  displayName: 'Albert Einstein',
  id: 'einstein',
  shareType,
  attributes: []
})

describe('InviteCollaborator RecipientContainer', () => {
  describe('renders a recipient with a deselect button', () => {
    it.each(ShareTypes.authenticated)('different recipients for different shareTypes', (type) => {
      const recipient = getRecipient(type.value)
      const { wrapper } = getMountedWrapper(recipient)
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  it('displays an avatar image if capability is present', async () => {
    const recipient = getRecipient()
    const { wrapper } = getMountedWrapper(recipient, true)
    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('emits an event if deselect button is clicked', async () => {
    const recipient = getRecipient()
    const { wrapper } = getMountedWrapper(recipient, true)
    const spyOnDeselect = (wrapper.vm.deselect as any).mockImplementation()
    const button = wrapper.find('.files-share-invite-recipient-btn-remove')
    await button.trigger('click')
    expect(spyOnDeselect).toHaveBeenCalledTimes(1)
  })
})

function getMountedWrapper(recipient: CollaboratorAutoCompleteItem, avatarsEnabled = false) {
  const capabilities = {
    files_sharing: {
      user: {
        profile_picture: avatarsEnabled
      }
    }
  } satisfies Partial<CapabilityStore['capabilities']>

  return {
    wrapper: mount(RecipientContainer, {
      props: {
        recipient,
        deselect: vi.fn()
      },
      global: {
        plugins: [...defaultPlugins({ piniaOptions: { capabilityState: { capabilities } } })],
        stubs: {
          OcIcon: true
        }
      }
    })
  }
}
