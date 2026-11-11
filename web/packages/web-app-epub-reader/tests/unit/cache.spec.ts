import type { Resource } from '@ownclouders/web-client'
import { describe, expect, it } from 'vitest'
import { resourceCacheId, resourceFingerprint } from '../../src/utils/cache'

function resource(overrides: Partial<Resource> = {}): Resource {
  return {
    id: 'resource-id',
    name: 'book.epub',
    path: '/Books/book.epub',
    type: 'file',
    ...overrides
  } as Resource
}

describe('EPUB metadata cache keys', () => {
  it('prefers the stable file id', () => {
    expect(resourceCacheId(resource({ fileId: 'file-id' }))).toBe('file-id')
  })

  it('changes the fingerprint when file metadata changes', () => {
    const original = resourceFingerprint(resource({ etag: 'one', size: 100 }))
    const changed = resourceFingerprint(resource({ etag: 'two', size: 100 }))

    expect(changed).not.toBe(original)
  })
})
