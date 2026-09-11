import { Page, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { objects } from '../../../index'
import { editor } from '../utils'
import { config } from '../../../../config'

const mermaidEditorRoot = '#mermaid-editor'
const mermaidEditorContent = `${mermaidEditorRoot} .cm-content`
const mermaidPreviewDiagram = `${mermaidEditorRoot} .mermaid-preview-pane-diagram svg`
const mermaidPreviewError = `${mermaidEditorRoot} .mermaid-preview-pane-error`
const saveMermaidFileButton = '#app-save-action:visible'

export const selectors = {
  mermaidEditorRoot,
  mermaidEditorContent,
  mermaidPreviewDiagram,
  mermaidPreviewError
}

// The mermaid editor is checked against WCAG 2.2 from day one via its own tag
// list, kept local to this file rather than added to the shared a11y helper
// (tests/e2e/support/objects/a11y/actions.ts). That helper backs every other
// existing e2e a11y check in the repo and is still scoped to WCAG 2.1 -
// widening it to 2.2 for the whole app is a separate, repo-wide decision with
// its own triage cost for pre-existing violations, out of scope here.
const wcag22Tags = ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa', 'best-practice']

const assertNoSevereWcag22Violations = async (page: Page, label: string): Promise<void> => {
  if (config.skipA11y) {
    return
  }
  const { violations } = await new AxeBuilder({ page })
    .withTags(wcag22Tags)
    .include(mermaidEditorRoot)
    .analyze()
  const severe = violations.filter((v) => v.impact === 'critical' || v.impact === 'serious')
  expect(severe, `Found ${severe.length} severe WCAG 2.2 violations in ${label}`).toHaveLength(0)
}

export const editMermaidDocument = async ({
  page,
  content
}: {
  page: Page
  content: string
}): Promise<void> => {
  await page.locator(mermaidEditorRoot).waitFor()
  await page.locator(mermaidEditorContent).fill(content)
  // the preview is debounced (see web-app-mermaid-editor's App.vue), so wait for
  // it to settle into either a rendered diagram or an inline error before
  // scanning - never mid-debounce against a transiently empty pane
  await Promise.race([
    page.locator(mermaidPreviewDiagram).waitFor(),
    page.locator(mermaidPreviewError).waitFor()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [mermaidEditorRoot],
    'mermaid editor'
  )
  await assertNoSevereWcag22Violations(page, 'mermaid editor')
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 204 && resp.request().method() === 'PUT'),
    page.waitForResponse((resp) => resp.status() === 207 && resp.request().method() === 'PROPFIND'),
    page.locator(saveMermaidFileButton).click()
  ])
  await editor.close(page)
}
