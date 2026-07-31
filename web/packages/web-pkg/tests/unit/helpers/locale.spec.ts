import { getLocaleFromLanguage } from '../../../src/helpers'

describe('getLocaleFromLanguage', () => {
  describe('returns empty string for invalid input', () => {
    it.each([
      [null, ''],
      [undefined, ''],
      ['', '']
    ])('invalid input "%s"', (input: string, expected: string) => {
      expect(getLocaleFromLanguage(input)).toBe(expected)
    })
  })
  describe('returns single segment locale from input', () => {
    it.each([
      ['de', 'de'],
      ['de_DE', 'de']
    ])('input "%s"', (input, expected) => {
      expect(getLocaleFromLanguage(input)).toBe(expected)
    })
  })
})
