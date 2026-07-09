import { ConfigStore, useAuthStore, useConfigStore } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { Router } from 'vue-router'
import { AuthService, isNetworkError } from '../../../../src/services/auth/authService'
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

  describe('isNetworkError', () => {
    it.each([
      ['axios ERR_NETWORK (offline / connection refused)', { code: 'ERR_NETWORK' }, true],
      ['axios ECONNABORTED (timeout)', { code: 'ECONNABORTED' }, true],
      ['fetch TypeError - Safari "Load failed"', new TypeError('Load failed'), true],
      ['fetch TypeError - Chromium "Failed to fetch"', new TypeError('Failed to fetch'), true],
      [
        'fetch TypeError - Firefox "NetworkError ..."',
        new TypeError('NetworkError when attempting to fetch resource.'),
        true
      ],
      ['an answered request carrying an HTTP response (401)', { response: { status: 401 } }, false],
      [
        'an ERR_NETWORK code that nonetheless carries a response',
        { code: 'ERR_NETWORK', response: { status: 500 } },
        false
      ],
      ['an unrelated TypeError (real bug)', new TypeError('foo is not a function'), false],
      ['an error with an unknown code', { code: 'ERR_BAD_REQUEST' }, false],
      ['a null error', null, false],
      ['a plain string error', 'boom', false]
    ])('returns %s => %s', (_desc, error, expected) => {
      expect(isNetworkError(error)).toBe(expected)
    })
  })

  describe('initializeContext auth-error classification', () => {
    const buildAuthService = (rejection: unknown) => {
      const authService = new AuthService()
      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          getAccessToken: vi.fn().mockResolvedValue('access-token'),
          getUser: vi.fn().mockResolvedValue(mock<User>({ expires_in: 3600 })),
          updateContext: vi.fn().mockRejectedValue(rejection)
        })
      })
      const router = createRouter()
      initAuthService({ authService, router })
      const handleAuthErrorSpy = vi
        .spyOn(authService, 'handleAuthError')
        .mockResolvedValue(undefined)
      return { authService, handleAuthErrorSpy }
    }

    it('preserves the session (no handleAuthError) when updateContext fails with a network error on reload', async () => {
      const { authService, handleAuthErrorSpy } = buildAuthService({ code: 'ERR_NETWORK' })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(handleAuthErrorSpy).not.toHaveBeenCalled()
    })

    it('triggers handleAuthError when updateContext fails with a 401 on reload', async () => {
      const { authService, handleAuthErrorSpy } = buildAuthService({ response: { status: 401 } })

      await authService.initializeContext(mock<RouteLocation>({}))

      expect(handleAuthErrorSpy).toHaveBeenCalled()
    })

    it('addUserLoaded handler preserves the session on a network error but logs out on an auth rejection', async () => {
      const authService = new AuthService()
      let userLoadedCb: (user: User) => Promise<void>
      const updateContext = vi.fn()
      Object.defineProperty(authService, 'userManager', {
        value: mock<UserManager>({
          // skip the reload branch, we drive the addUserLoaded handler directly
          getAccessToken: vi.fn().mockResolvedValue(null),
          // force event-handler registration (mock<> would otherwise make this truthy)
          areEventHandlersRegistered: false,
          updateContext,
          events: {
            addAccessTokenExpired: vi.fn(),
            addAccessTokenExpiring: vi.fn(),
            addUserLoaded: vi.fn().mockImplementation((cb) => {
              userLoadedCb = cb
            }),
            addUserUnloaded: vi.fn(),
            addSilentRenewError: vi.fn()
          } as unknown as UserManager['events']
        })
      })
      const router = createRouter()
      initAuthService({ authService, router })
      const handleAuthErrorSpy = vi
        .spyOn(authService, 'handleAuthError')
        .mockResolvedValue(undefined)

      await authService.initializeContext(mock<RouteLocation>({}))

      // network error -> session preserved
      updateContext.mockRejectedValueOnce({ code: 'ERR_NETWORK' })
      await userLoadedCb(mock<User>({ access_token: 'tok', expires_in: 3600 }))
      expect(handleAuthErrorSpy).not.toHaveBeenCalled()

      // auth rejection -> existing logout flow
      updateContext.mockRejectedValueOnce({ response: { status: 401 } })
      await userLoadedCb(mock<User>({ access_token: 'tok', expires_in: 3600 }))
      expect(handleAuthErrorSpy).toHaveBeenCalledTimes(1)
    })
  })
})
