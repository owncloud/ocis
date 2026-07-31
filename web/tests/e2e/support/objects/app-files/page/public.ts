import { Download, FrameLocator, Page } from '@playwright/test'
import { File } from '../../../types'
import util from 'util'
import path from 'path'
import * as po from '../resource/actions'
import { objects } from '../../../index'

const passwordInput = 'input[type="password"]'
const fileUploadInput = '//input[@id="files-file-upload-input"]'
const dropUploadResourceSelector = '.upload-info-items [data-test-resource-name="%s"]'
const uploadInfoSuccessLabelSelector = '.upload-info-success'
const publicLinkAuthorizeButton = '.oc-login-authorize-button'
const folderModalIframe = '#iframe-folder-view'
const passwordProtectedPublicLinkForm =
  '//span[contains(text(),"password-protected")]/ancestor::form'
const publicLinkErrorMessage = 'div.oc-link-resolve-error-title'

export class Public {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async open({ url }: { url: string }): Promise<void> {
    await Promise.all([
      this.#page.waitForResponse(
        (res) =>
          res.url().includes('/public-files/') &&
          res.request().method() === 'PROPFIND' &&
          res.status() >= 207
      ),
      this.#page.goto(url)
    ])
    if (
      !(await this.#page.locator(passwordProtectedPublicLinkForm).isVisible()) &&
      !(await this.#page.locator(publicLinkErrorMessage).isVisible())
    ) {
      // wait for redirection to complete
      await this.#page.waitForURL('**/public/**')
    }
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      'public link page'
    )
  }

  async authenticate({
    password,
    passwordProtectedFolder = false,
    expectToSucceed = true
  }: {
    password: string
    passwordProtectedFolder?: boolean
    expectToSucceed?: boolean
  }): Promise<void> {
    let page: Page | FrameLocator = this.#page
    if (passwordProtectedFolder) {
      page = this.#page.frameLocator(folderModalIframe)
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        this.#page,
        ['folderViewModal'],
        'password protected folder modal'
      )
    } else {
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        this.#page,
        ['body'],
        'public link authenticate page'
      )
    }
    await page.locator(passwordInput).fill(password)
    await page.locator(publicLinkAuthorizeButton).click()
    if (expectToSucceed) {
      await page.locator('#web-content').waitFor()
    }
  }

  async dropUpload({ resources }: { resources: File[] }): Promise<void> {
    const startUrl = this.#page.url()
    await this.#page.locator(fileUploadInput).setInputFiles(resources.map((file) => file.path))
    const names = resources.map((file) => path.basename(file.name))
    await this.#page.locator(uploadInfoSuccessLabelSelector).waitFor()
    await Promise.all(
      names.map((name) =>
        this.#page.locator(util.format(dropUploadResourceSelector, name)).waitFor()
      )
    )
    await this.#page.goto(startUrl)
  }

  async reload(): Promise<void> {
    await this.#page.reload()
  }

  async download(args: Omit<po.downloadResourcesArgs, 'page'>): Promise<Download[]> {
    const startUrl = this.#page.url()
    const downloads = await po.downloadResources({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    return downloads
  }

  async rename(args: Omit<po.renameResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.renameResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async upload(args: Omit<po.uploadResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.uploadResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    await this.#page.locator('body').click()
  }

  async uploadInternal(
    args: Omit<po.uploadResourceArgs, 'page'> & { link: string }
  ): Promise<void> {
    // link is the public link url
    const { link } = args
    delete args.link
    await po.uploadResource({ ...args, page: this.#page })
    await this.#page.goto(link)
  }

  async delete(args: Omit<po.deleteResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.deleteResource({ ...args, page: this.#page, isPublicLink: true })
    await this.#page.goto(startUrl)
  }

  async expectThatLinkIsDeleted({ url }: { url: string }): Promise<void> {
    await po.expectThatPublicLinkIsDeleted({ page: this.#page, url })
  }

  async getDocumentContent({ page, editor }: { page: Page; editor: string }): Promise<string> {
    return await po.getDocumentContent({ page, editor })
  }

  async fillDocumentContent({
    page,
    text,
    editor
  }: {
    page: Page
    text: string
    editor: string
  }): Promise<void> {
    return await po.fillDocumentContent({ page, text, editor })
  }
}
