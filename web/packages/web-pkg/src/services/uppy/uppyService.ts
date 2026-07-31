import Uppy, { BasePlugin, UnknownPlugin, UppyFile } from '@uppy/core'
import Tus from '@uppy/tus'
import { TusOptions } from '@uppy/tus'
import XHRUpload, { XHRUploadOptions } from '@uppy/xhr-upload'
import { Language } from 'vue3-gettext'
import { eventBus } from '../eventBus'
import DropTarget from '@uppy/drop-target'
import { Resource, urlJoin } from '@ownclouders/web-client'
import { Body, generateFileID, MinimalRequiredUppyFile, NetworkError } from '@uppy/utils'

type UppyServiceTopics =
  | 'uploadStarted'
  | 'uploadCancelled'
  | 'uploadCompleted'
  | 'uploadSuccess'
  | 'uploadError'
  | 'filesSelected'
  | 'progress'
  | 'addedForUpload'
  | 'upload-progress'
  | 'drag-over'
  | 'drag-out'
  | 'drop'

export type uppyHeaders = {
  [name: string]: string | number
}

export type UploadResult = {
  successful: OcUppyFile[]
  failed: OcUppyFile[]
  uploadID?: string
}

interface UppyServiceOptions {
  language: Language
}

type FileWithPath = File & {
  readonly path?: string
  readonly relativePath?: string
}

// FIXME: tus error types seem to be wrong in Uppy, we need the type of the tus client lib
type TusClientError = Error & { originalResponse: any }

// IMPORTANT: must only contain primitive types, complex types won't be serialized properly!
export type OcUppyMeta = {
  retry?: boolean
  name?: string
  mtime?: number
  // current space & folder
  spaceId: string
  spaceName: string
  driveAlias: string
  driveType: string
  currentFolder: string // current folder path during upload initiation
  currentFolderId?: string
  fileId?: string
  // upload data
  uppyId?: string
  relativeFolder: string
  relativePath: string
  tusEndpoint: string
  uploadId: string
  topLevelFolderId?: string
  // route data
  routeName?: string
  routeDriveAliasAndItem?: string
  routeShareId?: string

  isFolder: boolean
}
export type OcUppyBody = Body
export type OcUppyFile = UppyFile<OcUppyMeta, OcUppyBody> & { isFolder?: boolean; spaceId: string }
type OcUppyPlugin = typeof BasePlugin<any, OcUppyMeta, OcUppyBody>
export type OcMinimalUppyFile = MinimalRequiredUppyFile<OcUppyMeta, OcUppyBody>

export type OcTusOptions = TusOptions<OcUppyMeta, OcUppyBody>

/** `OmitFirstArg<typeof someArray>` is the type of the returned value of `someArray.slice(1)`. */
type OmitFirstArg<T> = T extends [any, ...infer U] ? U : never

export class UppyService {
  uppy: Uppy<OcUppyMeta, OcUppyBody>
  uploadInputs: HTMLInputElement[] = []
  uploadFolderMap: Record<string, Resource> = {}

  constructor({ language }: UppyServiceOptions) {
    const { $gettext } = language
    this.uppy = new Uppy<OcUppyMeta, OcUppyBody>({
      autoProceed: false,
      onBeforeFileAdded: (file, files) => {
        if (file.id in files) {
          file.meta.retry = true
        }
        file.meta.relativePath = this.getRelativeFilePath({ ...file, spaceId: file.meta.spaceId })
        // id needs to be generated after the relative path has been set.
        file.id = generateFileID(file, this.uppy.getID())
        return file
      }
    })

    // FIXME: move to importer plugin, as strings are only visible in the dashboard anyhow
    this.uppy.setOptions({
      locale: {
        strings: {
          addedNumFiles: $gettext('Added %{numFiles} file(s)'), // for some reason this string is required and missing in uppy
          authenticate: $gettext('Connect'),
          authenticateWith: $gettext('Connect to %{pluginName}'),
          authenticateWithTitle: $gettext('Please authenticate with %{pluginName} to select files'),
          cancel: $gettext('Cancel'),
          companionError: $gettext('Connection with Companion failed'),
          loadedXFiles: $gettext('Loaded %{numFiles} files'),
          loading: $gettext('Loading...'),
          logOut: $gettext('Log out'),
          pluginWebdavInputLabel: $gettext('Public link without password protection'),
          selectX: {
            0: $gettext('Select %{smart_count}'),
            1: $gettext('Select %{smart_count}')
          },
          signInWithGoogle: $gettext('Sign in with Google')
        }
      }
    })

    this.setUpEvents()
  }

  getRelativeFilePath = (file: OcUppyFile): string | undefined => {
    const relativePath =
      (file.data as FileWithPath).relativePath || (file.data as File).webkitRelativePath
    return relativePath ? urlJoin(relativePath) : undefined
  }

  addPlugin<T extends OcUppyPlugin>(
    Plugin: T,
    // We want to let the plugin decide whether `opts` is optional or not
    // so we spread the argument rather than defining `opts:` ourselves.
    ...args: OmitFirstArg<ConstructorParameters<T>>
  ) {
    this.uppy.use(Plugin, ...args)
  }

  removePlugin(plugin: UnknownPlugin<OcUppyMeta, OcUppyBody>) {
    this.uppy.removePlugin(plugin)
  }

  getPlugin<
    T extends UnknownPlugin<OcUppyMeta, OcUppyBody> = UnknownPlugin<OcUppyMeta, OcUppyBody>
  >(name: string): T | undefined {
    return this.uppy.getPlugin(name)
  }

  useTus({
    chunkSize,
    overridePatchMethod,
    uploadDataDuringCreation,
    onBeforeRequest,
    headers
  }: TusOptions<OcUppyMeta, OcUppyBody>) {
    const tusPluginOptions: TusOptions<OcUppyMeta, OcUppyBody> = {
      chunkSize,
      removeFingerprintOnSuccess: true,
      overridePatchMethod,
      retryDelays: [0, 500, 1000],
      uploadDataDuringCreation,
      limit: 5,
      headers,
      onBeforeRequest,
      onShouldRetry: (err, retryAttempt, options, next) => {
        // status code 5xx means the upload is gone on the server side
        if ((err as TusClientError)?.originalResponse?.getStatus() >= 500) {
          return false
        }
        if ((err as TusClientError)?.originalResponse?.getStatus() === 401) {
          return true
        }
        return next(err)
      },
      onAfterResponse(_, res) {
        const status = res.getStatus()
        if (status >= 500 && res.getHeader('content-type')?.includes('text/html')) {
          throw new NetworkError(`Server error (${status}) - Please try again`)
        }
      }
    }

    const xhrPlugin = this.uppy.getPlugin('XHRUpload')
    if (xhrPlugin) {
      this.uppy.removePlugin(xhrPlugin)
    }

    const tusPlugin = this.uppy.getPlugin('Tus')
    if (tusPlugin) {
      tusPlugin.setOptions(tusPluginOptions)
      return
    }

    this.uppy.use(Tus, tusPluginOptions)
  }

  useXhr({ headers, timeout, endpoint }: XHRUploadOptions<OcUppyMeta, OcUppyBody>) {
    const xhrPluginOptions: XHRUploadOptions<OcUppyMeta, OcUppyBody> = {
      endpoint,
      method: 'put',
      headers,
      formData: false,
      timeout,
      getResponseData() {
        return {}
      },
      onAfterResponse(xhr) {
        if (xhr.status >= 500 && xhr.getResponseHeader('content-type')?.includes('text/html')) {
          throw new NetworkError(`Server error (${xhr.status}) - Please try again`, xhr)
        }
      }
    }

    const tusPlugin = this.uppy.getPlugin('Tus')
    if (tusPlugin) {
      this.uppy.removePlugin(tusPlugin)
    }

    const xhrPlugin = this.uppy.getPlugin('XHRUpload')
    if (xhrPlugin) {
      xhrPlugin.setOptions(xhrPluginOptions)
      return
    }

    this.uppy.use(XHRUpload, xhrPluginOptions)
  }

  tusActive() {
    return !!this.uppy.getPlugin('Tus')
  }

  useDropTarget({ targetSelector }: { targetSelector: string }) {
    if (this.uppy.getPlugin('DropTarget')) {
      return
    }
    this.uppy.use(DropTarget, {
      target: targetSelector,
      onDragOver: (event) => {
        this.publish('drag-over', event)
      },
      onDragLeave: (event) => {
        this.publish('drag-out', event)
      },
      onDrop: (event) => {
        this.publish('drop', event)
      }
    })
  }

  removeDropTarget() {
    const dropTargetPlugin = this.uppy.getPlugin('DropTarget')
    if (dropTargetPlugin) {
      this.uppy.removePlugin(dropTargetPlugin)
    }
  }

  subscribe<T>(topic: UppyServiceTopics, callback: (data?: T) => void): string {
    return eventBus.subscribe(topic, callback)
  }

  unsubscribe(topic: UppyServiceTopics, token: string): void {
    eventBus.unsubscribe(topic, token)
  }

  publish(topic: UppyServiceTopics, data?: unknown): void {
    eventBus.publish(topic, data)
  }

  private setUpEvents() {
    this.uppy.on('progress', (value) => {
      this.publish('progress', value)
    })
    this.uppy.on('upload-progress', (file, progress) => {
      this.publish('upload-progress', { file, progress })
    })
    this.uppy.on('cancel-all', () => {
      this.publish('uploadCancelled')
      this.clearInputs()
    })
    this.uppy.on('complete', (result) => {
      if (!result || result.successful.length === 0) {
        return
      }

      this.publish('uploadCompleted', result)
      result.successful.forEach((file: any) => {
        this.uppy.removeFile(file.id)
      })
      this.clearInputs()
    })
    this.uppy.on('upload-success', (file) => {
      this.publish('uploadSuccess', file)
    })
    this.uppy.on('upload-error', (file, error) => {
      this.publish('uploadError', { file, error })
    })
  }

  registerUploadInput(el: HTMLInputElement) {
    const listenerRegistered = el.getAttribute('listener')
    if (listenerRegistered !== 'true') {
      el.setAttribute('listener', 'true')
      el.addEventListener('change', (event) => {
        const target = event.target as HTMLInputElement
        const files = Array.from(target.files)
        this.addFiles(files)
      })
      this.uploadInputs.push(el)
    }
  }

  removeUploadInput(el: HTMLInputElement) {
    this.uploadInputs = this.uploadInputs.filter((input) => input !== el)
  }

  generateUploadId(uppyFile: OcUppyFile): string {
    return generateFileID(uppyFile, this.uppy.getID())
  }

  addFiles(files: OcMinimalUppyFile[] | File[]) {
    // uppy types say they do not accept File[] but they are wrong
    this.uppy.addFiles(files as OcMinimalUppyFile[])
  }

  uploadFiles() {
    return this.uppy.upload()
  }

  retryAllUploads() {
    return this.uppy.retryAll()
  }

  pauseAllUploads() {
    return this.uppy.pauseAll()
  }

  resumeAllUploads() {
    return this.uppy.resumeAll()
  }

  cancelAllUploads() {
    return this.uppy.cancelAll()
  }

  getCurrentUploads(): Record<string, unknown> {
    return this.uppy.getState().currentUploads
  }

  isRemoteUploadInProgress(): boolean {
    return this.uppy.getFiles().some((f) => f.isRemote && !f.error)
  }

  clearInputs() {
    this.uploadInputs.forEach((item) => {
      item.value = null
    })
  }

  /**
   * Set a specific upload folder for an upload. The HandleUpload plugin
   * checks, if a specific folder has been specified as upload destination.
   * If not, it falls back to the current folder.
   * The uploadId needs to be set within the meta object of the upload files
   * for the plugin to connect an upload to its destination folder.
   **/
  setUploadFolder(uploadId: string, folder: Resource) {
    this.uploadFolderMap[uploadId] = folder
  }

  removeUploadFolder(uploadId: string) {
    if (this.uploadFolderMap[uploadId]) {
      delete this.uploadFolderMap[uploadId]
    }
  }
}
