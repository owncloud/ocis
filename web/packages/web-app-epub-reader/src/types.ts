import type { Resource, SpaceResource } from '@ownclouders/web-client'

export interface EpubMetadata {
  title: string
  authors: string[]
  description: string
  language: string
  publisher: string
  published: string
  subjects: string[]
  spineItemCount: number
  coverUrl?: string
  coverBlob?: Blob
}

export interface LibraryBook extends EpubMetadata {
  id: string
  resource: Resource
  space: SpaceResource
  loadingMetadata: boolean
  metadataError?: string
  favorite: boolean
  readingStatus: ReadingStatus
  shelfIds: string[]
  readingProgress?: number
  hasReadingPosition: boolean
}

export type LibrarySort = 'recent' | 'title' | 'author'
export type ReadingStatus = 'unread' | 'reading' | 'finished'
export type LibraryViewMode = 'grid' | 'list'

export interface LibraryShelf {
  id: string
  name: string
}

export interface BookPreferences {
  favorite: boolean
  readingStatus: ReadingStatus
  shelfIds: string[]
}
