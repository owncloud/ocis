import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('space trash bin permissions', () => {
  test('viewer cannot list the trash bin of a space', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })

    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    await ui.userLogsIn({ stepUser: 'Alice' })

    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    await api.userHasCreatedFilesInsideSpace({
      stepUser: 'Alice',
      files: [{ name: 'lorem.txt', space: 'team', content: 'lorem ipsum' }]
    })

    await api.userHasDeletedResourcesInSpace({
      stepUser: 'Alice',
      spaceName: 'team',
      resources: ['lorem.txt']
    })

    await api.userHasAddedMembersToSpace({
      stepUser: 'Alice',
      space: 'team',
      sharee: [{ user: 'Brian', role: 'Can view', shareType: 'user' }]
    })

    await ui.userLogsOut({ stepUser: 'Alice' })

    await ui.userLogsIn({ stepUser: 'Brian' })

    await ui.userNavigatesToTrashbinOfSpace({ stepUser: 'Brian', space: 'team.1' })

    await ui.userShouldSeeDeletedFilesPermissionDeniedMessage({ stepUser: 'Brian' })

    await ui.userLogsOut({ stepUser: 'Brian' })
  })
})
