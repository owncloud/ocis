import { Capabilities, GetCapabilitiesFactory } from './capabilities'
import { FetchClient } from '../http'
import { SignUrlPayload, UrlSign } from './urlSign'

export * from './capabilities'

export interface OCS {
  getCapabilities: () => Promise<Capabilities>
  signUrl: (payload: SignUrlPayload) => Promise<string>
}

export const ocs = (baseURI: string, httpClient: FetchClient): OCS => {
  const url = new URL(baseURI)
  url.pathname = [...url.pathname.split('/'), 'ocs', 'v2.php'].filter(Boolean).join('/')
  const ocsV2BaseURI = url.href

  const capabilitiesFactory = GetCapabilitiesFactory(ocsV2BaseURI, httpClient)

  const urlSign = new UrlSign({ baseURI, httpClient })

  return {
    getCapabilities: () => {
      return capabilitiesFactory.getCapabilities()
    },
    signUrl: (payload: SignUrlPayload) => {
      return urlSign.signUrl(payload)
    }
  }
}
