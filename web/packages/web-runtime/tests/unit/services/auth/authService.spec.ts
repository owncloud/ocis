import { ConfigStore, useAuthStore, useConfigStore } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { Router } from 'vue-router'
import { AuthService } from '../../../../src/services/auth/authService'
import { UserManager } from '../../../../src/services/auth/userManager'
import { RouteLocation, createRouter, createTestingPinia } from '@ownclouders/web-test-helpers'
import { User } from 'oidc-client-ts'

const mockUpdateContext = vi.fn()
console.debug = vi.fn()

vi.mock('../../../../src/services/auth/userManager')

const initAuthService = ({
  authService,
  configStore = null,
  router = null
}: {
  authService: AuthService
  configStore?: ConfigStore
  router?: Router
}) => {
  createTestingPinia()
  const authStore = useAuthStore()
  configStore = configStore || useConfigStore()

  authService.initialize(configStore, null, router, null, null, null, authStore, null, null)
}

describe('AuthService', () => {
  describe('signInCallback', () => {
    it.each([
      ['/', '/', {}],
      ['/?details=sharing', '/', { details: 'sharing' }],
      [
        '/external?contextRouteName=files-spaces-personal&fileId=0f897576',
        '/external',
        {
          contextRouteName: 'files-spaces-personal',
          fileId: '0f897576'
        }
      ]
    ])(
      'parses query params and passes them explicitly to router.replace: %s => %s %s',
      async (url, path, query: Record<string, string>) => {
        const authService = new AuthService()

        Object.defineProperty(authService, 'userManager', {
          value: {
            signinRedirectCallback: vi.fn(),
            getUser: vi.fn().mockResolvedValue(null),
            getAndClearPostLoginRedirectUrl: () => url
          }
        })

        const router = createRouter()
        const replaceSpy = vi.spyOn(router, 'replace')

        initAuthService({ authService, router })
        await authService.signInCallback()

        expect(replaceSpy).toHaveBeenCalledWith({
          path,
          query
        })
      }
    )
  })

  describe('initializeContext', () => {
    it('when embed mode is disabled and access_token is present, should call updateContext', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getAccessToken: vi.fn().mockResolvedValue('access-token'),
          getUser: vi.fn().mockResolvedValue(mock<User>({ expires_in: 3600 })),
          updateContext: mockUpdateContext
        })
      })

      initAuthService({ authService })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(mockUpdateContext).toHaveBeenCalledWith('access-token', true)
    })

    it('when embed mode is disabled and access_token is not present, should not call updateContext', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getAccessToken: vi.fn().mockResolvedValue(null),
          updateContext: mockUpdateContext
        })
      })

      initAuthService({ authService })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(mockUpdateContext).not.toHaveBeenCalled()
    })

    it('when embed mode is enabled, access_token is present but auth is not delegated, should call updateContext', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getAccessToken: vi.fn().mockResolvedValue('access-token'),
          getUser: vi.fn().mockResolvedValue(mock<User>({ expires_in: 3600 })),
          updateContext: mockUpdateContext
        })
      })

      initAuthService({ authService })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(mockUpdateContext).toHaveBeenCalledWith('access-token', true)
    })

    it('when embed mode is enabled, access_token is present and auth is delegated, should not call updateContext', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getAccessToken: vi.fn().mockResolvedValue('access-token'),
          updateContext: mockUpdateContext
        })
      })

      const configStore = useConfigStore()
      configStore.options = { embed: { enabled: true, delegateAuthentication: true } }
      initAuthService({ authService, configStore })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(mockUpdateContext).not.toHaveBeenCalled()
    })

    it('when embed mode is disabled, access_token is present and auth is delegated, should call updateContext', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getAccessToken: vi.fn().mockResolvedValue('access-token'),
          getUser: vi.fn().mockResolvedValue(mock<User>({ expires_in: 3600 })),
          updateContext: mockUpdateContext
        })
      })

      initAuthService({ authService })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(mockUpdateContext).toHaveBeenCalledWith('access-token', true)
    })
  })

  describe('handleDelegatedTokenUpdate', () => {
    const buildMessageEvent = (origin: string) =>
      mock<MessageEvent>({
        origin,
        data: { name: 'owncloud-embed:update-token', data: { access_token: 'attacker-token' } }
      })

    it('when delegateAuthenticationOrigin is not configured, should reject the message regardless of its origin', () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({ updateContext: mockUpdateContext })
      })

      const configStore = useConfigStore()
      configStore.options = { embed: { enabled: true, delegateAuthentication: true } }
      initAuthService({ authService, configStore })
      ;(authService as any).handleDelegatedTokenUpdate(
        buildMessageEvent('https://attacker.example')
      )

      expect(mockUpdateContext).not.toHaveBeenCalled()
    })

    it('when the message origin does not match the configured delegateAuthenticationOrigin, should reject the message', () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({ updateContext: mockUpdateContext })
      })

      const configStore = useConfigStore()
      configStore.options = {
        embed: {
          enabled: true,
          delegateAuthentication: true,
          delegateAuthenticationOrigin: 'https://trusted.example'
        }
      }
      initAuthService({ authService, configStore })
      ;(authService as any).handleDelegatedTokenUpdate(
        buildMessageEvent('https://attacker.example')
      )

      expect(mockUpdateContext).not.toHaveBeenCalled()
    })

    it('when the message origin matches the configured delegateAuthenticationOrigin, should update the context', () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({ updateContext: mockUpdateContext })
      })

      const configStore = useConfigStore()
      configStore.options = {
        embed: {
          enabled: true,
          delegateAuthentication: true,
          delegateAuthenticationOrigin: 'https://trusted.example'
        }
      }
      initAuthService({ authService, configStore })
      ;(authService as any).handleDelegatedTokenUpdate(buildMessageEvent('https://trusted.example'))

      expect(mockUpdateContext).toHaveBeenCalled()
    })

    describe('when dispatched through the window message listener', () => {
      let authService: AuthService

      const signInDelegated = async () => {
        authService = new AuthService()

        Object.defineProperty(authService, 'userManager', {
          value: mock<UserManager>({
            getUser: vi.fn().mockResolvedValue(null),
            getAndClearPostLoginRedirectUrl: vi.fn().mockReturnValue('/'),
            updateContext: mockUpdateContext
          })
        })

        const configStore = useConfigStore()
        configStore.options = {
          embed: {
            enabled: true,
            delegateAuthentication: true,
            delegateAuthenticationOrigin: 'https://host.example.org'
          }
        }
        initAuthService({ authService, configStore, router: createRouter() })

        await authService.signInCallback('initial-token')
        mockUpdateContext.mockClear()
      }

      afterEach(() => {
        window.removeEventListener(
          'message',
          (authService as any).handleDelegatedTokenUpdate as EventListener
        )
      })

      it('ignores "owncloud-embed:update-token" messages from unexpected origins', async () => {
        await signInDelegated()

        window.dispatchEvent(
          new MessageEvent('message', {
            data: { name: 'owncloud-embed:update-token', data: { access_token: 'renewed-token' } },
            origin: 'https://attacker.example.org'
          })
        )

        expect(mockUpdateContext).not.toHaveBeenCalled()
      })

      it('ignores "owncloud-embed:update-token" messages without an access token', async () => {
        await signInDelegated()

        window.dispatchEvent(
          new MessageEvent('message', {
            data: { name: 'owncloud-embed:update-token', data: {} },
            origin: 'https://host.example.org'
          })
        )

        expect(mockUpdateContext).not.toHaveBeenCalled()
      })

      it('updates the user context with the access token from "owncloud-embed:update-token" messages', async () => {
        await signInDelegated()

        window.dispatchEvent(
          new MessageEvent('message', {
            data: { name: 'owncloud-embed:update-token', data: { access_token: 'renewed-token' } },
            origin: 'https://host.example.org'
          })
        )

        expect(mockUpdateContext).toHaveBeenCalledWith('renewed-token', false)
      })
    })
  })

  describe('acr', () => {
    const mockSignInRedirect = vi.fn()

    it('when user is not authenticated, should redirect to login page', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getUser: vi.fn().mockResolvedValue(null),
          signinRedirect: mockSignInRedirect
        })
      })

      await authService.requireAcr('advanced', '/')
      expect(mockSignInRedirect).toHaveBeenCalledWith({ acr_values: 'advanced' })
    })

    it('when user is authenticated and acr is not the one required, should redirect to login page', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getUser: vi
            .fn()
            .mockResolvedValue(mock<User>({ profile: { acr: 'regular' }, expired: false })),
          signinRedirect: mockSignInRedirect
        })
      })

      await authService.requireAcr('advanced', '/')
      expect(mockSignInRedirect).toHaveBeenCalledWith({ acr_values: 'advanced' })
    })

    it('when user is authenticated and acr is the one required but access token is expired, should redirect to login page', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getUser: vi
            .fn()
            .mockResolvedValue(mock<User>({ profile: { acr: 'advanced' }, expired: true })),
          signinRedirect: mockSignInRedirect
        })
      })

      await authService.requireAcr('advanced', '/')
      expect(mockSignInRedirect).toHaveBeenCalledWith({ acr_values: 'advanced' })
    })

    it('when user is authenticated and acr is the one required, should not redirect to login page', async () => {
      const authService = new AuthService()

      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getUser: vi
            .fn()
            .mockResolvedValue(mock<User>({ profile: { acr: 'advanced' }, expired: false })),
          signinRedirect: mockSignInRedirect
        })
      })

      await authService.requireAcr('advanced', '/')
      expect(mockSignInRedirect).not.toHaveBeenCalled()
    })
  })
})
