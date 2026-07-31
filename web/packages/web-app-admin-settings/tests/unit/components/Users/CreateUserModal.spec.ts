import CreateUserModal from '../../../../src/components/Users/CreateUserModal.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  mockAxiosReject,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Modal, eventBus, useMessages } from '@ownclouders/web-pkg'
import { useUserSettingsStore } from '../../../../src/composables/stores/userSettings'
import { User } from '@ownclouders/web-client/graph/generated'

describe('CreateUserModal', () => {
  describe('computed method "isFormInvalid"', () => {
    it('should be true if any data set is invalid', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).formData.userName.valid = false
      expect((wrapper.vm as any).isFormInvalid).toBeTruthy()
    })
  })
  it('should be false if no data set is invalid', () => {
    const { wrapper } = getWrapper()
    Object.keys((wrapper.vm as any).formData).forEach((key) => {
      ;(wrapper.vm as any).formData[key].valid = true
    })
    expect((wrapper.vm as any).isFormInvalid).toBeFalsy()
  })

  describe('method "validateUserName"', () => {
    it('should be false when userName is empty', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.onPremisesSamAccountName = ''
      expect(await (wrapper.vm as any).validateUserName()).toBeFalsy()
    })
    it('should be false when userName is longer than 255 characters', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'n'.repeat(256)
      expect(await (wrapper.vm as any).validateUserName()).toBeFalsy()
    })
    it('should be false when userName contains white spaces', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'jan owncCloud'
      expect(await (wrapper.vm as any).validateUserName()).toBeFalsy()
    })
    it('should be false when userName starts with a numeric value', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.onPremisesSamAccountName = '1moretry'
      expect(await (wrapper.vm as any).validateUserName()).toBeFalsy()
    })
    it('should be false when userName is already existing', async () => {
      const { wrapper, mocks } = getWrapper()
      const graphMock = mocks.$clientService.graphAuthenticated
      const getUserStub = graphMock.users.getUser.mockResolvedValue(
        mock<User>({ onPremisesSamAccountName: 'jan' })
      )
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'jan'
      expect(await (wrapper.vm as any).validateUserName()).toBeFalsy()
      expect(getUserStub).toHaveBeenCalled()
    })
    it('should be true when userName is valid', async () => {
      const { wrapper, mocks } = getWrapper()
      const graphMock = mocks.$clientService.graphAuthenticated
      const getUserStub = graphMock.users.getUser.mockRejectedValue(() => mockAxiosReject())
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'jana'
      expect(await (wrapper.vm as any).validateUserName()).toBeTruthy()
      expect(getUserStub).toHaveBeenCalled()
    })
    it('should be true when userName is an email address', async () => {
      const { wrapper, mocks } = getWrapper()
      const graphMock = mocks.$clientService.graphAuthenticated
      const getUserStub = graphMock.users.getUser.mockRejectedValue(() => mockAxiosReject())
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'sk@domain.tld'
      expect(await (wrapper.vm as any).validateUserName()).toBeTruthy()
      expect(getUserStub).toHaveBeenCalled()
    })
  })

  describe('method "validateDisplayName"', () => {
    it('should be false when displayName is empty', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.displayName = ''
      expect((wrapper.vm as any).validateDisplayName()).toBeFalsy()
    })
    it('should be false when displayName is longer than 255 characters', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.displayName = 'n'.repeat(256)
      expect(await (wrapper.vm as any).validateDisplayName()).toBeFalsy()
    })
    it('should be true when displayName is valid', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.displayName = 'jana'
      expect((wrapper.vm as any).validateDisplayName()).toBeTruthy()
    })
  })

  describe('method "validateEmail"', () => {
    it('should be false when email is invalid', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.mail = 'jana@'
      expect((wrapper.vm as any).validateEmail()).toBeFalsy()
    })

    it('should be true when email is valid', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.mail = 'jana@owncloud.com'
      expect((wrapper.vm as any).validateEmail()).toBeTruthy()
    })
  })

  describe('method "validatePassword"', () => {
    it('should be false when password is empty', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.passwordProfile.password = ''
      expect((wrapper.vm as any).validatePassword()).toBeFalsy()
    })

    it('should be true when password is valid', () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).user.passwordProfile.password = 'asecret'
      expect((wrapper.vm as any).validatePassword()).toBeTruthy()
    })
  })
  describe('method "onConfirm"', () => {
    it('should not create user if form is invalid', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)
      const { wrapper } = getWrapper()

      const eventSpy = vi.spyOn(eventBus, 'publish')
      try {
        await wrapper.vm.onConfirm()
      } catch {}

      const { showMessage } = useMessages()
      expect(showMessage).not.toHaveBeenCalled()
      expect(eventSpy).not.toHaveBeenCalled()
    })
    it('should create user on success', async () => {
      const { wrapper, mocks } = getWrapper()
      mocks.$clientService.graphAuthenticated.users.getUser.mockRejectedValueOnce(new Error(''))
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'foo'
      await (wrapper.vm as any).validateUserName()
      ;(wrapper.vm as any).user.displayName = 'foo bar'
      ;(wrapper.vm as any).validateDisplayName()
      ;(wrapper.vm as any).user.mail = 'foo@bar.com'
      ;(wrapper.vm as any).validateEmail()
      ;(wrapper.vm as any).user.passwordProfile.password = 'asecret'
      ;(wrapper.vm as any).validatePassword()

      mocks.$clientService.graphAuthenticated.users.createUser.mockResolvedValue(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )
      mocks.$clientService.graphAuthenticated.users.getUser.mockResolvedValueOnce(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )

      await wrapper.vm.onConfirm()

      const { upsertUser } = useUserSettingsStore()
      expect(upsertUser).toHaveBeenCalled()
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
    })

    it('should show message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const { wrapper, mocks } = getWrapper()
      mocks.$clientService.graphAuthenticated.users.getUser.mockRejectedValue(new Error(''))
      ;(wrapper.vm as any).user.onPremisesSamAccountName = 'foo'
      await (wrapper.vm as any).validateUserName()
      ;(wrapper.vm as any).user.displayName = 'foo bar'
      ;(wrapper.vm as any).validateDisplayName()
      ;(wrapper.vm as any).user.mail = 'foo@bar.com'
      ;(wrapper.vm as any).validateEmail()
      ;(wrapper.vm as any).user.passwordProfile.password = 'asecret'
      ;(wrapper.vm as any).validatePassword()

      mocks.$clientService.graphAuthenticated.users.createUser.mockResolvedValue(
        mock<User>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )
      const eventSpy = vi.spyOn(eventBus, 'publish')
      await wrapper.vm.onConfirm()

      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
      expect(eventSpy).not.toHaveBeenCalled()
    })
  })
})

function getWrapper() {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(CreateUserModal, {
      props: {
        modal: mock<Modal>()
      },
      global: {
        mocks,
        provide: mocks,
        plugins: [...defaultPlugins()]
      }
    })
  }
}
