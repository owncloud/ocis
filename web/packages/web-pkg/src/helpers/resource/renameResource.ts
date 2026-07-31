import { basename, join } from 'path'
import { Resource, SpaceResource, extractExtensionFromFile } from '@ownclouders/web-client'

export function renameResource(space: SpaceResource, resource: Resource, newPath: string) {
  resource.name = basename(newPath)
  resource.path = newPath
  resource.webDavPath = join(space.webDavPath, newPath)
  resource.extension = extractExtensionFromFile(resource)
  return resource
}
