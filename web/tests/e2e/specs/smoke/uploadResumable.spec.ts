import { test } from '../../environment/test'
import * as ui from '../../steps/ui/index'
import * as api from '../../steps/api/api'
import { userCreatesAFileOfSizeInTempUploadDir } from '../../steps/runtimeFs'
import { resourcePage } from '../../environment/constants'

test.describe('Upload large resources', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following user using API
    // | id    |
    // | Alice |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
  })

  test('Upload large resources in personal space', async () => {
    // Given the user creates a file "largefile.txt" of "1GB" size in the temp upload directory
    await userCreatesAFileOfSizeInTempUploadDir({ fileName: 'largefile.txt', fileSize: '1GB' })

    // Given "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // When "Alice" starts uploading the following large resources from the temp upload directory
    // | resource      |
    // | largefile.txt |
    await ui.userStartsUploadingFileFromTheTempUploadDir({
      stepUser: 'Alice',
      file: 'largefile.txt'
    })

    // When "Alice" pauses the file upload
    await ui.userPausesUpload({ stepUser: 'Alice' })

    // When "Alice" cancels the file upload
    await ui.userCancelsUpload({ stepUser: 'Alice' })

    // Then following resources should not be displayed in the files list for user "Alice"
    // | resource      |
    // | largefile.txt |
    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['largefile.txt']
    })

    // When "Alice" starts uploading the following large resources from the temp upload directory
    // | resource      |
    // | largefile.txt |
    await ui.userStartsUploadingFileFromTheTempUploadDir({
      stepUser: 'Alice',
      file: 'largefile.txt'
    })

    // When "Alice" pauses the file upload
    await ui.userPausesUpload({ stepUser: 'Alice' })

    // When "Alice" resumes the file upload
    await ui.userResumesUpload({ stepUser: 'Alice' })

    // Then following resources should be displayed in the files list for user "Alice"
    // | resource      |
    // | largefile.txt |
    await ui.userShouldSeeResources({
      listType: resourcePage.filesList,
      stepUser: 'Alice',
      resources: ['largefile.txt']
    })

    // When "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
