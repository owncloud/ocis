import { Page } from '@playwright/test'
import util from 'util'
import { locatorUtils } from '../../../utils'
import { objects } from '../../../index'

const sidebarPanel = '#app-sidebar'
const contextMenuButton =
  '//span[@data-test-resource-name="%s"]/ancestor::tr[contains(@class, "oc-tbody-tr")]//button[contains(@class, "resource-table-btn-action-dropdown")]'
const contextMenuContainer = '#oc-files-context-menu'
const folderModalIframe = '#iframe-folder-view'
const actionMenuForCurrentFolderSelector = '#oc-breadcrumb-contextmenu-trigger'
const closeSidebarRootPanelBtn = `${sidebarPanel} .is-active-root-panel .header__close:visible`
const closeSidebarSubPanelBtn = `${sidebarPanel} .is-active-sub-panel .header__close:visible`

const openForResource = async ({
  page,
  resource,
  resourceType
}: {
  page: Page
  resource: string
  resourceType: string
}): Promise<void> => {
  if (resourceType === 'passwordProtectedFolder') {
    await page.frameLocator(folderModalIframe).locator(actionMenuForCurrentFolderSelector).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      [folderModalIframe],
      'account page'
    )
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'account page'
    )
    await page
      .frameLocator(folderModalIframe)
      .locator('.oc-files-actions-show-details-trigger')
      .click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      [folderModalIframe],
      'account page'
    )
  } else {
    await page.locator(util.format(contextMenuButton, resource)).waitFor()
    await page.locator(util.format(contextMenuButton, resource)).click()
    await page.locator(contextMenuContainer).waitFor()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'account page'
    )
    await page
      .locator(contextMenuContainer)
      .locator('.oc-files-actions-show-details-trigger')
      .click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'account page'
    )
  }
}

export const openPanelForResource = async ({
  page,
  resource,
  panel
}: {
  page: Page
  resource: string
  panel: string
}): Promise<void> => {
  await page.locator(util.format(contextMenuButton, resource)).waitFor()
  await page.locator(util.format(contextMenuButton, resource)).click()
  await page.locator(contextMenuContainer).waitFor()
  await page
    .locator(contextMenuContainer)
    .locator(`.oc-files-actions-show-${panel}-trigger`)
    .click()
}

const openGlobal = async ({ page }: { page: Page }): Promise<void> => {
  await page.locator('#files-toggle-sidebar').click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [sidebarPanel],
    'sidebar panel'
  )
}

export const open = async ({
  page,
  resource,
  resourceType
}: {
  page: Page
  resource?: string
  resourceType?: string
}): Promise<void> => {
  if (resourceType === 'passwordProtectedFolder') {
    if (await page.frameLocator(folderModalIframe).locator(sidebarPanel).count()) {
      await closePasswordProtectedFolder({ page })
    }
  } else {
    if (await page.locator(sidebarPanel).count()) {
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        [sidebarPanel],
        'sidebar panel'
      )
      await Promise.all([
        page.locator(sidebarPanel).waitFor({ state: 'detached' }),
        close({ page })
      ])
    }
  }

  resource ? await openForResource({ page, resource, resourceType }) : await openGlobal({ page })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [sidebarPanel],
    'sidebar panel opening'
  )
}

export const close = async ({ page }: { page: Page }): Promise<void> => {
  // await sidebar transitions
  await new Promise((resolve) => setTimeout(resolve, 250))
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [sidebarPanel],
    'sidebar panel'
  )
  const isSubPanelActive = await page.locator(closeSidebarSubPanelBtn).isVisible()
  if (isSubPanelActive) {
    await page.locator(closeSidebarSubPanelBtn).click()
  } else {
    await page.locator(closeSidebarRootPanelBtn).click()
  }
}

export const openPanel = async ({
  page,
  name,
  resourceType
}: {
  page: Page
  name: string
  resourceType?: string
}): Promise<void> => {
  const currentPanel = page.locator('.sidebar-panel.is-active')
  const backButton = currentPanel.locator('.header__back')

  if (await backButton.count()) {
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['sidebarPanel'],
      'sidebar panel'
    )
    await Promise.all([
      locatorUtils.waitForEvent(currentPanel, 'transitionend'),
      backButton.click()
    ])
  }
  if (resourceType === 'passwordProtectedFolder') {
    const panelSelector = page
      .frameLocator(folderModalIframe)
      .locator(`#sidebar-panel-${name}-select`)
    const nextPanel = page.frameLocator(folderModalIframe).locator(`#sidebar-panel-${name}`)
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      [folderModalIframe],
      'sidebar panel'
    )
    await Promise.all([nextPanel.waitFor(), panelSelector.click()])
  } else {
    const panelSelector = page.locator(`#sidebar-panel-${name}-select`)
    const nextPanel = page.locator(`#sidebar-panel-${name}`)
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      [sidebarPanel],
      'sidebar panel'
    )
    await Promise.all([nextPanel.waitFor(), panelSelector.click()])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'activity sidebar panel with no activities'
    )
  }
}

export const closePasswordProtectedFolder = async ({ page }: { page: Page }): Promise<void> => {
  const isSubPanelActive = await page
    .frameLocator(folderModalIframe)
    .locator(closeSidebarSubPanelBtn)
    .isVisible()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [sidebarPanel],
    'sidebar panel'
  )
  if (isSubPanelActive) {
    await page.frameLocator(folderModalIframe).locator(closeSidebarSubPanelBtn).click()
  } else {
    await page.frameLocator(folderModalIframe).locator(closeSidebarRootPanelBtn).click()
  }
}
