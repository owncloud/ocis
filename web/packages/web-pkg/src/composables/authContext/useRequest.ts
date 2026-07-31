import { useClientService } from '../clientService'
import type { Router, RouteLocationNormalizedLoaded } from 'vue-router'
import type { Method, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ClientService } from '../../services'
import { AuthStore, useAuthStore } from '../piniaStores'

interface RequestOptions {
  router?: Router
  authStore?: AuthStore
  clientService?: ClientService
  currentRoute?: RouteLocationNormalizedLoaded
}

export interface RequestResult {
  makeRequest(method: Method, url: string, config?: AxiosRequestConfig): Promise<AxiosResponse>
}

export function useRequest(options: RequestOptions = {}): RequestResult {
  const clientService = options.clientService ?? useClientService()
  const authStore = options.authStore ?? useAuthStore()

  const makeRequest = (
    method: Method,
    url: string,
    config: AxiosRequestConfig = {}
  ): Promise<AxiosResponse> => {
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

    config.method = method
    config.url = url

    return httpClient.request(config)
  }

  return {
    makeRequest
  }
}
