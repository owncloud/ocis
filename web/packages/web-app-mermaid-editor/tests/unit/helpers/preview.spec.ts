import { PREVIEW_SIZE_LIMIT, isPreviewTooLarge } from '../../../src/helpers/preview'

describe('preview size guard', () => {
  it('flags content above the limit', () => {
    expect(isPreviewTooLarge('a'.repeat(PREVIEW_SIZE_LIMIT + 1))).toBe(true)
  })

  it('allows content at or below the limit', () => {
    expect(isPreviewTooLarge('a'.repeat(PREVIEW_SIZE_LIMIT))).toBe(false)
    expect(isPreviewTooLarge('')).toBe(false)
  })

  it('treats nullish content as not too large', () => {
    expect(isPreviewTooLarge(undefined as unknown as string)).toBe(false)
  })
})
