import isEqual from 'lodash-es/isEqual'
import { cacheService } from '../cache'
import { ClientService } from '../client'
import { encodePath } from '../../utils'
import { isPublicSpaceResource } from '@ownclouders/web-client'
import { BuildQueryStringOptions, LoadPreviewOptions } from '.'
import { AuthStore, CapabilityStore, ConfigStore, UserStore } from '../../composables'

// @ts-ignore
import { stringify } from 'qs'

const DEFAULT_SUPPORTED_MIME_TYPES = [
  'image/gif',
  'image/png',
  'image/jpeg',
  'text/plain',
  'image/tiff',
  'image/bmp',
  'image/x-ms-bmp',
  'application/vnd.geogebra.slides',
  'application/vnd.geogebra.pinboard'
]

export class PreviewService {
  clientService: ClientService
  configStore: ConfigStore
  userStore: UserStore
  authStore: AuthStore
  capabilityStore: CapabilityStore

  constructor({
    clientService,
    userStore,
    authStore,
    capabilityStore,
    configStore
  }: {
    clientService: ClientService
    userStore: UserStore
    authStore: AuthStore
    capabilityStore: CapabilityStore
    configStore: ConfigStore
  }) {
    this.clientService = clientService
    this.userStore = userStore
    this.authStore = authStore
    this.configStore = configStore
    this.capabilityStore = capabilityStore
  }

  private get available(): boolean {
    if (this.capabilityStore.filesThumbnail) {
      return !!this.capabilityStore.filesThumbnail.version
    }

    return true
  }

  private get supportedMimeTypes() {
    return this.capabilityStore.filesThumbnail?.supportedMimeTypes || DEFAULT_SUPPORTED_MIME_TYPES
  }

  private get user() {
    return this.userStore.user
  }

  public isMimetypeSupported(mimeType: string, onlyImages = false) {
    if (!this.supportedMimeTypes.length) {
      return true
    }
    const mimeTypes = this.getSupportedMimeTypes(onlyImages ? 'image/' : null)
    return mimeTypes.includes(mimeType)
  }

  public getSupportedMimeTypes(filter?: string) {
    if (!filter) {
      return this.supportedMimeTypes
    }
    return this.supportedMimeTypes.filter((mimeType) => mimeType.startsWith(filter))
  }

  public async loadPreview(
    options: LoadPreviewOptions,
    cached = false,
    silenceErrors = true,
    signal?: AbortSignal
  ): Promise<string | undefined> {
    const { space, resource } = options
    const serverSupportsPreview = this.available && this.isMimetypeSupported(resource.mimeType)
    const resourceSupportsPreview =
      resource.type !== 'folder' && resource.extension && resource.canDownload()
    if (!serverSupportsPreview || !resourceSupportsPreview) {
      return undefined
    }

    const isPublic = isPublicSpaceResource(space)
    if (
      !isPublic &&
      (!this.configStore.serverUrl ||
        !this.user.onPremisesSamAccountName ||
        !this.authStore.accessToken)
    ) {
      return undefined
    }

    if (isPublic) {
      return this.publicPreviewUrl(options, signal)
    }
    try {
      return await this.privatePreviewBlob(options, cached, silenceErrors, signal)
    } catch (e) {
      if (silenceErrors) {
        return undefined
      }
      throw e
    }
  }

  private async cacheFactory(
    options: LoadPreviewOptions,
    silenceErrors: boolean,
    signal?: AbortSignal
  ): Promise<string> {
    const { resource, dimensions } = options
    const hit = cacheService.filePreview.get(resource.id.toString())

    if (hit && hit.etag === resource.etag && isEqual(dimensions, hit.dimensions)) {
      return hit.src
    }
    try {
      const src = await this.privatePreviewBlob(options, false, true, signal)
      return cacheService.filePreview.set(
        resource.id.toString(),
        { src, etag: resource.etag, dimensions },
        0
      ).src
    } catch (e) {
      if (silenceErrors) {
        return
      }
      throw e
    }
  }

  private buildQueryString(options: BuildQueryStringOptions): string {
    return stringify({
      scalingup: options.scalingup || 0,
      preview: Object.hasOwnProperty.call(options, 'preview') ? options.preview : 1,
      a: Object.hasOwnProperty.call(options, 'a') ? options.a : 1,
      ...(options.processor && { processor: options.processor }),
      ...(options.etag && { c: options.etag.replaceAll('"', '') }),
      ...(options.dimensions && options.dimensions[0] && { x: options.dimensions[0] }),
      ...(options.dimensions && options.dimensions[1] && { y: options.dimensions[1] })
    })
  }

  private async privatePreviewBlob(
    options: LoadPreviewOptions,
    cached = false,
    silenceErrors = true,
    signal?: AbortSignal
  ): Promise<string> {
    const { resource, dimensions, processor } = options
    if (cached) {
      return this.cacheFactory(options, silenceErrors, signal)
    }

    const url = [
      this.configStore.serverUrl,
      'dav',
      encodePath(resource.webDavPath),
      '?',
      this.buildQueryString({ etag: resource.etag, dimensions, processor })
    ].join('')

    try {
      const { data } = await this.clientService.httpAuthenticated.get<Blob>(url, {
        responseType: 'blob',
        signal
      })
      return window.URL.createObjectURL(data)
    } catch (e) {
      if ([425, 429].includes(e.status)) {
        const retryAfter = e.response?.headers?.['retry-after'] || 5
        await new Promise((resolve) => setTimeout(resolve, retryAfter * 1000))
        return this.privatePreviewBlob(options, cached, silenceErrors, signal)
      }

      throw e
    }
  }

  private async publicPreviewUrl(
    options: LoadPreviewOptions,
    signal?: AbortSignal
  ): Promise<string> {
    const { resource, dimensions, processor } = options
    // In a public context, i.e. public shares, the downloadURL contains a pre-signed url to
    // download the file.
    const [url, signedQuery] = resource.downloadURL.split('?')

    // Since the pre-signed url contains query parameters and the caller of this method
    // can also provide query parameters we have to combine them.
    const combinedQuery = [
      this.buildQueryString({ etag: resource.etag, dimensions, processor }),
      signedQuery
    ]
      .filter(Boolean)
      .join('&')

    const previewUrl = [url, combinedQuery].filter(Boolean).join('?')

    try {
      const { status } = await this.clientService.httpUnAuthenticated.head(previewUrl, { signal })

      if (status !== 404) {
        return previewUrl
      }
    } catch (e) {
      if ([425, 429].includes(e.status)) {
        const retryAfter = e.response?.headers?.['retry-after'] || 5
        await new Promise((resolve) => setTimeout(resolve, retryAfter * 1000))
        return this.publicPreviewUrl(options, signal)
      }

      throw e
    }
  }
}
