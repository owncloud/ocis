import {
  DavHttpError,
  PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION,
  Resource,
  SpaceResource,
  urlJoin
} from '@ownclouders/web-client'
import { captureException } from '@sentry/vue'
import { unref } from 'vue'
import { storeToRefs } from 'pinia'
import { useClientService } from '../../clientService'
import { useGetMatchingSpace } from '../../spaces'
import { useSpacesStore } from '../../piniaStores'

/**
 * The `.psec` pointer file together with the space it lives in. The space is needed
 * so callers can move/rename the file, since it may reside in a different space than
 * the real folder (the folder always lives in the personal space).
 */
export type ResolvedPsecFile = { psecFile: Resource; space: SpaceResource }

/**
 * A password protected folder is represented by two independent resources:
 * the `.psec` pointer file (visible to the owner) and the real folder living under
 * `/.PasswordProtectedFolders/projects/<space>/...` (accessible via the public link).
 * This resolves one from the other so callers can keep both in sync.
 */
export const useResolvePasswordProtectedFolderCounterpart = () => {
  const clientService = useClientService()
  const { getMatchingSpace } = useGetMatchingSpace()
  const spacesStore = useSpacesStore()
  const { getSpacesByName } = spacesStore
  const { personalSpace } = storeToRefs(spacesStore)

  const getPasswordProtectedFolder = async (psecFile: Resource): Promise<Resource | null> => {
    try {
      const matchingSpace = getMatchingSpace(psecFile)

      const folderPath = urlJoin(
        '/.PasswordProtectedFolders/projects/',
        matchingSpace.name,
        psecFile.path.replace(psecFile.name, ''),
        psecFile.name.replace(`.${PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION}`, '')
      )

      return await clientService.webdav.getFileInfo(unref(personalSpace), { path: folderPath })
    } catch (error) {
      if (error instanceof DavHttpError && error.statusCode === 404) {
        return null
      }

      console.error(error)
      captureException(error)
      return null
    }
  }

  const getPsecFile = async (folder: Resource): Promise<ResolvedPsecFile | null> => {
    const [spaceName, ...psecFilePathParts] = folder.path
      .replace('/.PasswordProtectedFolders/projects/', '')
      .split('/')
    const matchingSpaces = getSpacesByName(spaceName)

    if (matchingSpaces.length < 1) {
      return null
    }

    for (const space of matchingSpaces) {
      try {
        const psecFile = await clientService.webdav.getFileInfo(space, {
          path: urlJoin(...psecFilePathParts) + '.' + PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION
        })
        return { psecFile, space }
      } catch (error) {
        if (error instanceof DavHttpError && error.statusCode === 404) {
          continue
        }

        console.error(error)
        captureException(error)
        continue
      }
    }

    return null
  }

  return { getPasswordProtectedFolder, getPsecFile }
}
