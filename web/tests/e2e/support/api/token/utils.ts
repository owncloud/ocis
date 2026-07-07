import fetch, { Response } from 'node-fetch'
import { checkResponseStatus, request } from '../http'
import { TokenEnvironmentFactory } from '../../environment'
import { config } from '../../../config'
import { User } from '../../types'

const logonUrl = '/signin/v1/identifier/_/logon'
const redirectUrl = '/oidc-callback.html'
const tokenUrl = '/konnect/v1/token'

interface Token {
  access_token: string
  refresh_token: string
}

const logonRequest = (username: string, password: string): Promise<Response> => {
  return fetch(config.baseUrl + logonUrl, {
    method: 'POST',
    headers: {
      'Kopano-Konnect-XSRF': '1',
      Referer: config.baseUrl,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      params: [username, password, '1'],
      hello: {
        scope: 'openid profile email',
        client_id: 'web',
        redirect_uri: config.baseUrl + redirectUrl,
        flow: 'oidc'
      }
    })
  })
}

const getAuthorizedEndPoint = async (user: User): Promise<Array<string>> => {
  const timeout = 5000 // 5 seconds timeout
  const startTime = Date.now()
  const usernames = Array.from(new Set([user.id, user.originalId].filter(Boolean))) as string[]
  let logonResponse: Response | undefined

  for (const username of usernames) {
    let retry = true

    // server may return 502 Bad Gateway if the request is too early
    // retry until timeout
    while (retry) {
      logonResponse = await logonRequest(username, user.password)
      const elapsedTime = Date.now() - startTime
      retry = elapsedTime < timeout

      if (logonResponse.status === 200) {
        const cookies = logonResponse.headers.raw()['set-cookie']?.[0] || ''
        const data = (await logonResponse.json()) as { hello: { continue_uri: string } }
        return [data.hello.continue_uri, cookies]
      }

      if (logonResponse.status === 502 && retry) {
        console.info('[INFO] Failed with 502 Bad Gateway. Retrying logon request...')
        // wait for 1 second before retrying
        await new Promise((resolve) => setTimeout(resolve, 1000))
        continue
      }

      break
    }
  }

  throw new Error(
    `Logon failed for all candidate usernames: ${usernames.join(', ')}. Last status: ${logonResponse?.status} ${logonResponse?.statusText}`
  )
}

const getCode = async ({
  continueUrl,
  cookies
}: {
  continueUrl: string
  cookies: string
}): Promise<string> => {
  const params = new URLSearchParams({
    client_id: 'web',
    prompt: 'none',
    redirect_uri: config.baseUrl + redirectUrl,
    response_mode: 'query',
    response_type: 'code',
    scope: 'openid profile offline_access email'
  })
  const authorizeResponse = await fetch(`${continueUrl}?${params.toString()}`, {
    method: 'GET',
    redirect: 'manual',
    headers: {
      Cookie: cookies
    }
  })
  if (authorizeResponse.status !== 302) {
    throw new Error(
      `Authorization failed: Expected status code be 302 but received ${authorizeResponse.status} Message: ${authorizeResponse.statusText}`
    )
  }

  const locationHeader = authorizeResponse.headers.get('location')
  const urlParams = new URLSearchParams(new URL(locationHeader).search)

  if (locationHeader.includes('error=login_required')) {
    const errorDescription = urlParams.get('error_description')
    throw new Error(`Redirection failed. ${errorDescription}`)
  }
  return urlParams.get('code')
}

const getToken = async (code: string): Promise<Response> => {
  const response = await fetch(config.baseUrl + tokenUrl, {
    method: 'POST',
    body: new URLSearchParams({
      client_id: 'web',
      code: code,
      redirect_uri: config.baseUrl + redirectUrl,
      grant_type: 'authorization_code'
    })
  })
  if (response.status !== 200) {
    throw new Error(
      `Request failed: Expected status code be 200 but received ${response.status} Message: ${response.statusText}`
    )
  }
  return response
}

export const setAccessAndRefreshToken = async (user: User) => {
  const [authorizedUrl, cookies] = await getAuthorizedEndPoint(user)
  const code = await getCode({ continueUrl: authorizedUrl, cookies })
  const response = await getToken(code)
  const tokenList = (await response.json()) as Token

  const tokenEnvironment = TokenEnvironmentFactory()
  tokenEnvironment.setToken({
    user: { ...user },
    token: {
      userId: user.id,
      accessToken: tokenList.access_token,
      refreshToken: tokenList.refresh_token
    }
  })
}

export const refreshAccessToken = async (user: User): Promise<void> => {
  const tokenEnvironment = TokenEnvironmentFactory()

  const body = new URLSearchParams()
  body.append('client_id', 'web')
  body.append('grant_type', 'refresh_token')
  body.append('refresh_token', tokenEnvironment.getToken({ user }).refreshToken)

  const response = await request({
    method: 'POST',
    path: tokenUrl,
    body,
    header: { 'Content-Type': 'application/x-www-form-urlencoded' },
    user
  })
  checkResponseStatus(response, 'Failed refresh access token')

  const tokenList = (await response.json()) as Token

  // update tokens
  tokenEnvironment.setToken({
    user: { ...user },
    token: {
      userId: user.id,
      accessToken: tokenList.access_token,
      refreshToken: tokenList.refresh_token
    }
  })
}
