import { Resource, SpaceResource, urlJoin } from '@ownclouders/web-client'
import { stringify } from 'qs'
import { useCapabilityStore, useClientService, useRequest } from '@ownclouders/web-pkg'

export const useCreateFileHandler = () => {
  const capabilityStore = useCapabilityStore()
  const clientService = useClientService()
  const { makeRequest } = useRequest({ clientService })

  const createFileHandler = async ({
    fileName,
    space,
    currentFolder
  }: {
    fileName: string
    space: SpaceResource
    currentFolder: Resource
  }) => {
    if (fileName === '') {
      return
    }

    const query = stringify({
      parent_container_id: currentFolder.fileId,
      filename: fileName
    })
    const url = `${capabilityStore.filesAppProviders[0].new_url}?${query}`
    const response = await makeRequest('POST', url)
    if (response.status !== 200) {
      throw new Error(`An error has occurred: ${response.status}`)
    }

    const path = urlJoin(currentFolder.path, fileName) || ''
    return clientService.webdav.getFileInfo(space, { path })
  }

  return {
    createFileHandler
  }
}
