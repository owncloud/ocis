import { checkResponseStatus, request } from '../http'
import { Group, Me, User } from '../../types'
import join from 'join-path'
import { config } from '../../../config'
import { getApplicationEntity } from './utils'
import { userRoleStore } from '../../store'
import { UsersEnvironment } from '../../environment'
import { setAccessAndRefreshToken } from '../token'

interface GroupResponse {
  value: Group[]
}

export const me = async ({ user }: { user: User }): Promise<Me> => {
  const response = await request({
    method: 'GET',
    path: join('graph', 'v1.0', 'me'),
    user
  })

  return (await response.json()) as Me
}

export const createUser = async ({ user, admin }: { user: User; admin: User }): Promise<User> => {
  const body = JSON.stringify({
    displayName: user.displayName,
    mail: user.email,
    onPremisesSamAccountName: user.id,
    passwordProfile: { password: user.password }
  })

  const response = await request({
    method: 'POST',
    path: join('graph', 'v1.0', 'users'),
    body,
    user: admin
  })

  checkResponseStatus(response, 'Failed while creating user')

  const usersEnvironment = new UsersEnvironment()
  const resBody = (await response.json()) as User
  usersEnvironment.storeCreatedUser(user.originalId || user.id, { ...user, uuid: resBody.id })
  await setAccessAndRefreshToken(user)
  return user
}

export const getMeInfo = async (user: User): Promise<User> => {
  const response = await request({
    method: 'GET',
    path: join('graph', 'v1.0', 'me'),
    user
  })
  checkResponseStatus(response, `Failed get user: ${user.id}`)

  const resBody = (await response.json()) as User
  return resBody
}

export const deleteUser = async ({ user, admin }: { user: User; admin: User }): Promise<User> => {
  const response = await request({
    method: 'DELETE',
    path: join('graph', 'v1.0', 'users', user.id),
    user: admin
  })
  // do not throw error if user is not found
  if (response.status !== 204 && response.status !== 404) {
    throw Error(`Failed to delete user: ${user.id}, Status: ${response.status}`)
  }
  const usersEnvironment = new UsersEnvironment()
  usersEnvironment.removeCreatedUser({ key: user.id })
  return user
}

export const getUserId = async ({ user, admin }: { user: User; admin: User }): Promise<string> => {
  let userId = ''
  const response = await request({
    method: 'GET',
    path: join('graph', 'v1.0', 'users', user.id),
    user: admin
  })
  if (response.ok) {
    const resBody = (await response.json()) as User
    userId = resBody.id
  }
  return userId
}

export const createGroup = async ({
  group,
  admin
}: {
  group: Group
  admin: User
}): Promise<Group> => {
  const body = JSON.stringify({
    displayName: group.displayName
  })

  const response = await request({
    method: 'POST',
    path: join('graph', 'v1.0', 'groups'),
    body,
    user: admin
  })

  checkResponseStatus(response, 'Failed while creating group')

  const usersEnvironment = new UsersEnvironment()
  const resBody = (await response.json()) as Group
  // Store with originalId for parallel safety
  usersEnvironment.storeCreatedGroup({
    group: { ...group, uuid: resBody.id, originalId: group.id }
  })
  return { ...group, uuid: resBody.id, originalId: group.id }
}

export const deleteGroup = async ({
  group,
  admin
}: {
  group: Group
  admin: User
}): Promise<Group> => {
  const usersEnvironment = new UsersEnvironment()
  const groupId = usersEnvironment.getCreatedGroup({ key: group.originalId || group.id }).uuid

  await request({
    method: 'DELETE',
    path: join('graph', 'v1.0', 'groups', groupId),
    user: admin
  })
  return group
}

export const addUserToGroup = async ({
  user,
  group,
  admin
}: {
  user: User
  group: Group
  admin: User
}): Promise<void> => {
  const usersEnvironment = new UsersEnvironment()
  const userId = usersEnvironment.getCreatedUser({ key: user.originalId || user.id }).uuid
  const groupId = usersEnvironment.getCreatedGroup({ key: group.originalId || group.id }).uuid
  const body = JSON.stringify({
    '@odata.id': join(config.baseUrl, 'graph', 'v1.0', 'users', userId)
  })

  const response = await request({
    method: 'POST',
    path: join('graph', 'v1.0', 'groups', groupId, 'members', '$ref'),
    body,
    user: admin
  })

  checkResponseStatus(response, `Failed to add user ${user.id} to group ${group.id}`)
}

export const assignRole = async (admin: User, id: string, role: string): Promise<void> => {
  const applicationEntity = await getApplicationEntity(admin)

  if (!userRoleStore.has(role)) {
    applicationEntity.appRoles.forEach((role) => {
      userRoleStore.set(role.displayName, role.id)
    })
  }
  if (userRoleStore.get(role) === undefined) {
    throw new Error(`unknown role "${role}"`)
  }
  const response = await request({
    method: 'POST',
    path: join('graph', 'v1.0', 'users', id, 'appRoleAssignments'),
    user: admin,
    body: JSON.stringify({
      principalId: id,
      appRoleId: userRoleStore.get(role),
      resourceId: applicationEntity.id
    })
  })
  checkResponseStatus(response, 'Failed while assigning role to the user')
}

export const getGroups = async (adminUser: User): Promise<Group[]> => {
  const response = await request({
    method: 'GET',
    path: join('graph', 'v1.0', 'groups'),
    user: adminUser
  })
  const data = (await response.json()) as GroupResponse
  return data.value
}
