import { DavProperty } from '@ownclouders/web-client/webdav'
import { webdav, HttpError, Resource, SpaceResource, urlJoin } from '@ownclouders/web-client'

import { WorkerTopic } from '../../piniaStores/webWorkers'
import { resolveFileNameDuplicate } from '../../../helpers/resource/conflictHandling/conflictUtils'
import { PDFRenderer } from './renderer'

/**
 * Data structure for messages sent to the PDF export worker.
 *
 * Contains all necessary information for exporting markdown content as PDF,
 * including authentication, destination information, and content to process.
 */
type MessageData = {
  /** Access token for authentication */
  accessToken?: string
  /** Base URL of the ownCloud instance */
  baseUrl?: string
  /** HTTP headers to include in requests */
  headers?: Record<string, string>
  /** Destination folder resource where the PDF will be saved */
  destinationFolder?: Resource
  /** Space resource containing the destination folder */
  space?: SpaceResource
  /** Name for the output PDF file */
  fileName?: string
  /** Markdown content to convert to PDF */
  content?: string
}

/**
 * Message structure for communication with the PDF export worker.
 *
 * Used for both PDF export requests and token update notifications.
 */
type Message = {
  /** Type of operation to perform */
  topic: WorkerTopic
  /** Payload data for the operation */
  data: MessageData
}

let storedHeaders: Record<string, string>

/**
 * Main message handler for the PDF export worker.
 *
 * This worker handles two types of operations:
 * 1. Token updates - Updates stored authentication headers when tokens change
 * 2. PDF export - Converts markdown content to PDF and uploads it to the specified location
 *
 * For PDF export operations, the worker:
 * - Establishes a WebDAV connection using provided credentials
 * - Checks for filename conflicts and resolves duplicates
 * - Renders the markdown content to PDF using the PDFRenderer
 * - Uploads the generated PDF to the destination folder
 * - Posts back success/failure results
 *
 * @param event - The message event containing the operation details
 */
self.onmessage = async (event: MessageEvent) => {
  const { topic, data } = JSON.parse(event.data) as Message

  if (topic === 'tokenUpdate' && storedHeaders) {
    const existingToken = storedHeaders.Authorization

    if (existingToken?.toString().startsWith('Bearer')) {
      storedHeaders.Authorization = data.accessToken
    }

    return
  }

  const { baseUrl, headers, destinationFolder, space, content } = data
  let { fileName } = data

  storedHeaders = headers
  const webdavService = webdav(
    baseUrl,
    () => {},
    () => storedHeaders
  )

  try {
    const { children: existingResources } = await webdavService.listFiles(
      space,
      {
        fileId: destinationFolder.fileId
      },
      { davProperties: [DavProperty.Name] }
    )

    const resourceAlreadyExists = existingResources.find(
      (existingResource) => existingResource.name === fileName
    )

    if (resourceAlreadyExists) {
      fileName = resolveFileNameDuplicate(fileName, 'pdf', existingResources)
    }

    const renderer = new PDFRenderer(content)
    const pdf = await renderer.renderAsArrayBuffer()

    const resource = await webdavService.putFileContents(space, {
      fileName,
      parentFolderId: destinationFolder.id,
      content: pdf,
      path: urlJoin(destinationFolder.path, fileName)
    })

    postMessage(JSON.stringify({ successful: [resource], failed: [] }))
  } catch (e) {
    const error = {
      message: e.message || 'Unexpected error',
      statusCode: e instanceof HttpError ? e.statusCode : 500,
      xReqId: e instanceof HttpError ? e.response?.headers?.get('x-request-id') : undefined
    }
    postMessage(
      JSON.stringify({
        successful: [],
        failed: [
          {
            resourceName: fileName,
            error
          }
        ]
      })
    )
  }
}
