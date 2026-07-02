import { useGettext } from 'vue3-gettext'
import { useGetMatchingSpace } from '../spaces'
import { useAppsStore, useResourcesStore, useSpacesStore } from '../piniaStores'
import { useClientService } from '../clientService'
import { EDITOR_MODE_EDIT, useFileActions } from './files'
import { storeToRefs } from 'pinia'
import { unref } from 'vue'
import { resolveFileNameDuplicate } from '../../helpers'
import { urlJoin } from '@ownclouders/web-client'

// open an editor with an empty file within the current folder
export const useOpenEmptyEditor = () => {
  const { getMatchingSpace } = useGetMatchingSpace()
  const spacesStore = useSpacesStore()
  const appsStore = useAppsStore()
  const resourcesStore = useResourcesStore()
  const clientService = useClientService()
  const { $gettext } = useGettext()
  const { openEditor } = useFileActions()
  const { resources, currentFolder } = storeToRefs(resourcesStore)

  const openEmptyEditor = async (appId: string, extension: string) => {
    let destinationSpace = unref(currentFolder) ? getMatchingSpace(unref(currentFolder)) : null
    let destinationFiles = unref(resources)
    let filePath = unref(currentFolder)?.path

    if (!destinationSpace || !unref(currentFolder).canCreate()) {
      destinationSpace = spacesStore.personalSpace
      destinationFiles = (await clientService.webdav.listFiles(destinationSpace)).children
      filePath = ''
    }

    let fileName = $gettext('New file') + `.${extension}`

    if (destinationFiles.some((f) => f.name === fileName)) {
      fileName = resolveFileNameDuplicate(fileName, extension, destinationFiles)
    }

    const emptyResource = await clientService.webdav.putFileContents(destinationSpace, {
      path: urlJoin(filePath, fileName)
    })

    const space = getMatchingSpace(emptyResource)
    const appFileExtension = appsStore.fileExtensions.find(
      ({ app, extension: ext }) => app === appId && ext === extension
    )

    openEditor(appFileExtension, space, emptyResource, EDITOR_MODE_EDIT)
  }

  return {
    openEmptyEditor
  }
}
