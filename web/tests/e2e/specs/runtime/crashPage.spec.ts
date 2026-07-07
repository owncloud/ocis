import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { CRASH_CODES } from '../../../../packages/web-pkg/src/errors/codes'

test.describe('crash page', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
    await ui.userLogsIn({ stepUser: 'Alice' })
  })

  test('when spaces loading fails, the crash page is displayed', async ({ world }) => {
    const { page } = world.actorsEnvironment.getActor({ key: 'Alice' })

    await page.route('**/me/drives*', (route) => {
      route.abort('failed')
    })
    await page.reload()

    await ui.expectCrashPageToBeVisible({ page })
    await ui.userLogsOut({ stepUser: 'Alice' })
  })

  test('the crash page does not have any accessibility violations', async ({ world }) => {
    const { page } = world.actorsEnvironment.getActor({ key: 'Alice' })
    await ui.openCrashPage({
      page,
      errorCode: CRASH_CODES.RUNTIME_BOOTSTRAP_SPACES_LOAD
    })
    await ui.expectCrashPageToBeVisible({ page })
    await ui.expectCrashPageHasNoAccessibilityViolations({ page })
  })
})
