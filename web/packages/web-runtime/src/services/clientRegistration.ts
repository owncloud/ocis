import { buildUrl, OpenIdConnectConfig } from '@ownclouders/web-pkg'
import { v4 as uuidV4 } from 'uuid'
import { router } from '../router'

async function get(url: string) {
  return await fetch(url, { headers: { 'X-Request-ID': uuidV4() } })
    .then((res) => {
      return res.json()
    })
    .catch((err) => {
      console.error('Error: ', err)
    })
}

async function post(url: string, data: unknown) {
  return await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Request-ID': uuidV4()
    },
    body: JSON.stringify(data)
  })
    .then((res) => {
      return res.json()
    })
    .catch((err) => {
      console.error('Error: ', err)
    })
}

export async function registerClient(openIdConfig: OpenIdConnectConfig) {
  const clientData = JSON.parse(sessionStorage.getItem('dynamicClientData'))
  if (clientData !== null) {
    const client_secret_expires_at = clientData.client_secret_expires_at || 0

    if (client_secret_expires_at === 0 || Date.now() < client_secret_expires_at * 1000) {
      return JSON.parse(clientData)
    }
  }
  sessionStorage.removeItem('dynamicClientData')
  const wellKnown = await get(`${openIdConfig.authority}/.well-known/openid-configuration`)
  const resp = await post(wellKnown.registration_endpoint, {
    redirect_uris: [buildUrl(router, '/oidc-callback.html')],
    client_name: `ownCloud Web on ${window.location.origin}`
  })
  sessionStorage.setItem('dynamicClientData', JSON.stringify(resp))
  return resp
}
