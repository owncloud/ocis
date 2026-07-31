import AutocompleteItem from '../../../../../../../src/components/SideBar/Shares/Collaborators/InviteCollaborator/AutocompleteItem.vue'
import { CollaboratorAutoCompleteItem, ShareTypes } from '@ownclouders/web-client'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

describe('AutocompleteItem component', () => {
  it.each(ShareTypes.all)('sets a class that reflects the share type', (shareType) => {
    const { wrapper } = createWrapper({ shareType: shareType.value })
    expect(wrapper.find('div').attributes('class')).toContain(
      `files-collaborators-search-${shareType.key}`
    )
  })
  it.each(ShareTypes.all)(
    'displays the correct image/icon according to the shareType',
    (shareType) => {
      const { wrapper } = createWrapper({ shareType: shareType.value })
      const isUserShareType = shareType.key === ShareTypes.user.key
      if (isUserShareType) {
        expect(wrapper.find('avatar-image-stub').exists()).toBeTruthy()
        expect(wrapper.find('oc-avatar-item-stub').exists()).toBeFalsy()
      } else {
        expect(wrapper.find('avatar-image-stub').exists()).toBeFalsy()
        expect(wrapper.find('oc-avatar-item-stub').exists()).toBeTruthy()
        expect(wrapper.find('oc-avatar-item-stub').attributes().icon).toEqual(shareType.icon)
      }
    }
  )
  describe('avatar image', () => {
    it('sets the userId', () => {
      const { wrapper } = createWrapper({
        shareType: ShareTypes.user.value,
        id: 'the-user-id'
      })
      expect(wrapper.find('avatar-image-stub').attributes('userid')).toEqual('the-user-id')
    })
    it('sets the user-name', () => {
      const { wrapper } = createWrapper({
        shareType: ShareTypes.user.value,
        displayName: 'the-user-name'
      })
      expect(wrapper.find('avatar-image-stub').attributes('user-name')).toEqual('the-user-name')
    })
  })
  describe('autocomplete text', () => {
    it('shows the user name', () => {
      const { wrapper } = createWrapper({ displayName: 'Alice Hansen' })
      expect(wrapper.find('.files-collaborators-autocomplete-username').text()).toEqual(
        'Alice Hansen'
      )
    })
    it.each([ShareTypes.user.value, ShareTypes.group.value])(
      'hides share type for users and groups',
      (shareType: number) => {
        const { wrapper } = createWrapper({ shareType })
        expect(wrapper.find('.files-collaborators-autocomplete-share-type').exists()).toBeFalsy()
      }
    )
    it('shows share type for guests', () => {
      const { wrapper } = createWrapper({ shareType: ShareTypes.guest.value })
      expect(wrapper.find('.files-collaborators-autocomplete-share-type').text()).toEqual('(Guest)')
    })
  })
  describe('additional info', () => {
    it('shows the attributes for a user if given', () => {
      const attributes = ['foo', 'bar']
      const { wrapper } = createWrapper({ shareType: ShareTypes.user.value, attributes })
      expect(wrapper.find('.files-collaborators-autocomplete-additionalInfo').text()).toEqual(
        attributes.join(' · ')
      )
    })
    it('shows the email for a user if given', () => {
      const mail = 'foo@bar.com'
      const { wrapper } = createWrapper({ shareType: ShareTypes.user.value, mail })
      expect(wrapper.find('.files-collaborators-autocomplete-additionalInfo').text()).toEqual(mail)
    })
    it('shows the onPremisesSamAccountName for a user if no mail given', () => {
      const onPremisesSamAccountName = 'fooBar'
      const { wrapper } = createWrapper({
        shareType: ShareTypes.user.value,
        onPremisesSamAccountName
      })
      expect(wrapper.find('.files-collaborators-autocomplete-additionalInfo').text()).toEqual(
        onPremisesSamAccountName
      )
    })
    it('does not show for group shares', () => {
      const { wrapper } = createWrapper({ shareType: ShareTypes.group.value })
      expect(wrapper.find('.files-collaborators-autocomplete-additionalInfo').exists()).toBeFalsy()
    })
  })
})

function createWrapper({
  shareType = ShareTypes.user.value,
  id = '',
  displayName = '',
  mail = '',
  onPremisesSamAccountName = '',
  attributes = []
}: Partial<CollaboratorAutoCompleteItem>) {
  return {
    wrapper: shallowMount(AutocompleteItem, {
      props: {
        item: { shareType, id, displayName, mail, onPremisesSamAccountName, attributes }
      },
      global: {
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins()],
        stubs: { 'avatar-image': true }
      }
    })
  }
}
