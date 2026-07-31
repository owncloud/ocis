import CreateGroupModal from '../../../../src/components/Groups/CreateGroupModal.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  mockAxiosReject,
  mockAxiosResolve,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Modal, eventBus, useMessages } from '@ownclouders/web-pkg'
import { useGroupSettingsStore } from '../../../../src/composables'
import { Group } from '@ownclouders/web-client/graph/generated'

describe('CreateGroupModal', () => {
  describe('computed method "isFormInvalid"', () => {
    it('should be true if any data set is invalid', () => {
      const { wrapper } = getWrapper()
      wrapper.vm.formData.displayName.valid = false
      expect(wrapper.vm.isFormInvalid).toBeTruthy()
    })
    it('should be false if no data set is invalid', () => {
      const { wrapper } = getWrapper()
      Object.keys(wrapper.vm.formData).forEach((key) => {
        wrapper.vm.formData[key].valid = true
      })
      expect(wrapper.vm.isFormInvalid).toBeFalsy()
    })
  })
  describe('method "validateDisplayName"', () => {
    it('should be false when displayName is empty', async () => {
      const { wrapper } = getWrapper()
      wrapper.vm.group.displayName = ''
      expect(await wrapper.vm.validateDisplayName()).toBeFalsy()
    })
    it('should be false when displayName is longer than 255 characters', async () => {
      const { wrapper } = getWrapper()
      wrapper.vm.group.displayName = 'n'.repeat(256)
      expect(await wrapper.vm.validateDisplayName()).toBeFalsy()
    })
    it('should be false when displayName is already existing', async () => {
      const { wrapper, mocks } = getWrapper()
      const graphMock = mocks.$clientService.graphAuthenticated
      wrapper.vm.group.displayName = 'admins'
      const getGroupSub = graphMock.groups.getGroup.mockResolvedValue(
        mock<Group>({ displayName: 'admins' })
      )
      expect(await wrapper.vm.validateDisplayName()).toBeFalsy()
      expect(getGroupSub).toHaveBeenCalled()
    })
    it('should be true when displayName is valid', async () => {
      const { wrapper, mocks } = getWrapper()
      const graphMock = mocks.$clientService.graphAuthenticated
      const getGroupSub = graphMock.groups.getGroup.mockRejectedValue(() => mockAxiosReject())
      wrapper.vm.group.displayName = 'users'
      expect(await wrapper.vm.validateDisplayName()).toBeTruthy()
      expect(getGroupSub).toHaveBeenCalled()
    })
  })
  describe('method "onConfirm"', () => {
    it('should not create group if form is invalid', async () => {
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
    it('should create group on success', async () => {
      const { wrapper, mocks } = getWrapper()
      mocks.$clientService.graphAuthenticated.groups.getGroup.mockRejectedValueOnce(new Error(''))

      wrapper.vm.group.displayName = 'foo bar'
      await wrapper.vm.validateDisplayName()

      mocks.$clientService.graphAuthenticated.groups.createGroup.mockResolvedValueOnce(
        mock<Group>({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )

      await wrapper.vm.onConfirm()

      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
      const { upsertGroup } = useGroupSettingsStore()
      expect(upsertGroup).toHaveBeenCalled()
    })

    it('should show message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const { wrapper, mocks } = getWrapper()
      mocks.$clientService.graphAuthenticated.groups.getGroup.mockRejectedValue(new Error(''))

      wrapper.vm.group.displayName = 'foo bar'
      await wrapper.vm.validateDisplayName()

      mocks.$clientService.graphAuthenticated.groups.createGroup.mockRejectedValue(
        mockAxiosResolve({ id: 'e3515ffb-d264-4dfc-8506-6c239f6673b5' })
      )
      await wrapper.vm.onConfirm()

      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
      const { upsertGroup } = useGroupSettingsStore()
      expect(upsertGroup).not.toHaveBeenCalled()
    })
  })
})

function getWrapper() {
  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: shallowMount(CreateGroupModal, {
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
