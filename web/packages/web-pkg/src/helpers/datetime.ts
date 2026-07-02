import { DateTime } from 'luxon'
import { getLocaleFromLanguage } from './locale'

export const formatDateFromDateTime = (
  date: DateTime,
  currentLanguage: string,
  format = DateTime.DATETIME_MED
) => {
  return date.setLocale(getLocaleFromLanguage(currentLanguage)).toLocaleString(format)
}
export const formatDateFromJSDate = (
  date: Date,
  currentLanguage: string,
  format = DateTime.DATETIME_MED
) => {
  return formatDateFromDateTime(DateTime.fromJSDate(date), currentLanguage, format)
}
export const formatDateFromHTTP = (
  date: string,
  currentLanguage: string,
  format = DateTime.DATETIME_MED
) => {
  return formatDateFromDateTime(DateTime.fromHTTP(date), currentLanguage, format)
}
export const formatDateFromISO = (
  date: string,
  currentLanguage: string,
  format = DateTime.DATETIME_MED
) => {
  return formatDateFromDateTime(DateTime.fromISO(date), currentLanguage, format)
}
export const formatDateFromRFC = (
  date: string,
  currentLanguage: string,
  format = DateTime.DATETIME_MED
) => {
  return formatDateFromDateTime(DateTime.fromRFC2822(date), currentLanguage, format)
}
export const formatRelativeDateFromDateTime = (date: DateTime, currentLanguage: string) => {
  return date.setLocale(getLocaleFromLanguage(currentLanguage)).toRelative()
}
export const formatRelativeDateFromJSDate = (date: Date, currentLanguage: string) => {
  return formatRelativeDateFromDateTime(DateTime.fromJSDate(date), currentLanguage)
}
export const formatRelativeDateFromHTTP = (date: string, currentLanguage: string) => {
  return formatRelativeDateFromDateTime(DateTime.fromHTTP(date), currentLanguage)
}
export const formatRelativeDateFromISO = (date: string, currentLanguage: string) => {
  return formatRelativeDateFromDateTime(DateTime.fromISO(date), currentLanguage)
}
export const formatRelativeDateFromRFC = (date: string, currentLanguage: string) => {
  return formatRelativeDateFromDateTime(DateTime.fromRFC2822(date), currentLanguage)
}
