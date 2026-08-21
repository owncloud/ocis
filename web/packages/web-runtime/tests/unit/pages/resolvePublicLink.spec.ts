import ResolvePublicLink from '../../../src/pages/resolvePublicLink.vue'
import { defaultPlugins, defaultComponentMocks, shallowMount } from '@ownclouders/web-test-helpers'
import { mockDeep } from 'vitest-mock-extended'
import { CapabilityStore, ClientService, useRouteParam, useRouteQuery } from '@ownclouders/web-pkg'
import { DavHttpError, SpaceResource } from '@ownclouders/web-client'
import { authService } from '../../../src/services/auth'
import { ref } from 'vue'
import { DavErrorCode } from '@ownclouders/web-client/webdav'

vi.mock('../../../src/services/auth')

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRouteParam: vi.fn(),
  useRouteQuery: vi.fn()
}))

const selectors = {
  cardFooter: '.oc-card-footer',
  ocSpinnerStub: 'oc-spinner-stub',
  submitButton: '.oc-login-authorize-button'
}

describe('resolvePublicLink', () => {
  it('should display the configuration theme general slogan as the login card footer', () => {
    const { wrapper } = getWrapper()
    const slogan = wrapper.find(selectors.cardFooter)
    expect(slogan.html()).toMatchSnapshot()
  })
  it('should display the loading spinner', () => {
    const { wrapper } = getWrapper({ passwordRequired: true })
    const loading = wrapper.find(selectors.ocSpinnerStub)
    expect(loading.exists()).toBeTruthy()
  })
  describe('password required form', () => {
    it('should display if password is required', async () => {
      const { wrapper } = getWrapper({ passwordRequired: true })
      await (wrapper.vm as any).loadPublicSpaceTask.last

      expect(wrapper.find('form').html()).toMatchSnapshot()
    })
    describe('submit button', () => {
      it('should be set as disabled if "password" is empty', async () => {
        const { wrapper } = getWrapper({ passwordRequired: true })
        await (wrapper.vm as any).loadPublicSpaceTask.last

        expect(wrapper.find(selectors.submitButton).attributes().disabled).toBe('true')
      })
      it('should be set as enabled if "password" is not empty', async () => {
        const { wrapper } = getWrapper({ passwordRequired: true })
        await (wrapper.vm as any).loadPublicSpaceTask.last
        ;(wrapper.vm as any).password = 'password'
        await wrapper.vm.$nextTick()

        expect(wrapper.find(selectors.submitButton).attributes().disabled).toBe('false')
      })
      it('should be disabled and showing spinner when submit button is clicked', async () => {
        const { wrapper } = getWrapper({ passwordRequired: true })
        ;(wrapper.vm as any).resolvePublicLinkTask.perform(true)
        await (wrapper.vm as any).loadPublicSpaceTask.last
        ;(wrapper.vm as any).password = 'password'
        await wrapper.vm.$nextTick()
        const submitButton = wrapper.find(selectors.submitButton)
        expect(submitButton.attributes().disabled).toBe('true')
        expect(submitButton.attributes().showspinner).toBe('true')
      })
      it('should resolve the public link on click', async () => {
        const resolvePublicLinkSpy = vi.spyOn(authService, 'resolvePublicLink')
        const { wrapper } = getWrapper({ passwordRequired: true })
        await (wrapper.vm as any).loadPublicSpaceTask.last
        ;(wrapper.vm as any).password = 'password'
        await wrapper.vm.$nextTick()
        await wrapper.find(selectors.submitButton).trigger('submit')
        await (wrapper.vm as any).resolvePublicLinkTask.last

        expect(resolvePublicLinkSpy).toHaveBeenCalled()
      })
    })
  })
  describe('error message', () => {
    it('should display an error message if the space cannot be resolved', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({ getFileInfoErrorStatusCode: 404 })

      try {
        await (wrapper.vm as any).loadPublicSpaceTask.last
      } catch {}

      expect(wrapper.find('.oc-link-resolve-error-message').text()).toContain(
        'The resource could not be located, it may not exist anymore.'
      )
    })
    it('should display an error message if the space cannot be resolved after entering password', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({
        passwordRequired: true,
        getFileInfoErrorStatusCode: 404
      }) as any
      await wrapper.vm.loadPublicSpaceTask.last
      await expect(wrapper.vm.resolvePublicLinkTask.perform(true)).rejects.toThrow()

      expect(wrapper.find('.oc-link-resolve-error-message').text()).toContain(
        'The resource could not be located, it may not exist anymore.'
      )
    })
    it('should display the blocked message if the link is temporarily blocked', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({ getFileInfoErrorStatusCode: 429 })

      try {
        await (wrapper.vm as any).loadPublicSpaceTask.last
      } catch {}

      expect(wrapper.find('.oc-link-resolve-error-message').text()).toContain(
        'Too many failed password attempts for this link. It has been temporarily blocked, please try again later.'
      )
    })
    it('should display the blocked message and no password form if blocked after entering password', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({
        passwordRequired: true,
        getFileInfoErrorStatusCode: 429
      }) as any
      await wrapper.vm.loadPublicSpaceTask.last
      await expect(wrapper.vm.resolvePublicLinkTask.perform(true)).rejects.toThrow()

      expect(wrapper.find('.oc-link-resolve-error-message').text()).toContain(
        'Too many failed password attempts for this link. It has been temporarily blocked, please try again later.'
      )
      expect(wrapper.find('form').exists()).toBe(false)
    })
    it('should display the server message for an unknown failure after entering password', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({
        passwordRequired: true,
        getFileInfoErrorStatusCode: 500,
        getFileInfoErrorMessage: 'Internal server error'
      }) as any
      await wrapper.vm.loadPublicSpaceTask.last
      await expect(wrapper.vm.resolvePublicLinkTask.perform(true)).rejects.toThrow()

      const message = wrapper.find('.oc-link-resolve-error-message').text()
      expect(message).toContain('Internal server error')
      expect(message).not.toContain('The resource could not be located, it may not exist anymore.')
    })
    it('should display a neutral message if the server sent no message', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({
        getFileInfoErrorStatusCode: 502,
        getFileInfoErrorMessage: 'Unknown error'
      })

      try {
        await (wrapper.vm as any).loadPublicSpaceTask.last
      } catch {}

      const message = wrapper.find('.oc-link-resolve-error-message').text()
      expect(message).toContain('An unexpected error occurred, please try again later.')
      expect(message).not.toContain('Unknown error')
    })
    it('should display a neutral message if the request failed without a status', async () => {
      console.error = vi.fn()
      const { wrapper } = getWrapper({
        getFileInfoError: new TypeError("Cannot read properties of undefined (reading 'status')")
      })

      try {
        await (wrapper.vm as any).loadPublicSpaceTask.last
      } catch {}

      const message = wrapper.find('.oc-link-resolve-error-message').text()
      expect(message).toContain('An unexpected error occurred, please try again later.')
      expect(message).not.toContain('Cannot read properties of undefined')
    })
  })
  describe('internal link', () => {
    it('redirects the user to the login page', async () => {
      const { wrapper, mocks } = getWrapper({ isInternalLink: true })
      await (wrapper.vm as any).loadPublicSpaceTask.last

      expect(mocks.$router.push).toHaveBeenCalledWith({
        name: 'login',
        query: { redirectUrl: '/i/token' }
      })
    })
  })
})

function getWrapper({
  passwordRequired = false,
  isInternalLink = false,
  getFileInfoErrorStatusCode = null,
  getFileInfoErrorMessage = '',
  getFileInfoError = null
}: {
  passwordRequired?: boolean
  isInternalLink?: boolean
  getFileInfoErrorStatusCode?: number
  getFileInfoErrorMessage?: string
  getFileInfoError?: Error
} = {}) {
  const $clientService = mockDeep<ClientService>()
  const spaceResource = mockDeep<SpaceResource>({ driveType: 'public' })

  // loadPublicSpaceTask response
  if (passwordRequired) {
    $clientService.webdav.getFileInfo.mockRejectedValueOnce(
      new DavHttpError('', 'ERR_MISSING_BASIC_AUTH', undefined, 401)
    )
  } else if (isInternalLink) {
    $clientService.webdav.getFileInfo.mockRejectedValueOnce(
      new DavHttpError('', 'ERR_MISSING_BEARER_AUTH', undefined, 401)
    )
  }

  if (getFileInfoError) {
    $clientService.webdav.getFileInfo.mockRejectedValueOnce(getFileInfoError)
  } else if (getFileInfoErrorStatusCode) {
    $clientService.webdav.getFileInfo.mockRejectedValueOnce(
      new DavHttpError(
        getFileInfoErrorMessage,
        'ERR_UNKNOWN' as DavErrorCode,
        undefined,
        getFileInfoErrorStatusCode
      )
    )
  } else {
    $clientService.webdav.getFileInfo.mockResolvedValueOnce(spaceResource)
  }

  const mocks = { ...defaultComponentMocks(), $clientService }

  const capabilities = {
    files_sharing: { federation: { incoming: true, outgoing: true } }
  } satisfies Partial<CapabilityStore['capabilities']>

  vi.mocked(useRouteParam).mockReturnValue(ref('token'))
  vi.mocked(useRouteQuery).mockReturnValue(ref('redirectUrl'))

  return {
    mocks,
    wrapper: shallowMount(ResolvePublicLink, {
      global: {
        plugins: [...defaultPlugins({ piniaOptions: { capabilityState: { capabilities } } })],
        mocks,
        provide: mocks
      }
    })
  }
}
