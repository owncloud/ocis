import { cacheService, ClientService } from '@ownclouders/web-pkg'
import { ImageDimension } from '@ownclouders/web-pkg'

interface AvatarUrlOptions {
  clientService: ClientService
  server: string
  username: string
  size?: number
}

export const avatarUrl = async (options: AvatarUrlOptions, cached = false): Promise<string> => {
  const size = options.size || ImageDimension.Avatar

  if (cached) {
    return cacheFactory({ ...options, size })
  }

  const url = [options.server, 'dav/avatars/', options.username, `/${size}.png`].join('')

  const { status, statusText } = await options.clientService.httpAuthenticated.head(url)

  if (status !== 200) {
    throw new Error(statusText)
  }

  return options.clientService.ocs.signUrl({ url, username: options.username })
}

const cacheFactory = async (options: AvatarUrlOptions): Promise<string> => {
  const hit = cacheService.avatarUrl.get(options.username)
  if (hit && hit.size === options.size) {
    return hit.src
  }

  try {
    const src = await avatarUrl(options)
    return cacheService.avatarUrl.set(options.username, { src, size: options.size }, 0).src
  } catch {}
}
