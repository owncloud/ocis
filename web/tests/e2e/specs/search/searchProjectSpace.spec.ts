import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { resourcePage, searchScope } from '../../environment/constants'

test.describe('Search in the project space', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })

    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    await ui.userLogsIn({ stepUser: 'Alice' })

    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    await ui.userCreatesResources({
      stepUser: 'Alice',
      resources: [{ name: 'folder(WithSymbols:!;_+-&)', type: 'folder' }]
    })

    await ui.userUploadsResources({
      stepUser: 'Alice',
      resources: [{ name: "new-'single'quotes.txt", to: 'folder(WithSymbols:!;_+-&)' }]
    })

    await ui.userNavigatesToPersonalSpacePage({ stepUser: 'Alice' })
  })

  test('Search in the project spaces', async () => {
    // search for project space objects
    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: "-'s",
      filter: searchScope.allFiles
    })

    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ["new-'single'quotes.txt"]
    })

    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder(WithSymbols:!;_+-&)']
    })

    await ui.userSearchesGloballyWithFilter({
      stepUser: 'Alice',
      keyword: '!;_+-&)',
      filter: searchScope.allFiles
    })

    await ui.userShouldSeeResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ['folder(WithSymbols:!;_+-&)']
    })

    await ui.userShouldNotSeeTheResources({
      listType: resourcePage.searchList,
      stepUser: 'Alice',
      resources: ["new-'single'quotes.txt"]
    })

    await ui.userLogsOut({ stepUser: 'Alice' })
  })
})
