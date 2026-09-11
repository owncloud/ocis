import { Page, expect } from '@playwright/test'

const externalEditorIframe = '[name="app-iframe"]'
// Collabora
const collaboraDocPermissionModeSelector = '#permissionmode-container'
const collaboraEditorSaveSelector = '.notebookbar-shortcuts-bar #save'
const collaboraDocTextAreaSelector = '#clipboard-area'
const collaboraCanvasEditorSelector = '.leaflet-layer'
const collaboraWelcomeModal = '.iframe-welcome-modal'
const collaboraWelcomeSlide = '.slider > a'
const collaboraWelcomeClose = '//button[text()="Close"]'
// OnlyOffice
const onlyOfficeInnerFrameSelector = '[name="frameEditor"]'
const onlyOfficeSaveButtonSelector = '#slot-btn-dt-save > button'
const onlyofficeDocTextAreaSelector = '#area_id'
const onlyOfficeCanvasEditorSelector = '#id_viewer_overlay'

export const removeCollaboraWelcomeModal = async (page: Page): Promise<void> => {
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  await editorMainFrame.locator('#document-header').waitFor()
  const versionSet = await editorMainFrame.locator('body').evaluate(() => {
    return localStorage.getItem('WSDWelcomeVersion')
  })
  if (!versionSet) {
    await editorMainFrame.locator(collaboraWelcomeModal).waitFor()
    const welcomeModal = editorMainFrame.frameLocator(collaboraWelcomeModal)
    await welcomeModal.locator(collaboraWelcomeSlide).last().click()
    await welcomeModal.locator(collaboraWelcomeClose).click()
  }
}

export const waitForCollaboraEditor = async (page: Page): Promise<void> => {
  await removeCollaboraWelcomeModal(page)
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  await editorMainFrame.locator(collaboraDocTextAreaSelector).waitFor()
}

export const waitForOnlyOfficeEditor = async (page: Page): Promise<void> => {
  const editorMainFrame = page
    .frameLocator(externalEditorIframe)
    .frameLocator(onlyOfficeInnerFrameSelector)
  await editorMainFrame.locator(onlyofficeDocTextAreaSelector).waitFor()
}

export const focusCollaboraEditor = async (page: Page): Promise<void> => {
  await removeCollaboraWelcomeModal(page)
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  await editorMainFrame.locator(collaboraCanvasEditorSelector).click()
}

export const focusOnlyOfficeEditor = async (page: Page): Promise<void> => {
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  const innerFrame = editorMainFrame.frameLocator(onlyOfficeInnerFrameSelector)
  await innerFrame.locator(onlyOfficeCanvasEditorSelector).click()
}

export const getOfficeDocumentContent = async (page: Page): Promise<string> => {
  // navigator.clipboard calls throw "Document is not focused" unless this page's tab is the
  // OS-focused one - with multiple actors (e.g. Alice + Brian) each owning their own page/tab,
  // the one we're about to read from isn't necessarily the foregrounded one.
  await page.bringToFront()
  // clear the clipboard
  await page.evaluate("navigator.clipboard.writeText('')")
  // copying and getting the value with keyboard requires some time
  await page.keyboard.press('ControlOrMeta+A', { delay: 200 })
  await page.keyboard.press('ControlOrMeta+C', { delay: 200 })
  return page.evaluate(() => navigator.clipboard.readText())
}

export const fillCollaboraDocumentContent = async (page: Page, content: string): Promise<void> => {
  await removeCollaboraWelcomeModal(page)

  const editorMainFrame = page.frameLocator(externalEditorIframe)
  await editorMainFrame.locator(collaboraDocTextAreaSelector).focus()
  await page.keyboard.press('ControlOrMeta+A')
  await editorMainFrame.locator(collaboraDocTextAreaSelector).fill(content)
  const saveLocator = editorMainFrame.locator(collaboraEditorSaveSelector)
  await expect(saveLocator).toHaveAttribute('class', /.*savemodified.*/)
  await saveLocator.click()
  await expect(saveLocator).not.toHaveAttribute('class', /.*savemodified.*/)
  // allow the document to save
  await page.waitForTimeout(500)
}

export const fillOnlyOfficeDocumentContent = async (page: Page, content: string): Promise<void> => {
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  const innerIframe = editorMainFrame.frameLocator(onlyOfficeInnerFrameSelector)
  const textAreaLocator = innerIframe.locator(onlyofficeDocTextAreaSelector)
  const saveButtonDisabledLocator = innerIframe.locator(onlyOfficeSaveButtonSelector)

  // right after the document opens, the hidden text area can already be focusable even though
  // OnlyOffice hasn't finished initializing internally yet, so Ctrl+A/paste can be a silent
  // no-op that leaves the document content unchanged - verify it actually landed and retry if
  // it didn't, instead of trusting a single attempt.
  let actualContent = ''
  for (let attempt = 1; attempt <= 5; attempt++) {
    await textAreaLocator.focus()
    await page.keyboard.press('ControlOrMeta+A')
    await textAreaLocator.fill(content)
    await expect(saveButtonDisabledLocator).toHaveAttribute('disabled', 'disabled')

    actualContent = (await getOfficeDocumentContent(page)).trim()
    if (actualContent === content.trim()) {
      return
    }
    await page.waitForTimeout(1000)
  }
  expect(actualContent).toBe(content.trim())
}

export const canEditCollaboraDocument = async (page: Page): Promise<boolean> => {
  await removeCollaboraWelcomeModal(page)
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  const collaboraDocPermissionModeLocator = editorMainFrame.locator(
    collaboraDocPermissionModeSelector
  )
  const permissionMode = (await collaboraDocPermissionModeLocator.innerText()).trim()
  return permissionMode === 'Edit'
}

export const canEditOnlyOfficeDocument = async (page: Page): Promise<boolean> => {
  const editorMainFrame = page.frameLocator(externalEditorIframe)
  const innerFrame = editorMainFrame.frameLocator(onlyOfficeInnerFrameSelector)
  try {
    await expect(innerFrame.locator(onlyOfficeSaveButtonSelector)).toBeVisible()
    return true
  } catch {
    return false
  }
}
