import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import {
  createResourceTypes,
  verifyTopbarResourceName
} from '../../support/objects/app-files/resource/actions'

test.describe('Open resources in a new tab', { tag: '@predefined-users' }, () => {
  test('opens files with # and ? in their name in a new tab', async () => {
    const files: { name: string; type: createResourceTypes; content: string }[] = [
      { name: '#.txt', type: 'txtFile', content: 'file with # in its name' },
      { name: '?.txt', type: 'txtFile', content: 'file with ? in its name' },
      { name: 'a#b.txt', type: 'txtFile', content: 'file with a#b in its name' },
      { name: 'a?b.txt', type: 'txtFile', content: 'file with a?b in its name' },
      { name: 'a#b?c.txt', type: 'txtFile', content: 'file with both # and ? in its name' }
    ]

    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    await ui.userLogsIn({ stepUser: 'Alice' })

    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: files
    })

    for (const file of files) {
      const newPage = await ui.userOpensResourceInNewTab({
        stepUser: 'Alice',
        resource: file.name
      })

      await verifyTopbarResourceName({
        page: newPage,
        expectedName: file.name
      })

      await newPage.close()
    }

    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
