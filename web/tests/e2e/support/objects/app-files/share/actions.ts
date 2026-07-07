import { Page, Locator } from '@playwright/test'
import util from 'util'
import Collaborator, { ICollaborator, IAccessDetails } from './collaborator'
import { sidebar } from '../utils'
import { clickResource } from '../resource/actions'
import { User } from '../../../types'
import { locatorUtils } from '../../../utils'
import { objects } from '../../../index'
import { a11y } from '../../index'
import { fileAction } from '../../../../environment/constants'

const invitePanel = '//*[@id="oc-files-sharing-sidebar"]'
const quickShareButton =
  '//*[@data-test-resource-name="%s"]/ancestor::tr//button[contains(@class, "files-quick-action-show-shares")]'
const actionMenuDropdownButton =
  '//*[@data-test-resource-name="%s"]/ancestor::tr//button[contains(@class, "resource-table-btn-action-dropdown")]'
const actionsTriggerButton =
  '//*[@data-test-resource-name="%s"]/ancestor::tr//button[contains(@class, "oc-files-actions-%s-trigger")]'
const selecAllCheckbox = '#resource-table-select-all'
const acceptButton = '.oc-files-actions-enable-sync-trigger'
const pendingShareItem =
  '//div[@id="files-shared-with-me-pending-section"]//tr[contains(@class,"oc-tbody-tr")]'
const showMoreOptionsButton = '#show-more-share-options-btn'
const calendarDatePickerId = 'recipient-datepicker-btn'
const informMessage = '//div[contains(@class,"oc-notification-message-title")]'
const showMoreBtn = '.toggle-shares-list-btn:has-text("Show more")'
const userTypeFilter = '.invite-form-share-role-type'
const userTypeSelector = '.invite-form-share-role-type-item'

export interface ShareArgs {
  page: Page
  resource: string
  recipients: ICollaborator[]
  expirationDate?: string
}

export const openSharingPanel = async function (
  page: Page,
  resource: string,
  via: ActionViaType = fileAction.sideBarPanel
): Promise<void> {
  const folderPaths = resource.split('/')
  const item = folderPaths.pop()

  if (folderPaths.length) {
    await clickResource({ page, path: folderPaths.join('/') })
  }

  switch (via) {
    case fileAction.quickAction:
      await page.locator(util.format(quickShareButton, item)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['appSidebar'],
        'account page'
      )
      break

    case fileAction.sideBarPanel:
      await sidebar.open({ page, resource: item })
      await sidebar.openPanel({ page, name: 'sharing' })
      await page.locator(invitePanel).waitFor()
      break
  }

  // always click on the “Show more” button if it exists
  const showMore = page.locator(showMoreBtn)
  if ((await showMore.count()) > 0) {
    await showMore.click()
  }
}

export type ActionViaType =
  | typeof fileAction.sideBarPanel
  | typeof fileAction.quickAction
  | typeof fileAction.urlNavigation

export interface createShareArgs extends ShareArgs {
  via: ActionViaType
}

export const createShare = async (args: createShareArgs): Promise<void> => {
  const { page, resource, recipients, via } = args

  if (via !== fileAction.urlNavigation) {
    await openSharingPanel(page, resource, via)

    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'app sidebar'
    )
  }
  const expirationDate = recipients[0].expirationDate

  if (expirationDate) {
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['tippyBox'], 'app sidebar')
    await page.locator(showMoreOptionsButton).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['tippyBox'], 'app sidebar')
    await Promise.all([
      locatorUtils.waitForEvent(page.locator(invitePanel), 'transitionend'),
      page.getByTestId(calendarDatePickerId).click()
    ])
    await Collaborator.setExpirationDate(page, expirationDate)
  }
  const federatedShare = recipients[0].shareType
  if (federatedShare) {
    // --- WHY THIS WORKAROUND EXISTS ---
    // The "External users" filter chip (OcFilterChip → OcDrop → Tippy.js) teleports its
    // dropdown content to document.body, outside the Vue component tree. Tippy's own
    // `close-on-click` handler fires on the toggle BEFORE the bubbled event reaches Vue's
    // @click="selectShareRoleType(option)" in InviteCollaboratorForm.vue.
    // With the original page.locator(userTypeFilter).click(), isExternalShareRoleType stayed
    // false: the invite input searched all users (not Federated), Brian Murphy was not found.
    //
    // Fix: open the dropdown via Tippy's JS API (btn._tippy.show()) — bypasses Playwright's
    // synthetic event path. Then fire the item click via dispatchEvent({bubbles:true})
    // directly on the DOM node so it reaches Vue's handler before Tippy can intercept.
    //
    await page.evaluate(() => {
      const btn = document.querySelector(
        '.invite-form-share-role-type .oc-filter-chip-button'
      ) as any
      btn?._tippy?.show()
    })
    await page.locator(userTypeSelector).filter({ hasText: federatedShare }).waitFor()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['tippyBox'], 'app sidebar')
    // Use JS dispatchEvent — Playwright .click() on teleported Tippy content may not reach Vue handlers
    await page.evaluate((shareType) => {
      const items = Array.from(document.querySelectorAll('.invite-form-share-role-type-item'))
      const target = items.find((el) =>
        el.textContent?.toLowerCase().includes(shareType.toLowerCase())
      ) as HTMLElement
      target?.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }))
    }, federatedShare)
    // Small wait for Vue reactivity to update
    await page.waitForTimeout(200)
  }

  await Collaborator.inviteCollaborators({ page, collaborators: recipients })
  await sidebar.close({ page })
}

/**/

export interface ShareStatusArgs extends Omit<ShareArgs, 'recipients'> {
  via?: 'STATUS' | typeof fileAction.contextMenu
}

export const enableSync = async (args: ShareStatusArgs): Promise<void> => {
  const { resource, page } = args
  await clickActionInContextMenu({ page, resource }, 'enable-sync')
}

export const syncAllShares = async ({ page }: { page: Page }): Promise<void> => {
  await page.locator(selecAllCheckbox).click()
  const numberOfPendingShares = await page.locator(pendingShareItem).count()
  const checkResponses = []
  for (let i = 0; i < numberOfPendingShares; i++) {
    checkResponses.push(
      page.waitForResponse(
        (resp) =>
          resp.url().includes('root/children') &&
          resp.status() === 201 &&
          resp.request().method() === 'POST'
      )
    )
  }
  await Promise.all([...checkResponses, page.locator(acceptButton).click()])
  await a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['files'],
    'Files view after sync is enabled for all shared resources'
  )
}

export const disableSync = async (args: ShareStatusArgs): Promise<void> => {
  const { page, resource } = args
  await clickActionInContextMenu({ page, resource }, 'disable-sync')
}

export const clickActionInContextMenu = async (
  args: ShareStatusArgs,
  action: string
): Promise<void> => {
  const { page, resource } = args
  await page.locator(util.format(actionMenuDropdownButton, resource)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['.tippy-box'],
    'context menu dropdown'
  )
  switch (action) {
    case 'enable-sync':
      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().includes('root/children') &&
            resp.status() === 201 &&
            resp.request().method() === 'POST'
        ),
        page.locator(util.format(actionsTriggerButton, resource, action)).click()
      ])
      break
    case 'disable-sync':
      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().includes('drives') &&
            resp.status() === 204 &&
            resp.request().method() === 'DELETE'
        ),
        page.locator(util.format(actionsTriggerButton, resource, action)).click()
      ])
      break
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['#files-shared-with-me-view'],
    'Shared with me file list'
  )
}

export const changeShareeRole = async (args: ShareArgs): Promise<void> => {
  const { page, resource, recipients } = args
  await openSharingPanel(page, resource)

  for (const collaborator of recipients) {
    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().includes('permissions') &&
          resp.status() === 200 &&
          resp.request().method() === 'PATCH'
      ),
      Collaborator.changeCollaboratorRole({ page, collaborator })
    ])
  }
}

/**/

export interface removeShareeArgs extends ShareArgs {
  removeOwnSpaceAccess?: boolean
}

export const removeSharee = async (args: removeShareeArgs): Promise<void> => {
  const { page, resource, recipients, removeOwnSpaceAccess } = args
  await openSharingPanel(page, resource)

  for (const collaborator of recipients) {
    await Collaborator.removeCollaborator({ page, collaborator, removeOwnSpaceAccess })
  }
}

/**/

export const checkSharee = async (args: ShareArgs): Promise<void> => {
  const { resource, page, recipients } = args
  await openSharingPanel(page, resource)

  for (const collaborator of recipients) {
    await Collaborator.checkCollaborator({ page, collaborator })
  }
}

export const addExpirationDate = async (args: {
  page: Page
  resource: string
  collaborator: Omit<ICollaborator, 'role'>
  expirationDate: string
}): Promise<void> => {
  const { page, resource, collaborator, expirationDate } = args
  await openSharingPanel(page, resource)

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('drives') &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    Collaborator.setExpirationDateForCollaborator({ page, collaborator, expirationDate })
  ])
}

export const getAccessDetails = async (args: {
  page: Page
  resource: string
  collaborator: Omit<ICollaborator, 'role'>
}): Promise<IAccessDetails> => {
  const { page, resource, collaborator } = args
  await openSharingPanel(page, resource)

  return Collaborator.getAccessDetails(page, collaborator)
}

export const getMessage = ({ page }: { page: Page }): Promise<string> => {
  return page.locator(informMessage).textContent()
}

export const changeRoleLocator = (args: { page: Page; recipient: User }): Locator => {
  const { page, recipient } = args
  const recipientRow = Collaborator.getCollaboratorUserOrGroupSelector(recipient, 'user')
  return page.locator(util.format(Collaborator.collaboratorRoleDropdownButton, recipientRow))
}

export const changeShareLocator = (args: { page: Page; recipient: User }): Locator => {
  const { page, recipient } = args
  const recipientRow = Collaborator.getCollaboratorUserOrGroupSelector(recipient, 'user')
  return page.locator(util.format(Collaborator.collaboratorEditDropdownButton, recipientRow))
}
