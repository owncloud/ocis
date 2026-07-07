import { test } from '../../environment/test'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'
import { fileAction } from '../../environment/constants'

test.describe('check files pagination in project space', () => {
  test('pagination', async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    //   | Carol |
    //   | David |
    //   | Edith |
    await api.usersHaveBeenCreated({
      stepUser: 'Admin',
      users: ['Alice', 'Brian', 'Carol', 'David', 'Edith']
    })

    // And "Admin" creates following group using API
    //   | id       |
    //   | sales    |
    //   | security |
    await api.groupsHaveBeenCreated({ groupIds: ['sales', 'security'], stepUser: 'Admin' })

    // And "Admin" adds user to the group using API
    //   | user  | group    |
    //   | David | sales    |
    //   | Edith | security |
    await api.usersHaveBeenAddedToGroup({
      stepUser: 'Admin',
      usersToAdd: [
        { user: 'David', group: 'sales' },
        { user: 'Edith', group: 'security' }
      ]
    })

    // And "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" creates the following project space using API
    //   | name | id     |
    //   | team | team.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    // And "Alice" adds following users to the project space
    //   | user     | role     | kind  |
    //   | Brian    | Can edit | user  |
    //   | Carol    | Can view | user  |
    //   | sales    | Can view | group |
    //   | security | Can edit | group |
    await ui.userAddsMembersToSpace({
      stepUser: 'Alice',
      members: [
        { user: 'Brian', role: 'Can edit', kind: 'user' },
        { user: 'Carol', role: 'Can view', kind: 'user' },
        { user: 'sales', role: 'Can view', kind: 'group' },
        { user: 'security', role: 'Can edit', kind: 'group' }
      ]
    })

    // When "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Brian', space: 'team.1' })

    // And "Brian" creates the following resources
    //   | resource | type   |
    //   | parent   | folder |
    await ui.userCreatesResources({
      stepUser: 'Brian',
      resources: [{ name: 'parent', type: 'folder' }]
    })

    // And "Brian" uploads the following resources
    //   | resource  | to     |
    //   | lorem.txt | parent |
    await ui.userUploadsResources({
      stepUser: 'Brian',
      resources: [{ name: 'lorem.txt', to: 'parent' }]
    })

    // When "David" logs in
    await ui.userLogsIn({ stepUser: 'David' })

    // And "David" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'David', space: 'team.1' })

    // Then "David" should not be able to edit folder "parent"
    await ui.userShouldNotBeAbleToEditResource({ stepUser: 'David', resource: 'parent' })

    // And "David" logs out
    await ui.userLogsOut({ stepUser: 'David' })

    // When "Edith" logs in
    await ui.userLogsIn({ stepUser: 'Edith' })

    // And "Edith" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Edith', space: 'team.1' })

    // And "Edith" creates the following resources
    //   | resource | type   |
    //   | edith    | folder |
    await ui.userCreatesResources({
      stepUser: 'Edith',
      resources: [{ name: 'edith', type: 'folder' }]
    })

    // And "Edith" uploads the following resources
    //   | resource  | to    |
    //   | lorem.txt | edith |
    await ui.userUploadsResources({
      stepUser: 'Edith',
      resources: [{ name: 'lorem.txt', to: 'edith' }]
    })

    // And "Edith" logs out
    await ui.userLogsOut({ stepUser: 'Edith' })

    // When "Carol" logs in
    await ui.userLogsIn({ stepUser: 'Carol' })

    // And "Carol" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Carol', space: 'team.1' })

    // Then "Carol" should not be able to edit folder "parent"
    await ui.userShouldNotBeAbleToEditResource({ stepUser: 'Carol', resource: 'parent' })

    // And "Alice" creates a public link of following resource using the sidebar panel
    //   | resource | role     | password |
    //   | parent   | Can edit | %public% |
    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'parent',
      role: 'Can edit',
      password: '%public%'
    })

    // And "Anonymous" opens the public link "Unnamed link"
    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'Unnamed link' })

    // And "Anonymous" unlocks the public link with password "%public%"
    await ui.userUnlocksPublicLink({ password: '%public%', stepUser: 'Anonymous' })

    // And "Anonymous" uploads the following resources in public link page
    //   | resource     |
    //   | textfile.txt |
    await ui.userUploadsResourcesInPublicLink({
      stepUser: 'Anonymous',
      resources: [{ name: 'textfile.txt' }]
    })

    // And "Anonymous" deletes the following resources from public link using sidebar panel
    //   | resource  | from |
    //   | lorem.txt |      |
    await ui.userDeletesResourcesFromPublicLink({
      stepUser: 'Anonymous',
      actionType: fileAction.sideBarPanel,
      resources: [{ resource: 'lorem.txt' }]
    })

    // When "Brian" deletes the following resources using the sidebar panel
    //   | resource     | from   |
    //   | textfile.txt | parent |
    await ui.userDeletesResources({
      stepUser: 'Brian',
      actionType: fileAction.sideBarPanel,
      resources: [{ name: 'textfile.txt', from: 'parent' }]
    })

    // When "Carol" navigates to the trashbin of the project space "team.1"
    await ui.userNavigatesToTrashbinOfSpace({ stepUser: 'Carol', space: 'team.1' })

    // Then "Carol" should not be able to delete following resources from the trashbin
    //   | resource            |
    //   | parent/lorem.txt    |
    //   | parent/textfile.txt |
    await ui.userShouldNotBeAbleToDeleteResourceFromTrashbin({
      stepUser: 'Carol',
      resources: ['parent/textfile.txt', 'parent/lorem.txt']
    })
    // And "Carol" should not be able to restore following resources from the trashbin
    //   | resource            |
    //   | parent/lorem.txt    |
    //   | parent/textfile.txt |
    await ui.userShouldNotBeAbleToRestoreResourceFromTrashbin({
      stepUser: 'Carol',
      resources: ['parent/lorem.txt', 'parent/textfile.txt']
    })
    // When "Brian" navigates to the trashbin of the project space "team.1"
    await ui.userNavigatesToTrashbinOfSpace({ stepUser: 'Brian', space: 'team.1' })

    // Then "Brian" should be able to restore following resource from the trashbin
    //   | resource         |
    //   | parent/lorem.txt |
    await ui.userShouldBeAbleToRestoreResourceFromTrashbin({
      stepUser: 'Brian',
      resources: ['parent/lorem.txt']
    })

    // And "Brian" should not be able to delete following resource from the trashbin
    //   | resource            |
    //   | parent/textfile.txt |
    await ui.userShouldNotBeAbleToDeleteResourceFromTrashbin({
      stepUser: 'Brian',
      resources: ['parent/textfile.txt']
    })

    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    // And "Alice" removes access to following users from the project space
    //   | user  |
    //   | Brian |
    await ui.userRemovesAccessToMember({ stepUser: 'Alice', reciver: 'Brian', role: 'role' })

    // Then "Brian" should not see space "team.1"
    await ui.userShouldNotSeeSpace({ stepUser: 'Brian', space: 'team.1' })

    // // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })

    // When "Alice" changes the roles of the following users in the project space
    //   | user  | role       |
    //   | Carol | Can manage |
    await ui.userChangesMemberRole({ stepUser: 'Alice', role: 'Can manage', sharee: 'Carol' })

    // And "Carol" navigates to the trashbin of the project space "team.1"
    await ui.userNavigatesToTrashbinOfSpace({ stepUser: 'Carol', space: 'team.1' })

    // Then "Carol" should be able to delete following resource from the trashbin
    //   | resource            |
    //   | parent/textfile.txt |
    await ui.userShouldBeAbleToDeleteResourceFromTrashbin({
      stepUser: 'Carol',
      resources: ['parent/textfile.txt']
    })

    // And "Carol" logs out
    await ui.userLogsOut({ stepUser: 'Carol' })

    // And "Alice" as project manager removes their own access to the project space
    await ui.userRemovesAccessToMember({ stepUser: 'Alice', reciver: 'Alice', role: 'Can manage' })

    // Then "Alice" should not see space "team.1"
    await ui.userShouldNotSeeSpace({ stepUser: 'Alice', space: 'team.1' })

    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
