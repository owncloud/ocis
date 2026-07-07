import { ref, unref } from 'vue'
import { ErrorTimeout } from 'oidc-client-ts'
import { AuthServiceInterface } from '../../authContext'
import { WebWorker, useWebWorkersStore } from '../../piniaStores/webWorkers'
import TokenWorker from './worker?worker'

export type TokenTimerWorkerTopic = 'set' | 'reset'

export const useTokenTimerWorker = ({ authService }: { authService: AuthServiceInterface }) => {
  const { createWorker } = useWebWorkersStore()

  const worker = ref<WebWorker>()

  const startWorker = () => {
    worker.value = createWorker(TokenWorker as unknown as string)

    unref(unref(worker).worker).onmessage = () => {
      authService.signinSilent().catch(async (error) => {
        if (error instanceof ErrorTimeout) {
          console.warn('token renewal timed out, retrying in 5 seconds...')
          unref(worker).post(JSON.stringify({ topic: 'set', expiry: 5, expiryThreshold: 0 }))
          return
        }

        console.error('token renewal error:', error)

        // log out user if they don't have a refresh token
        const refreshToken = await authService.getRefreshToken()
        if (!refreshToken) {
          return authService.logoutUser()
        }
      })
    }
  }

  const setTokenTimer = ({
    expiry,
    expiryThreshold
  }: {
    expiry: number
    expiryThreshold: number
  }) => {
    if (!unref(worker)) {
      console.error('token timer worker is not running')
      return
    }

    unref(worker).post(JSON.stringify({ topic: 'set', expiry, expiryThreshold }))
  }

  const resetTokenTimer = () => {
    if (!unref(worker)) {
      console.error('token timer worker is not running')
      return
    }

    unref(worker).post(JSON.stringify({ topic: 'reset' }))
  }

  return { startWorker, setTokenTimer, resetTokenTimer }
}
