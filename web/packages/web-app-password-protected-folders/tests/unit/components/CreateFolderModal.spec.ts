import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import CreateFolderModal from '../../../src/components/CreateFolderModal.vue'
import { useCreateFileHandler } from '../../../src/composables/useCreateFileHandler'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { VueWrapper } from '@vue/test-utils'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { PasswordPolicyService, usePasswordPolicyService } from '@ownclouders/web-pkg'
import { PasswordPolicy } from '../../../../design-system/src/helpers/types'
import LinkRoleDropdown from '@ownclouders/web-pkg/src/components/LinkRoleDropdown.vue'

vi.mock('../../../src/composables/useCreateFileHandler', () => ({
  useCreateFileHandler: vi.fn().mockReturnValue({ createFileHandler: vi.fn() })
}))

vi.mock('../../../../web-pkg/src/composables/passwordPolicyService')

const currentFolder = mock<Resource>()
const currentSpace = mock<SpaceResource>({ driveType: 'personal' })

const SELECTORS = Object.freeze({
  inputFolderName: '#input-folder-name',
  inputFolderPassword: '#input-folder-password'
})

describe('CreateFolderModal', () => {
  it('should call "createFileHandler" when form is valid', () => {
    const { wrapper } = getWrapper({ passwordPolicyFulfilled: true })

    const folderNameInput = wrapper.findComponent(SELECTORS.inputFolderName) as VueWrapper
    const passwordInput = wrapper.findComponent(SELECTORS.inputFolderPassword) as VueWrapper
    const permissionsInput = wrapper.findComponent(LinkRoleDropdown) as VueWrapper

    folderNameInput.vm.$emit('update:modelValue', 'name')
    passwordInput.vm.$emit('update:modelValue', 'password')
    permissionsInput.vm.$emit('update:modelValue', SharingLinkType.Edit)

    wrapper.vm.onConfirm()

    expect(useCreateFileHandler().createFileHandler).toHaveBeenCalledWith({
      fileName: 'name',
      password: 'password',
      personalSpace: currentSpace,
      currentSpace: currentSpace,
      currentFolder,
      type: SharingLinkType.Edit
    })
  })

  it('should not call "createFileHandler" when form is invalid', () => {
    const { wrapper } = getWrapper({ passwordPolicyFulfilled: false })

    const folderNameInput = wrapper.findComponent(SELECTORS.inputFolderName) as VueWrapper
    folderNameInput.vm.$emit('update:modelValue', 'name')

    expect(wrapper.vm.onConfirm()).rejects.toThrow()
    expect(useCreateFileHandler().createFileHandler).not.toHaveBeenCalled()
  })
})

function getWrapper({
  passwordPolicyFulfilled = true
}: { passwordPolicyFulfilled?: boolean } = {}) {
  vi.mocked(usePasswordPolicyService).mockReturnValue(
    mock<PasswordPolicyService>({
      getPolicy: () => mock<PasswordPolicy>({ check: () => passwordPolicyFulfilled })
    })
  )

  const mocks = defaultComponentMocks()

  return {
    wrapper: shallowMount(CreateFolderModal, {
      global: {
        plugins: defaultPlugins({
          piniaOptions: {
            resourcesStore: { currentFolder },
            spacesState: { currentSpace, spaces: [currentSpace] }
          }
        }),
        mocks,
        provide: mocks
      }
    }),
    mocks
  }
}
