import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application } from '../../environment/constants'

test.describe('Markdown editor inlines images as base64', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" creates the following resources
    //   | resource   | type   | content      |
    //   | notes.md   | mdFile | # image test |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'notes.md', type: 'mdFile', content: '# image test' }]
    })

    // And "Alice" opens the following file in texteditor
    //   | resource |
    //   | notes.md |
    // (creating the file saves and closes the editor, so it has to be reopened)
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'notes.md',
      viewer: application.textEditor
    })
  })

  test.afterEach(async () => {
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('picked image is inlined into the document', async () => {
    // When "Alice" picks the image "testavatar.png" in the markdown editor
    await ui.userPicksImageInMarkdownEditor({ stepUser: 'Alice', image: 'testavatar.png' })

    // Then "Alice" should see the image inlined in the markdown editor
    await ui.userShouldSeeInlinedImageInMarkdownEditor({ stepUser: 'Alice' })
  })

  test('image above the per-image limit is rejected', async () => {
    // When "Alice" picks an image of 2500000 bytes in the markdown editor
    await ui.userPicksOversizedImageInMarkdownEditor({ stepUser: 'Alice', sizeInBytes: 2500000 })

    // Then "Alice" should see the image rejected in the markdown editor
    //   | image         | limit |
    //   | oversized.png | 2 MB  |
    await ui.userShouldSeeImageRejectionInMarkdownEditor({
      stepUser: 'Alice',
      image: 'oversized.png',
      limit: '2 MB'
    })

    // And "Alice" should not see an image inlined in the markdown editor
    await ui.userShouldNotSeeInlinedImageInMarkdownEditor({ stepUser: 'Alice' })
  })

  test('crop and upload is not offered', async () => {
    // Then "Alice" should not see the crop option in the markdown editor
    await ui.userShouldNotSeeCropOptionInMarkdownEditor({ stepUser: 'Alice' })
  })
})
