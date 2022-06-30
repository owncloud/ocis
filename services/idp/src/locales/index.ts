import allLocales from './locales.json';

interface Locale {
  locale: string,
  name: string,
  nativeName: string,
}

function enableLocales(locales: Locale[]): Locale[] {
  return locales;
}

export const locales = enableLocales(allLocales);

export default locales;
