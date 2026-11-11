import type { Resource } from '@ownclouders/web-client'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { getReaderProgress } from '../../src/utils/readerProgress'

const values = new Map<string, string>()
const resource = { id: 'book-id' } as Resource

beforeEach(() => {
  values.clear()
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => values.get(key) ?? null
  })
})

describe('EPUB reader progress', () => {
  it('reads the location saved by the built-in EPUB reader', () => {
    values.set(
      'oc_epubReader_resource_book-id',
      JSON.stringify({
        currentLocation: {
          start: { cfi: 'epubcfi(/6/4!/4/2)', percentage: 0.42 }
        }
      })
    )

    expect(getReaderProgress(resource)).toEqual({
      finished: false,
      hasPosition: true,
      percentage: 42
    })
  })

  it('returns an empty state for books that have not been opened', () => {
    expect(getReaderProgress(resource)).toEqual({ finished: false, hasPosition: false })
  })

  it('estimates progress when EPUB.js leaves its percentage at zero', () => {
    values.set(
      'oc_epubReader_resource_book-id',
      JSON.stringify({
        currentLocation: {
          atStart: false,
          start: {
            cfi: 'epubcfi(/6/8!/4/2)',
            index: 2,
            percentage: 0,
            displayed: { page: 3, total: 10 }
          }
        }
      })
    )

    expect(getReaderProgress(resource, 10)).toEqual({
      finished: false,
      hasPosition: true,
      percentage: 22
    })
  })

  it('marks a location at the end as finished', () => {
    values.set(
      'oc_epubReader_resource_book-id',
      JSON.stringify({
        currentLocation: {
          atEnd: true,
          start: { cfi: 'epubcfi(/6/20!/4/2)', index: 9, percentage: 0 }
        }
      })
    )

    expect(getReaderProgress(resource, 10)).toEqual({
      finished: true,
      hasPosition: true,
      percentage: 100
    })
  })
})
