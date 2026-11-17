import { dirname } from 'path'
import { extractNameWithoutExtension, Resource, SpaceResource } from '@ownclouders/web-client'
import type { TransferFailure } from './types'

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

/**
 * What exists at the target after a transfer, so it can be shown in the files list: the resources
 * transferred completely, and the folders copied incompletely (they exist, but miss some of their
 * content).
 */
export const getResourcesAtTarget = (
  successful: Resource[],
  failed: TransferFailure[]
): { path: string; fileId?: string }[] => [
  ...successful,
  ...failed
    .filter(({ incomplete }) => incomplete)
    .map(({ incomplete }) => ({ path: incomplete.path }))
]
