import PQueue from 'p-queue'
import { join } from 'path'
import { webdav as _webdav } from '@ownclouders/web-client'
import type { WorkerTopic } from '../../piniaStores/webWorkers'
import { TransferType } from '../../../helpers/resource/conflictHandling/types'
import type { TransferData } from '../../../helpers/resource/conflictHandling'
import { PasteWorkerReturnData } from './usePasteWorker'

type MessageData = {
  baseUrl?: string
  accessToken?: string
  headers?: Record<string, string>
  transferData?: TransferData[]
}

type Message = {
  topic: WorkerTopic
  data: MessageData
}

let storedHeaders: Record<string, string>

self.onmessage = async (e: MessageEvent) => {
  const { topic, data } = JSON.parse(e.data) as Message

  if (topic === 'tokenUpdate' && storedHeaders) {
    const existingToken = storedHeaders.Authorization

    // token must only be updated for bearer tokens, not on public links
    if (existingToken?.toString().startsWith('Bearer')) {
      storedHeaders.Authorization = data.accessToken
    }

    return
  }

  const { baseUrl, headers, transferData } = data

  storedHeaders = headers
  const webdav = _webdav(
    baseUrl,
    () => {},
    () => storedHeaders
  )

  const successful: PasteWorkerReturnData['successful'] = []
  const failed: PasteWorkerReturnData['failed'] = []
  const queue = new PQueue({ concurrency: 4 })

  const doCopy = (transferData: TransferData) =>
    webdav.copyFiles(
      transferData.sourceSpace,
      transferData.resource,
      transferData.targetSpace,
      { path: transferData.path },
      { overwrite: transferData.overwrite }
    )

  const doMove = (transferData: TransferData) =>
    webdav.moveFiles(
      transferData.sourceSpace,
      transferData.resource,
      transferData.targetSpace,
      { path: transferData.path },
      { overwrite: transferData.overwrite }
    )

  const promises = transferData.map((data) => {
    return queue.add(async () => {
      const resource = data.resource

      try {
        if (
          data.transferType === TransferType.COPY ||
          data.transferType === TransferType.DUPLICATE
        ) {
          await doCopy(data)
          resource.id = undefined
          resource.fileId = undefined
        } else {
          await doMove(data)
        }

        resource.path = join(data.targetFolder.path, resource.name)
        resource.webDavPath = join(data.targetFolder.webDavPath, resource.name)

        successful.push(resource)
      } catch (e) {
        console.error(e)
        failed.push({
          resourceName: resource.name,
          message: e.message,
          statusCode: e.statusCode,
          xReqId: e.response.headers?.get('x-request-id')
        })
      }
    })
  })

  await Promise.allSettled(promises)

  postMessage(JSON.stringify({ successful, failed }))
}
