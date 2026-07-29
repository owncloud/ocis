import type { Resource } from '@ownclouders/web-client'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  getBookPreferences,
  loadShelves,
  loadViewMode,
  saveBookPreferences,
  saveShelves,
  saveViewMode
} from '../../src/utils/preferences'

const values = new Map<string, string>()

beforeEach(() => {
  values.clear()
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => values.get(key) ?? null,
    setItem: (key: string, value: string) => values.set(key, value)
  })
})

describe('library preferences', () => {
  const resource = { id: 'book-1', path: '/Books/book.epub' } as Resource

  it('persists per-book organization state', () => {
    saveBookPreferences(resource, {
      favorite: true,
      readingStatus: 'reading',
      shelfIds: ['classics']
    })

    expect(getBookPreferences(resource)).toEqual({
      favorite: true,
      readingStatus: 'reading',
      shelfIds: ['classics']
    })
  })

  it('persists shelves and view mode', () => {
    saveShelves([{ id: 'classics', name: 'Classics' }])
    saveViewMode('list')

    expect(loadShelves()).toEqual([{ id: 'classics', name: 'Classics' }])
    expect(loadViewMode()).toBe('list')
  })
})
