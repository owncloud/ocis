import { HttpClient } from '../../http'
import { graph, ocs, webdav } from '@ownclouders/web-client'
import { Graph } from '@ownclouders/web-client/graph'
import { OCS } from '@ownclouders/web-client/ocs'
import { AuthParameters } from './auth'
import { FetchClient, type OnResponseArgs } from '@ownclouders/web-client'
import { v4 as uuidV4 } from 'uuid'
import { WebDAV } from '@ownclouders/web-client/webdav'
import { Language } from 'vue3-gettext'
import { FetchEventSourceInit } from '@microsoft/fetch-event-source'
import { sse } from '@ownclouders/web-client/sse'
import { AuthStore, ConfigStore } from '../../composables'
import { shouldResponseTriggerMaintenance } from '@ownclouders/web-client'

const createFetchOptions = (authParams: AuthParameters, language: string): FetchEventSourceInit => {
  return {
    headers: {
      Authorization: `Bearer ${authParams.accessToken}`,
      'Accept-Language': language,
      'X-Request-ID': uuidV4(),
      'X-Requested-With': 'XMLHttpRequest'
    }
  }
}

export interface ClientServiceOptions {
  configStore: ConfigStore
  language: Language
  authStore: AuthStore
}

export class ClientService {
  private configStore: ConfigStore
  private language: Language
  private authStore: AuthStore

  private httpAuthenticatedClient: HttpClient
  private httpUnAuthenticatedClient: HttpClient

  private graphClient: Graph
  private graphHttpClient: FetchClient
  private ocsClient: OCS
  private webDavClient: WebDAV

  public initiatorId = uuidV4()
  public lastSuccessfulRequestTime: number | null = null

  private staticHeaders: Record<string, string> = {
    'Initiator-ID': this.initiatorId,
    'X-Requested-With': 'XMLHttpRequest'
  }

  constructor(options: ClientServiceOptions) {
    this.configStore = options.configStore
    this.language = options.language
    this.authStore = options.authStore

    this.initGraphClient(this.configStore.isInVault)
    this.initOcsClient(this.configStore.isInVault)
    this.initWebDavClient()

    this.httpAuthenticatedClient = new HttpClient({
      baseUrl: this.configStore.serverUrl,
      staticHeaders: this.staticHeaders,
      headers: () => this.getDynamicHeaders(),
      onResponse: (args) => this.handleResponse(args)
    })
    this.httpUnAuthenticatedClient = new HttpClient({
      baseUrl: this.configStore.serverUrl,
      staticHeaders: this.staticHeaders,
      headers: () => this.getDynamicHeaders({ useAuth: false }),
      onResponse: (args) => this.handleResponse(args)
    })
  }

  public get httpAuthenticated() {
    return this.httpAuthenticatedClient
  }

  public get httpUnAuthenticated() {
    return this.httpUnAuthenticatedClient
  }

  public get graphAuthenticated() {
    return this.graphClient
  }

  public reinitializeGraphClient(isInVault: boolean) {
    this.initGraphClient(isInVault)
  }

  public reinitializeOcsClient(isInVault: boolean) {
    this.initOcsClient(isInVault)
  }

  public get sseAuthenticated(): EventSource {
    return sse(
      this.configStore.serverUrl,
      createFetchOptions({ accessToken: this.authStore.accessToken }, this.currentLanguage)
    )
  }

  public get ocs() {
    return this.ocsClient
  }

  public get webdav() {
    return this.webDavClient
  }

  get currentLanguage() {
    return this.language.current
  }

  private initGraphClient(isInVault: boolean) {
    if (!this.graphHttpClient) {
      this.graphHttpClient = new FetchClient({
        staticHeaders: this.staticHeaders,
        headers: () => this.getDynamicHeaders(),
        onResponse: (args) => this.handleResponse(args)
      })
    }

    this.graphClient = graph(
      isInVault ? `${this.configStore.serverUrl}vault` : this.configStore.serverUrl,
      this.graphHttpClient
    )
  }

  private initOcsClient(isInVault: boolean) {
    const httpClient = new FetchClient({
      staticHeaders: this.staticHeaders,
      headers: () => this.getDynamicHeaders(),
      onResponse: (args) => this.handleResponse(args)
    })

    const baseUrl = isInVault
      ? `${this.configStore.serverUrl}?vault=true`
      : this.configStore.serverUrl
    this.ocsClient = ocs(baseUrl, httpClient)
  }

  private initWebDavClient() {
    this.webDavClient = webdav(
      this.configStore.serverUrl,
      this.configStore.setMaintenanceMode,
      () => {
        const headers = { ...this.staticHeaders, ...this.getDynamicHeaders() }

        if (this.authStore.publicLinkToken) {
          headers['public-token'] = this.authStore.publicLinkToken
        }

        if (this.authStore.publicLinkPassword) {
          headers['Authorization'] =
            'Basic ' +
            Buffer.from(['public', this.authStore.publicLinkPassword].join(':')).toString('base64')
        }

        return headers
      }
    )
  }

  /**
   * Dynamic headers that should be provided via callback or interceptor because they may
   * change during the lifetime of the application (e.g. token renewal).
   */
  private getDynamicHeaders({ useAuth = true }: { useAuth?: boolean } = {}): Record<
    string,
    string
  > {
    return {
      'Accept-Language': this.currentLanguage,
      'X-Request-ID': uuidV4(),
      ...(useAuth &&
        this.authStore.accessToken && { Authorization: 'Bearer ' + this.authStore.accessToken })
    }
  }

  /**
   * Called for every response the client receives. The asymmetry is deliberate: only a
   * successful response clears maintenance mode, and a non-2xx response never clears it.
   *
   * `args.requestUrl` is the caller's URL, not `response.url` — the maintenance
   * allow-list is matched against relative paths. `args.status` is 500 when the transport
   * failed and there is no response at all.
   */
  public handleResponse({ response, status, requestUrl }: OnResponseArgs): void {
    if (response?.ok) {
      this.configStore.setMaintenanceMode(false)
      this.lastSuccessfulRequestTime = Math.floor(Date.now() / 1000)
      return
    }

    if (shouldResponseTriggerMaintenance(status, requestUrl)) {
      this.configStore.setMaintenanceMode(true)
    }
  }
}
