export const getLocaleFromLanguage = (currentLanguage: string) => {
  return (currentLanguage || '').split('_')[0]
}
