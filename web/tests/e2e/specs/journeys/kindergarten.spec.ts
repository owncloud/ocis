import { test } from '../../environment/test'
import { fileAction } from '../../environment/constants'
import * as api from '../../steps/api/api.js'
import * as ui from '../../steps/ui/index'

test.describe('Kindergarten can use web to organize a day', () => {
  test.beforeEach(async () => {
    // Given "Admin" creates following users using API
    //   | id    |
    //   | Alice |
    //   | Brian |
    //   | Carol |
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian', 'Carol'] })

    // And "Admin" creates following group using API
    //   | id       |
    //   | sales    |
    //   | security |
    await api.groupsHaveBeenCreated({ groupIds: ['sales', 'security'], stepUser: 'Admin' })

    // And "Admin" adds user to the group using API
    //   | user  | group |
    //   | Brian | sales |
    await api.usersHaveBeenAddedToGroup({
      stepUser: 'Admin',
      usersToAdd: [{ user: 'Brian', group: 'sales' }]
    })
  })

  test('Alice can share this weeks meal plan with all parents', async () => {
    // This journey performs 16 downloads plus shares and deletes, each cycling the
    // sidebar with a11y scans - it consistently runs beyond the default 180s budget.
    test.slow()

    // When "Alice" logs in
    await ui.userLogsIn({ stepUser: 'Alice' })

    // And "Alice" navigates to the personal space page
    await ui.userNavigatesToPersonalSpacePage({ stepUser: 'Alice' })

    // And "Alice" creates the following resources
    //   | resource                             | type   |
    //   | groups/Kindergarten Koalas/meal plan | folder |
    //   | groups/Pre-Schools Pirates/meal plan | folder |
    //   | groups/Teddy Bear Daycare/meal plan  | folder |
    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [
        { name: 'groups/Kindergarten Koalas/meal plan', type: 'folder' },
        { name: 'groups/Pre-Schools Pirates/meal plan', type: 'folder' },
        { name: 'groups/Teddy Bear Daycare/meal plan', type: 'folder' }
      ]
    })

    // And "Alice" uploads the following resources
    //   | resource          | to                                   |
    //   | PARENT/parent.txt | groups/Kindergarten Koalas/meal plan |
    //   | lorem.txt         | groups/Kindergarten Koalas/meal plan |
    //   | lorem-big.txt     | groups/Kindergarten Koalas/meal plan |
    //   | data.zip          | groups/Pre-Schools Pirates/meal plan |
    //   | lorem.txt         | groups/Pre-Schools Pirates/meal plan |
    //   | lorem-big.txt     | groups/Pre-Schools Pirates/meal plan |
    //   | data.zip          | groups/Teddy Bear Daycare/meal plan  |
    //   | lorem.txt         | groups/Teddy Bear Daycare/meal plan  |
    //   | lorem-big.txt     | groups/Teddy Bear Daycare/meal plan  |
    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [
        { name: 'PARENT/parent.txt', to: 'groups/Kindergarten Koalas/meal plan' },
        { name: 'lorem.txt', to: 'groups/Kindergarten Koalas/meal plan' },
        { name: 'lorem-big.txt', to: 'groups/Kindergarten Koalas/meal plan' },
        { name: 'data.zip', to: 'groups/Pre-Schools Pirates/meal plan' },
        { name: 'lorem.txt', to: 'groups/Pre-Schools Pirates/meal plan' },
        { name: 'lorem-big.txt', to: 'groups/Pre-Schools Pirates/meal plan' },
        { name: 'data.zip', to: 'groups/Teddy Bear Daycare/meal plan' },
        { name: 'lorem.txt', to: 'groups/Teddy Bear Daycare/meal plan' },
        { name: 'lorem-big.txt', to: 'groups/Teddy Bear Daycare/meal plan' }
      ]
    })

    // Implementation of sharing with different roles is currently broken
    // since we switched to bulk creating of shares with a single dropdown

    // And "Alice" shares the following resources using the sidebar panel
    //   | resource                                           | recipient | type  | role                   | resourceType |
    //   | groups/Pre-Schools Pirates/meal plan               | Brian     | user  | Can edit with trashbin | folder       |
    //   | groups/Pre-Schools Pirates/meal plan               | Carol     | user  | Can edit with trashbin | folder       |
    //   | groups/Pre-Schools Pirates/meal plan/lorem-big.txt | sales     | group | Can view               | file         |
    //   | groups/Pre-Schools Pirates/meal plan/lorem-big.txt | Carol     | user  | Can view               | file         |
    //   | groups/Kindergarten Koalas/meal plan               | sales     | group | Can view               | folder       |
    //   | groups/Kindergarten Koalas/meal plan               | security  | group | Can edit with trashbin | folder       |
    //   | groups/Kindergarten Koalas/meal plan/lorem.txt     | sales     | group | Can view               | file         |
    //   | groups/Kindergarten Koalas/meal plan/lorem.txt     | security  | group | Can view               | file         |
    //   | groups/Teddy Bear Daycare/meal plan                | Brian     | user  | Can edit with trashbin | folder       |
    //   | groups/Teddy Bear Daycare/meal plan                | Carol     | user  | Can edit with trashbin | folder       |
    //   | groups/Teddy Bear Daycare/meal plan/data.zip       | Brian     | user  | Can edit with trashbin | file         |
    //   | groups/Teddy Bear Daycare/meal plan/data.zip       | Carol     | user  | Can edit with trashbin | file         |
    await ui.userSharesResources({
      actionType: fileAction.sideBarPanel,
      stepUser: 'Alice',
      shares: [
        {
          resource: 'groups/Pre-Schools Pirates/meal plan',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Pre-Schools Pirates/meal plan',
          recipient: 'Carol',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Pre-Schools Pirates/meal plan/lorem-big.txt',
          recipient: 'sales',
          type: 'group',
          role: 'Can view',
          resourceType: 'file'
        },
        {
          resource: 'groups/Pre-Schools Pirates/meal plan/lorem-big.txt',
          recipient: 'Carol',
          type: 'user',
          role: 'Can view',
          resourceType: 'file'
        },
        {
          resource: 'groups/Kindergarten Koalas/meal plan',
          recipient: 'sales',
          type: 'group',
          role: 'Can view',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Kindergarten Koalas/meal plan',
          recipient: 'security',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Kindergarten Koalas/meal plan/lorem.txt',
          recipient: 'sales',
          type: 'group',
          role: 'Can view',
          resourceType: 'file'
        },
        {
          resource: 'groups/Kindergarten Koalas/meal plan/lorem.txt',
          recipient: 'security',
          type: 'group',
          role: 'Can view',
          resourceType: 'file'
        },
        {
          resource: 'groups/Teddy Bear Daycare/meal plan',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Teddy Bear Daycare/meal plan',
          recipient: 'Carol',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Teddy Bear Daycare/meal plan/data.zip',
          recipient: 'Brian',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        },
        {
          resource: 'groups/Teddy Bear Daycare/meal plan/data.zip',
          recipient: 'Carol',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        }
      ]
    })

    // update share
    // And "Alice" updates following sharee role
    //   | resource                                           | recipient | type  | role                   | resourceType |
    //   | groups/Pre-Schools Pirates/meal plan               | Carol     | user  | Can view               | folder       |
    //   | groups/Pre-Schools Pirates/meal plan/lorem-big.txt | sales     | group | Can edit with trashbin | file         |
    //   | groups/Kindergarten Koalas/meal plan               | sales     | group | Can edit with trashbin | folder       |
    //   | groups/Teddy Bear Daycare/meal plan/data.zip       | Carol     | user  | Can edit with trashbin | file         |
    await ui.userUpdatesShareeRoles({
      stepUser: 'Alice',
      roleUpdates: [
        {
          resource: 'groups/Pre-Schools Pirates/meal plan',
          recipient: 'Carol',
          type: 'user',
          role: 'Can view',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Pre-Schools Pirates/meal plan/lorem-big.txt',
          recipient: 'sales',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        },
        {
          resource: 'groups/Kindergarten Koalas/meal plan',
          recipient: 'sales',
          type: 'group',
          role: 'Can edit with trashbin',
          resourceType: 'folder'
        },
        {
          resource: 'groups/Teddy Bear Daycare/meal plan/data.zip',
          recipient: 'Carol',
          type: 'user',
          role: 'Can edit with trashbin',
          resourceType: 'file'
        }
      ]
    })
    // Then what do we check for to be confident that the above things done by Alice have worked?
    // When "Brian" logs in
    await ui.userLogsIn({ stepUser: 'Brian' })

    // And "Brian" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Brian' })

    // And "Brian" downloads the following resources using the sidebar panel
    //   | resource | from      | type |
    //   | data.zip | meal plan | file |
    await ui.userDownloadsResource({
      stepUser: 'Brian',
      resourceToDownload: [{ resource: 'data.zip', from: 'meal plan', type: 'file' }],
      actionType: fileAction.sideBarPanel
    })

    // Then what do we check for to be confident that the above things done by Brian have worked?
    // Then the downloaded zip should contain... ?
    // When "Carol" logs in
    await ui.userLogsIn({ stepUser: 'Carol' })

    // And "Carol" navigates to the shared with me page
    await ui.userNavigatesToSharedWithMePage({ stepUser: 'Carol' })
    // And "Carol" downloads the following resources using the sidebar panel
    //   | resource      | from      | type   |
    //   | data.zip      | meal plan | file   |
    //   | lorem.txt     | meal plan | file   |
    //   | lorem-big.txt | meal plan | file   |
    //   | meal plan     |           | folder |
    // Then what do we check for to be confident that the above things done by Carol have worked?
    // Then the downloaded files should have content "abc..."
    await ui.userDownloadsResource({
      stepUser: 'Carol',
      resourceToDownload: [
        { resource: 'data.zip', from: 'meal plan', type: 'file' },
        { resource: 'lorem.txt', from: 'meal plan', type: 'file' },
        { resource: 'lorem-big.txt', from: 'meal plan', type: 'file' },
        { resource: 'meal plan', type: 'folder' }
      ],
      actionType: fileAction.sideBarPanel
    })
    // And "Carol" logs out
    await ui.userLogsOut({ stepUser: 'Carol' })

    // When "Brian" downloads the following resources using the sidebar panel
    //   | resource      | from      | type   |
    //   | lorem.txt     | meal plan | file   |
    //   | lorem-big.txt | meal plan | file   |
    //   | meal plan     |           | folder |
    // Then what do we check for to be confident that the above things done by Brian have worked?
    // Then the downloaded files should have content "abc..."
    await ui.userDownloadsResource({
      stepUser: 'Brian',
      resourceToDownload: [
        { resource: 'lorem.txt', from: 'meal plan', type: 'file' },
        { resource: 'lorem-big.txt', from: 'meal plan', type: 'file' },
        { resource: 'meal plan', type: 'folder' }
      ],
      actionType: fileAction.sideBarPanel
    })

    // And "Brian" logs out
    await ui.userLogsOut({ stepUser: 'Brian' })
    // And "Alice" downloads the following resources using the sidebar panel
    //   | resource            | from                                 | type   |
    //   | parent.txt          | groups/Kindergarten Koalas/meal plan | file   |
    //   | lorem.txt           | groups/Kindergarten Koalas/meal plan | file   |
    //   | lorem-big.txt       | groups/Kindergarten Koalas/meal plan | file   |
    //   | data.zip            | groups/Pre-Schools Pirates/meal plan | file   |
    //   | lorem.txt           | groups/Pre-Schools Pirates/meal plan | file   |
    //   | lorem-big.txt       | groups/Pre-Schools Pirates/meal plan | file   |
    //   | data.zip            | groups/Teddy Bear Daycare/meal plan  | file   |
    //   | lorem.txt           | groups/Teddy Bear Daycare/meal plan  | file   |
    //   | lorem-big.txt       | groups/Teddy Bear Daycare/meal plan  | file   |
    //   | meal plan           | groups/Kindergarten Koalas           | folder |
    //   | meal plan           | groups/Pre-Schools Pirates           | folder |
    //   | meal plan           | groups/Teddy Bear Daycare            | folder |
    //   | Kindergarten Koalas | groups                               | folder |
    //   | Pre-Schools Pirates | groups                               | folder |
    //   | Teddy Bear Daycare  | groups                               | folder |
    //   | groups              |                                      | folder |
    await ui.userDownloadsResource({
      stepUser: 'Alice',
      resourceToDownload: [
        {
          resource: 'parent.txt',
          from: 'groups/Kindergarten Koalas/meal plan',
          type: 'file'
        },
        {
          resource: 'lorem.txt',
          from: 'groups/Kindergarten Koalas/meal plan',
          type: 'file'
        },
        {
          resource: 'lorem-big.txt',
          from: 'groups/Kindergarten Koalas/meal plan',
          type: 'file'
        },
        {
          resource: 'data.zip',
          from: 'groups/Pre-Schools Pirates/meal plan',
          type: 'file'
        },
        {
          resource: 'lorem.txt',
          from: 'groups/Pre-Schools Pirates/meal plan',
          type: 'file'
        },
        {
          resource: 'lorem-big.txt',
          from: 'groups/Pre-Schools Pirates/meal plan',
          type: 'file'
        },
        {
          resource: 'data.zip',
          from: 'groups/Teddy Bear Daycare/meal plan',
          type: 'file'
        },
        {
          resource: 'lorem.txt',
          from: 'groups/Teddy Bear Daycare/meal plan',
          type: 'file'
        },
        {
          resource: 'lorem-big.txt',
          from: 'groups/Teddy Bear Daycare/meal plan',
          type: 'file'
        },
        {
          resource: 'meal plan',
          from: 'groups/Kindergarten Koalas',
          type: 'folder'
        },
        {
          resource: 'meal plan',
          from: 'groups/Pre-Schools Pirates',
          type: 'folder'
        },
        {
          resource: 'meal plan',
          from: 'groups/Teddy Bear Daycare',
          type: 'folder'
        },
        {
          resource: 'Kindergarten Koalas',
          from: 'groups',
          type: 'folder'
        },
        {
          resource: 'Pre-Schools Pirates',
          from: 'groups',
          type: 'folder'
        },
        {
          resource: 'Teddy Bear Daycare',
          from: 'groups',
          type: 'folder'
        },
        {
          resource: 'groups',
          type: 'folder'
        }
      ],
      actionType: fileAction.sideBarPanel
    })
    // And "Alice" deletes the following resources using the batch action
    //   | resource            | from                                 |
    //   | lorem.txt           | groups/Kindergarten Koalas/meal plan |
    //   | lorem-big.txt       | groups/Kindergarten Koalas/meal plan |
    //   | data.zip            | groups/Pre-Schools Pirates/meal plan |
    //   | lorem.txt           | groups/Pre-Schools Pirates/meal plan |
    //   | lorem-big.txt       | groups/Pre-Schools Pirates/meal plan |
    //   | data.zip            | groups/Teddy Bear Daycare/meal plan  |
    //   | lorem.txt           | groups/Teddy Bear Daycare/meal plan  |
    //   | lorem-big.txt       | groups/Teddy Bear Daycare/meal plan  |
    //   | Kindergarten Koalas | groups                               |
    //   | Pre-Schools Pirates | groups                               |
    //   | Teddy Bear Daycare  | groups                               |
    // # Then what do we check for to be confident that the above things done by Alice have worked?
    // # Then the downloaded files should have content "abc..."
    await ui.userDeletesResources({
      stepUser: 'Alice',
      actionType: fileAction.sideBarPanel,
      resources: [
        { name: 'lorem.txt', from: 'groups/Kindergarten Koalas/meal plan' },
        { name: 'lorem-big.txt', from: 'groups/Kindergarten Koalas/meal plan' },
        { name: 'data.zip', from: 'groups/Pre-Schools Pirates/meal plan' },
        { name: 'lorem.txt', from: 'groups/Pre-Schools Pirates/meal plan' },
        { name: 'lorem-big.txt', from: 'groups/Pre-Schools Pirates/meal plan' },
        { name: 'data.zip', from: 'groups/Teddy Bear Daycare/meal plan' },
        { name: 'lorem.txt', from: 'groups/Teddy Bear Daycare/meal plan' },
        { name: 'lorem-big.txt', from: 'groups/Teddy Bear Daycare/meal plan' },
        { name: 'Kindergarten Koalas', from: 'groups' },
        { name: 'Pre-Schools Pirates', from: 'groups' },
        { name: 'Teddy Bear Daycare', from: 'groups' }
      ]
    })
    // And "Alice" logs out
    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
