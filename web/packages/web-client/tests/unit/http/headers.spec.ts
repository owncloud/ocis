import { httpHeaders } from '../../../src/http'

describe('httpHeaders', () => {
  const headers = () =>
    httpHeaders(new Headers({ ETag: '"abc"', 'Retry-After': '30', 'Lock-Token': '<token>' }))

  it('reads a header through get, as a native Headers does', () => {
    expect(headers().get('etag')).toBe('"abc"')
    expect(headers().get('nope')).toBeNull()
  })

  it('reads a header through bracket access, as axios allowed', () => {
    expect(headers()['etag']).toBe('"abc"')
    expect(headers()['retry-after']).toBe('30')
  })

  it('ignores case in bracket access', () => {
    expect(headers()['Lock-Token']).toBe('<token>')
    expect(headers()['LOCK-TOKEN']).toBe('<token>')
  })

  it('resolves an absent header to undefined rather than null', () => {
    expect(headers()['nope']).toBeUndefined()
  })

  it('enumerates the lowercased header names', () => {
    expect(Object.keys(headers()).sort()).toEqual(['etag', 'lock-token', 'retry-after'])
    expect({ ...headers() }).toEqual({
      etag: '"abc"',
      'lock-token': '<token>',
      'retry-after': '30'
    })
  })

  it('answers the in operator for a header and for a Headers method', () => {
    expect('etag' in headers()).toBe(true)
    expect('nope' in headers()).toBe(false)
    expect('get' in headers()).toBe(true)
  })

  it('stays a Headers', () => {
    const wrapped = headers()

    expect(wrapped).toBeInstanceOf(Headers)
    expect(wrapped.has('etag')).toBe(true)
    expect([...wrapped]).toHaveLength(3)
  })

  it('is not mistaken for a thenable', async () => {
    await expect(Promise.resolve(headers())).resolves.toBeInstanceOf(Headers)
  })
})
