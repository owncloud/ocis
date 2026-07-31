import { ResourcesStore } from '@ownclouders/web-pkg'
import { extractNodeId } from '@ownclouders/web-client'
import { eventSchema, SseEventWrapperOptions } from './types'

export const sseEventWrapper = (options: SseEventWrapperOptions) => {
  const { topic, msg, method, ...sseEventOptions } = options
  try {
    const sseData = eventSchema.parse(JSON.parse(msg.data))
    console.debug(`SSE event '${topic}'`, sseData)

    return method({ ...sseEventOptions, sseData })
  } catch (e) {
    console.error(`Unable to process sse event ${topic}`, e)
  }
}
export const isItemInCurrentFolder = ({
  resourcesStore,
  parentFolderId
}: {
  resourcesStore: ResourcesStore
  parentFolderId: string
}) => {
  const currentFolder = resourcesStore.currentFolder
  if (!currentFolder || !currentFolder.id) {
    return false
  }

  if (!extractNodeId(currentFolder.id)) {
    // if we don't have a nodeId here, we have a space (root) as current folder and can only check against the storageId
    const spaceNodeId = currentFolder.id.split('$')[1]
    if (`${currentFolder.id}!${spaceNodeId}` !== parentFolderId) {
      return false
    }
  } else {
    if (currentFolder.id !== parentFolderId) {
      return false
    }
  }

  return true
}
