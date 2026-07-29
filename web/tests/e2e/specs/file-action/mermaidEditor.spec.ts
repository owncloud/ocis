import { expect } from '@playwright/test'
import util from 'util'
import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { selectors as mermaidSelectors } from '../../support/objects/app-files/resource/mermaidEditor'
import { resourceNameSelector } from '../../support/objects/app-files/resource/actions'

const diagramRowSelector = util.format(resourceNameSelector, 'diagram.mmd')

test.describe('Mermaid diagram editor', { tag: '@predefined-users' }, () => {
  test('renders a diagram, surfaces invalid syntax, and switches view modes', async ({ world }) => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // When "Alice" creates the following resource
    //   | resource    | type    | content             |
    //   | diagram.mmd | mmdFile | graph TD; A --> B   |
    // (editMermaidDocument, invoked for this resource type, also asserts the
    // resulting diagram renders and has no severe WCAG 2.1 / WCAG 2.2 violations)
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'diagram.mmd', type: 'mmdFile', content: 'graph TD\n  A --> B' }]
    })

    const { page } = world.actorsEnvironment.getActor({ key: 'Alice' })

    // And "Alice" edits the following resource
    //   | resource    | content                    |
    //   | diagram.mmd | graph TD; A --> B --> C    |
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [
        { name: 'diagram.mmd', type: 'mmdFile', content: 'graph TD\n  A --> B --> C' }
      ]
    })

    // Then re-opening the file shows the rendered diagram and the view-mode
    // toggle works (split -> preview-only -> split)
    await page.locator(diagramRowSelector).click()
    await page.locator(mermaidSelectors.mermaidEditorRoot).waitFor()
    await expect(page.locator(mermaidSelectors.mermaidPreviewDiagram)).toBeVisible()
    await expect(
      page.locator(mermaidSelectors.mermaidEditorRoot).locator('.mermaid-editor-body-split')
    ).toBeVisible()

    await page.locator(`${mermaidSelectors.mermaidEditorRoot} .mermaid-editor-viewmode-preview`).click()
    await expect(
      page.locator(mermaidSelectors.mermaidEditorRoot).locator('.mermaid-editor-body-preview-only')
    ).toBeVisible()

    await page.locator(`${mermaidSelectors.mermaidEditorRoot} .mermaid-editor-viewmode-split`).click()
    await expect(
      page.locator(mermaidSelectors.mermaidEditorRoot).locator('.mermaid-editor-body-split')
    ).toBeVisible()

    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // When "Alice" edits the resource with invalid Mermaid syntax
    // Then the editor surfaces an inline error instead of a blank/crashed pane
    // (editMermaidDocument waits for either the diagram or the error state, and
    // re-asserts WCAG 2.1 / WCAG 2.2 conformance on that error state too)
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [{ name: 'diagram.mmd', type: 'mmdFile', content: 'this is not }} a diagram' }]
    })

    await page.locator(diagramRowSelector).click()
    await expect(page.locator(mermaidSelectors.mermaidPreviewError)).toBeVisible()
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
