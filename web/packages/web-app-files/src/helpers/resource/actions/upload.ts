import { Language } from 'vue3-gettext'
import { Resource } from '@ownclouders/web-client'
import { extractExtensionFromFile } from '@ownclouders/web-client'
import {
  ConflictDialog,
  OcUppyFile,
  ResolveConflict,
  resolveFileNameDuplicate,
  ResolveStrategy,
  ResourceConflictModal,
  ResourcesStore,
  useModals
} from '@ownclouders/web-pkg'

interface ConflictedResource {
  name: string
  type: string
}

export class UploadResourceConflict extends ConflictDialog {
  resourcesStore: ResourcesStore

  constructor(resourcesStore: ResourcesStore, language: Language) {
    const { $gettext, $ngettext } = language
    super($gettext, $ngettext)

    this.resourcesStore = resourcesStore
  }

  resolveFileExists(
    resource: Resource,
    conflictCount: number,
    suggestMerge = false,
    separateSkipHandling = false // separate skip-handling between files and folders
  ): Promise<ResolveConflict> {
    const { dispatchModal } = useModals()

    return new Promise<ResolveConflict>((resolve) => {
      dispatchModal({
        variation: 'danger',
        title: resource.isFolder
          ? this.$gettext('Folder already exists')
          : this.$gettext('File already exists'),
        hideActions: true,
        customComponent: ResourceConflictModal,
        customComponentAttrs: () => ({
          confirmSecondaryTextOverwrite: resource.isFolder
            ? this.$gettext('Merge')
            : this.$gettext('Replace'),
          resource,
          conflictCount,
          suggestMerge,
          separateSkipHandling,
          callbackFn: (conflict: ResolveConflict) => {
            resolve(conflict)
          }
        })
      })
    })
  }

  getConflicts(files: OcUppyFile[]): ConflictedResource[] {
    const conflicts: ConflictedResource[] = []
    for (const file of files) {
      const relativeFilePath = file.meta.relativePath
      if (relativeFilePath) {
        // Logic for folders, applies to all files inside folder and subfolders
        const rootFolder = relativeFilePath.replace(/^\/+/, '').split('/')[0]
        const exists = this.resourcesStore.resources.find((f) => f.name === rootFolder)
        if (exists) {
          if (conflicts.some((conflict) => conflict.name === rootFolder)) {
            continue
          }
          conflicts.push({ name: rootFolder, type: 'folder' })
          continue
        }
      }
      // Logic for files
      const exists = this.resourcesStore.resources.find(
        (f) => f.name === file.name && !file.meta.relativeFolder
      )
      if (exists) {
        conflicts.push({ name: file.name, type: 'file' })
      }
    }
    return conflicts
  }

  async displayOverwriteDialog(
    files: OcUppyFile[],
    conflicts: ConflictedResource[]
  ): Promise<OcUppyFile[]> {
    let fileCount = 0
    let folderCount = 0
    const resolvedFileConflicts: { name: string; strategy: ResolveStrategy }[] = []
    const resolvedFolderConflicts: { name: string; strategy: ResolveStrategy }[] = []
    let doForAllConflicts = false
    let allConflictsStrategy
    let doForAllConflictsFolders = false
    let allConflictsStrategyFolders

    for (const conflict of conflicts) {
      const isFolder = conflict.type === 'folder'
      const conflictArray = isFolder ? resolvedFolderConflicts : resolvedFileConflicts

      if (doForAllConflicts && !isFolder) {
        conflictArray.push({
          name: conflict.name,
          strategy: allConflictsStrategy
        })
        continue
      }
      if (doForAllConflictsFolders && isFolder) {
        conflictArray.push({
          name: conflict.name,
          strategy: allConflictsStrategyFolders
        })
        continue
      }

      const conflictsLeft =
        conflicts.filter((c) => c.type === conflict.type).length -
        (isFolder ? folderCount : fileCount)

      const resolvedConflict: ResolveConflict = await this.resolveFileExists(
        { name: conflict.name, isFolder } as Resource,
        conflictsLeft,
        isFolder,
        true
      )
      isFolder ? folderCount++ : fileCount++
      if (resolvedConflict.doForAllConflicts) {
        if (isFolder) {
          doForAllConflictsFolders = true
          allConflictsStrategyFolders = resolvedConflict.strategy
        } else {
          doForAllConflicts = true
          allConflictsStrategy = resolvedConflict.strategy
        }
      }

      conflictArray.push({
        name: conflict.name,
        strategy: resolvedConflict.strategy
      })
    }
    const filesToSkip = resolvedFileConflicts
      .filter((e) => e.strategy === ResolveStrategy.SKIP)
      .map((e) => e.name)
    const foldersToSkip = resolvedFolderConflicts
      .filter((e) => e.strategy === ResolveStrategy.SKIP)
      .map((e) => e.name)

    files = files.filter((e) => !filesToSkip.includes(e.name))
    files = files.filter(
      (file) =>
        !foldersToSkip.some((folderName) => file.meta.relativeFolder.split('/')[1] === folderName)
    )

    const filesToKeepBoth = resolvedFileConflicts
      .filter((e) => e.strategy === ResolveStrategy.KEEP_BOTH)
      .map((e) => e.name)
    const foldersToKeepBoth = resolvedFolderConflicts
      .filter((e) => e.strategy === ResolveStrategy.KEEP_BOTH)
      .map((e) => e.name)

    for (const fileName of filesToKeepBoth) {
      const file = files.find((e) => e.name === fileName && !e.meta.relativeFolder)
      const extension = extractExtensionFromFile({ name: fileName } as Resource)
      file.name = resolveFileNameDuplicate(fileName, extension, this.resourcesStore.resources)
      file.meta.name = file.name
      if (file.xhrUpload?.endpoint) {
        const endpoint =
          typeof file.xhrUpload.endpoint === 'function'
            ? await file.xhrUpload.endpoint(file)
            : file.xhrUpload.endpoint

        file.xhrUpload.endpoint = endpoint.replace(
          new RegExp(`/${encodeURIComponent(fileName)}`),
          `/${encodeURIComponent(file.name)}`
        )
      }
    }
    for (const folder of foldersToKeepBoth) {
      const filesInFolder = files.filter((e) => e.meta.relativeFolder.split('/')[1] === folder)
      for (const file of filesInFolder) {
        const newFolderName = resolveFileNameDuplicate(folder, '', this.resourcesStore.resources)
        file.meta.relativeFolder = file.meta.relativeFolder.replace(
          new RegExp(`/${folder}`),
          `/${newFolderName}`
        )
        file.meta.relativePath = file.meta.relativePath.replace(
          new RegExp(`/${folder}/`),
          `/${newFolderName}/`
        )
        file.meta.tusEndpoint = file.meta.tusEndpoint.replace(
          new RegExp(`/${encodeURIComponent(folder)}$`),
          `/${encodeURIComponent(newFolderName)}`
        )
        if (file.xhrUpload?.endpoint) {
          const endpoint =
            typeof file.xhrUpload.endpoint === 'function'
              ? await file.xhrUpload.endpoint(file)
              : file.xhrUpload.endpoint

          file.xhrUpload.endpoint = endpoint.replace(
            new RegExp(`/${encodeURIComponent(folder)}$`),
            `/${encodeURIComponent(newFolderName)}`
          )
        }
        if (file.tus?.endpoint) {
          file.tus.endpoint = file.tus.endpoint.replace(
            new RegExp(`/${encodeURIComponent(folder)}$`),
            `/${encodeURIComponent(newFolderName)}`
          )
        }
      }
    }
    return files
  }
}
