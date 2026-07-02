import account from '../../../src/pages/account.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  mockAxiosReject,
  mockAxiosResolve,
  mount
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  Extension,
  ExtensionPoint,
  OptionsConfig,
  useExtensionRegistry,
  useMessages,
  useSharesStore,
  useResourcesStore
} from '@ownclouders/web-pkg'
import { LanguageOption, SettingsBundle, SettingsValue } from '../../../src/helpers/settings'
import { User } from '@ownclouders/web-client/graph/generated'
import { VueWrapper } from '@vue/test-utils'
import { SpaceResource } from '@ownclouders/web-client'
import { Capabilities } from '@ownclouders/web-client/ocs'

const $route = {
  meta: {
    title: 'Some Title'
  }
}

const selectors = {
  pageTitle: '.oc-page-title',
  loaderStub: 'oc-spinner-stub',
  editUrlButton: '[data-testid="account-page-edit-url-btn"]',
  editPasswordButton: '[data-testid="account-page-edit-password-btn"]',
  logoutButton: '[data-testid="account-page-logout-url-btn"]',
  accountPageInfo: '.account-page-info',
  groupNames: '[data-testid="group-names"]',
  groupNamesEmpty: '[data-testid="group-names-empty"]',
  gdprExport: '[data-testid="gdpr-export"]',
  extensionsSection: '.account-page-extensions',
  crossInstanceReferenceRow: '[data-testid="account-cross-instance-reference-row"]',
  crossInstanceReferenceCopyButton: '[data-testid="account-cross-instance-reference-copy-btn"]'
}

describe('account page', () => {
  beforeEach(() => {
    navigator.clipboard.write = vi.fn()
  })

  describe('public link context', () => {
    it('should render a limited view', async () => {
      const { wrapper } = getWrapper({ isUserContext: false, isPublicLinkContext: true })
      await blockLoadingState(wrapper)

      expect(wrapper.html()).toMatchSnapshot()
    })
  })

  describe('header section', () => {
    describe('edit url button', () => {
      it('should be displayed if defined via config', async () => {
        const { wrapper } = getWrapper({
          accountEditLink: { href: '/' }
        })
        await blockLoadingState(wrapper)

        const editUrlButton = wrapper.find(selectors.editUrlButton)
        expect(editUrlButton.html()).toMatchSnapshot()
      })
      it('should not be displayed if not defined via config', async () => {
        const { wrapper } = getWrapper()
        await blockLoadingState(wrapper)

        const editUrlButton = wrapper.find(selectors.editUrlButton)
        expect(editUrlButton.exists()).toBeFalsy()
      })
    })
  })

  describe('account information section', () => {
    it('displays basic user information', async () => {
      const { wrapper } = getWrapper({
        user: mock<User>({
          onPremisesSamAccountName: 'some-username',
          displayName: 'some-displayname',
          mail: 'some-email',
          memberOf: []
        })
      })
      await blockLoadingState(wrapper)

      const accountPageInfo = wrapper.find(selectors.accountPageInfo)
      expect(accountPageInfo.html()).toMatchSnapshot()
    })

    describe('group membership', () => {
      it('displays message if not member of any groups', async () => {
        const { wrapper } = getWrapper()
        await blockLoadingState(wrapper)

        const groupNamesEmpty = wrapper.find(selectors.groupNamesEmpty)
        expect(groupNamesEmpty.exists()).toBeTruthy()
      })
      it('displays group names', async () => {
        const { wrapper } = getWrapper({
          user: mock<User>({
            memberOf: [{ displayName: 'one' }, { displayName: 'two' }, { displayName: 'three' }]
          })
        })
        await blockLoadingState(wrapper)

        const groupNames = wrapper.find(selectors.groupNames)
        expect(groupNames.html()).toMatchSnapshot()
      })
    })

    describe('Logout from all devices link', () => {
      it('should render the logout from active devices if logoutUrl is provided', async () => {
        const { wrapper } = getWrapper()
        await blockLoadingState(wrapper)

        expect(wrapper.find('[data-testid="logout"]').exists()).toBe(true)
      })
      it("shouldn't render the logout from active devices if logoutUrl isn't provided", async () => {
        const { wrapper } = getWrapper()
        await blockLoadingState(wrapper)

        wrapper.vm.logoutUrl = undefined
        expect(wrapper.find('[data-testid="logout"]').exists()).toBe(true)
      })
      it('should use url from configuration manager', async () => {
        const { wrapper } = getWrapper()
        await blockLoadingState(wrapper)

        const logoutButton = wrapper.find(selectors.logoutButton)
        expect(logoutButton.attributes('href')).toBe('https://account-manager/logout')
      })
    })

    it('should render the cross-instance reference if set', async () => {
      const { wrapper } = getWrapper({
        user: mock<User>({
          onPremisesSamAccountName: 'some-username',
          displayName: 'some-displayname',
          mail: 'some-email',
          memberOf: [],
          crossInstanceReference: 'some-cross-instance-reference'
        })
      })
      await blockLoadingState(wrapper)
      expect(wrapper.find(selectors.crossInstanceReferenceRow).exists()).toBe(true)
    })

    it('should not render the cross-instance reference if not set', async () => {
      const { wrapper } = getWrapper()
      await blockLoadingState(wrapper)

      expect(wrapper.find(selectors.crossInstanceReferenceRow).exists()).toBeFalsy()
    })

    it('should copy the cross-instance reference to the clipboard when the copy button is clicked', async () => {
      const { wrapper } = getWrapper({
        user: mock<User>({
          onPremisesSamAccountName: 'some-username',
          displayName: 'some-displayname',
          mail: 'some-email',
          memberOf: [],
          crossInstanceReference: 'some-cross-instance-reference'
        })
      })

      await blockLoadingState(wrapper)
      await wrapper.find(selectors.crossInstanceReferenceCopyButton).trigger('click')
      expect(navigator.clipboard.write).toHaveBeenCalledWith([
        new ClipboardItem({
          'text/plain': new Blob(['some-cross-instance-reference'], { type: 'text/plain' })
        })
      ])
    })
  })

  describe('Preferences section', () => {
    describe('change password button', () => {
      it('should be displayed if not disabled via capability', async () => {
        const { wrapper } = getWrapper({
          capabilities: {
            graph: { users: { change_password_self_disabled: false }, tags: { max_tag_length: 30 } }
          }
        })
        await blockLoadingState(wrapper)

        const editPasswordButton = wrapper.find(selectors.editPasswordButton)
        expect(editPasswordButton.exists()).toBeTruthy()
      })
      it('should not be displayed if disabled via capability', async () => {
        const { wrapper } = getWrapper({
          capabilities: {
            graph: { users: { change_password_self_disabled: true }, tags: { max_tag_length: 30 } }
          }
        })
        await blockLoadingState(wrapper)

        const editPasswordButton = wrapper.find(selectors.editPasswordButton)
        expect(editPasswordButton.exists()).toBeFalsy()
      })
    })
  })

  describe('Method "updateDisableEmailNotifications', () => {
    it('should show a message on success', async () => {
      const { wrapper, mocks } = getWrapper()
      await blockLoadingState(wrapper)

      mocks.$clientService.httpAuthenticated.post.mockResolvedValueOnce(
        mockAxiosResolve({ value: { id: 'settings-language' } })
      )
      await wrapper.vm.updateDisableEmailNotifications(true)
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
    })
    it('should show a message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const { wrapper, mocks } = getWrapper()
      await blockLoadingState(wrapper)

      mocks.$clientService.httpAuthenticated.post.mockImplementation(() => mockAxiosReject('err'))
      await wrapper.vm.updateDisableEmailNotifications(true)
      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
    })
  })

  describe('Method "updateSelectedLanguage', () => {
    it('should show a message on success', async () => {
      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.graphAuthenticated.users.editMe.mockResolvedValueOnce(undefined)
      await wrapper.vm.updateSelectedLanguage({ value: 'en' } as LanguageOption)
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
    })

    it('should fetch share roles', async () => {
      const { wrapper } = getWrapper({})
      await blockLoadingState(wrapper)
      const sharesStore = useSharesStore()

      await wrapper.vm.updateSelectedLanguage({ value: 'en' } as LanguageOption)

      expect(sharesStore.fetchShareRolesDefinitions).toHaveBeenCalled()

      expect(sharesStore.setGraphRoles).toHaveBeenCalled()
    })

    it('should show a message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.graphAuthenticated.users.editMe.mockRejectedValue(new Error())
      await wrapper.vm.updateSelectedLanguage({ value: 'en' } as LanguageOption)
      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
    })

    it('should refetch settings bundles', async () => {
      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.graphAuthenticated.users.editMe.mockResolvedValueOnce(undefined)
      await wrapper.vm.updateSelectedLanguage({ value: 'en' } as LanguageOption)
      expect(mocks.$clientService.httpAuthenticated.post).toHaveBeenCalledWith(
        '/api/v0/settings/bundles-list',
        {},
        { signal: expect.any(AbortSignal) }
      )
    })
  })

  describe('Method "updateViewOptionsWebDavDetails', () => {
    it('should show a message on success', async () => {
      const { wrapper } = getWrapper({})
      await blockLoadingState(wrapper)

      await wrapper.vm.updateViewOptionsWebDavDetails(true)
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()

      const { setAreWebDavDetailsShown } = useResourcesStore()
      expect(setAreWebDavDetailsShown).toHaveBeenCalled()
    })
  })

  describe('Extensions section', () => {
    it('should be hidden if no extension points offer preferences', async () => {
      const { wrapper } = getWrapper({})
      await blockLoadingState(wrapper)

      expect(wrapper.find(selectors.extensionsSection).exists()).toBeFalsy()
    })

    it('should be hidden if an extension point only has 1 or less extensions', async () => {
      const extensionPointMock = mock<ExtensionPoint<Extension>>({
        userPreference: {
          label: 'example-extension-point'
        }
      })
      const { wrapper } = getWrapper({
        extensionPoints: [extensionPointMock]
      })
      await blockLoadingState(wrapper)

      expect(wrapper.find(selectors.extensionsSection).exists()).toBeFalsy()
    })

    it('should be visible if an extension point has at least 2 extensions', async () => {
      const extensionPoint = mock<ExtensionPoint<Extension>>({
        id: 'test-extension-point',
        multiple: false,
        defaultExtensionId: 'foo-2',
        userPreference: {
          label: 'Foo container'
        }
      })
      const extensions = [
        mock<Extension>({
          id: 'foo-1',
          userPreference: {
            optionLabel: 'Foo 1'
          }
        }),
        mock<Extension>({
          id: 'foo-2',
          userPreference: {
            optionLabel: 'Foo 2'
          }
        })
      ]
      const { wrapper } = getWrapper({
        extensionPoints: [extensionPoint],
        extensions
      })
      await blockLoadingState(wrapper)

      expect(wrapper.find(selectors.extensionsSection).exists()).toBeTruthy()
    })
  })

  describe('Method "updateMultiChoiceSettingsValue"', () => {
    it('should show a message on success', async () => {
      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.httpAuthenticated.post.mockResolvedValueOnce(
        mockAxiosResolve({
          value: { identifier: { setting: 'setting-id' }, value: { id: 'value-id' } }
        })
      )
      await wrapper.vm.updateMultiChoiceSettingsValue('setting-id', 'setting-key', true)
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
    })

    it('should show a message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.httpAuthenticated.post.mockImplementation(() => mockAxiosReject('err'))
      await wrapper.vm.updateMultiChoiceSettingsValue('setting-id', 'setting-key', true)
      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
    })
  })

  describe('Method "updateSingleChoiceValue"', () => {
    it('should show a message on success', async () => {
      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.httpAuthenticated.post.mockResolvedValueOnce(
        mockAxiosResolve({
          value: { identifier: { setting: 'setting-id' }, value: { id: 'value-id' } }
        })
      )
      await wrapper.vm.updateSingleChoiceValue('setting-id', {
        displayValue: 'Daily',
        value: { stringValue: 'daily' }
      })
      const { showMessage } = useMessages()
      expect(showMessage).toHaveBeenCalled()
    })

    it('should show a message on error', async () => {
      vi.spyOn(console, 'error').mockImplementation(() => undefined)

      const { wrapper, mocks } = getWrapper({})
      await blockLoadingState(wrapper)

      mocks.$clientService.httpAuthenticated.post.mockImplementation(() => mockAxiosReject('err'))
      await wrapper.vm.updateSingleChoiceValue('setting-id', {
        displayValue: 'Daily',
        value: { stringValue: 'daily' }
      })
      const { showErrorMessage } = useMessages()
      expect(showErrorMessage).toHaveBeenCalled()
    })
  })
})

const blockLoadingState = async (wrapper: VueWrapper<any, any>) => {
  await wrapper.vm.loadAccountBundleTask.last
  await wrapper.vm.loadValuesListTask.last
  await wrapper.vm.loadGraphUserTask.last
}

function getWrapper({
  user = mock<User>({ memberOf: [] }),
  capabilities = {},
  accountEditLink = undefined,
  spaces = [],
  isPublicLinkContext = false,
  isUserContext = true,
  extensionPoints = [],
  extensions = []
}: {
  user?: User
  capabilities?: Partial<Capabilities['capabilities']>
  accountEditLink?: OptionsConfig['accountEditLink']
  spaces?: SpaceResource[]
  isPublicLinkContext?: boolean
  isUserContext?: boolean
  extensionPoints?: ExtensionPoint<Extension>[]
  extensions?: Extension[]
} = {}) {
  const plugins = defaultPlugins({
    piniaOptions: {
      userState: { user },
      authState: {
        userContextReady: isUserContext,
        publicLinkContextReady: isPublicLinkContext
      },
      spacesState: { spaces },
      capabilityState: { capabilities },
      configState: {
        options: {
          logoutUrl: 'https://account-manager/logout',
          ...(accountEditLink && { accountEditLink })
        }
      }
    }
  })

  const { getExtensionPoints, requestExtensions } = useExtensionRegistry()
  vi.mocked(getExtensionPoints).mockReturnValue(extensionPoints)
  vi.mocked(requestExtensions).mockReturnValue(extensions)

  const mocks = {
    ...defaultComponentMocks(),
    $route
  }

  mocks.$clientService.httpAuthenticated.post.mockImplementation((url) => {
    let response = {}

    if (url.endsWith('bundles-list')) {
      response = { bundles: [mock<SettingsBundle>()] }
    }
    if (url.endsWith('values-list')) {
      response = { values: [mock<SettingsValue>()] }
    }

    return Promise.resolve(mockAxiosResolve(response))
  })
  mocks.$clientService.graphAuthenticated.users.getMe.mockResolvedValue(mock<User>({ id: '1' }))

  return {
    mocks,
    wrapper: mount(account, {
      global: {
        plugins,
        mocks,
        provide: mocks,
        stubs: {
          'extension-preference': true
        }
      }
    })
  }
}
