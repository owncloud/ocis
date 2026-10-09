import { UserManager } from './userManager'
import { PublicLinkManager } from './publicLinkManager'
import {
  AuthStore,
  ClientService,
  UserStore,
  CapabilityStore,
  ConfigStore,
  useTokenTimerWorker,
  useMessages,
  AuthServiceInterface
} from '@ownclouders/web-pkg'
import { RouteLocation, RouteLocationRaw, Router } from 'vue-router'
import {
  base,
  extractPublicLinkToken,
  isAnonymousContext,
  isIdpContextRequired,
  isPublicLinkContextRequired,
  isUserContextRequired
} from '../../router'
import { unref } from 'vue'
import { Ability } from '@ownclouders/web-client'
import { Language } from 'vue3-gettext'
import { PublicLinkType } from '@ownclouders/web-client'
import { WebWorkersStore } from '@ownclouders/web-pkg'
import { isSilentRedirectRoute } from '../../helpers/silentRedirect'
import { vaultStepUpFailedKey } from '../../helpers/vaultStepUp'

/**
 * Custom OIDC state sent with an MFA step-up request. oidc-client-ts stores it with the
 * request and returns it from the sign-in callback, so a step-up is only treated as answered
 * when the IdP actually redirected back for it.
 */
interface MfaStepUpState {
  mfaStepUpTarget: string
}

const isMfaStepUpState = (state: unknown): state is MfaStepUpState =>
  typeof (state as MfaStepUpState)?.mfaStepUpTarget === 'string'

export class AuthService implements AuthServiceInterface {
  private clientService: ClientService
  private configStore: ConfigStore
  private router: Router
  private userManager: UserManager
  private publicLinkManager: PublicLinkManager
  private ability: Ability
  private language: Language
  private userStore: UserStore
  private authStore: AuthStore
  private capabilityStore: CapabilityStore
  private webWorkersStore: WebWorkersStore

  private tokenTimerWorker: ReturnType<typeof useTokenTimerWorker>
  private tokenTimerInitialized = false

  // number of seconds before an access token is to expire to raise the accessTokenExpiring event
  private accessTokenExpiryThreshold = 10

  // Target of an MFA step-up the IdP just answered (set in `signInCallback`) ...
  private pendingStepUpReturnTarget: string | null = null
  // ... and moved here for exactly the navigation that follows the callback.
  private stepUpReturnTarget: string | null = null

  public hasAuthErrorOccurred: boolean

  public initialize(
    configStore: ConfigStore,
    clientService: ClientService,
    router: Router,
    ability: Ability,
    language: Language,
    userStore: UserStore,
    authStore: AuthStore,
    capabilityStore: CapabilityStore,
    webWorkersStore: WebWorkersStore
  ): void {
    this.configStore = configStore
    this.clientService = clientService
    this.router = router
    this.hasAuthErrorOccurred = false
    this.ability = ability
    this.language = language
    this.userStore = userStore
    this.authStore = authStore
    this.capabilityStore = capabilityStore
    this.webWorkersStore = webWorkersStore
  }

  /**
   * Initialize publicLinkContext and userContext (whichever is available, respectively).
   *
   * FIXME: at the moment the order "publicLink first, user second" is important, because we trigger the `ready` hook of all applications
   * as soon as any context is ready. This works well for user context pages, because they can't have a public link context at the same time.
   * Public links on the other hand could have a logged in user as well, thus we need to make sure that the public link context is loaded first.
   * For the moment this is fine. In the future we might want to wait triggering the `ready` hook of applications until all available contexts
   * are loaded.
   *
   * @param to {Route}
   */
  public async initializeContext(to: RouteLocation): Promise<RouteLocationRaw | false | void> {
    // A step-up answer is only valid for the navigation right after the callback.
    this.stepUpReturnTarget = this.pendingStepUpReturnTarget
    this.pendingStepUpReturnTarget = null

    this.showPendingVaultStepUpFailure()

    if (!this.publicLinkManager) {
      this.publicLinkManager = new PublicLinkManager({
        clientService: this.clientService,
        authStore: this.authStore,
        capabilityStore: this.capabilityStore
      })
    }

    if (isPublicLinkContextRequired(this.router, to)) {
      const publicLinkToken = extractPublicLinkToken(to)
      if (publicLinkToken) {
        await this.publicLinkManager.updateContext(publicLinkToken)
      }
    } else if (to.name !== 'resolvePublicLink') {
      // no need to clear public context if we're routing the to public link resolving page
      this.publicLinkManager.clearContext()
    }

    if (!this.userManager) {
      this.userManager = new UserManager({
        clientService: this.clientService,
        configStore: this.configStore,
        ability: this.ability,
        language: this.language,
        userStore: this.userStore,
        authStore: this.authStore,
        capabilityStore: this.capabilityStore,
        webWorkersStore: this.webWorkersStore,
        accessTokenExpiryThreshold: this.accessTokenExpiryThreshold
      })

      // don't load worker in the silent redirect iframe
      const isSilentRedirect = isSilentRedirectRoute()
      if (!this.tokenTimerWorker && !isSilentRedirect) {
        const { options } = this.configStore

        if (!options.embed?.enabled || !options.embed?.delegateAuthentication) {
          this.tokenTimerWorker = useTokenTimerWorker({ authService: this })
          this.tokenTimerWorker.startWorker()
        }
      }
    }

    if (to.params.scope === 'vault') {
      // Permissions known: send an unentitled user to accessDenied instead of the IdP.
      if (this.authStore.userContextReady && !this.ability.can('read-all', 'Vault')) {
        return { name: 'accessDenied' }
      }

      // Capabilities only — the full context loads vault spaces, which 401 without MFA and crash.
      if (!this.capabilityStore.isInitialized) {
        this.capabilityStore.setCapabilities(await this.clientService.ocs.getCapabilities())
      }

      const requiredAcr = this.capabilityStore.authMfaRequiredLevelname
      const user = await this.userManager.getUser()

      // Address-bar navigation: permissions unloaded — load them and deny before MFA.
      if (user && !user.expired && !this.authStore.userContextReady) {
        try {
          await this.userManager.loadUserAbilities()
        } catch (e) {
          // Can't confirm access: deny, don't hand off to MFA or leave a blank page.
          console.error('failed to load vault abilities on cold load, denying access:', e)
          return { name: 'accessDenied' }
        }
        if (!this.ability.can('read-all', 'Vault')) {
          return { name: 'accessDenied' }
        }
      }

      if (!user || user.expired || user.profile.acr !== requiredAcr) {
        // `acr_values` is a voluntary claim: an IdP that can't reach the required level
        // (no second factor available/enrolled) returns a lower `acr` instead of an error.
        // Redirecting again would loop forever, so give up once the IdP answered the step-up.
        if (user && !user.expired && this.isStepUpReturn(to.fullPath)) {
          console.warn(
            `[authService:initializeContext] - MFA step-up returned acr "${user.profile.acr}", required "${requiredAcr}". Not retrying.`
          )
          this.leaveVaultAfterFailedStepUp()
          // cancel the vault navigation, the page reloads outside the vault
          return false
        }

        this.userManager.setPostLoginRedirectUrl(to.fullPath)
        await this.signinRedirectForStepUp(requiredAcr, to.fullPath)
        // redirecting to the IdP, don't establish the user context below
        return
      }
    }

    if (isPublicLinkContextRequired(this.router, to)) {
      const user = await this.userManager.getUser()

      if (user?.expired) {
        try {
          await this.userManager.signinSilent()
        } catch {
          await this.userManager.removeUser()
        }
      }
    }

    if (!isAnonymousContext(this.router, to)) {
      const fetchUserData = !isIdpContextRequired(this.router, to)

      if (!this.userManager.areEventHandlersRegistered) {
        this.userManager.events.addAccessTokenExpired((...args): void => {
          const handleExpirationError = () => {
            console.error('AccessToken Expired：', ...args)
            this.handleAuthError(unref(this.router.currentRoute), { forceLogout: true })
          }

          // retry silent signin once, force logout if it fails
          this.userManager.signinSilent().catch(handleExpirationError)
        })

        this.userManager.events.addAccessTokenExpiring((...args) => {
          console.debug('AccessToken Expiring：', ...args)
        })

        this.userManager.events.addUserLoaded(async (user) => {
          this.tokenTimerWorker?.setTokenTimer({
            expiry: user.expires_in,
            expiryThreshold: this.accessTokenExpiryThreshold
          })

          console.debug(
            `New User Loaded. access_token： ${user.access_token}, refresh_token: ${user.refresh_token}`
          )
          try {
            await this.userManager.updateContext(user.access_token, fetchUserData)
          } catch (e) {
            console.error(e)
            await this.handleAuthError(unref(this.router.currentRoute))
          }
        })

        this.userManager.events.addUserUnloaded(() => {
          console.log('user unloaded…')
          this.tokenTimerWorker?.resetTokenTimer()
          this.resetStateAfterUserLogout()

          if (this.userManager.unloadReason === 'authError') {
            this.hasAuthErrorOccurred = true
            return this.router.push({
              name: 'accessDenied',
              query: { redirectUrl: unref(this.router.currentRoute)?.fullPath }
            })
          }

          // handle redirect after logout
          if (this.configStore.isOAuth2) {
            const oAuth2 = this.configStore.oAuth2
            if (oAuth2.logoutUrl) {
              return (window.location = oAuth2.logoutUrl as any)
            }
          }
        })
        this.userManager.events.addSilentRenewError(async (error) => {
          console.error('Silent Renew Error：', error)
          await this.handleAuthError(unref(this.router.currentRoute))
        })

        this.userManager.areEventHandlersRegistered = true
      }

      // This is to prevent issues in embed mode when the expired token is still saved but already expired
      // If the following code gets executed, it would toggle errorOccurred var which would then lead to redirect to the access denied screen
      if (
        this.configStore.options.embed?.enabled &&
        this.configStore.options.embed.delegateAuthentication
      ) {
        return
      }

      // relevant for page reload: token is already in userStore
      // no userLoaded event and no signInCallback gets triggered
      const accessToken = await this.userManager.getAccessToken()
      if (accessToken) {
        console.debug('[authService:initializeContext] - updating context with saved access_token')

        try {
          await this.userManager.updateContext(accessToken, fetchUserData)

          if (!this.tokenTimerInitialized) {
            const user = await this.userManager.getUser()
            this.tokenTimerWorker?.setTokenTimer({
              expiry: user.expires_in,
              expiryThreshold: this.accessTokenExpiryThreshold
            })

            this.tokenTimerInitialized = true
          }
        } catch (e) {
          console.error(e)
          await this.handleAuthError(unref(this.router.currentRoute))
        }
      }
    }
  }

  public loginUser(redirectUrl?: string) {
    this.userManager.setPostLoginRedirectUrl(redirectUrl)
    return this.userManager.signinRedirect()
  }

  public signinSilent() {
    return this.userManager.signinSilent()
  }

  /**
   * Sign in callback gets called from the IDP after initial login.
   */
  public async signInCallback(accessToken?: string) {
    try {
      if (
        this.configStore.options.embed.enabled &&
        this.configStore.options.embed.delegateAuthentication &&
        accessToken
      ) {
        console.debug('[authService:signInCallback] - setting access_token and fetching user')
        await this.userManager.updateContext(accessToken, true)

        // Setup a listener to handle token refresh
        console.debug('[authService:signInCallback] - adding listener to update-token event')
        window.addEventListener('message', this.handleDelegatedTokenUpdate)
      } else {
        const callbackUser = await this.userManager.signinRedirectCallback(
          this.buildSignInCallbackUrl()
        )
        if (isMfaStepUpState(callbackUser?.state)) {
          this.pendingStepUpReturnTarget = callbackUser.state.mfaStepUpTarget
        }
      }

      const redirectRoute = this.router.resolve(this.userManager.getAndClearPostLoginRedirectUrl())
      return this.router.replace({
        path: redirectRoute.path,
        ...(redirectRoute.query && { query: redirectRoute.query })
      })
    } catch (e) {
      console.warn('error during authentication:', e)
      return this.handleAuthError(unref(this.router.currentRoute))
    }
  }

  /**
   * Sign in silent callback gets called with OIDC during access token renewal when no `refresh_token`
   * is present (`refresh_token` exists when `offline_access` is present in scopes).
   *
   * The oidc-client lib emits a userLoaded event internally, which already handles the token update
   * in web.
   */
  public async signInSilentCallback() {
    await this.userManager.signinSilentCallback(this.buildSignInCallbackUrl())
  }

  /**
   * craft a url that the parser in oidc-client-ts can handle…
   */
  private buildSignInCallbackUrl() {
    const currentQuery = unref(this.router.currentRoute).query
    return '/?' + new URLSearchParams(currentQuery as Record<string, string>).toString()
  }

  public async handleAuthError(
    route: RouteLocation,
    { forceLogout = false }: { forceLogout?: boolean } = {}
  ) {
    if (isPublicLinkContextRequired(this.router, route)) {
      const token = extractPublicLinkToken(route)
      this.publicLinkManager.clear(token)
      return this.router.push({
        name: 'resolvePublicLink',
        params: { token },
        query: { redirectUrl: route.fullPath }
      })
    }
    if (isUserContextRequired(this.router, route) || isIdpContextRequired(this.router, route)) {
      if (forceLogout) {
        this.tokenTimerWorker?.resetTokenTimer()
        await this.logoutUser()
        return
      }

      const user = await this.userManager.getUser()
      if (user?.expires_in !== undefined && user.expires_in < 0) {
        // token expired, simply return and let the regular auth flow do its thing
        return
      }

      await this.userManager.removeUser('authError')
      this.tokenTimerWorker?.resetTokenTimer()
      return
    }
    // authGuard is taking care of redirecting the user to the
    // accessDenied page if hasAuthErrorOccurred is set to true
    // we can't push the route ourselves, see authGuard for details.
    this.hasAuthErrorOccurred = true
  }

  public async resolvePublicLink(
    token: string,
    passwordRequired: boolean,
    password: string,
    type: PublicLinkType
  ) {
    this.publicLinkManager.setPasswordRequired(token, passwordRequired)
    this.publicLinkManager.setPassword(token, password)
    this.publicLinkManager.setResolved(token, true)
    this.publicLinkManager.setType(token, type)

    await this.publicLinkManager.updateContext(token)
  }

  public async logoutUser() {
    const endSessionEndpoint = await this.userManager.metadataService?.getEndSessionEndpoint()
    if (!endSessionEndpoint) {
      await this.userManager.removeUser()
      return this.router.push({ name: 'logout' })
    }

    const u = await this.userManager.getUser()
    if (u && u.id_token) {
      return this.userManager.signoutRedirect({ id_token_hint: u.id_token })
    }

    return await this.userManager.removeUser()
  }

  private resetStateAfterUserLogout() {
    // TODO: create UserUnloadTask interface and allow registering unload-tasks in the authService
    this.userStore.reset()
    this.authStore.clearUserContext()
  }

  public async getRefreshToken() {
    const user = await this.userManager.getUser()
    return user?.refresh_token
  }

  private handleDelegatedTokenUpdate = (event: MessageEvent) => {
    if (event.origin !== this.configStore.options.embed?.delegateAuthenticationOrigin) {
      return
    }

    if (event.data?.name !== 'owncloud-embed:update-token') {
      return
    }

    const accessToken = event.data.data?.access_token
    if (!accessToken) {
      return
    }

    console.debug('[authService:handleDelegatedTokenUpdate] - going to update the access_token')
    return this.userManager.updateContext(accessToken, false)
  }

  /**
   * Ensures the current user has authenticated with the given `acr` (e.g. MFA), redirecting
   * to the IdP for a step-up if needed.
   *
   * @param acrValue - The ACR value to require.
   * @param redirectUrl - The URL to redirect to after login.
   *
   * @returns false if the IdP just answered a step-up for `redirectUrl` without the required
   * `acr`. The caller must not grant access and should navigate elsewhere.
   */
  public async requireAcr(acrValue: string, redirectUrl: string): Promise<boolean> {
    const user = await this.userManager.getUser()
    const isAuthenticated = user && !user.expired

    if (isAuthenticated && user.profile.acr === acrValue) {
      return true
    }

    // `acr_values` is a voluntary claim, see `initializeContext`: don't redirect again
    if (isAuthenticated && this.isStepUpReturn(redirectUrl)) {
      console.warn(
        `[authService:requireAcr] - MFA step-up returned acr "${user.profile.acr}", required "${acrValue}". Not retrying.`
      )
      this.showStepUpFailedMessage()
      return false
    }

    this.userManager.setPostLoginRedirectUrl(redirectUrl)
    await this.signinRedirectForStepUp(acrValue, redirectUrl)
    return true
  }

  private signinRedirectForStepUp(acrValue: string, target: string) {
    const state: MfaStepUpState = { mfaStepUpTarget: target }
    return this.userManager.signinRedirect({ acr_values: acrValue, state })
  }

  /**
   * Whether the current navigation is the one right after the IdP answered a step-up for `target`.
   */
  private isStepUpReturn(target: string): boolean {
    return this.stepUpReturnTarget !== null && this.stepUpReturnTarget === target
  }

  /**
   * Vault mode is fixed for the lifetime of the page (see `useVault`), so leaving it
   * requires a full page load. An in-app navigation would keep the vault clients and
   * load vault spaces with a non-MFA token.
   */
  private leaveVaultAfterFailedStepUp() {
    try {
      sessionStorage.setItem(vaultStepUpFailedKey, 'true')
    } catch {
      // message can't be shown after the reload, leaving the vault still has to happen
    }
    this.navigateOutsideVault(base?.href || `${window.location.origin}/`)
  }

  private navigateOutsideVault(url: string) {
    window.location.assign(url)
  }

  private showPendingVaultStepUpFailure() {
    try {
      if (!sessionStorage.getItem(vaultStepUpFailedKey)) {
        return
      }
      sessionStorage.removeItem(vaultStepUpFailedKey)
    } catch {
      return
    }
    this.showVaultStepUpFailedMessage()
  }

  private showVaultStepUpFailedMessage() {
    const { $pgettext } = this.language
    useMessages().showErrorMessage({
      title: $pgettext(
        'Error message title shown when the multi-factor authentication step-up required to open the vault failed',
        'Multi-factor authentication required'
      ),
      desc: $pgettext(
        'Error message shown when the multi-factor authentication step-up required to open the vault failed, e.g. because the user has no second factor at hand',
        'The vault requires multi-factor authentication, which could not be completed. Please set up a second factor or contact your administrator.'
      )
    })
  }

  private showStepUpFailedMessage() {
    const { $pgettext } = this.language
    useMessages().showErrorMessage({
      title: $pgettext(
        'Error message title shown when the multi-factor authentication step-up required to open a page (e.g. admin settings) failed',
        'Multi-factor authentication required'
      ),
      desc: $pgettext(
        'Error message shown when the multi-factor authentication step-up required to open a page (e.g. admin settings) failed, e.g. because the user has no second factor at hand',
        'This page requires multi-factor authentication, which could not be completed. Please set up a second factor or contact your administrator.'
      )
    })
  }
}

export const authService = new AuthService()
