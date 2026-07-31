import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application } from '../../environment/constants'

test.describe('Different file viewers', { tag: '@predefined-users' }, () => {
  test('file viewers', async () => {
    // Given "Admin" creates following user using API
    //   | id    |
    //   | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // When "Alice" creates the following resources
    //   | resource  | type    | content   |
    //   | lorem.txt | txtFile | some text |
    //   | lorem.md  | mdFile  | readme    |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'lorem.txt', type: 'txtFile', content: 'some text' },
        { name: 'lorem.md', type: 'mdFile', content: 'readme' }
      ]
    })

    // And "Alice" edits the following resources
    //   | resource  | content                   |
    //   | lorem.txt | new content edited        |
    //   | lorem.md  | new readme content edited |
    await ui.userEditsResources({
      stepUser: 'Alice',
      resources: [
        { name: 'lorem.txt', content: 'new content edited' },
        { name: 'lorem.md', content: 'new readme content edited' }
      ]
    })

    // And "Alice" uploads the following resource
    //   | resource        |
    //   | simple.pdf      |
    //   | sampleGif.gif   |
    //   | testimage.mp3   |
    //   | sampleOgg.ogg   |
    //   | sampleWebm.webm |
    //   | test_video.mp4  |
    //   | testavatar.jpeg |
    //   | testavatar.png  |
    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [
        { name: 'simple.pdf' },
        { name: 'sampleGif.gif' },
        { name: 'testimage.mp3' },
        { name: 'sampleOgg.ogg' },
        { name: 'sampleWebm.webm' },
        { name: 'test_video.mp4' },
        { name: 'testavatar.jpeg' },
        { name: 'testavatar.png' }
      ]
    })

    // Then "Alice" should see thumbnail and preview for file "sampleGif.gif"
    await ui.userShouldSeeThumbnailAndPreviewForFile({
      stepUser: 'Alice',
      resource: 'sampleGif.gif'
    })

    // And "Alice" should see thumbnail and preview for file "testavatar.jpeg"
    await ui.userShouldSeeThumbnailAndPreviewForFile({
      stepUser: 'Alice',
      resource: 'testavatar.jpeg'
    })

    // And "Alice" should see thumbnail and preview for file "testavatar.png"
    await ui.userShouldSeeThumbnailAndPreviewForFile({
      stepUser: 'Alice',
      resource: 'testavatar.png'
    })

    // When "Alice" opens a file "testavatar.png" in the media-viewer using the sidebar panel
    await ui.userOpensMediaUsingSidebarPanel({ stepUser: 'Alice', resource: 'testavatar.png' })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" opens the following file in mediaviewer
    //   | resource        |
    //   | testavatar.jpeg |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'testavatar.jpeg',
      viewer: application.mediaViewer
    })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // When "Alice" navigates to the next media resource
    await ui.userNavigatesToMediaResource({ stepUser: 'Alice', navigationType: 'next' })

    // And "Alice" navigates to the previous media resource
    await ui.userNavigatesToMediaResource({ stepUser: 'Alice', navigationType: 'previous' })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" opens the following file in mediaviewer
    //   | resource      |
    //   | sampleGif.gif |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'sampleGif.gif',
      viewer: application.mediaViewer
    })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" opens the following file in mediaviewer
    //   | resource      |
    //   | testimage.mp3 |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'testimage.mp3',
      viewer: application.mediaViewer
    })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" opens the following file in mediaviewer
    //   | resource      |
    //   | sampleOgg.ogg |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'sampleOgg.ogg',
      viewer: application.mediaViewer
    })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" opens the following file in mediaviewer
    //   | resource        |
    //   | sampleWebm.webm |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'sampleWebm.webm',
      viewer: application.mediaViewer
    })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // When "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" opens the following file in mediaviewer
    //   | resource       |
    //   | test_video.mp4 |
    await ui.userOpensResourceInViewer({
      stepUser: 'Alice',
      resource: 'test_video.mp4',
      viewer: application.mediaViewer
    })

    // Then "Alice" is in a media-viewer
    await ui.userShouldBeInFileViewer({
      stepUser: 'Alice',
      fileViewerType: application.mediaViewer
    })

    // And "Alice" closes the file viewer
    await ui.userClosesFileViewer({ stepUser: 'Alice' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
