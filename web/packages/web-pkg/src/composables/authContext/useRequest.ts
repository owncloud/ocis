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

    config.headers = config.headers || {}

    if (authStore.publicLinkContextReady) {
      if (authStore.publicLinkPassword) {
        config.headers.Authorization =
          'Basic ' +
          Buffer.from(['public', authStore.publicLinkPassword].join(':')).toString('base64')
      }
      if (authStore.publicLinkToken) {
        config.headers['public-token'] = authStore.publicLinkToken
      }
    }

    return httpClient.request({ ...config, method, url })
  }

  return {
    makeRequest
  }
}
