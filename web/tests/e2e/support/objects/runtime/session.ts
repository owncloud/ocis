import { Page } from '@playwright/test'
import { User } from '../../types'
import { config } from '../../../config'
import { TokenEnvironmentFactory } from '../../environment'
import { objects } from '../../index'

export class Session {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  signIn(username: string, password: string): Promise<void> {
    if (config.keycloak) {
      return this.keycloakSignIn(username, password)
    }
    return this.idpSignIn(username, password)
  }

  async idpSignIn(username: string, password: string): Promise<void> {
    await this.#page.locator('button[type="submit"]').waitFor()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      'login page'
    )
    await this.#page.locator('//input[@type="text" or @placeholder="Username"]').fill(username)
    await this.#page.locator('//input[@type="password" or @placeholder="Password"]').fill(password)
    await this.#page.locator('button[type="submit"]').click()
  }

  async keycloakSignIn(username: string, password: string): Promise<void> {
    await this.#page.locator('#username').waitFor()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      'login page'
    )
    await this.#page.locator('#username').fill(username)
    await this.#page.locator('#password').fill(password)
    await this.#page.locator('#kc-login').click()
  }

  async keycloakOTPSignIn(otp: string): Promise<void> {
    await this.#page.locator('#otp').fill(otp)
    await this.#page.locator('#kc-login').click()
  }

  async login(user: User): Promise<void> {
    if (config.keycloak) {
      return this.keycloakLogin(user)
    }

    const [response] = await Promise.all([
      this.#page.waitForResponse(
        (resp) =>
          resp.url().endsWith('/token') &&
          resp.status() === 200 &&
          resp.request().method() === 'POST'
      ),
      this.signIn(user.id, user.password)
    ])

    if (config.predefinedUsers) {
      const tokenRes = await response.json()
      const tokenEnvironment = TokenEnvironmentFactory()
      tokenEnvironment.setToken({
        user: { ...user },
        token: {
          userId: user.id,
          accessToken: tokenRes.access_token,
          refreshToken: tokenRes.refresh_token
        }
      })
    }
  }

  private async keycloakLogin(user: User): Promise<void> {
    // Retry OIDC login: admin user is shared across all workers, so multiple
    // workers may authenticate it concurrently. Keycloak can return the login
    // form (200) instead of a redirect (302). Each retry starts a fresh OIDC
    // session, naturally staggering the requests.
    const MAX_RETRIES = 3
    for (let attempt = 0; attempt < MAX_RETRIES; attempt++) {
      try {
        if (attempt > 0) {
          await this.#page.goto(config.baseUrl)
        }
        const [response] = await Promise.all([
          this.#page.waitForResponse(
            (resp) =>
              resp.url().endsWith('/token') &&
              resp.status() === 200 &&
              resp.request().method() === 'POST',
            { timeout: config.tokenTimeout * 1000 }
          ),
          this.signIn(user.id, user.password)
        ])
        if (config.predefinedUsers) {
          const tokenRes = await response.json()
          const tokenEnvironment = TokenEnvironmentFactory()
          tokenEnvironment.setToken({
            user: { ...user },
            token: {
              userId: user.id,
              accessToken: tokenRes.access_token,
              refreshToken: tokenRes.refresh_token
            }
          })
        }
        return
      } catch (e) {
        if (attempt === MAX_RETRIES - 1) throw e
        await new Promise((r) => setTimeout(r, config.minTimeout * 1000))
      }
    }
  }

  async logout(): Promise<void> {
    await this.#page.locator('#_userMenuButton').click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['#account-info-container'],
      'files page'
    )
    await this.#page.locator('#oc-topbar-account-logout').click()
  }
}
