import type { Resource } from '@ownclouders/web-client'
import type { EpubMetadata } from '../types'

const DATABASE_NAME = 'ocis-epub-library'
const DATABASE_VERSION = 1
const STORE_NAME = 'metadata'

interface CachedMetadataRecord {
  id: string
  fingerprint: string
  metadata: Omit<EpubMetadata, 'coverBlob' | 'coverUrl'>
  coverBlob?: Blob
}

export function resourceCacheId(resource: Resource): string {
  return resource.fileId || resource.id || `${resource.spaceId}:${resource.path}`
}

export function resourceFingerprint(resource: Resource): string {
  return [resource.etag || '', resource.mdate || '', String(resource.size || 0)].join(':')
}

function openDatabase(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DATABASE_NAME, DATABASE_VERSION)
    request.onupgradeneeded = () => {
      if (!request.result.objectStoreNames.contains(STORE_NAME)) {
        request.result.createObjectStore(STORE_NAME, { keyPath: 'id' })
      }
    }
    request.onsuccess = () => resolve(request.result)
    request.onerror = () => reject(request.error)
  })
}

export async function getCachedMetadata(resource: Resource): Promise<EpubMetadata | null> {
  if (typeof indexedDB === 'undefined') return null
  try {
    const database = await openDatabase()
    const record = await new Promise<CachedMetadataRecord | undefined>((resolve, reject) => {
      const request = database
        .transaction(STORE_NAME, 'readonly')
        .objectStore(STORE_NAME)
        .get(resourceCacheId(resource))
      request.onsuccess = () => resolve(request.result as CachedMetadataRecord | undefined)
      request.onerror = () => reject(request.error)
    })
    database.close()
    if (
      !record ||
      record.fingerprint !== resourceFingerprint(resource) ||
      typeof record.metadata.spineItemCount !== 'number'
    ) {
      return null
    }
    return {
      ...record.metadata,
      coverBlob: record.coverBlob,
      coverUrl: record.coverBlob ? URL.createObjectURL(record.coverBlob) : undefined
    }
  } catch {
    return null
  }
}

export async function setCachedMetadata(resource: Resource, metadata: EpubMetadata): Promise<void> {
  if (typeof indexedDB === 'undefined') return
  try {
    const database = await openDatabase()
    const persistedMetadata: CachedMetadataRecord['metadata'] = {
      authors: metadata.authors,
      description: metadata.description,
      language: metadata.language,
      published: metadata.published,
      publisher: metadata.publisher,
      subjects: metadata.subjects,
      spineItemCount: metadata.spineItemCount,
      title: metadata.title
    }
    const record: CachedMetadataRecord = {
      id: resourceCacheId(resource),
      fingerprint: resourceFingerprint(resource),
      metadata: persistedMetadata,
      coverBlob: metadata.coverBlob
    }
    await new Promise<void>((resolve, reject) => {
      const transaction = database.transaction(STORE_NAME, 'readwrite')
      transaction.objectStore(STORE_NAME).put(record)
      transaction.oncomplete = () => resolve()
      transaction.onerror = () => reject(transaction.error)
    })
    database.close()
  } catch (e) {
    console.error('Failed to set cached metadata', e)
  }
}
