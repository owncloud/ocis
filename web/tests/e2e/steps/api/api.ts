import { config } from '../../config'
import { api } from '../../support'
import { ResourceType } from '../../support/api/share/share'
import { Space } from '../../support/types'
import fs from 'fs'
import { getWorld } from '../../environment/world'

export async function usersHaveBeenCreated({
  stepUser,
  users
}: {
  stepUser: string
  users: Array<string>
}): Promise<void> {
  const world = getWorld()
  const admin = world.usersEnvironment.getUser({ key: stepUser })
  for (const userToBeCreated of users) {
    const user = world.usersEnvironment.getUser({ key: userToBeCreated })
    // do not try to create users when using predefined users
    if (!config.predefinedUsers) {
      await api.provision.createUser({ user, admin })
    }
  }
}

export async function userHasCreatedFolder({
  stepUser,
  folderName
}: {
  stepUser: string
  folderName: string
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  await api.dav.createFolderInsidePersonalSpace({ user, folder: folderName })
}

export async function userHasCreatedFolders({
  stepUser,
  folderNames
}: {
  stepUser: string
  folderNames: string[]
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const folderName of folderNames) {
    await api.dav.createFolderInsidePersonalSpace({
      user,
      folder: folderName
    })
  }
}

export async function userHasCreatedFiles({
  stepUser,
  files
}: {
  stepUser: string
  files: { pathToFile: string; content: string; mtimeDeltaDays?: string }[]
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const file of files) {
    await api.dav.uploadFileInPersonalSpace({
      user,
      pathToFile: file.pathToFile,
      content: file.content,
      mtimeDeltaDays: file.mtimeDeltaDays
    })
  }
}

export async function userHasSharedResources({
  stepUser,
  shares
}: {
  stepUser: string
  shares: {
    resource: string
    recipient: string
    type: string
    role: string
    resourceType?: string
  }[]
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const resource of shares) {
    await api.share.createShare({
      user,
      path: resource.resource,
      shareType: resource.type,
      shareWith: resource.recipient,
      role: resource.role,
      resourceType: resource.resourceType as ResourceType
    })
  }
}

export async function userHasCreatedPublicLinkOfResource({
  stepUser,
  resource,
  role,
  name,
  password,
  space
}: {
  stepUser: string
  resource: string
  role?: string
  name?: string
  password?: string
  space?: 'Personal'
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })

  await api.share.createLinkShare({
    user,
    path: resource,
    password: password,
    name: name ? name : 'Unnamed link',
    role: role,
    spaceName: space
  })
}

export async function userHasCreatedPublicLinkOfSpace({
  stepUser,
  space,
  password,
  role,
  name
}: {
  stepUser: string
  space: string
  password: string
  role?: string
  name?: string
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })

  await api.share.createSpaceLinkShare({
    user,
    spaceName: space,
    password: password,
    name: name ? name : 'Unnamed link',
    role: role
  })
}

export async function userHasAssignedRolesToUsers({
  stepUser,
  users
}: {
  stepUser: string
  users: { id: string; role: string }[]
}) {
  const world = getWorld()
  const admin = world.usersEnvironment.getUser({ key: stepUser })
  for (const { id, role } of users) {
    const user = world.usersEnvironment.getUser({ key: id })
    /**
     The oCIS API request for assigning roles allows only one role per user,
      whereas the Keycloak API request can assign multiple roles to a user.
      If multiple roles are assigned to a user in Keycloak,
      oCIS map the highest priority role among Keycloak assigned roles.
      Therefore, we need to unassign the previous role before
      assigning a new one when using the Keycloak API.
    */
    await api.provision.unAssignRole({ admin, user })
    await api.provision.assignRole({ admin, user, role })
  }
}

export async function userHasCreatedProjectSpaces({
  stepUser,
  spaces
}: {
  stepUser: string
  spaces: Array<{ name: string; id: string }>
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const space of spaces) {
    const spaceId = await api.graph.createSpace({
      user,
      space: { id: space.id, name: space.name } as unknown as Space
    })
    world.spacesEnvironment.createSpace({
      key: space.id || space.name,
      space: { name: space.name, id: spaceId }
    })
  }
}

export async function userHasUploadedFilesInPersonalSpace({
  stepUser,
  filesToUpload
}: {
  stepUser: string
  filesToUpload: { localFile: string; to: string }[]
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const file of filesToUpload) {
    const fileInfo = world.filesEnvironment.getFile({ name: file.localFile.split('/').pop()! })
    const content = fs.readFileSync(fileInfo.path)
    await api.dav.uploadFileInPersonalSpace({
      user,
      pathToFile: file.to,
      content
    })
  }
}

export async function userHasCreatedFoldersInSpace({
  stepUser,
  spaceName,
  folders
}: {
  stepUser: string
  spaceName: string
  folders: Array<string>
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const folder of folders) {
    await api.dav.createFolderInsideSpaceBySpaceName({
      user,
      folder,
      spaceName
    })
  }
}

export async function userHasCreatedFilesInsideSpace({
  stepUser,
  files
}: {
  stepUser: string
  files: { name: string; space: string; content?: string }[]
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const file of files) {
    await api.dav.uploadFileInsideSpaceBySpaceName({
      user,
      pathToFile: file.name,
      spaceName: file.space,
      content: file.content
    })
  }
}

export async function usersHaveBeenAddedToGroup({
  stepUser,
  usersToAdd
}: {
  stepUser: string
  usersToAdd: { user: string; group: string }[]
}) {
  const world = getWorld()
  const admin = world.usersEnvironment.getUser({ key: stepUser })
  for (const info of usersToAdd) {
    const group = world.usersEnvironment.getGroup({ key: info.group })
    const user = world.usersEnvironment.getUser({ key: info.user })
    await api.graph.addUserToGroup({ user, group, admin })
  }
}

export async function userHasDeletedGroup({
  stepUser,
  name
}: {
  stepUser: string
  name: string
}): Promise<void> {
  const world = getWorld()
  const admin = world.usersEnvironment.getUser({ key: stepUser })
  const group = world.usersEnvironment.getGroup({ key: name })
  await api.graph.deleteGroup({ group, admin })
}

export async function userHasAddedMembersToSpace({
  stepUser,
  space,
  sharee
}: {
  stepUser: string
  space: string
  sharee: Array<{ user: string; shareType: string; role: string }>
}) {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const share of sharee) {
    await api.share.addMembersToTheProjectSpace({
      user,
      spaceName: space,
      shareType: share.shareType,
      shareWith: share.user,
      role: share.role
    })
  }
}

export async function groupsHaveBeenCreated({
  groupIds,
  stepUser
}: {
  groupIds: string[]
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const admin = world.usersEnvironment.getUser({ key: stepUser })
  for (const groupId of groupIds) {
    const group = world.usersEnvironment.getGroup({ key: groupId })
    await api.graph.createGroup({ group, admin })
  }
}

export async function userHasAddedTagsToResources({
  stepUser,
  tags
}: {
  stepUser: string
  tags: { resource: string; tags: string }[]
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  for (const resource of tags) {
    await api.dav.addTagToResource({ user, resource: resource.resource, tags: resource.tags })
  }
}

export async function userHasDisabledAutoAcceptingShare({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const user = world.usersEnvironment.getUser({ key: stepUser })
  await api.settings.configureAutoAcceptShare({ user, state: false })
}
