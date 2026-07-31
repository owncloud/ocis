import { Capabilities, GetCapabilitiesFactory } from './capabilities'
import { AxiosInstance } from 'axios'
import { SignUrlPayload, UrlSign } from './urlSign'

export * from './capabilities'

export interface OCS {
  getCapabilities: () => Promise<Capabilities>
  signUrl: (payload: SignUrlPayload) => Promise<string>
}

export const ocs = (baseURI: string, axiosClient: AxiosInstance): OCS => {
  const url = new URL(baseURI)
  url.pathname = [...url.pathname.split('/'), 'ocs', 'v2.php'].filter(Boolean).join('/')
  const ocsV2BaseURI = url.href

  const capabilitiesFactory = GetCapabilitiesFactory(ocsV2BaseURI, axiosClient)

  const urlSign = new UrlSign({ baseURI, axiosClient })

  return {
    getCapabilities: () => {
      return capabilitiesFactory.getCapabilities()
    },
    signUrl: (payload: SignUrlPayload) => {
      return urlSign.signUrl(payload)
    }
  }
}
