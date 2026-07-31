import { ref, unref } from 'vue'
import { WebWorker, useWebWorkersStore } from '../../piniaStores/webWorkers'
import MfaWorker from './worker?worker'

export type MfaExpiryWorkerTopic = 'set' | 'reset'

const MFA_WARNING_THRESHOLD_SECONDS = 300 // 5 minutes before expiry

export const useMfaExpiryWorker = ({ onExpiring }: { onExpiring: () => void }) => {
  const { createWorker } = useWebWorkersStore()

  const worker = ref<WebWorker>()

  const startWorker = () => {
    worker.value = createWorker(MfaWorker as unknown as string)

    unref(unref(worker).worker).onmessage = () => {
      onExpiring()
    }
  }

  const setMfaTimer = ({ expiresAt }: { expiresAt: number }) => {
    if (!unref(worker)) {
      console.error('mfa expiry worker is not running')
      return
    }

    unref(worker).post(
      JSON.stringify({
        topic: 'set',
        expiresAt,
        warningThreshold: MFA_WARNING_THRESHOLD_SECONDS
      })
    )
  }

  const resetMfaTimer = () => {
    if (!unref(worker)) {
      console.error('mfa expiry worker is not running')
      return
    }

    unref(worker).post(JSON.stringify({ topic: 'reset' }))
  }

  return { startWorker, setMfaTimer, resetMfaTimer }
}
