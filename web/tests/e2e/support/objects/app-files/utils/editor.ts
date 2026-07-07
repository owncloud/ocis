import { Locator, Page } from '@playwright/test'
import { objects } from '../../../index'
import { application } from '../../../../environment/constants'

const closeTextEditorOrViewerButton = '#app-top-bar-close'
const saveTextEditorOrViewerButton = '#app-save-action'
const texEditor = '#text-editor'
const pdfViewer = '#pdf-viewer'
const imageViewer = '.stage'

export const close = async (page: Page): Promise<void> => {
  await Promise.all([
    page.waitForURL(/.*\/files\/(spaces|shares|link|search)\/.*/),
    page.locator(closeTextEditorOrViewerButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['body'], 'Personal Page')
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
