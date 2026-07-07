import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { fileAction } from '../../environment/constants'

test.describe('Notifications', () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    //   | Carol |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian', 'Carol'] })

    // And "Admin" assigns following roles to the users using API
    //   | id    | role        |
    //   | Alice | Space Admin |
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })
  })

  test.afterEach(async () => {
    await ui.userLogsOut({ stepUser: 'Carol' })
    await ui.userLogsOut({ stepUser: 'Brian' })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('user should be able to read and dismiss notifications', async () => {
    // Given "Admin" creates following groups using API
    //   | id    |
    //   | sales |
    await api.groupsHaveBeenCreated({ groupIds: ['sales'], stepUser: 'Admin' })

    // And "Admin" adds user to the group using API
    //   | user  | group |
    //   | Alice | sales |
    //   | Brian | sales |
    await api.usersHaveBeenAddedToGroup({
      stepUser: 'Admin',
      usersToAdd: [
        { user: 'Alice', group: 'sales' },
        { user: 'Brian', group: 'sales' }
      ]
    })

    // And "Alice" creates the following folder in personal space using API
    //   | name             |
    //   | folder_to_shared |
    //   | share_to_group   |
    await api.userHasCreatedFolders({
      stepUser: 'Alice',
      folderNames: ['folder_to_shared', 'share_to_group']
    })

    // And "Alice" creates the following project space using API
    //   | name | id     |
    //   | team | team.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // When "Alice" shares the following resource using the sidebar panel
    //   | resource         | recipient | type  | role                      | resourceType |
    //   | folder_to_shared | Brian     | user  | Can edit without versions | folder       |
    //   | share_to_group   | sales     | group | Can edit without versions | folder       |
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
        },
        {
          resource: 'share_to_group',
          recipient: 'sales',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        }
      ]
    })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // Then "Brian" should see the following notifications
    // | message                                                   |
    // | %user_alice_displayName% shared folder_to_shared with you |
    // | %user_alice_displayName% shared share_to_group with you   |
    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: [
        '%user_alice_displayName% shared folder_to_shared with you',
        '%user_alice_displayName% shared share_to_group with you'
      ]
    })

    // And "Brian" marks all notifications as read
    await ui.userMarksAllNotificationsAsRead({ stepUser: 'Brian' })

    // When "Alice" removes following sharee
    //   | resource         | recipient |
    //   | folder_to_shared | Brian     |
    await ui.userRemovesSharees({
      stepUser: 'Alice',
      sharees: [
        {
          resource: 'folder_to_shared',
          recipient: 'Brian'
        }
      ]
    })

    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    // And "Alice" adds following users to the project space
    //   | user  | role     | kind |
    //   | Brian | Can edit | user |
    //   | Carol | Can edit | user |
    await ui.userAddsMembersToSpace({
      stepUser: 'Alice',
      members: [
        { user: 'Brian', role: 'Can edit with versions and trash bin', kind: 'user' },
        { user: 'Carol', role: 'Can edit with versions and trash bin', kind: 'user' }
      ]
    })

    // Then "Alice" should see no notifications
    await ui.userShouldSeeNoNotifications({ stepUser: 'Alice' })

    // And "Brian" should see the following notifications
    //   | message                                         |
    //   | %user_alice_displayName% unshared folder_to_shared with you |
    //   | %user_alice_displayName% added you to Space team            |
    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: [
        '%user_alice_displayName% unshared folder_to_shared with you',
        '%user_alice_displayName% added you to Space team'
      ]
    })

    // And "Brian" marks all notifications as read
    await ui.userMarksAllNotificationsAsRead({ stepUser: 'Brian' })

    // When "Alice" removes access to following users from the project space
    //   | user  | role                      | kind |
    //   | Carol | Can edit without versions | user |
    await ui.userRemovesAccessToMember({
      stepUser: 'Alice',
      reciver: 'Carol',
      role: 'Can edit with trashbin'
    })

    // And "Carol" logs in
    await ui.userLogsIn({ stepUser: 'Carol' })

    // And "Carol" should see the following notifications
    //   | message                                              |
    //   | %user_alice_displayName% added you to Space team     |
    //   | %user_alice_displayName% removed you from Space team |
    await ui.userShouldSeeNotifications({
      stepUser: 'Carol',
      expectedMessages: [
        '%user_alice_displayName% added you to Space team',
        '%user_alice_displayName% removed you from Space team'
      ]
    })

    // When "Alice" opens the "admin-settings" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })

    // And "Alice" navigates to the project spaces management page
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })

    // And "Alice" disables the space "team.1" using the context-menu
    await ui.userManagesSpaceUsingContexMenu({
      stepUser: 'Alice',
      action: 'disables',
      space: 'team.1'
    })

    // Then "Brian" should see the following notifications
    //  | message                          |
    //  | %user_alice_displayName% disabled Space team |
    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: ['%user_alice_displayName% disabled Space team']
    })

    // When "Alice" deletes the space "team.1" using the context-menu
    await ui.userManagesSpaceUsingContexMenu({
      stepUser: 'Alice',
      action: 'deletes',
      space: 'team.1'
    })

    // Then "Brian" should see the following notifications
    //   | message                         |
    //   | %user_alice_displayName% deleted Space team |
    await ui.userShouldSeeNotifications({
      stepUser: 'Brian',
      expectedMessages: ['%user_alice_displayName% deleted Space team']
    })

    await api.userHasDeletedGroup({ stepUser: 'Admin', name: 'sales' })
  })

  test('user should not get any notification when notification is disabled', async () => {
    // Given "Alice" creates the following folder in personal space using API
    //   | name             |
    //   | folder_to_shared |
    await api.userHasCreatedFolder({ stepUser: 'Alice', folderName: 'folder_to_shared' })

    // And "Alice" creates the following project space using API
    //   | name | id     |
    //   | team | team.1 |
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    // And "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" opens the user menu
    await ui.userOpensAccountPage({ stepUser: 'Brian' })

    // When "Brian" disables notification for the following events
    //   | event                   |
    //   | Share Received          |
    //   | Share Removed           |
    //   | Added as space member   |
    //   | Removed as space member |
    await ui.userDisablesNotificationEvents({
      stepUser: 'Brian',
      events: [
        'Share Received',
        'Share Removed',
        'Added as space member',
        'Removed as space member'
      ]
    })

    // And "Carol" logs in
    await ui.userLogsIn({ stepUser: 'Carol' })

    // And "Carol" opens the user menu
    await ui.userOpensAccountPage({ stepUser: 'Carol' })

    // And "Carol" disables notification for the following events
    //   | event          |
    //   | Space disabled |
    //   | Space deleted  |
    await ui.userDisablesNotificationEvents({
      stepUser: 'Brian',
      events: ['Space disabled', 'Space deleted']
    })

    // And "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" shares the following resource using the sidebar panel
    //   | resource         | recipient | type  | role                      | resourceType |
    //   | folder_to_shared | Brian     | user  | Can edit without versions | folder       |
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

    // When "Alice" removes following sharee
    //   | resource         | recipient |
    //   | folder_to_shared | Brian     |
    await ui.userRemovesSharees({
      stepUser: 'Alice',
      sharees: [
        {
          resource: 'folder_to_shared',
          recipient: 'Brian'
        }
      ]
    })

    // And "Alice" navigates to the project space "team.1"
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    // And "Alice" adds following users to the project space
    //   | user  | role     | kind |
    //   | Brian | Can edit | user |
    //   | Carol | Can edit | user |
    await ui.userAddsUsersToProjectSpace({
      stepUser: 'Alice',
      space: 'team.1',
      members: [
        { reciver: 'Brian', role: 'Can edit with versions and trash bin', kind: 'user' },
        { reciver: 'Carol', role: 'Can edit with versions and trash bin', kind: 'user' }
      ]
    })

    // And "Alice" removes access to following users from the project space
    //   | user  | role     | kind |
    //   | Brian | Can edit | user |
    await ui.userRemovesAccessToMember({
      stepUser: 'Alice',
      reciver: 'Carol',
      role: 'Can edit with versions and trash bin'
    })

    // Then "Alice" should see no notifications
    await ui.userShouldSeeNoNotifications({ stepUser: 'Brian' })

    // When "Alice" opens the "admin-settings" app
    await ui.userOpensApplication({ stepUser: 'Alice', name: 'admin-settings' })

    // And "Alice" navigates to the project spaces management page
    await ui.userNavigatesToProjectSpaceManagementPage({ stepUser: 'Alice' })

    // And "Alice" disables the space "team.1" using the context-menu
    await ui.userManagesSpaceUsingContexMenu({
      stepUser: 'Alice',
      action: 'disables',
      space: 'team.1'
    })

    // When "Alice" deletes the space "team.1" using the context-menu
    await ui.userManagesSpaceUsingContexMenu({
      stepUser: 'Alice',
      action: 'deletes',
      space: 'team.1'
    })

    // Then "Carol" should see no notifications
    await ui.userShouldSeeNoNotifications({ stepUser: 'Carol' })
  })
})
