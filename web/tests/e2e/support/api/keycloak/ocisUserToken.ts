import fetch from 'node-fetch'
import { TokenEnvironmentFactory } from '../../environment'
import { config } from '../../../config'
import { User } from '../../types'

interface ocisTokenForKeycloak {
  access_token: string
  refresh_token: string
}

const authorizationEndpoint = config.keycloakUrl + '/realms/oCIS/protocol/openid-connect/auth'
const tokenEndpoint = config.keycloakUrl + '/realms/oCIS/protocol/openid-connect/token'
const redirectUrl = config.baseUrl + '/oidc-callback.html'

async function getAuthorizationEndPoint() {
  const loginParams = {
    client_id: 'web',
    redirect_uri: redirectUrl,
    response_mode: 'query',
    response_type: 'code',
    scope: 'openid profile email'
  }
  const queryString = new URLSearchParams(loginParams).toString()
  const authorizationUrl = `${authorizationEndpoint}?${queryString}`

  const authorizationResponse = await fetch(authorizationUrl, {
    method: 'GET',
    redirect: 'manual'
  })

  if (authorizationResponse.status === 302) {
    const locationHeader = authorizationResponse.headers.get('location')
    const urlParams = new URLSearchParams(new URL(locationHeader).search)
    const errorDescription = urlParams.get('error_description')
    throw new Error(`Unexpected redirection. ${errorDescription}`)
  } else if (authorizationResponse.status !== 200) {
    throw new Error(
      `Authorization failed: Expected status code to be 200 but received ${authorizationResponse.status}. \nMessage: ${authorizationResponse.statusText}`
    )
  }

  const cookies = authorizationResponse.headers.raw()['set-cookie']?.[0]
  const htmlData = await authorizationResponse.text()

  // authorization url for login is send back from server in the HTML body.
  const match = htmlData.match(/action="([^"]+)"/)
  if (!match) {
    throw new Error('No authorization url found in the HTML response body.')
  }
  const auhorizationUrl = match[1]
  return [auhorizationUrl, cookies]
}

const getCode = async ({
  user,
  auhorizationUrl,
  cookies
}: {
  user: User
  auhorizationUrl: string
  cookies: string
}) => {
  const authCodeResponse = await fetch(auhorizationUrl, {
    method: 'POST',
    body: new URLSearchParams({
      username: user.id,
      password: user.password
    }),
    redirect: 'manual',
    headers: {
      Cookie: cookies
    }
  })

  if (authCodeResponse.status !== 302) {
    throw new Error(
      `Login failed: Expected status code to be 302 but received ${authCodeResponse.status}. \nMessage: ${authCodeResponse.statusText}`
    )
  }

  const locationHeader = authCodeResponse.headers.get('location')
  const urlParams = new URLSearchParams(new URL(locationHeader).search)
  return urlParams.get('code')
}

const getToken = async (authorizationCode: string) => {
  const tokenResponse = await fetch(tokenEndpoint, {
    method: 'POST',
    body: new URLSearchParams({
      client_id: 'web',
      code: authorizationCode,
      redirect_uri: redirectUrl,
      grant_type: 'authorization_code'
    })
  })

  if (tokenResponse.status !== 200) {
    throw new Error(
      `Failed to retrieve token: Expected status code to be 200 but received ${tokenResponse.status}. \nMessage: ${tokenResponse.statusText}`
    )
  }

  return tokenResponse
}

export const setAccessTokenForKeycloakOcisUser = async (user: User) => {
  const tokenEnvironment = TokenEnvironmentFactory()

  // admin's Keycloak username is 'admin', not the world-transformed user.id
  const loginId =
    (user.originalId || user.id).toLowerCase() === config.keycloakAdminUser.toLowerCase()
      ? user.originalId
      : user.id

  // Retry OIDC auth code flow: admin user is shared across all workers, so
  // multiple workers may authenticate it concurrently. Keycloak can return the
  // login form (200) instead of a redirect (302) under concurrent load. Each
  // retry starts a fresh OIDC session, naturally staggering the requests.
  const MAX_RETRIES = 3
  for (let attempt = 0; attempt < MAX_RETRIES; attempt++) {
    try {
      const [auhorizationUrl, cookies] = await getAuthorizationEndPoint()
      const authorizationCode = await getCode({
        user: { ...user, id: loginId },
        auhorizationUrl,
        cookies
      })
      const tokenResponse = await getToken(authorizationCode)
      const token = (await tokenResponse.json()) as ocisTokenForKeycloak
      tokenEnvironment.setToken({
        user: { ...user },
        token: {
          userId: user.id,
          accessToken: token.access_token,
          refreshToken: token.refresh_token
        }
      })
      return
    } catch (e) {
      if (attempt === MAX_RETRIES - 1) {
        throw e
      }
      await new Promise((r) => setTimeout(r, config.minTimeout * 1000))
    }
  }
}

export const refreshAccessTokenForKeycloakOcisUser = async (user: User) => {
  const tokenEnvironment = TokenEnvironmentFactory()
  const refreshToken = tokenEnvironment.getToken({ user }).refreshToken
  const tokenResponse = await fetch(tokenEndpoint, {
    method: 'POST',
    body: new URLSearchParams({
      client_id: 'web',
      refresh_token: refreshToken,
      grant_type: 'refresh_token'
    })
  })
  if (tokenResponse.status !== 200) {
    throw new Error(
      `Failed to retrieve token: Expected status code to be 200 but received ${tokenResponse.status}. \nMessage: ${tokenResponse.statusText}`
    )
  }

  const token = (await tokenResponse.json()) as ocisTokenForKeycloak
  tokenEnvironment.setToken({
    user: { ...user },
    token: {
      userId: user.id,
      accessToken: token.access_token,
      refreshToken: token.refresh_token
    }
  })
}
