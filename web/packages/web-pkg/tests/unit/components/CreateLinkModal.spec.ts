import CreateLinkModal from '../../../src/components/CreateLinkModal.vue'
import {
  ComponentProps,
  defaultComponentMocks,
  defaultPlugins,
  mount
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { PasswordPolicyService } from '../../../src/services'
import { usePasswordPolicyService } from '../../../src/composables/passwordPolicyService'
import { AbilityRule, LinkShare, Resource, ShareRole, SpaceResource } from '@ownclouders/web-client'
import { PasswordPolicy } from '@ownclouders/design-system/helpers'
import { useEmbedMode } from '../../../src/composables/embedMode'
import { useLinkTypes } from '../../../src/composables/links'
import { ref } from 'vue'
import { CapabilityStore, Modal, useSharesStore } from '../../../src/composables/piniaStores'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { describe } from 'vitest'

vi.mock('../../../src/composables/embedMode')
vi.mock('../../../src/composables/passwordPolicyService')
vi.mock('../../../src/composables/links', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useLinkTypes: vi.fn()
}))

const selectors = {
  passwordInput: '.link-modal-password-input',
  roleElements: '.role-dropdown-list li',
  contextMenuToggle: '#link-modal-context-menu-toggle',
  confirmBtn: '.link-modal-confirm',
  cancelBtn: '.link-modal-cancel',
  linkRoleDropDownToggle: '.link-role-dropdown-toggle'
}

describe('CreateLinkModal', () => {
  describe('password input', () => {
    it('should be rendered', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find(selectors.passwordInput).exists()).toBeTruthy()
    })
  })
  describe('datepicker', () => {
    it('should be rendered', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.findComponent({ name: 'oc-datepicker' }).exists()).toBeTruthy()
    })
  })
  describe('link role drop', () => {
    it('lists all types as roles', async () => {
      const availableLinkTypes = [SharingLinkType.View, SharingLinkType.Edit]
      const { wrapper } = getWrapper({ availableLinkTypes })
      await wrapper.find(selectors.linkRoleDropDownToggle).trigger('click')

      expect(wrapper.findAll(selectors.roleElements).length).toBe(availableLinkTypes.length)
    })
  })
  describe('method "confirm"', () => {
    it('creates links for all resources', async () => {
      const resources = [mock<Resource>({ isFolder: false }), mock<Resource>({ isFolder: false })]
      const { wrapper } = getWrapper({ resources })
      await wrapper.vm.onConfirm()

      const { addLink } = useSharesStore()
      expect(addLink).toHaveBeenCalledTimes(resources.length)
    })
    it('emits event in embed mode including the created links', async () => {
      const resources = [mock<Resource>({ isFolder: false })]
      const { wrapper, mocks } = getWrapper({ resources, embedModeEnabled: true })
      const link = mock<LinkShare>({ webUrl: 'someurl' })

      const { addLink } = useSharesStore()
      vi.mocked(addLink).mockResolvedValue(link)
      await wrapper.vm.onConfirm()
      expect(mocks.postMessageMock).toHaveBeenCalledWith('owncloud-embed:share', [link.webUrl])
    })
    it('shows error messages for links that failed to be created', async () => {
      const consoleMock = vi.fn(() => undefined)
      vi.spyOn(console, 'error').mockImplementation(consoleMock)
      const resources = [mock<Resource>({ isFolder: false })]
      const { wrapper } = getWrapper({ resources })

      const { addLink } = useSharesStore()
      vi.mocked(addLink).mockRejectedValue({ response: {} })
      await wrapper.vm.onConfirm()

      expect(consoleMock).toHaveBeenCalledTimes(1)
    })
    it('calls clipboard after link creation', async () => {
      global.ClipboardItem = vi.fn().mockImplementation(function (data) {
        return { data }
      }) as any
      global.navigator = {
        ...global.navigator,
        clipboard: {
          write: vi.fn().mockResolvedValue(undefined)
        }
      } as any
      const resources = [mock<Resource>({ isFolder: false })]
      const { wrapper } = getWrapper({ resources })
      const link = mock<LinkShare>({ webUrl: 'https://example.com/link' })

      const { addLink } = useSharesStore()
      vi.mocked(addLink).mockResolvedValue(link)

      await wrapper.vm.onConfirm()

      expect(global.navigator.clipboard.write).toHaveBeenCalled()
      expect(global.ClipboardItem).toHaveBeenCalled()
    })
  })
  describe('action buttons', () => {
    describe('confirm button', () => {
      it('is disabled when password policy is not fulfilled', async () => {
        const { wrapper } = getWrapper({ passwordPolicyFulfilled: false })
        expect(wrapper.find(selectors.confirmBtn).attributes('disabled')).toBeTruthy()
      })
    })
  })
})

function getWrapper({
  resources = [],
  space = undefined,
  defaultLinkType = SharingLinkType.View,
  userCanCreatePublicLinks = true,
  passwordEnforced = false,
  passwordPolicyFulfilled = true,
  embedModeEnabled = false,
  callbackFn = undefined,
  availableLinkTypes = [SharingLinkType.View]
}: {
  resources?: Resource[]
  space?: SpaceResource
  defaultLinkType?: SharingLinkType
  userCanCreatePublicLinks?: boolean
  passwordEnforced?: boolean
  passwordPolicyFulfilled?: boolean
  embedModeEnabled?: boolean
  callbackFn?: ComponentProps<typeof CreateLinkModal>['callbackFn']
  availableLinkTypes?: SharingLinkType[]
} = {}) {
  vi.mocked(usePasswordPolicyService).mockReturnValue(
    mock<PasswordPolicyService>({
      getPolicy: () => mock<PasswordPolicy>({ check: () => passwordPolicyFulfilled })
    })
  )
  vi.mocked(useLinkTypes).mockReturnValue(
    mock<ReturnType<typeof useLinkTypes>>({
      defaultLinkType: ref(defaultLinkType),
      getAvailableLinkTypes: () => availableLinkTypes,
      getLinkRoleByType: () => mock<ShareRole>({ description: 'role', displayName: 'role' }),
      isPasswordEnforcedForLinkType: () => passwordEnforced
    })
  )

  const postMessageMock = vi.fn()
  vi.mocked(useEmbedMode).mockReturnValue(
    mock<ReturnType<typeof useEmbedMode>>({
      isEnabled: ref(embedModeEnabled),
      postMessage: postMessageMock
    })
  )

  const mocks = { ...defaultComponentMocks(), postMessageMock }

  const abilities = [] as AbilityRule[]
  if (userCanCreatePublicLinks) {
    abilities.push({ action: 'create-all', subject: 'PublicLink' })
  }

  const capabilities = {
    files_sharing: {
      public: {
        can_edit: true,
        can_contribute: true,
        alias: true,
        password: { enforced_for: { read_only: passwordEnforced } }
      }
    }
  } satisfies Partial<CapabilityStore['capabilities']>

  return {
    mocks,
    wrapper: mount(CreateLinkModal, {
      props: {
        space,
        resources,
        callbackFn,
        modal: mock<Modal>()
      },
      global: {
        plugins: [
          ...defaultPlugins({ abilities, piniaOptions: { capabilityState: { capabilities } } })
        ],
        mocks,
        provide: mocks,
        stubs: { OcTextInput: true, OcDatepicker: true, OcButton: true }
      }
    })
  }
}
