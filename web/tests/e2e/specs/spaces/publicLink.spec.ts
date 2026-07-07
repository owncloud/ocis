import { test } from '../../environment/test'
import * as api from '../../steps/api/api'
import * as ui from '../../steps/ui/index'
import { application, fileAction } from '../../environment/constants'

test.describe('spaces public link', () => {
  test.beforeEach(async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Alice'] })
  })

  test('public link for space', async () => {
    await api.usersHaveBeenCreated({ stepUser: 'Admin', users: ['Brian', 'Carol', 'David'] })

    await api.userHasAssignedRolesToUsers({
      stepUser: 'Admin',
      users: [{ id: 'Alice', role: 'Space Admin' }]
    })

    await ui.userLogsIn({ stepUser: 'Alice' })

    await api.userHasCreatedProjectSpaces({
      stepUser: 'Alice',
      spaces: [{ name: 'team', id: 'team.1' }]
    })

    await api.userHasCreatedFoldersInSpace({
      stepUser: 'Alice',
      spaceName: 'team',
      folders: ['spaceFolder/subFolder']
    })

    await api.userHasCreatedFilesInsideSpace({
      stepUser: 'Alice',
      files: [
        { name: 'spaceFolder/shareToBrian.txt', space: 'team', content: 'some text' },
        { name: 'spaceFolder/subFolder/shareToBrian.md', space: 'team', content: 'readme' }
      ]
    })

    await api.userHasAddedMembersToSpace({
      stepUser: 'Alice',
      space: 'team',
      sharee: [
        { user: 'Brian', role: 'Can edit with versions and trash bin', shareType: 'user' },
        { user: 'Carol', role: 'Can view', shareType: 'user' },
        { user: 'David', role: 'Can manage', shareType: 'user' }
      ]
    })

    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    await ui.userUploadsResourcesViaDragNDrop({
      stepUser: 'Alice',
      resourceNames: ['simple.pdf', 'testavatar.jpg']
    })

    await ui.userCreatesPublicLinkOfSpaceWithPassword({ stepUser: 'Alice', password: '%public%' })

    await ui.userRenamesTheMostRecentlyCreatedPublicLinkOfSpace({
      stepUser: 'Alice',
      newName: 'spaceLink'
    })

    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'spaceFolder',
      password: '%public%'
    })

    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'spaceFolder',
      newName: 'folderLink'
    })

    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'spaceFolder/shareToBrian.txt',
      password: '%public%'
    })

    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'spaceFolder/shareToBrian.txt',
      newName: 'textLink'
    })

    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'spaceFolder/subFolder/shareToBrian.md',
      password: '%public%'
    })

    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'spaceFolder/subFolder/shareToBrian.md',
      newName: 'markdownLink'
    })

    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'simple.pdf',
      password: '%public%'
    })

    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'simple.pdf',
      newName: 'pdfLink'
    })

    await ui.userCreatesPublicLink({
      stepUser: 'Alice',
      resource: 'testavatar.jpg',
      password: '%public%'
    })

    await ui.userRenamesMostRecentlyCreatedPublicLinkOfResource({
      stepUser: 'Alice',
      resource: 'testavatar.jpg',
      newName: 'imageLink'
    })

    await ui.userLogsOut({ stepUser: 'Alice' })

    await ui.userLogsIn({ stepUser: 'Brian' })

    await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'spaceLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Brian', password: '%public%' })

    await ui.userShouldNotBeAbleToEditThePublicLink({ stepUser: 'Brian', linkName: 'spaceLink' })

    await ui.userShouldNotBeAbleToEditThePublicLink({ stepUser: 'Brian', linkName: 'folderLink' })

    await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'textLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Brian', password: '%public%' })

    await ui.userShouldBeInFileViewer({ stepUser: 'Brian', fileViewerType: application.textEditor })

    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    await ui.userOpensPublicLink({ stepUser: 'Brian', name: 'markdownLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Brian', password: '%public%' })

    await ui.userShouldBeInFileViewer({ stepUser: 'Brian', fileViewerType: application.textEditor })

    await ui.userClosesFileViewer({ stepUser: 'Brian' })

    await ui.userLogsOut({ stepUser: 'Brian' })

    await ui.userLogsIn({ stepUser: 'Carol' })

    await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'spaceLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Carol', password: '%public%' })

    await ui.userShouldNotBeAbleToEditThePublicLink({ stepUser: 'Carol', linkName: 'spaceLink' })

    await ui.userShouldNotBeAbleToEditThePublicLink({ stepUser: 'Carol', linkName: 'folderLink' })

    await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'folderLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Carol', password: '%public%' })

    await ui.userShouldNotBeAbleToEditResource({ stepUser: 'Carol', resource: 'subFolder' })

    await ui.userOpensPublicLink({ stepUser: 'Carol', name: 'pdfLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Carol', password: '%public%' })

    await ui.userShouldBeInFileViewer({ stepUser: 'Carol', fileViewerType: application.pdfViewer })

    await ui.userClosesFileViewer({ stepUser: 'Carol' })

    await ui.userLogsOut({ stepUser: 'Carol' })

    await ui.userLogsIn({ stepUser: 'David' })

    await ui.userOpensPublicLink({ stepUser: 'David', name: 'spaceLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'David', password: '%public%' })

    await ui.userEditsThePublicLinkOfSpaceChangingRole({
      stepUser: 'David',
      linkName: 'spaceLink',
      role: 'Can edit'
    })

    await ui.userChangesRoleOfPublicLinkOfResource({
      stepUser: 'David',
      resource: 'spaceFolder',
      linkName: 'folderLink',
      newRole: 'Can edit'
    })

    await ui.userOpensPublicLink({ stepUser: 'David', name: 'imageLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'David', password: '%public%' })

    await ui.userShouldBeInFileViewer({
      stepUser: 'David',
      fileViewerType: application.mediaViewer
    })

    await ui.userClosesFileViewer({ stepUser: 'David' })

    await ui.userLogsOut({ stepUser: 'David' })
  })

  test('crud operation to public link for space', async () => {
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
      files: [{ name: 'example.txt', space: 'team', content: 'some text' }]
    })

    await ui.userNavigatesToSpace({ stepUser: 'Alice', space: 'team.1' })

    await ui.userCreatesPublicLinkOfSpaceWithPassword({ stepUser: 'Alice', password: '%public%' })

    await ui.userRenamesTheMostRecentlyCreatedPublicLinkOfSpace({
      stepUser: 'Alice',
      newName: 'spaceLink'
    })

    await ui.userEditsThePublicLinkOfSpaceChangingRole({
      stepUser: 'Alice',
      linkName: 'spaceLink',
      role: 'Can edit'
    })

    await ui.userLogsOut({ stepUser: 'Alice' })

    await ui.userOpensPublicLink({ stepUser: 'Anonymous', name: 'spaceLink' })

    await ui.userUnlocksPublicLink({ stepUser: 'Anonymous', password: '%public%' })

    await ui.userDownloadsThePublicLinkResources({
      stepUser: 'Anonymous',
      actionType: fileAction.sideBarPanel,
      resources: [{ resource: 'example.txt', type: 'file' }]
    })

    await ui.userUploadsResourcesInPublicLink({
      stepUser: 'Anonymous',
      resources: [{ name: 'new-lorem.txt' }]
    })

    await ui.userRenamesPublicLinkResources({
      stepUser: 'Anonymous',
      resources: [{ resource: 'example.txt', newName: 'renamed.txt' }]
    })

    await ui.userEditsResources({
      stepUser: 'Anonymous',
      resources: [{ name: 'renamed.txt', content: 'new content' }]
    })

    await ui.userDeletesResources({
      stepUser: 'Anonymous',
      actionType: fileAction.sideBarPanel,
      resources: [{ name: 'renamed.txt' }, { name: 'new-lorem.txt' }]
    })
  })
})
