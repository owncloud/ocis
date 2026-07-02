import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction } from '../../environment/constants'
import { test } from '../../environment/test'

test.describe('spaces.personal', () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian', 'Carol'] })
    // And "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('unstructured collection of testable space interactions,', async () => {
    // When "Alice" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'files' })
    // And "Alice" navigates to the projects space page
    await ui.userNavigatesToSpacesPage({ stepUser: 'Alice' })
    // And "Alice" creates the following project spaces
    //   | name  | id     |
    //   | team  | team.1 |
    //   | team2 | team.2 |
    await ui.userCreatesProjectSpaces({
      stepUser: 'Alice',
      spaces: [
        { name: 'team', id: 'team.1' },
        { name: 'team2', id: 'team.2' }
      ]
    })

    // team.1
    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })
    // And "Alice" updates the space "team.1"
    await ui.userUpdatesSpace({
      stepUser: 'Alice',
      key: 'team.1',
      updates: [
        { attribute: 'name', value: 'developer team' },
        { attribute: 'subtitle', value: 'developer team - subtitle' },
        { attribute: 'description', value: 'developer team - description' },
        { attribute: 'quota', value: '50' },
        { attribute: 'image', value: 'testavatar.png' }
      ]
    })

    // shared examples
    // And "Alice" creates the following resources
    //   | resource         | type   |
    //   | folderPublic     | folder |
    //   | folder_to_shared | folder |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'folderPublic', type: 'folder' },
        { name: 'folder_to_shared', type: 'folder' }
      ]
    })
    // And "Alice" uploads the following resources
    //   | resource  | to               |
    //   | lorem.txt | folderPublic     |
    //   | lorem.txt | folder_to_shared |
    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [
        { name: 'lorem.txt', to: 'folderPublic' },
        { name: 'lorem.txt', to: 'folder_to_shared' }
      ]
    })
    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource     | role             | password |
    //   | folderPublic | Secret File Drop | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'folderPublic',
      role: 'Secret File Drop',
      password: '%public%'
    })
    // And "Alice" renames the most recently created public link of resource "folderPublic" to "team.1"
    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      newName: 'team.1'
    })
    // And "Alice" sets the expiration date of the public link named "team.1" of resource "folderPublic" to "+5 days"
    await ui.userSetsExpirationDateOfThePublicLinkOfResource({
      stepUser: 'Alice',
      linkName: 'team.1',
      resource: 'folderPublic',
      expireDate: '+5 days'
    })

    // borrowed from share.feature
    // When "Alice" shares the following resource using the sidebar panel
    //   | resource         | recipient | type | role                   | resourceType |
    //   | folder_to_shared | Brian     | user | Can edit with trashbin | folder       |
    await ui.userSharesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      shares: [
        {
          resource: 'folder_to_shared',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        }
      ]
    })

    // team.2
    // And "Alice" navigates to the project space "team.2"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.2' })
    // And "Alice" updates the space "team.2"
    await ui.userUpdatesSpace({
      stepUser: 'Alice',
      key: 'team.2',
      updates: [
        { attribute: 'name', value: 'management team' },
        { attribute: 'subtitle', value: 'management team - subtitle' },
        { attribute: 'description', value: 'management team - description' },
        { attribute: 'quota', value: '500' },
        { attribute: 'image', value: 'sampleGif.gif' }
      ]
    })
    // And "Alice" creates the following resources
    //   | resource     | type   |
    //   | folderPublic | folder |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'folderPublic', type: 'folder' }]
    })
    // And "Alice" uploads the following resources
    //   | resource  | to           |
    //   | lorem.txt | folderPublic |
    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [{ name: 'lorem.txt', to: 'folderPublic' }]
    })
    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource     | password |
    //   | folderPublic | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'folderPublic',
      password: '%public%'
    })
    // And "Alice" renames the most recently created public link of resource "folderPublic" to "team.2"
    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      newName: 'team.2'
    })
    // And "Alice" edits the public link named "team.2" of resource "folderPublic" changing role to "Secret File Drop"
    await ui.userChangesRoleOfPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'folderPublic',
      linkName: 'team.2',
      newRole: 'Secret File Drop'
    })
    // And "Alice" sets the expiration date of the public link named "team.2" of resource "folderPublic" to "+5 days"
    await ui.userSetsExpirationDateOfThePublicLinkOfResource({
      stepUser: 'Alice',
      linkName: 'team.2',
      resource: 'folderPublic',
      expireDate: '+5 days'
    })
    // And "Alice" changes the password of the public link named "team.2" of resource "folderPublic" to "new-strongPass1"
    await ui.userChangesThePasswordOfPublicLink({
      stepUser: 'Alice',
      linkName: 'team.2',
      resource: 'folderPublic',
      newPassword: 'new-strongPass1'
    })

    // borrowed from link.feature, all existing resource actions can be reused
    // When "Anonymous" opens the public link "team.1"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'team.1' })
    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })
    // And "Anonymous" drop uploads following resources
    //   | resource     |
    //   | textfile.txt |
    await ui.userDropUploadsResources({ stepUser: 'Anonymous', resources: ['textfile.txt'] })

    // borrowed from share.feature
    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" opens the "files" app
    await ui.userOpensApplication({ stepUser: 'Brian', name: 'files' })
    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })
    // And "Brian" renames the following resource
    //   | resource                   | as            |
    //   | folder_to_shared/lorem.txt | lorem_new.txt |
    await ui.userRenamesResource({
      stepUser: 'Brian',
      resource: 'folder_to_shared/lorem.txt',
      newResourceName: 'lorem_new.txt'
    })
    // And "Brian" uploads the following resource
    //   | resource   | to               |
    //   | simple.pdf | folder_to_shared |
    await ui.userUploadsResources({
      stepUser: 'Brian',
      resources: [{ name: 'simple.pdf', to: 'folder_to_shared' }]
    })
    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })
    // And "Alice" updates the space "team.1" image to "testavatar.jpeg"
    await ui.userUpdatesSpace({
      stepUser: 'Alice',
      key: 'team.1',
      updates: [{ attribute: 'image', value: 'testavatar.jpeg' }]
    })
    // And "Alice" uploads the following resource
    //   | resource          | to               | option  |
    //   | PARENT/simple.pdf | folder_to_shared | replace |
    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [{ name: 'PARENT/simple.pdf', to: 'folder_to_shared', option: 'replace' }]
    })
    // And "Brian" should not see the version panel for the file
    //   | resource   | to               |
    //   | simple.pdf | folder_to_shared |
    await ui.userShouldNotSeeVersionPanelForFiles({
      stepUser: 'Brian',
      file: 'simple.pdf',
      to: 'folder_to_shared'
    })

    // When "Alice" deletes the following resources using the sidebar panel
    //   | resource         | from             |
    //   | lorem_new.txt    | folder_to_shared |
    //   | folder_to_shared |                  |
    await ui.userDeletesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [
        {
          name: 'lorem_new.txt',
          from: 'folder_to_shared'
        },
        {
          name: 'folder_to_shared'
        }
      ]
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
    // alice is done
    // When "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })

    // borrowed from link.feature, all existing resource actions can be reused
    // When "Anonymous" opens the public link "team.2"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'team.2' })
    // And "Anonymous" unlocks the public link with password "new-strongPass1"
    await ui.userUnlocksPublicLink({ password: 'new-strongPass1', stepUser: 'Anonymous' })
    // And "Anonymous" drop uploads following resources
    //   | resource     |
    //   | textfile.txt |
    await ui.userDropUploadsResources({ stepUser: 'Anonymous', resources: ['textfile.txt'] })
  })

  test('members of the space can control the versions of the files', async () => {
    // And "Alice" creates the following project space using API
    //   | name | id     |
    //   | team | team.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })
    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })
    // And "Alice" creates the following resources
    //   | resource            | type    | content             |
    //   | parent/textfile.txt | txtFile | some random content |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        {
          name: 'parent/textfile.txt',
          type: 'txtFile',
          content: 'some random content'
        }
      ]
    })
    // When "Alice" uploads the following resources
    //   | resource     | to     | option  |
    //   | textfile.txt | parent | replace |
    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [{ name: 'textfile.txt', to: 'parent', option: 'replace' }]
    })
    // And "Alice" adds following users to the project space
    //   | user  | role                                | kind |
    //   | Carol | Can view                            | user |
    //   | Brian | Can edit with versions and trashbin | user |
    await ui.userAddsMembersToSpace({
      stepUser: 'Alice',
      members: [
        {
          user: 'Carol',
          role: 'Can view',
          kind: 'user'
        },
        {
          user: 'Brian',
          role: 'Can edit with versions and trash bin',
          kind: 'user'
        }
      ]
    })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
    // When "Carol" logs in
    await ui.userLogsIn({ stepUser: 'Carol' })
    // And "Carol" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Carol', space: 'team.1' })
    // And "Carol" should not see the version panel for the file
    //   | resource     | to     |
    //   | textfile.txt | parent |
    await ui.userShouldNotSeeVersionPanelForFiles({
      stepUser: 'Carol',
      file: 'textfile.txt',
      to: 'parent'
    })
    // And "Carol" logs out
    await ui.userLogsOut({ stepUser: 'Carol' })

    // When "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })
    // And "Brian" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'team.1' })
    // And "Brian" downloads old version of the following resource
    //   | resource     | to     |
    //   | textfile.txt | parent |
    await ui.userDownloadsPreviousVersionOfResource({
      stepUser: 'Brian',
      resource: 'textfile.txt',
      to: 'parent'
    })
    // And "Brian" restores following resources version
    //   | resource     | to     | version | openDetailsPanel |
    //   | textfile.txt | parent | 1       | true             |
    await ui.userRestoresResourceVersion({
      stepUser: 'Brian',
      file: 'textfile.txt',
      to: 'parent',
      version: 1,
      openDetailsPanel: true
    })
    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
