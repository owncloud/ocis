import { useClientService } from '../clientService'
import { triggerDownloadWithFilename } from '../../../src/helpers'
import { useGettext } from 'vue3-gettext'
import { ClientService } from '../../services'
import { useCapabilityStore, useMessages, useUserStore } from '../piniaStores'
import { Resource, SpaceResource } from '@ownclouders/web-client'

export interface DownloadFileOptions {
  clientService?: ClientService
}

export const useDownloadFile = (options?: DownloadFileOptions) => {
  const { showErrorMessage } = useMessages()
  const clientService = options?.clientService || useClientService()
  const { $gettext } = useGettext()
  const capabilityStore = useCapabilityStore()
  const userStore = useUserStore()

  const downloadFile = async (space: SpaceResource, file: Resource, version: string = null) => {
    try {
      const url = await clientService.webdav.getFileUrl(space, file, {
        version,
        doHeadRequest: true,
        isUrlSigningEnabled: capabilityStore.supportUrlSigning,
        username: userStore.user?.onPremisesSamAccountName
      })
      triggerDownloadWithFilename(url, file.name)
    } catch (e) {
      console.error(e)
      showErrorMessage({
        title: $gettext('Download failed'),
        desc: $gettext('File could not be located'),
        errors: [e]
      })
    }
  }

  return {
    downloadFile
  }
}
