import { SpaceAction, SpaceActionOptions } from '../types'
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'

import { useOpenWithDefaultApp } from '../useOpenWithDefaultApp'
import { getRelativeSpecialFolderSpacePath, Resource, SpaceResource } from '@ownclouders/web-client'
import { useClientService } from '../../clientService'
import { useSharesStore, useSpacesStore, useUserStore } from '../../piniaStores'
import { useCreateSpace, useSpaceHelpers } from '../../spaces'

export const useSpaceActionsEditReadmeContent = () => {
  const clientService = useClientService()
  const { openWithDefaultApp } = useOpenWithDefaultApp()
  const { createDefaultMetaFolder } = useCreateSpace()
  const userStore = useUserStore()
  const spacesStore = useSpacesStore()
  const sharesStore = useSharesStore()
  const { $gettext } = useGettext()
  const { getDefaultMetaFolder } = useSpaceHelpers()

  const createReadme = async (space: SpaceResource, metaFolder: Resource) => {
    // FIXME: remove path as soon as we make the full switch to id-based dav requests
    const markdownResource = await clientService.webdav.putFileContents(space, {
      path: '.space/readme.md',
      parentFolderId: metaFolder.id,
      fileName: 'readme.md'
    })

    const updatesSpace = await clientService.graphAuthenticated.drives.updateDrive(
      space.id,
      {
        name: space.name,
        special: [{ specialFolder: { name: 'readme' }, id: markdownResource.id }]
      },
      sharesStore.graphRoles
    )

    spacesStore.updateSpaceField({
      id: space.id,
      field: 'spaceReadmeData',
      value: updatesSpace.spaceReadmeData
    })

    return markdownResource
  }

  const handler = async ({ resources }: SpaceActionOptions) => {
    let markdownResource: Resource = null

    let metaFolder = await getDefaultMetaFolder(resources[0])
    if (!metaFolder) {
      metaFolder = await createDefaultMetaFolder(resources[0])
      markdownResource = await createReadme(resources[0], metaFolder)
    }

    if (!markdownResource) {
      const path = getRelativeSpecialFolderSpacePath(resources[0], 'readme')
      if (path) {
        markdownResource = await clientService.webdav.getFileInfo(resources[0], { path })
      } else {
        markdownResource = await createReadme(resources[0], metaFolder)
      }
    }

    openWithDefaultApp({ space: resources[0], resource: markdownResource })
  }

  const actions = computed((): SpaceAction[] => [
    {
      name: 'editReadmeContent',
      icon: 'article',
      label: () => {
        return $gettext('Edit description')
      },
      handler,
      isVisible: ({ resources }) => {
        if (resources.length !== 1) {
          return false
        }

        return resources[0].canEditReadme({ user: userStore.user })
      },
      class: 'oc-files-actions-edit-readme-content-trigger'
    }
  ])

  return {
    actions
  }
}
