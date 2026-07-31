import { defineStore } from 'pinia'
import { Ref, computed, ref, unref } from 'vue'
import {
  extractStorageId,
  isMountPointSpaceResource,
  isProjectSpaceResource,
  SpaceResource,
  type Resource
} from '@ownclouders/web-client'
import { getParentPaths } from '../../helpers'
import { AncestorMetaData, AncestorMetaDataValue } from '../../types'
import { DavProperty, WebDAV } from '@ownclouders/web-client/webdav'
import { useSpacesStore } from './spaces'
import { useUserStore } from './user'
import { useConfigStore } from './config'

export const useResourcesStore = defineStore('resources', () => {
  const configStore = useConfigStore()
  const spacesStore = useSpacesStore()
  const userStore = useUserStore()

  const resources = ref<Resource[]>([]) as Ref<Resource[]>
  const currentFolder = ref<Resource>()
  const ancestorMetaData = ref<AncestorMetaData>({})
  const deleteQueue = ref<string[]>([])

  const activeResources = computed(() => {
    let res = unref(resources)

    if (!unref(areHiddenFilesShown)) {
      res = res.filter((file) => !file.name.startsWith('.'))
    }

    return res
  })

  const totalResourcesCount = computed(() => {
    const fileCount = unref(resources).filter(({ type }) => type === 'file').length
    const hiddenFileCount = unref(resources).filter(
      ({ type, name }) => type === 'file' && name.startsWith('.')
    ).length
    const folderCount = unref(resources).filter(({ type }) => type === 'folder').length
    const hiddenFolderCount = unref(resources).filter(
      ({ type, name }) => type === 'folder' && name.startsWith('.')
    ).length
    const spaceCount = unref(resources).filter(isProjectSpaceResource).length

    return {
      files: fileCount,
      hiddenFiles: hiddenFileCount,
      folders: folderCount,
      hiddenFolders: hiddenFolderCount,
      spaces: spaceCount
    }
  })

  const setResources = (file: Resource[]) => {
    resources.value = file
  }

  const removeResources = (values: Resource[]) => {
    resources.value = unref(resources).filter((file) => !values.find(({ id }) => id === file.id))
  }

  const clearResources = () => {
    resources.value = []
  }

  const upsertResource = (resource: Resource) => {
    const existing = unref(resources).find(({ id }) => id === resource.id)
    if (existing) {
      Object.assign(existing, resource)
      return
    }
    unref(resources).push(resource)
  }

  const upsertResources = (values: Resource[]) => {
    const other = unref(resources).filter((f) => !values.some((r) => r.path === f.path))
    resources.value = [...other, ...values]
  }

  const updateResourceField = <T extends Resource>({
    id,
    field,
    value
  }: {
    id: T['id']
    field: keyof T
    value: T[keyof T]
  }) => {
    const resource = unref(resources).find((resource) => id === resource.id) as T
    if (resource) {
      resource[field] = value
    }
  }

  const setCurrentFolder = (value: Resource) => {
    currentFolder.value = value
  }

  const clearCurrentFolder = () => {
    currentFolder.value = undefined
  }

  const initResourceList = <T extends Resource>(data: { resources: T[]; currentFolder: T }) => {
    resources.value = data.resources
    currentFolder.value = data.currentFolder
  }

  const clearResourceList = () => {
    resources.value = []
    currentFolder.value = undefined
    selectedIds.value = []
  }

  const selectedIds = ref<string[]>([])
  const latestSelectedId = ref<string>(null)

  const selectedResources = computed(() => {
    return unref(resources).filter((f) => unref(selectedIds).includes(f.id))
  })

  const setSelection = (ids: string[]) => {
    const latestSelected = ids.find((id) => !unref(selectedIds).includes(id))

    if (latestSelected) {
      latestSelectedId.value = latestSelected
    }
    selectedIds.value = ids
  }

  const addSelection = (id: string) => {
    latestSelectedId.value = id
    if (!unref(selectedIds).includes(id)) {
      unref(selectedIds).push(id)
    }
  }

  const removeSelection = (id: string) => {
    latestSelectedId.value = id
    if (unref(selectedIds).includes(id)) {
      selectedIds.value = unref(selectedIds).filter((i) => i !== id)
    }
  }

  const toggleSelection = (id: string) => {
    if (unref(selectedIds).includes(id)) {
      removeSelection(id)
    } else {
      addSelection(id)
    }
  }

  const resetSelection = () => {
    selectedIds.value = []
  }

  const setLastSelectedId = (id: string) => {
    latestSelectedId.value = id
  }

  const shouldShowFlatList = ref(false)
  const areHiddenFilesShown = ref(true)
  const areFileExtensionsShown = ref(true)
  const areWebDavDetailsShown = ref(false)

  const setAreHiddenFilesShown = (value: boolean) => {
    areHiddenFilesShown.value = value
    window.localStorage.setItem('oc_hiddenFilesShown', value.toString())
  }
  const setShouldShowFlatList = (value: boolean) => {
    shouldShowFlatList.value = value
    window.localStorage.setItem('oc_flatList', value.toString())
  }
  const setAreFileExtensionsShown = (value: boolean) => {
    areFileExtensionsShown.value = value
    window.localStorage.setItem('oc_fileExtensionsShown', value.toString())
  }
  const setAreWebDavDetailsShown = (value: boolean) => {
    areWebDavDetailsShown.value = value
    window.localStorage.setItem('oc_webDavDetailsShown', value.toString())
  }

  const setAncestorMetaData = (value: AncestorMetaData) => {
    ancestorMetaData.value = value
  }

  const updateAncestorField = <
    T extends AncestorMetaDataValue,
    K extends keyof AncestorMetaDataValue
  >({
    path,
    field,
    value
  }: {
    path: T['path']
    field: K
    value: T[K]
  }) => {
    const resource = unref(ancestorMetaData)[path] ?? null
    if (resource) {
      resource[field] = value
    }
  }

  const loadAncestorMetaData = ({
    folder,
    space,
    client,
    signal
  }: {
    folder: Resource
    space: SpaceResource
    client: WebDAV
    signal?: AbortSignal
  }) => {
    const data: AncestorMetaData = {
      [folder.path]: {
        id: folder.fileId,
        shareTypes: folder.shareTypes,
        parentFolderId: folder.parentFolderId,
        spaceId: space.id,
        path: folder.path
      }
    }
    const promises = []
    const davProperties = [DavProperty.FileId, DavProperty.ShareTypes, DavProperty.FileParent]
    const parentPaths = getParentPaths(folder.path)
    const spaces = spacesStore.spaces

    const getMountPoints = () =>
      spaces.filter(
        (s) =>
          isMountPointSpaceResource(s) && extractStorageId(s.root.remoteItem.rootId) === space.id
      )

    let fullyAccessibleSpace = true
    if (configStore.options.routing.fullShareOwnerPaths) {
      // keep logic in sync with "isResourceAccessible" from useGetMatchingSpace
      const projectSpace = spaces.find((s) => isProjectSpaceResource(s) && s.id === space.id)
      fullyAccessibleSpace = space.isOwner(userStore.user) || projectSpace?.isMember(userStore.user)
    }

    for (const path of parentPaths) {
      const cachedData = unref(ancestorMetaData)[path] ?? null
      if (cachedData?.spaceId === space.id) {
        data[path] = cachedData
        continue
      }

      // keep logic in sync with "isResourceAccessible" from useGetMatchingSpace
      if (
        !fullyAccessibleSpace &&
        !getMountPoints().find((m) => path.startsWith(m.root.remoteItem.path))
      ) {
        // no access to the parent resource
        break
      }

      promises.push(
        client
          .listFiles(space, { path }, { depth: 0, davProperties, signal })
          .then(({ resource }) => {
            data[path] = {
              id: resource.fileId,
              shareTypes: resource.shareTypes,
              parentFolderId: resource.parentFolderId,
              spaceId: space.id,
              path
            }
          })
      )
    }

    return Promise.all(promises).then(() => {
      if (!Object.keys(data).includes('/')) {
        // add space as root element
        const cachedRoot = unref(ancestorMetaData)['/']

        if (cachedRoot?.spaceId === space.id) {
          data['/'] = cachedRoot
        } else {
          const { parentFolderId } = Object.values(data)[0]

          if (parentFolderId) {
            const space = spacesStore.spaces.find(({ id }) => parentFolderId.startsWith(id))

            if (space) {
              data['/'] = {
                id: space.id,
                shareTypes: space.shareTypes,
                parentFolderId: space.id,
                spaceId: space.id,
                path: '/'
              }
            }
          }
        }
      }

      setAncestorMetaData(data)
    })
  }

  const getAncestorById = (id: string) => {
    return Object.values(unref(ancestorMetaData)).find((a) => id === a.id)
  }

  const addResourcesIntoDeleteQueue = (ids: string[]): void => {
    deleteQueue.value = deleteQueue.value.concat(
      ids.filter((id) => !unref(deleteQueue).includes(id))
    )
  }

  const removeResourcesFromDeleteQueue = (ids: string[]): void => {
    deleteQueue.value = deleteQueue.value.filter((id) => !ids.includes(id))
  }

  return {
    resources,
    currentFolder,
    activeResources,
    totalResourcesCount,

    setResources,
    removeResources,
    clearResources,
    upsertResource,
    upsertResources,
    updateResourceField,

    setCurrentFolder,
    clearCurrentFolder,

    initResourceList,
    clearResourceList,

    selectedIds,
    latestSelectedId,
    selectedResources,
    setSelection,
    addSelection,
    removeSelection,
    toggleSelection,
    resetSelection,
    setLastSelectedId,

    shouldShowFlatList,
    setShouldShowFlatList,
    areHiddenFilesShown,
    areFileExtensionsShown,
    areWebDavDetailsShown,
    setAreHiddenFilesShown,
    setAreFileExtensionsShown,
    setAreWebDavDetailsShown,

    ancestorMetaData,
    setAncestorMetaData,
    updateAncestorField,
    loadAncestorMetaData,
    getAncestorById,

    deleteQueue,
    addResourcesIntoDeleteQueue,
    removeResourcesFromDeleteQueue
  }
})

export type ResourcesStore = ReturnType<typeof useResourcesStore>
