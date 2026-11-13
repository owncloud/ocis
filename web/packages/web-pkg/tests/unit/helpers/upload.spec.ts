import {
  getUploadMtimeSeconds,
  getUploadMtimeHeaderValue,
  getUploadMtimeMeta,
  getUploadMtimeHeader
} from '../../../src/helpers'

describe('getUploadMtimeSeconds', () => {
  it.each([
    [{ lastModified: 1000 } as File, 1],
    [{ lastModified: 1234000 } as File, 1234],
    [{ lastModified: 1 } as File, 0.001]
  ])('returns lastModified %s in seconds as %s', (data, expected) => {
    expect(getUploadMtimeSeconds(data)).toEqual(expected)
  })

  it.each([
    [undefined, 'undefined'],
    [new Blob(), 'a Blob without lastModified'],
    [{ lastModified: NaN } as File, 'NaN'],
    [{ lastModified: 0 } as File, 'zero'],
    [{ lastModified: Infinity } as File, 'Infinity'],
    [{ lastModified: -1000 } as File, 'negative']
  ])('returns undefined when lastModified is %s (%s)', (data) => {
    expect(getUploadMtimeSeconds(data)).toBeUndefined()
  })
})

describe('getUploadMtimeHeaderValue', () => {
  it.each([
    [{ lastModified: 1000 } as File, '1'],
    [{ lastModified: 1234567 } as File, '1235'],
    [{ lastModified: 1500 } as File, '2']
  ])('returns lastModified %s as the rounded integer-seconds string %s', (data, expected) => {
    expect(getUploadMtimeHeaderValue(data)).toEqual(expected)
  })

  it.each([
    [undefined, 'undefined'],
    [new Blob(), 'a Blob without lastModified'],
    [{ lastModified: NaN } as File, 'NaN'],
    [{ lastModified: 0 } as File, 'zero'],
    [{ lastModified: Infinity } as File, 'Infinity'],
    [{ lastModified: -1000 } as File, 'negative']
  ])('returns undefined when lastModified is %s (%s)', (data) => {
    expect(getUploadMtimeHeaderValue(data)).toBeUndefined()
  })
})

describe('getUploadMtimeMeta', () => {
  it('returns a mtime key in seconds for a File with a valid lastModified', () => {
    expect(getUploadMtimeMeta({ lastModified: 1000 } as File)).toEqual({ mtime: 1 })
  })

  it.each([
    [undefined, 'undefined'],
    [new Blob(), 'a Blob without lastModified'],
    [{ lastModified: 0 } as File, 'zero']
  ])('returns an object without a mtime key when lastModified is %s (%s)', (data) => {
    const result = getUploadMtimeMeta(data)
    expect(result).toEqual({})
    expect('mtime' in result).toBe(false)
  })
})

describe('getUploadMtimeHeader', () => {
  it("returns an 'x-oc-mtime' key with the rounded integer-seconds string for a File with a valid lastModified", () => {
    expect(getUploadMtimeHeader({ lastModified: 1000 } as File)).toEqual({ 'x-oc-mtime': '1' })
  })

  it.each([
    [undefined, 'undefined'],
    [new Blob(), 'a Blob without lastModified'],
    [{ lastModified: 0 } as File, 'zero']
  ])("returns an object without an 'x-oc-mtime' key when lastModified is %s (%s)", (data) => {
    const result = getUploadMtimeHeader(data)
    expect(result).toEqual({})
    expect('x-oc-mtime' in result).toBe(false)
  })
})
