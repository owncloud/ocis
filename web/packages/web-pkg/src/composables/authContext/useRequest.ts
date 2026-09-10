import { useClientService } from '../clientService'
import type { Router, RouteLocationNormalizedLoaded } from 'vue-router'
import type { HttpResponse } from '@ownclouders/web-client'
import type { RequestConfig } from '../../http'
import { ClientService } from '../../services'
import { AuthStore, useAuthStore } from '../piniaStores'

interface RequestOptions {
  router?: Router
  authStore?: AuthStore
  clientService?: ClientService
  currentRoute?: RouteLocationNormalizedLoaded
}

export interface RequestResult {
  makeRequest(method: string, url: string, config?: RequestConfig): Promise<HttpResponse>
}

export function useRequest(options: RequestOptions = {}): RequestResult {
  const clientService = options.clientService ?? useClientService()
  const authStore = options.authStore ?? useAuthStore()

  const makeRequest = (
    method: string,
    url: string,
    config: RequestConfig = {}
  ): Promise<HttpResponse> => {
    const httpClient = authStore.accessToken
      ? clientService.httpAuthenticated
      : clientService.httpUnAuthenticated

    // A Headers rather than a plain object so that these entries replace a caller's own
    // spelling of the same name instead of being appended alongside it.
    const headers = new Headers(config.headers)

    if (authStore.publicLinkContextReady) {
      if (authStore.publicLinkPassword) {
        headers.set(
          'Authorization',
          'Basic ' +
            Buffer.from(['public', authStore.publicLinkPassword].join(':')).toString('base64')
        )
      }
      if (authStore.publicLinkToken) {
        headers.set('public-token', authStore.publicLinkToken)
      }
    }

    return httpClient.request({ ...config, headers, method, url })
  }

  return {
    makeRequest
  }
}
