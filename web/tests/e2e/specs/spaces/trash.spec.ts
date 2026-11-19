import { test, expect } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'

test.describe('space trash permissions', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice', 'Brian'] })
    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [
        { id: 'Alice', role: 'Admin' },
        { id: 'Brian', role: 'Space Admin' }
      ]
    })
    await ui.userLogsIn({ stepUser: 'Alice' })
    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'Marketing', id: 'marketing.1' }]
    })
    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'marketing.1' })

    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [{ name: 'simple.pdf' }]
    })

    await ui.userAddsMembersToSpace({
      stepUser: 'Alice',
      members: [{ user: 'Brian', role: 'Can edit', kind: 'user' }]
    })
    await ui.userNavigatesToTrashbin({ stepUser: 'Alice' })
  })

  test('Show an error when a member cannot access the space trash', async ({ page }) => {
    await expect(page.getByTestId('files-trashbin-forbidden')).toContainText(
      'Unable to load deleted files. You may not have permission to view them.'
    )

    // And the empty-trash message should not be displayed
    await expect(page.getByText('Space has no deleted files')).not.toBeVisible()
  })

  test('Show the empty-trash message when the Manager opens an empty trash', async ({ page }) => {
    await expect(page.getByTestId('files-trashbin-empty')).toContainText(
      'Space has no deleted files'
    )
  })
})
