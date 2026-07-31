import { OwnCloudServer } from './types'
import { ClientService } from '@ownclouders/web-pkg'
import { urlJoin } from '@ownclouders/web-client'

interface OwnCloudInstancesResponse {
  subject: string
  links: OwnCloudServer[]
}

const OWNCLOUD_REL = 'http://webfinger.owncloud/rel/server-instance'

export class WebfingerDiscovery {
  private serverUrl: string
  private clientService: ClientService

  constructor(serverUrl: string, clientService: ClientService) {
    this.serverUrl = serverUrl
    this.clientService = clientService
  }

  public async discoverOwnCloudServers(): Promise<OwnCloudServer[]> {
    const client = this.clientService.httpAuthenticated
    const url =
      urlJoin(this.serverUrl, '.well-known', 'webfinger') + `?resource=${encodeURI(this.serverUrl)}`
    const response = (await client.get<OwnCloudInstancesResponse>(url)).data
    return response.links.filter((o) => o.rel === OWNCLOUD_REL)
  }
}
