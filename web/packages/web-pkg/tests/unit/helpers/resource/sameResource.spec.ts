import { isSameResource } from '../../../../src'

describe('isSameResource', () => {
  test('evaluates to false if one of the resources is nullish', () => {
    expect(isSameResource(null, null)).toBe(false)
    expect(isSameResource(undefined, undefined)).toBe(false)
    expect(isSameResource(null, { id: '1', path: '', webDavPath: '', spaceId: '1' })).toBe(false)
    expect(isSameResource(undefined, { id: '1', path: '', webDavPath: '', spaceId: '1' })).toBe(
      false
    )
    expect(isSameResource({ id: '1', path: '', webDavPath: '', spaceId: '1' }, null)).toBe(false)
    expect(isSameResource({ id: '1', path: '', webDavPath: '', spaceId: '1' }, undefined)).toBe(
      false
    )
  })
  test('evaluates to false if ids are different values', () => {
    expect(
      isSameResource(
        { id: '1', path: '', webDavPath: '', spaceId: '1' },
        { id: '2', path: '', webDavPath: '', spaceId: '1' }
      )
    ).toBe(false)
  })
  test('evaluates to true if ids are the same', () => {
    expect(
      isSameResource(
        { id: '1', path: '', webDavPath: '', spaceId: '1' },
        { id: '1', path: '', webDavPath: '', spaceId: '1' }
      )
    ).toBe(true)
  })
})
