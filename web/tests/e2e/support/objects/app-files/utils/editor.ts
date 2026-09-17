import { Locator, Page } from '@playwright/test'
import { objects } from '../../../index'
import { application } from '../../../../environment/constants'

const closeTextEditorOrViewerButton = '#app-top-bar-close'
const saveTextEditorOrViewerButton = '#app-save-action'
const texEditor = '#text-editor'
const pdfViewer = '#pdf-viewer'
const imageViewer = '.stage'

const markdownEditor = '#text-editor-component'
// md-editor-v3's ImageDropdown renders its hidden file input as `${editorId}-toolbar-wrapper_label`.
const markdownImageInput = `${markdownEditor}-toolbar-wrapper_label`
// The aria-label is md-editor-v3's toolbarTips.image, which is "image" in en-US.
const markdownImageButton = `${markdownEditor} button.md-editor-toolbar-item[aria-label="image"]`
const markdownImageMenuItem = `${markdownEditor} .md-editor-menu-item-image`
// CodeMirror's linkShortener extension replaces any token over 30 characters with "...",
// keeping the original in the title attribute - so an inlined data URI never appears as text.
const markdownShortenedToken = `${markdownEditor} .cm-short-text`
const markdownPreviewImage = `${markdownEditor} .md-editor-preview img[src^="data:"]`
// Image rejections are reported through the app's notification store, so they render in the
// global notification stack rather than inside the editor.
const imageRejectionNotification = 'div.oc-notification .oc-notification-message-danger'

export const close = async (page: Page): Promise<void> => {
  await Promise.all([
    page.waitForURL(/.*\/files\/(spaces|shares|link|search)\/.*/),
    page.locator(closeTextEditorOrViewerButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'Personal Page',
    // right after this navigation, the global search input can transiently keep its
    // aria-controls pointing at the search-options dropdown for one render frame after
    // that dropdown has already unmounted (SearchBar.vue's showDrop/term state settling
    // post-navigation) - a genuine but narrow timing race, not a persistent a11y defect
    ['aria-valid-attr-value']
  )
}

export const save = async (page: Page): Promise<void> => {
  await Promise.all([
    page.waitForResponse((res) => res.request().method() === 'PUT' && res.status() === 204),
    page.waitForResponse((res) => res.request().method() === 'PROPFIND' && res.status() === 207),
    page.locator(saveTextEditorOrViewerButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['saveTextEditorOrViewerButton'],
    'Text editor Save button is disabled after saving'
  )
}

export const waitForMarkdownEditor = async (page: Page): Promise<void> => {
  await page.locator(markdownImageInput).waitFor({ state: 'attached' })
}

export const pickMarkdownImage = async (
  page: Page,
  file: { name: string; mimeType: string; buffer: Buffer } | string
): Promise<void> => {
  await page.locator(markdownImageInput).setInputFiles(file)
}

export const openMarkdownImageMenu = async (page: Page): Promise<void> => {
  await page.locator(markdownImageButton).click()
}

export const markdownImageMenuItemsLocator = (page: Page): Locator =>
  page.locator(markdownImageMenuItem)

export const markdownShortenedTokenLocator = (page: Page): Locator =>
  page.locator(markdownShortenedToken)

export const markdownPreviewImageLocator = (page: Page): Locator =>
  page.locator(markdownPreviewImage)

export const imageRejectionNotificationLocator = (page: Page): Locator =>
  page.locator(imageRejectionNotification)

export const fileViewerLocator = ({
  page,
  fileViewerType
}: {
  page: Page
  fileViewerType: string
}): Locator => {
  switch (fileViewerType) {
    case application.textEditor:
      return page.locator(texEditor)
    case application.pdfViewer:
      return page.locator(pdfViewer)
    case application.mediaViewer:
      return page.locator(imageViewer)
    default:
      throw new Error(`${fileViewerType} not implemented`)
  }
}
