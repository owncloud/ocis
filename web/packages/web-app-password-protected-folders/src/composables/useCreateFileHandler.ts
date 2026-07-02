import { Resource, SpaceResource, urlJoin } from '@ownclouders/web-client'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { useClientService, useResourcesStore, useSharesStore } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

export const useCreateFileHandler = () => {
  const clientService = useClientService()
  const { upsertResource } = useResourcesStore()
  const { addLink } = useSharesStore()
  const { $pgettext } = useGettext()

  const createFileHandler = async ({
    fileName,
    personalSpace,
    currentSpace,
    currentFolder,
    password,
    type
  }: {
    fileName: string
    personalSpace: SpaceResource
    currentSpace: SpaceResource
    currentFolder: Resource
    password: string
    type: SharingLinkType
  }) => {
    if (fileName === '') {
      return
    }

    const folderPath = urlJoin(
      '/.PasswordProtectedFolders/projects/',
      currentSpace.name,
      currentFolder.path,
      fileName
    )
    const folder = await clientService.webdav.createFolder(personalSpace, {
      path: folderPath,
      recursive: true
    })

    try {
      if (currentSpace.id === personalSpace.id && currentFolder.path === '/') {
        const psecFolders = await clientService.webdav.getFileInfo(personalSpace, {
          path: '/.PasswordProtectedFolders'
        })
        upsertResource(psecFolders)
      }

      const share = await addLink({
        clientService,
        space: personalSpace,
        resource: folder,
        options: { password, type }
      })

      const shareUrl = new URL(share.webUrl)
      if (!['https:', 'http:'].includes(shareUrl.protocol)) {
        throw new Error(
          $pgettext(
            'Error shown when creating a password-protected folder fails because the generated share link has an unexpected format.',
            'The folder could not be created because the share link is invalid.'
          )
        )
      }

      const path = urlJoin(currentFolder.path, fileName + '.psec')

      const file = await clientService.webdav.putFileContents(currentSpace, {
        path,
        content: btoa(share.webUrl)
      })
      upsertResource(file)
    } catch (error) {
      await clientService.webdav.deleteFile(personalSpace, {
        path: folderPath
      })

      throw error
    }
  }

  return { createFileHandler }
}
