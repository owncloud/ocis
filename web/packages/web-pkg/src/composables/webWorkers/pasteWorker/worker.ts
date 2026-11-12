import PQueue from 'p-queue'
import { join } from 'path'
import { webdav as _webdav, type Resource, type SpaceResource } from '@ownclouders/web-client'
import type { WorkerTopic } from '../../piniaStores/webWorkers'
import { TransferType } from '../../../helpers/resource/conflictHandling/types'
import type { TransferData } from '../../../helpers/resource/conflictHandling'
import { PasteWorkerReturnData } from './usePasteWorker'

type MessageData = {
  baseUrl?: string
  accessToken?: string
  headers?: Record<string, string>
  transferData?: TransferData[]
}

type FolderContent = {
  files: Resource[]
  folders: { folder: Resource; content: FolderContent }[]
}

type Message = {
  topic: WorkerTopic
  data: MessageData
}

let storedHeaders: Record<string, string>

self.onmessage = async (e: MessageEvent) => {
  const { topic, data } = JSON.parse(e.data) as Message

  if (topic === 'tokenUpdate' && storedHeaders) {
    const existingToken = storedHeaders.Authorization

    // token must only be updated for bearer tokens, not on public links
    if (existingToken?.toString().startsWith('Bearer')) {
      storedHeaders.Authorization = data.accessToken
    }

    return
  }

  const { baseUrl, headers, transferData } = data

  storedHeaders = headers
  const webdav = _webdav(
    baseUrl,
    () => {},
    () => storedHeaders
  )

  const successful: PasteWorkerReturnData['successful'] = []
  const failed: PasteWorkerReturnData['failed'] = []
  const queue = new PQueue({ concurrency: 4 })

  const doCopy = (transferData: TransferData) =>
    webdav.copyFiles(
      transferData.sourceSpace,
      transferData.resource,
      transferData.targetSpace,
      { path: transferData.path },
      { overwrite: transferData.overwrite }
    )

  const doMove = (transferData: TransferData) =>
    webdav.moveFiles(
      transferData.sourceSpace,
      transferData.resource,
      transferData.targetSpace,
      { path: transferData.path },
      { overwrite: transferData.overwrite }
    )

  const getErrorData = (e: any) => ({
    message: e?.message,
    statusCode: e?.statusCode,
    xReqId: e?.response?.headers?.get('x-request-id')
  })

  const collectError = (errors: unknown[]) => (e: unknown) => {
    console.error(e)
    errors.push(e)
  }

  const listChildren = async (sourceSpace: SpaceResource, folder: Resource) => {
    const listing = await queue.add(() => webdav.listFiles(sourceSpace, { path: folder.path }))
    return listing.children ?? []
  }

  /**
   * Lists the given children of a folder and all of their subfolders. Errors are collected instead
   * of thrown, a subfolder that can't be listed counts as empty.
   */
  const listFolderContent = async (
    sourceSpace: SpaceResource,
    children: Resource[],
    errors: unknown[]
  ): Promise<FolderContent> => {
    const folders = await Promise.all(
      children
        .filter((child) => child.isFolder)
        .map(async (child) => {
          let grandChildren: Resource[] = []
          try {
            grandChildren = await listChildren(sourceSpace, child)
          } catch (e) {
            collectError(errors)(e)
          }
          return {
            folder: child,
            content: await listFolderContent(sourceSpace, grandChildren, errors)
          }
        })
    )
    return { files: children.filter((child) => !child.isFolder), folders }
  }

  /**
   * Copies listed folder content with one request per file and per subfolder. Errors are collected
   * instead of thrown, so one failing item doesn't stop the rest from being copied.
   */
  const copyFolderContent = async (
    sourceSpace: SpaceResource,
    content: FolderContent,
    targetSpace: SpaceResource,
    targetPath: string,
    errors: unknown[]
  ) => {
    const copyFile = (file: Resource) =>
      queue
        .add(() =>
          webdav.copyFiles(sourceSpace, file, targetSpace, { path: join(targetPath, file.name) })
        )
        .catch(collectError(errors))

    const copySubfolder = async ({ folder, content }: FolderContent['folders'][number]) => {
      const folderTargetPath = join(targetPath, folder.name)
      try {
        await queue.add(() =>
          webdav.createFolder(targetSpace, { path: folderTargetPath, fetchFolder: false })
        )
      } catch (e) {
        collectError(errors)(e)
        return
      }
      await copyFolderContent(sourceSpace, content, targetSpace, folderTargetPath, errors)
    }

    await Promise.all([...content.files.map(copyFile), ...content.folders.map(copySubfolder)])
  }

  /**
   * Copies a folder file by file. The server would copy a whole folder within one single request,
   * which gets cut off by request timeouts for big folders. Instead, every file and every subfolder
   * gets its own request here, so every request stays short.
   * The whole folder gets listed before anything is created, so folders created by this copy can
   * never show up in the listing (e.g. when the target is reached through another space).
   * Throws if the folder itself can't be listed or created, so nothing gets created then. Resolves
   * with the errors of its content otherwise.
   */
  const copyFolderFileByFile = async (data: TransferData) => {
    const errors: unknown[] = []
    const children = await listChildren(data.sourceSpace, data.resource)
    const content = await listFolderContent(data.sourceSpace, children, errors)

    await queue.add(() =>
      webdav.createFolder(data.targetSpace, { path: data.path, fetchFolder: false })
    )

    await copyFolderContent(data.sourceSpace, content, data.targetSpace, data.path, errors)
    return errors
  }

  /** The target lies inside the source folder, e.g. "/a" pasted into "/a/b". */
  const isCopyIntoItself = (data: TransferData) => {
    const withSlash = (path: string) => `${path.replace(/\/+$/, '')}/`
    return (
      data.sourceSpace.id === data.targetSpace.id &&
      withSlash(data.path).startsWith(withSlash(data.resource.path))
    )
  }

  const transfer = async (data: TransferData) => {
    const resource = data.resource
    const isCopy =
      data.transferType === TransferType.COPY || data.transferType === TransferType.DUPLICATE

    try {
      // these stay one request: replacing a folder (the server deletes the existing folder first)
      // and copying a folder into itself (the server rejects that)
      if (isCopy && resource.isFolder && !data.overwrite && !isCopyIntoItself(data)) {
        const errors = await copyFolderFileByFile(data)
        if (errors.length) {
          // the folder exists but misses some of its content
          failed.push({ resourceName: resource.name, incomplete: true, ...getErrorData(errors[0]) })
          return
        }
      } else if (isCopy) {
        await queue.add(() => doCopy(data))
      } else {
        await queue.add(() => doMove(data))
      }

      if (isCopy) {
        resource.id = undefined
        resource.fileId = undefined
      }

      resource.path = join(data.targetFolder.path, resource.name)
      resource.webDavPath = join(data.targetFolder.webDavPath, resource.name)

      successful.push(resource)
    } catch (e) {
      console.error(e)
      failed.push({ resourceName: resource.name, ...getErrorData(e) })
    }
  }

  // every request goes through the queue by itself, a folder copy never holds a slot of the queue
  // while it waits for its content, so several folders can't block each other
  const promises = transferData.map(transfer)

  await Promise.allSettled(promises)

  postMessage(JSON.stringify({ successful, failed }))
}
