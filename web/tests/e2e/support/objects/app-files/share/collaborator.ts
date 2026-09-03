import { Page } from '@playwright/test'
import { startCase } from 'lodash-es'
import util from 'util'
import { Group, User } from '../../../types'
import { getActualExpiryDate } from '../../../utils/datePicker'
import { locatorUtils } from '../../../utils'
import { objects } from '../../../index'

export interface ICollaborator {
  collaborator: User | Group
  role?: string
  type?: CollaboratorType
  resourceType?: string
  expirationDate?: string
  shareType?: string
}

export interface InviteCollaboratorsArgs {
  page: Page
  collaborators: ICollaborator[]
}

export interface CollaboratorArgs {
  page: Page
  collaborator: ICollaborator
}

export interface RemoveCollaboratorArgs extends Omit<CollaboratorArgs, 'collaborator'> {
  collaborator: Omit<ICollaborator, 'role'>
  removeOwnSpaceAccess?: boolean
}

export interface SetExpirationDateForCollaboratorArgs extends Omit<
  CollaboratorArgs,
  'collaborator'
> {
  collaborator: Omit<ICollaborator, 'role'>
  expirationDate: any
}

export interface RemoveExpirationDateFromCollaboratorArgs extends Omit<
  CollaboratorArgs,
  'collaborator'
> {
  collaborator: Omit<ICollaborator, 'role'>
}

export interface SetDenyShareForCollaboratorArgs extends Omit<CollaboratorArgs, 'collaborator'> {
  collaborator: Omit<ICollaborator, 'role'>

  deny: boolean
}

export interface IAccessDetails {
  Name?: string
  'Additional info'?: string
  Type?: string
  'Access expires'?: string
  'Shared on'?: string
  'Invited by'?: string
}

export type CollaboratorType = 'user' | 'group'

export default class Collaborator {
  private static readonly invitePanel = '//*[@id="oc-files-sharing-sidebar"]'
  private static readonly inviteInput = '#files-share-invite-input'
  private static readonly newCollaboratorRoleDropdown =
    '//*[@id="files-collaborators-role-button-new"]'
  private static readonly sendInvitationButton = '#new-collaborators-form-create-button'
  public static readonly collaboratorRoleDropdownButton =
    '%s//button[contains(@class,"files-recipient-role-select-btn")]'
  private static readonly collaboratorRoleItemSelector = '%s//button[contains(@id, "%s")]'
  private static readonly collaboratorRoleButton = '//button[contains(@id, "%s")]'
  public static readonly collaboratorEditDropdownButton =
    '%s//button[contains(@class,"collaborator-edit-dropdown-options-btn")]'
  private static readonly collaboratorUserSelector =
    '//*[starts-with(@data-testid,"collaborator-user-item-%s")]'
  private static readonly collaboratorGroupSelector =
    '//*[starts-with(@data-testid,"collaborator-group-item-%s")]'
  private static readonly collaboratorRoleSelector =
    '%s//button[contains(@class,"files-recipient-role-select-btn")]/span[text()="%s"]'
  private static readonly removeCollaboratorButton =
    '%s//ul[contains(@class,"collaborator-edit-dropdown-options-list")]//button[contains(@class,"remove-share")]'
  private static readonly denyShareCollaboratorButton =
    '%s//ul[contains(@class,"collaborator-edit-dropdown-options-list")]//span[contains(@class,"deny-share")]//button[contains(@aria-checked,"%s")]'
  private static readonly setExpirationDateCollaboratorButton =
    '%s//ul[contains(@class,"collaborator-edit-dropdown-options-list")]//button[contains(@class,"recipient-datepicker-btn")]'
  private static readonly removeExpirationDateCollaboratorButton =
    '%s//ul[contains(@class,"collaborator-edit-dropdown-options-list")]//button[contains(@class,"remove-expiration-date")]'
  private static readonly showAccessDetailsButton =
    '%s//ul[contains(@class,"collaborator-edit-dropdown-options-list")]//button[contains(@class,"show-access-details")]'
  private static readonly removeCollaboratorConfirmationButton = '.oc-modal-body-actions-confirm'
  private static readonly collaboratorExpirationDatepicker = '.oc-modal-body .oc-date-picker input'
  private static readonly collaboratorExpirationDatepickerConfirmButton =
    '.oc-modal-body-actions-confirm'
  private static readonly collaboratorDropdownItem =
    'div[data-testid="new-collaborators-form"] div[data-testid="recipient-autocomplete-item-%s"]'

  static async addCollaborator(args: CollaboratorArgs): Promise<void> {
    const {
      page,
      collaborator: { collaborator, shareType }
    } = args
    const collaboratorInputLocator = page.locator(Collaborator.inviteInput)
    await collaboratorInputLocator.click()
    let fillValue: string
    if (shareType === 'external') {
      fillValue = collaborator.displayName
    } else if ((collaborator as User).originalId) {
      fillValue = collaborator.id
    } else {
      fillValue = collaborator.displayName
    }
    await Promise.all([
      page.waitForResponse((resp) => resp.url().includes('users') && resp.status() === 200),
      collaboratorInputLocator.fill(fillValue)
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'account page'
    )
    await collaboratorInputLocator.focus()
    await page.locator('.vs--open').waitFor()
    await page
      .locator(util.format(Collaborator.collaboratorDropdownItem, collaborator.displayName))
      .first() // in CI, resolves to two elements
      .click()
  }

  static async sendInvitation(page: Page, collaborators: string[]): Promise<void> {
    const checkResponses = []
    for (let i = 0; i < collaborators.length; i++) {
      checkResponses.push(
        page.waitForResponse((resp) => {
          return (
            resp.url().endsWith('invite') &&
            resp.status() === 200 &&
            resp.request().method() === 'POST'
          )
        })
      )
    }
    // --- WHY THIS WORKAROUND EXISTS ---
    // The Share button (#new-collaborators-form-create-button) sits inside vue-select's
    // vs__actions div. Playwright's page.locator(...).click() did not reach Vue's
    // @click="share" handler — POST /graph/.../invite was never made, waitForResponse
    // timed out after 30 s. dispatchEvent fires the click directly on the DOM node.
    await Promise.all([
      ...checkResponses,
      page.evaluate(() => {
        const btn = document.querySelector('#new-collaborators-form-create-button') as HTMLElement
        btn?.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }))
      })
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'Shares panel after sending invitation'
    )
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['notificationContainer'],
      'notification popup after sending invitation'
    )
  }

  static async inviteCollaborators(args: InviteCollaboratorsArgs): Promise<void> {
    const { page, collaborators } = args
    // When adding multiple users/groups at once
    // the role of the first collaborator is used as the collaborators role
    const role = collaborators[0].role
    const resourceType = collaborators[0].resourceType
    const collaboratorNames = []
    for (const collaborator of collaborators) {
      await Collaborator.addCollaborator({ page, collaborator })
      collaboratorNames.push(collaborator.collaborator.displayName)
    }
    await Collaborator.setCollaboratorRole(page, role, resourceType)
    await Collaborator.sendInvitation(page, collaboratorNames)
  }

  static async setCollaboratorRole(
    page: Page,
    role: string,
    resourceType: string,
    dropdownSelector?: string,
    itemSelector?: string
  ): Promise<void> {
    if (!dropdownSelector) {
      dropdownSelector = Collaborator.newCollaboratorRoleDropdown
      itemSelector = Collaborator.collaboratorRoleButton
    }
    // --- WHY THIS WORKAROUND EXISTS ---
    // The role dropdown (#files-collaborators-role-button-new) is a Tippy toggle inside
    // vue-select's vs__actions container. Same mechanism as the filter chip in actions.ts:
    // Playwright .click() is intercepted by Tippy's close-on-click before Vue's
    // @option-change="collaboratorRoleChanged" fires — the role never changes, and the test
    // timed out waiting for //button[contains(@id,"fb6c3e19-...")] (the selected-role indicator).
    //
    // Fix: open via _tippy.show(), then fire via dispatchEvent on the role button.
    // Note: itemSelector is an XPath expression like '//button[contains(@id,"fb6c3e19-...")]'.
    // document.querySelector() rejects XPath syntax (throws DOMException) —
    // document.evaluate() is required to resolve XPath nodes.
    const toggleId = await page.locator(dropdownSelector).getAttribute('id')
    await page.evaluate((id) => {
      const btn = id
        ? document.getElementById(id)
        : (document.querySelector('[id^="files-collaborators-role-button"]') as any)
      btn?._tippy?.show()
    }, toggleId)
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['tippyBox'], 'tippy box')
    // dispatchEvent since Tippy close-on-click intercepts .click() before Vue handler fires
    const roleSelector = util.format(itemSelector, role)
    await page.locator(roleSelector).waitFor()
    await page.evaluate((xpathSelector) => {
      const btn = document.evaluate(
        xpathSelector,
        document,
        null,
        XPathResult.FIRST_ORDERED_NODE_TYPE,
        null
      ).singleNodeValue as HTMLElement
      btn?.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true }))
    }, roleSelector)
  }

  static async changeCollaboratorRole(args: CollaboratorArgs): Promise<void> {
    const {
      page,
      collaborator: { collaborator, type, role, resourceType }
    } = args

    const collaboratorRow = Collaborator.getCollaboratorUserOrGroupSelector(collaborator, type)
    const roleDropdownSelector = util.format(
      Collaborator.collaboratorRoleDropdownButton,
      collaboratorRow
    )
    const roleItemSelector = util.format(Collaborator.collaboratorRoleItemSelector, collaboratorRow)
    await Collaborator.setCollaboratorRole(
      page,
      role,
      resourceType,
      roleDropdownSelector,
      roleItemSelector
    )
  }

  static async removeCollaborator(args: RemoveCollaboratorArgs): Promise<void> {
    const {
      page,
      collaborator: { collaborator, type },
      removeOwnSpaceAccess
    } = args
    const collaboratorRow = Collaborator.getCollaboratorUserOrGroupSelector(collaborator, type)

    await page
      .locator(util.format(Collaborator.collaboratorEditDropdownButton, collaboratorRow))
      .first()
      .click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['tippyBox'], 'files modal')
    await page.locator(util.format(Collaborator.removeCollaboratorButton, collaboratorRow)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['removeUserModal'],
      'files modal'
    )

    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().includes('permissions') &&
          resp.status() === 204 &&
          resp.request().method() === 'DELETE'
      ),
      page.locator(Collaborator.removeCollaboratorConfirmationButton).click()
    ])
    if (removeOwnSpaceAccess) {
      await page.waitForURL(/.*\/files\/spaces.*/)
    }
  }

  static async checkCollaborator(args: CollaboratorArgs): Promise<void> {
    const {
      page,
      collaborator: { collaborator, type, role }
    } = args
    const collaboratorRow = Collaborator.getCollaboratorUserOrGroupSelector(collaborator, type)

    await page.locator(collaboratorRow).waitFor()

    if (role) {
      const parts = role.split(' ')
      const collaboratorRole = `${startCase(parts[0].toLowerCase())} ${
        parts[1] ? `${parts[1].toLowerCase()}` : ''
      }`
      const roleSelector = util.format(
        Collaborator.collaboratorRoleSelector,
        collaboratorRow,
        collaboratorRole
      )
      await page.locator(roleSelector).waitFor()
    }
  }

  static async setExpirationDateForCollaborator(
    args: SetExpirationDateForCollaboratorArgs
  ): Promise<void> {
    const {
      page,
      collaborator: { collaborator, type },
      expirationDate
    } = args
    const collaboratorRow = Collaborator.getCollaboratorUserOrGroupSelector(collaborator, type)
    await page.locator(collaboratorRow).waitFor()

    await page
      .locator(util.format(Collaborator.collaboratorEditDropdownButton, collaboratorRow))
      .click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'account page'
    )

    const panel = page.locator(Collaborator.invitePanel)
    await Promise.all([
      locatorUtils.waitForEvent(panel, 'transitionend'),
      page
        .locator(util.format(Collaborator.setExpirationDateCollaboratorButton, collaboratorRow))
        .click()
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['ocModal'], 'account page')

    await Collaborator.setExpirationDate(page, expirationDate)
  }

  static async setExpirationDate(page: Page, expirationDate: any): Promise<void> {
    const newExpiryDate = getActualExpiryDate(
      expirationDate.toLowerCase().match(/[dayrmonthwek]+/)[0],
      expirationDate
    )

    await page
      .locator(Collaborator.collaboratorExpirationDatepicker)
      .fill(newExpiryDate.toISOString().split('T')[0])
    await page.locator(Collaborator.collaboratorExpirationDatepickerConfirmButton).click()
  }

  static async removeExpirationDateFromCollaborator(
    args: RemoveExpirationDateFromCollaboratorArgs
  ): Promise<void> {
    const {
      page,
      collaborator: { collaborator, type }
    } = args
    const collaboratorRow = Collaborator.getCollaboratorUserOrGroupSelector(collaborator, type)
    await page.locator(collaboratorRow).waitFor()
    await page
      .locator(util.format(Collaborator.collaboratorEditDropdownButton, collaboratorRow))
      .click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'account page'
    )
    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().includes('permissions') &&
          resp.status() === 200 &&
          resp.request().method() === 'PATCH'
      ),
      page
        .locator(util.format(Collaborator.removeExpirationDateCollaboratorButton, collaboratorRow))
        .click()
    ])
  }

  static waitForInvitePanel(page: Page): Promise<void> {
    return page.locator(Collaborator.invitePanel).waitFor()
  }

  static getCollaboratorUserOrGroupSelector = (collaborator: User | Group, type = 'user') => {
    return type === 'group'
      ? util.format(Collaborator.collaboratorGroupSelector, collaborator.displayName)
      : util.format(Collaborator.collaboratorUserSelector, collaborator.displayName)
  }

  static async getAccessDetails(
    page: Page,
    recipient: Omit<ICollaborator, 'role'>
  ): Promise<IAccessDetails> {
    const { collaborator, type } = recipient
    const collaboratorRow = Collaborator.getCollaboratorUserOrGroupSelector(collaborator, type)
    await page
      .locator(util.format(Collaborator.collaboratorEditDropdownButton, collaboratorRow))
      .click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'Tippy box collaborator edit dropdown'
    )
    await page.locator(util.format(Collaborator.showAccessDetailsButton, collaboratorRow)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'Tippy box share access details'
    )

    return page.locator('.share-access-details-drop dl').evaluate((el) => {
      const nodes = el.childNodes
      const details: Record<string, string> = {}
      nodes.forEach((node) => {
        if (node.nodeName === 'DT') {
          details[node.textContent] = node.nextSibling.textContent
        }
      })
      return details
    })
  }
}
