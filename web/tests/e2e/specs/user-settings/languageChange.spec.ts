import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'

test.describe('language settings', { tag: '@predefined-users' }, () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    // Given "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('system language change', async () => {
    // And "Alice" creates the following folder in personal space using API
    //   | name          |
    //   | check_message |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'check_message' })
    // And "Alice" shares the following resource using API
    //   | resource      | recipient | type | role     | resourceType |
    //   | check_message | Brian     | user | Can edit | folder       |
    await api.userHasSharedResources({
      stepUser: 'Alice',
      shares: [
        {
          resource: 'check_message',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit',
          resourceType: 'folder'
        }
      ]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // When "Brian" opens the user menu
    await ui.userOpensAccountPage({ stepUser: 'Brian' })
    // And "Brian" changes the language to "Deutsch - German"
    await ui.userChangesLanguage({ stepUser: 'Brian', language: 'Deutsch - German' })
    // Then "Brian" should see the following account page title "Mein Konto"
    await ui.userShouldSeeAccountPageTitle({ stepUser: 'Brian', expectedTitle: 'Mein Konto' })
    // And "Brian" should see the following notifications
    // | Alice hat check_message mit Ihnen geteilt |
    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: ['Alice Hansen hat check_message mit Ihnen geteilt']
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })

  test('anonymous user language change', async () => {
    // And "Alice" creates the following folder in personal space using API
    //   | name         |
    //   | folderPublic |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'folderPublic' })
    // And "Alice" uploads the following local file into personal space using API
    //   | localFile                | to        |
    //   | filesForUpload/lorem.txt | lorem.txt |
    await api.userHasUploadedFilesInPersonalSpace({
      stepUser: 'Alice',
      filesToUpload: [{ localFile: 'filesForUpload/lorem.txt', to: 'lorem.txt' }]
    })
    // And "Alice" creates a public link of following resource using API
    //   | resource     | password |
    //   | folderPublic | %public% |
    await api.userHasCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'lorem.txt',
      password: '%public%'
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // When "Anonymous" opens the public link "Unnamed link"
    await ui.anonymousUserOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })
    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
    // And "Anonymous" opens the user menu
    await ui.userOpensAccountPage({ stepUser: 'Anonymous' })
    // And "Anonymous" changes the language to "Deutsch - German"
    await ui.userChangesLanguage({ stepUser: 'Anonymous', language: 'Deutsch - German' })
    // Then "Anonymous" should see the following account page title "Mein Konto"
    await ui.userShouldSeeAccountPageTitle({ stepUser: 'Anonymous', expectedTitle: 'Mein Konto' })
  })
})
