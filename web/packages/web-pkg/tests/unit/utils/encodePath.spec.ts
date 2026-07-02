import { encodePath } from '../../../src/utils'

describe('encodePath', () => {
  it('is empty if input is empty', () => {
    expect(encodePath()).toBe('')
    expect(encodePath('')).toBe('')
  })

  it('encodes a uri', () => {
    expect(encodePath('xyz.com')).toBe('xyz.com')
    expect(encodePath('xyz.com/')).toBe('xyz.com/')
    expect(encodePath('xyz.com/?a=12&b=55')).toBe('xyz.com/%3Fa%3D12%26b%3D55')
  })
})
