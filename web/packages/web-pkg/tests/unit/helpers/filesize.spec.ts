import { formatFileSize } from '../../../src/helpers'

describe('formatFileSize', () => {
  describe('converts numeric input to a human readable format', () => {
    it.each([
      [0, '0 B'],
      [1, '1 B'],
      [1023, '1 kB'],
      [1024, '1 kB'],
      [1287654323, '1.3 GB']
    ])('input "%s"', (input: number, expected: string) => {
      expect(formatFileSize(input, '')).toEqual(expected)
    })
  })
  describe('converts string input to a human readable format', () => {
    it.each([
      ['0', '0 B'],
      ['1287654323', '1.3 GB']
    ])('input "%s"', (input: string, expected: string) => {
      expect(formatFileSize(input, '')).toEqual(expected)
    })
  })
  describe('rounds values', () => {
    it('to integers if file size below 1 MB', () => {
      expect(formatFileSize(4321, '')).toBe('4 kB')
    })
    it('to 1 decimal if file size above 1 MB', () => {
      expect(formatFileSize(4321000, '')).toBe('4.3 MB')
    })
  })
  describe('handles invalid input', () => {
    it.each([
      [-1, '--'],
      ['nonNumericInput', '?'],
      [NaN, '?']
    ])('input "%s"', (input: any, expected: string) => {
      expect(formatFileSize(input, '')).toEqual(expected)
    })
  })
  describe('respects different locales', () => {
    it.each([
      ['en', 1287654323, '1.3 GB'],
      ['de', 1287654323, '1,3 GB']
    ])('language "%s"', (language: string, input: number, expected: string) => {
      expect(formatFileSize(input, language)).toEqual(expected)
    })
  })
})
