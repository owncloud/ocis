import join from 'join-path'
import { BodyInit, Response } from 'node-fetch'
import { request as httpRequest, checkResponseStatus } from '../http'
import { User } from '../../types'
import { TokenEnvironmentFactory } from '../../environment'
import { config } from '../../../config'

interface KeycloakToken {
  access_token: string
  refresh_token: string
}

export const realmBasePath = `admin/realms/${config.keycloakRealm}`

export const request = async (args: {
  method: 'POST' | 'DELETE' | 'PUT' | 'GET' | 'MKCOL' | 'PROPFIND' | 'PATCH'
  path: string
  body?: BodyInit
  user?: User
  header?: object
}): Promise<Response> => {
  return await httpRequest({ ...args, isKeycloakRequest: true })
}

export const getUserIdFromResponse = (response: Response): string => {
  return response.headers.get('location').split('/').pop()
}

export const refreshAccessTokenForKeycloakUser = async (user: User): Promise<void> => {
  const tokenEnvironment = TokenEnvironmentFactory('keycloak')

  const body = new URLSearchParams()
  // client-id `admin-cli` enables us to use password grant type to get access token
  body.append('client_id', 'admin-cli')
  body.append('grant_type', 'refresh_token')
  body.append('refresh_token', tokenEnvironment.getToken({ user }).refreshToken)

  const response = await request({
    method: 'POST',
    path: join('realms', 'master', 'protocol', 'openid-connect', 'token'),
    body,
    header: { 'Content-Type': 'application/x-www-form-urlencoded' },
    user
  })
  checkResponseStatus(response, 'Failed refresh access token')

  const resBody = (await response.json()) as KeycloakToken

  // update tokens
  tokenEnvironment.setToken({
    user: { ...user },
    token: {
      userId: user.id,
      accessToken: resBody.access_token,
      refreshToken: resBody.refresh_token
    }
  })
}

export const setAccessTokenForKeycloakUser = async (user: User): Promise<void> => {
  const keyCloakTokenUrl = config.keycloakUrl + '/realms/master/protocol/openid-connect/token'

  const response = await fetch(keyCloakTokenUrl, {
    method: 'POST',
    // password grant type is used to get keycloak token.
    // This approach is not recommended and used only for the test
    body: new URLSearchParams({
      client_id: 'admin-cli',
      username: config.keycloakAdminUser,
      password: config.keycloakAdminPassword,
      grant_type: 'password'
    })
  })

  const resBody = (await response.json()) as KeycloakToken
  const tokenEnvironment = TokenEnvironmentFactory('keycloak')

  tokenEnvironment.setToken({
    user: { ...user },
    token: {
      userId: user.id,
      accessToken: resBody.access_token,
      refreshToken: resBody.refresh_token
    }
  })
}

export const setupKeycloakAdminUser = (user: User) => {
  user.id = config.keycloakAdminUser
  user.password = config.keycloakAdminPassword
}
