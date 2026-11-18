import type { Resource, SpaceResource } from '@ownclouders/web-client'

export enum ResolveStrategy {
  SKIP,
  REPLACE,
  KEEP_BOTH,
  MERGE
}
export interface ResolveConflict {
  strategy: ResolveStrategy
  doForAllConflicts: boolean
}

export enum TransferType {
  COPY,
  MOVE,
  DUPLICATE
}

/**
 * What a folder copy that created the folder but missed some of its content reports.
 * The lists hold the original paths, starting with the original name of the copied folder.
 */
export type IncompleteCopy = {
  // where the incomplete copy is, in the target space
  path: string
  // items are the files and folders inside the copied folder, what's inside a folder that couldn't
  // be read is unknown and isn't counted
  copiedItems: number
  totalItems: number
  // the files and folders that weren't copied, folders with a trailing "/"
  notCopied: string[]
}

/**
 * A resource that couldn't be copied, moved or duplicated (completely). `error` isn't always an
 * `HttpError`, e.g. a refused overwrite of the parent folder is reported without a request.
 */
export type TransferFailure = {
  resourceName: string
  incomplete?: IncompleteCopy
  error: Error
}

export type TransferData = {
  resource: Resource
  sourceSpace: SpaceResource
  targetSpace: SpaceResource
  targetFolder: Resource
  path: string
  overwrite: boolean
  transferType: TransferType
}
