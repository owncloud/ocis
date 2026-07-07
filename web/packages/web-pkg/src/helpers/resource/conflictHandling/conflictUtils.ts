import { dirname } from 'path'
import { extractNameWithoutExtension, Resource, SpaceResource } from '@ownclouders/web-client'

export const resolveFileNameDuplicate = (
  name: string,
  extension: string,
  existingResources: Resource[],
  iteration = 1
): string => {
  let potentialName
  if (!extension) {
    potentialName = `${name} (${iteration})`
  } else {
    const nameWithoutExtension = extractNameWithoutExtension({ name, extension } as Resource)
    potentialName = `${nameWithoutExtension} (${iteration}).${extension}`
  }
  const hasConflict = existingResources.some((f) => f.name === potentialName)
  if (!hasConflict) {
    return potentialName
  }
  return resolveFileNameDuplicate(name, extension, existingResources, iteration + 1)
}

export const isResourceBeeingMovedToSameLocation = (
  sourceSpace: SpaceResource,
  resource: Resource,
  targetSpace: SpaceResource,
  targetFolder: Resource
) => {
  return sourceSpace.id === targetSpace.id && dirname(resource.path) === targetFolder.path
}
