import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import resourcesToBackend from 'i18next-resources-to-backend';
import LanguageDetector, { CustomDetector } from 'i18next-browser-languagedetector';

import queryString from 'query-string';

import locales from './locales';

const config = {
  uiLocalesQueryName: 'ui_locales', // Same as OIDC uses.
  uiLocaleCookieName: 'ui_locale',  // For domain wide syncing, not set here.
  uiLocaleLocalStorageName: 'lico.identifier_ui_locale', // Sufficiently unique, set here.
}

const supportedLanguages = locales.map((locale) => {
  return locale.locale;
});

const queryUiLocalesDetector: CustomDetector = {
  name: 'queryUiLocales',
  lookup: (options): string | string[] | undefined => {
    const query = queryString.parse(document.location.search);
    const ui_locales = query[config.uiLocalesQueryName];
    if (!ui_locales) {
      return;
    }
    if (Array.isArray(ui_locales)) {
      return ui_locales as string[];
    } else {
      return ui_locales.split(' ');
    }
  },
}

const languageDetector = new LanguageDetector();
languageDetector.addDetector(queryUiLocalesDetector);

i18n
  .use(resourcesToBackend((language, namespace, callback) => {
    import(
      /* webpackMode: "lazy-once" */
      /* webpackInclude: /\.json$/ */
      /* webpackChunkName: "all-i18n-data" */
      /* webpackPrefetch: true */
      `./locales/${language}/${namespace}.json`
    )
      .then((resources) => {
        callback(null, resources)
      })
      .catch((error) => {
        callback(error, null)
      })
  }))
  .use(languageDetector)
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    fallbackLng: 'en-GB',
    supportedLngs: [...supportedLanguages],
    cleanCode: true,
    returnEmptyString: false,
    debug: false,

    detection: {
      /* https://github.com/i18next/i18next-browser-languageDetector */
      order: [queryUiLocalesDetector.name, 'cookie', 'localStorage', 'navigator'],
      lookupCookie: config.uiLocaleCookieName,
      lookupLocalStorage: config.uiLocaleLocalStorageName,
      caches: ['localStorage'],
    },

    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
