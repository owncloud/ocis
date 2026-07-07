import { uniqueId } from './uniqueId'

describe('uniqueId', () => {
  it('increases numeric part of the id independent of the chosen prefix', () => {
    expect(uniqueId()).toBe('1')
    expect(uniqueId()).toBe('2')
    expect(uniqueId('prefix-')).toBe('prefix-3')
    expect(uniqueId('another-prefix-')).toBe('another-prefix-4')
    expect(uniqueId()).toBe('5')
  })
})
