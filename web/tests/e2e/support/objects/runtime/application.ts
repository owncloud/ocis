import { Page } from '@playwright/test'
import util from 'util'
import { config } from '../../../config'
import { objects } from '../../index'

const appSwitcherButton = '#_appSwitcherButton'
const appSelector = `//ul[contains(@class, "applications-list")]//*[@data-test-id="%s"]`
const notificationsBell = `#oc-notifications-bell`
const notificationsDrop = `#oc-notifications-drop`
const notificationsLoading = `#oc-notifications-drop .oc-notifications-loading`
const markNotificationsAsReadButton = `#oc-notifications-drop .oc-notifications-mark-all`
const notificationItemsMessages = `#oc-notifications-drop .oc-notifications-item .oc-notifications-message`
const closeSidebarRootPanelBtn = `#app-sidebar .is-active-root-panel .header__close`
const closeSidebarSubPanelBtn = `#app-sidebar .is-active-sub-panel .header__close`

export class Application {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async reloadPage(): Promise<void> {
    await this.#page.reload()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      'body after page reload'
    )
  }

  async open({ name }: { name: string }): Promise<void> {
    await this.#page.waitForTimeout(1000)
    await this.#page.locator(appSwitcherButton).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['appSwitcherDropdown'],
      'app switcher dropdown'
    )
    await this.#page.locator(util.format(appSelector, `app.${name}.menuItem`)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      `${name} app page`
    )
  }

  async getNotificationMessages(): Promise<string[]> {
    // reload will fetch notifications immediately
    // wait for the notifications to load
    await Promise.all([
      this.#page.waitForResponse(
        (resp) =>
          resp.url().endsWith('notifications') &&
          resp.status() === 200 &&
          resp.request().method() === 'GET'
      ),
      this.#page.reload()
    ])

    const dropIsOpen = await this.#page.locator(notificationsDrop).isVisible()
    if (!dropIsOpen) {
      await this.#page.locator(notificationsBell).click()
    }
    await this.#page.locator(notificationsLoading).waitFor({ state: 'detached' })
    const result = this.#page.locator(notificationItemsMessages)
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      [notificationItemsMessages],
      'notifications'
    )
    const messages = []
    const count = await result.count()
    for (let i = 0; i < count; i++) {
      messages.push(await result.nth(i).innerText())
    }
    return messages
  }

  async markNotificationsAsRead(): Promise<void> {
    const dropIsOpen = await this.#page.locator(notificationsDrop).isVisible()
    if (!dropIsOpen) {
      await this.#page.locator(notificationsBell).click()
    }
    await this.#page.locator(notificationsLoading).waitFor({ state: 'detached' })
    await this.#page.locator(markNotificationsAsReadButton).click()

    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      [notificationItemsMessages],
      'notifications'
    )
    await this.#page.locator(notificationsLoading).waitFor({ state: 'detached' })
  }

  async openUrl(url): Promise<void> {
    await this.#page.goto(url)
  }

  async closeSidebar(): Promise<void> {
    // await sidebar transitions
    await new Promise((resolve) => setTimeout(resolve, 250))
    const isSubPanelActive = await this.#page.locator(closeSidebarSubPanelBtn).isVisible()
    if (isSubPanelActive) {
      await this.#page.locator(closeSidebarSubPanelBtn).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        this.#page,
        ['body'],
        'body after closing sidebar sub panel'
      )
    } else {
      await this.#page.locator(closeSidebarRootPanelBtn).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        this.#page,
        ['body'],
        'body after closing sidebar root panel'
      )
    }
  }

  async waitForTokenRenewalViaRefreshToken(): Promise<void> {
    await Promise.all([
      this.#page.waitForResponse(
        (resp) =>
          resp.url().includes('/token') &&
          resp.status() === 200 &&
          resp.request().method() === 'POST' &&
          resp.request().postDataJSON().grant_type === 'refresh_token' &&
          resp.request().postDataJSON().hasOwnProperty('refresh_token') &&
          resp.request().postDataJSON().refresh_token &&
          resp.request().postDataJSON().hasOwnProperty('scope') &&
          resp.request().postDataJSON().scope.includes('offline_access'),
        { timeout: config.tokenTimeout * 1000 }
      )
    ])
  }

  async waitForTokenRenewalViaIframe(): Promise<void> {
    const waitForIframe = this.#page.evaluateHandle(() => {
      let iframe = null
      const observer = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
          if (mutation.type === 'childList') {
            mutation.addedNodes.forEach((node) => {
              if (node.nodeName === 'IFRAME') {
                iframe = node
              }
            })
          }
        })
      })
      observer.observe(document.body, { childList: true, subtree: true })
      return new Promise((resolve) => {
        const interval = setInterval(() => {
          if (iframe) {
            observer.disconnect()
            clearInterval(interval)
            resolve(iframe)
          }
        }, 1000)
      })
    })

    await Promise.all([
      this.#page.waitForResponse(
        (resp) =>
          resp.url().includes('/oidc-silent-redirect.html') &&
          resp.status() === 200 &&
          resp.request().method() === 'GET',
        { timeout: config.tokenTimeout * 1000 }
      ),
      this.#page.waitForResponse(
        (resp) =>
          resp.url().endsWith('/token') &&
          resp.status() === 200 &&
          resp.request().method() === 'POST' &&
          resp.request().postDataJSON().grant_type === 'authorization_code' &&
          resp.request().postDataJSON().hasOwnProperty('code') &&
          resp.request().postDataJSON().code,
        { timeout: config.tokenTimeout * 1000 }
      ),
      waitForIframe
    ])
  }
}
