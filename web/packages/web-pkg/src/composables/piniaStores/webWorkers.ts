import { defineStore } from 'pinia'
import { Ref, ref, unref } from 'vue'
import { v4 as uuidV4 } from 'uuid'
import { UseWebWorkerReturn, useWebWorker } from '@vueuse/core'

export interface WebWorker<T = any> extends UseWebWorkerReturn<T> {
  id: string
  needsTokenRenewal: boolean
}

export type WorkerTopic = 'startProcess' | 'tokenUpdate'

export const useWebWorkersStore = defineStore('webWorkers', () => {
  const workers = ref([]) as Ref<WebWorker[]>

  const createWorker = <T = any>(
    workerUrl: string,
    { needsTokenRenewal = false }: { needsTokenRenewal?: boolean } = {}
  ): WebWorker<T> => {
    const workerId = uuidV4()
    const result = useWebWorker(workerUrl, { type: 'module', name: workerId })
    const worker = { id: workerId, needsTokenRenewal, ...result }
    unref(workers).push(worker)
    return worker
  }

  const terminateWorker = (id: string) => {
    const worker = unref(workers).find((w) => id === w.id)
    if (worker) {
      worker.terminate()
      workers.value = unref(workers).filter((w) => id !== w.id)
    }
  }

  const terminateAllWorkers = () => {
    unref(workers).forEach(({ terminate }) => {
      terminate()
    })

    workers.value = []
  }

  const updateAccessTokens = (accessToken: string) => {
    unref(workers)
      .filter(({ needsTokenRenewal }) => needsTokenRenewal)
      .forEach(({ post }) => {
        post(
          JSON.stringify({ topic: 'tokenUpdate', data: { accessToken: `Bearer ${accessToken}` } })
        )
      })
  }

  return {
    createWorker,
    terminateWorker,
    terminateAllWorkers,
    updateAccessTokens
  }
})

export type WebWorkersStore = ReturnType<typeof useWebWorkersStore>
