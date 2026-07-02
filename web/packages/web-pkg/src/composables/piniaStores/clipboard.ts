import { defineStore } from 'pinia'
import { computed, ref, unref } from 'vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { ClipboardActions } from '../../helpers'
import { useGettext } from 'vue3-gettext'
import { useMessages } from './messages'
import { useConfigStore } from './config'
import { useGetMatchingSpace } from '../spaces'

const clipboardStorageKey = 'oc-clipboard'

export enum ClipboardMode {
  Default = 'default',
  Vault = 'vault'
}

type ClipboardResourceSnapshot = Pick<
  Resource,
  | 'id'
  | 'fileId'
  | 'extension'
  | 'isFolder'
  | 'name'
  | 'mimeType'
  | 'parentFolderId'
  | 'path'
  | 'remoteItemId'
  | 'remoteItemPath'
  | 'shareTypes'
  | 'spaceId'
  | 'storageId'
  | 'type'
  | 'webDavPath'
>

type ClipboardSpaceSnapshot = Pick<SpaceResource, 'driveType' | 'id' | 'storageId' | 'webDavPath'>

type PersistedClipboardPayload = {
  action: ClipboardActions
  resources: ClipboardResourceSnapshot[]
  sourceMode: ClipboardMode
  sourceSpaces: Record<string, ClipboardSpaceSnapshot>
}

export const useClipboardStore = defineStore('clipboard', () => {
  const { $gettext } = useGettext()
  const { showMessage } = useMessages()
  const configStore = useConfigStore()
  const { getMatchingSpace } = useGetMatchingSpace()

  const action = ref<ClipboardActions>()
  const resources = ref<Resource[]>([])
  const sourceMode = ref<ClipboardMode>()
  const sourceSpaces = ref<Record<string, ClipboardSpaceSnapshot>>({})

  const currentMode = computed(() =>
    unref(configStore.isInVault) ? ClipboardMode.Vault : ClipboardMode.Default
  )

  const toClipboardSnapshot = (resource: Resource): ClipboardResourceSnapshot => ({
    id: resource.id,
    fileId: resource.fileId,
    extension: resource.extension,
    isFolder: resource.isFolder,
    name: resource.name,
    mimeType: resource.mimeType,
    parentFolderId: resource.parentFolderId,
    path: resource.path,
    remoteItemId: resource.remoteItemId,
    remoteItemPath: resource.remoteItemPath,
    shareTypes: resource.shareTypes,
    spaceId: resource.spaceId,
    storageId: resource.storageId,
    type: resource.type,
    webDavPath: resource.webDavPath
  })

  const getClipboardSourceSpaceKey = (resource: Pick<Resource, 'spaceId' | 'storageId'>) => {
    return resource.storageId || resource.spaceId
  }

  const buildSourceSpaces = (spaces: SpaceResource[] = []) =>
    spaces.reduce<Record<string, ClipboardSpaceSnapshot>>((acc, space) => {
      const key = space && getClipboardSourceSpaceKey(space)
      if (key)
        acc[key] = {
          driveType: space.driveType,
          id: space.id,
          storageId: space.storageId,
          webDavPath: space.webDavPath
        }
      return acc
    }, {})

  const removePersistedClipboard = () => {
    try {
      sessionStorage.removeItem(clipboardStorageKey)
    } catch (e) {
      console.log('error removing session item: ', e)
    }
  }

  const persistClipboard = () => {
    if (!action.value || resources.value.length === 0) {
      removePersistedClipboard()
      return
    }

    const payload: PersistedClipboardPayload = {
      action: action.value,
      resources: resources.value.map(toClipboardSnapshot),
      sourceMode: sourceMode.value,
      sourceSpaces: sourceSpaces.value
    }

    try {
      sessionStorage.setItem(clipboardStorageKey, JSON.stringify(payload))
    } catch (e) {
      console.log('error setting session item: ', e)
    }
  }

  const hydrateClipboard = () => {
    try {
      const stored = sessionStorage.getItem(clipboardStorageKey)
      if (!stored) return
      const payload = JSON.parse(stored) as PersistedClipboardPayload
      if (!payload.action || !Array.isArray(payload.resources)) {
        removePersistedClipboard()
        return
      }
      action.value = payload.action
      resources.value = payload.resources as Resource[]
      sourceMode.value = payload.sourceMode
      sourceSpaces.value = payload.sourceSpaces || {}
    } catch {
      removePersistedClipboard()
    }
  }

  const copyResources = (r: Resource[]) => {
    if (!r.length || !r[0].canDownload?.()) {
      return
    }

    action.value = ClipboardActions.Copy
    resources.value = r
    sourceMode.value = unref(currentMode)
    sourceSpaces.value = buildSourceSpaces(r.map(getMatchingSpace))
    persistClipboard()

    showMessage({ title: $gettext('Copied to clipboard!'), status: 'success' })
  }

  const cutResources = (r: Resource[]) => {
    if (!r.length || !r[0].canDownload?.()) {
      return
    }

    action.value = ClipboardActions.Cut
    resources.value = r
    sourceMode.value = unref(currentMode)
    sourceSpaces.value = buildSourceSpaces(r.map(getMatchingSpace))
    persistClipboard()

    showMessage({ title: $gettext('Cut to clipboard!'), status: 'success' })
  }

  const clearClipboard = () => {
    action.value = undefined
    resources.value = []
    sourceMode.value = undefined
    sourceSpaces.value = {}
    removePersistedClipboard()
  }

  hydrateClipboard()

  return {
    action,
    resources,
    sourceMode,
    sourceSpaces,
    getClipboardSourceSpaceKey,

    copyResources,
    cutResources,
    clearClipboard
  }
})

export type ClipboardStore = ReturnType<typeof useClipboardStore>
