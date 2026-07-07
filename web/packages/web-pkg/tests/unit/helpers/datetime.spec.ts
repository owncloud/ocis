import { DateTime, Settings } from 'luxon'
import {
  formatDateFromDateTime,
  formatDateFromJSDate,
  formatDateFromHTTP,
  formatDateFromISO,
  formatDateFromRFC,
  formatRelativeDateFromDateTime,
  formatRelativeDateFromJSDate,
  formatRelativeDateFromISO,
  formatRelativeDateFromRFC
} from '../../../src/helpers/datetime'

describe('datetime helper', () => {
  const language = 'en'
  const dateFormat = DateTime.DATETIME_MED
  beforeEach(() => {
    Settings.defaultZone = 'utc'
  })
  describe('formatDateFromDateTime', () => {
    it('should give correct output', () => {
      expect(
        formatDateFromDateTime(DateTime.fromISO('2010-10-22T21:38:00'), language, dateFormat)
      ).toBe('Oct 22, 2010, 9:38\u202fPM')
    })
  })
  describe('formatDateFromJSDate', () => {
    it('should give correct output', () => {
      expect(formatDateFromJSDate(new Date('2010-10-22T21:38:00'), language, dateFormat)).toBe(
        'Oct 22, 2010, 9:38\u202fPM'
      )
    })
    it('should fail for null', () => {
      expect(formatDateFromJSDate(null, language, dateFormat)).toBe('Invalid DateTime')
    })
  })
  describe('formatDateFromHTTP', () => {
    it('should give correct output', () => {
      expect(formatDateFromHTTP('Tue, 15 Nov 1994 12:45:26 GMT', language, dateFormat)).toBe(
        'Nov 15, 1994, 12:45\u202fPM'
      )
    })
    it('should fail for invalid http date', () => {
      expect(formatDateFromHTTP('Some not http date 123', language, dateFormat)).toBe(
        'Invalid DateTime'
      )
    })
  })
  describe('formatDateFromISO', () => {
    it('should give correct output', () => {
      expect(formatDateFromISO('2010-10-22T21:38:00', language, dateFormat)).toBe(
        'Oct 22, 2010, 9:38\u202fPM'
      )
    })
    it('should fail for invalid iso date', () => {
      expect(formatDateFromISO('some invalid iso date 123', language, dateFormat)).toBe(
        'Invalid DateTime'
      )
    })
  })
  describe('formatDateFromRFC', () => {
    it('should give correct output', () => {
      expect(formatDateFromRFC('01 Jun 2016 14:31:46 -0700', language, dateFormat)).toBe(
        'Jun 1, 2016, 9:31\u202fPM'
      )
    })
    it('should fail for invalid rfc date', () => {
      expect(formatDateFromRFC('some invalid rfc 123', language, dateFormat)).toBe(
        'Invalid DateTime'
      )
    })
  })
  describe('formatRelativeDateFromDateTime', () => {
    it('should return correct relative time', () => {
      expect(formatRelativeDateFromDateTime(DateTime.now().minus({ years: 12 }), language)).toBe(
        '12 years ago'
      )
    })
  })
  describe('formatRelativeDateFromJSDate', () => {
    it('should return correct relative time', () => {
      expect(
        formatRelativeDateFromJSDate(DateTime.now().minus({ years: 12 }).toJSDate(), language)
      ).toBe('12 years ago')
    })
    it('should return null if date is null', () => {
      expect(formatRelativeDateFromJSDate(null, language)).toBe(null)
    })
  })
  describe('formatRelativeDateFromISO', () => {
    it('should return correct relative time', () => {
      expect(formatRelativeDateFromISO(DateTime.now().minus({ years: 12 }).toISO(), language)).toBe(
        '12 years ago'
      )
    })
    it('should return null if invalid iso', () => {
      expect(formatRelativeDateFromISO('some invalid iso 123', language)).toBe(null)
    })
  })
  describe('formatRelativeDateFromRFC', () => {
    it('should return correct relative time', () => {
      expect(
        formatRelativeDateFromRFC(DateTime.now().minus({ years: 12 }).toRFC2822(), language)
      ).toBe('12 years ago')
    })
    it('should return null if invalid rfc', () => {
      expect(formatRelativeDateFromRFC('some invalid rfc 123', language)).toBe(null)
    })
  })
})
