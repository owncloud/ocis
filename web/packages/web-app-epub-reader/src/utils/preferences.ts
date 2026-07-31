import type { Resource } from '@ownclouders/web-client'
import type { BookPreferences, LibraryShelf, LibraryViewMode } from '../types'
import { resourceCacheId } from './cache'

const STORAGE_KEY = 'ocis-epub-library-preferences-v1'

interface StoredPreferences {
  books: Record<string, BookPreferences>
  shelves: LibraryShelf[]
  viewMode: LibraryViewMode
}

const defaults = (): StoredPreferences => ({ books: {}, shelves: [], viewMode: 'grid' })

function read(): StoredPreferences {
  if (typeof localStorage === 'undefined') return defaults()
  try {
    return { ...defaults(), ...JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}') }
  } catch {
    return defaults()
  }
}

function write(value: StoredPreferences): void {
  if (typeof localStorage === 'undefined') return
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value))
  } catch {
    // Preferences remain available for this session when browser storage is disabled.
  }
}

export function getBookPreferences(resource: Resource): BookPreferences {
  return {
    favorite: false,
    readingStatus: 'unread',
    shelfIds: [],
    ...read().books[resourceCacheId(resource)]
  }
}

export function saveBookPreferences(resource: Resource, preferences: BookPreferences): void {
  const stored = read()
  stored.books[resourceCacheId(resource)] = preferences
  write(stored)
}

export function loadShelves(): LibraryShelf[] {
  return read().shelves
}

export function saveShelves(shelves: LibraryShelf[]): void {
  const stored = read()
  stored.shelves = shelves
  write(stored)
}

export function loadViewMode(): LibraryViewMode {
  return read().viewMode
}

export function saveViewMode(viewMode: LibraryViewMode): void {
  const stored = read()
  stored.viewMode = viewMode
  write(stored)
}
