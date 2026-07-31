import { ApplicationInformation } from '@ownclouders/web-pkg'
import { merge } from 'lodash-es'
import { Language, Translations } from 'vue3-gettext'

export const setCurrentLanguage = ({
  language,
  languageSetting = null
}: {
  language: Language
  languageSetting?: string
}): void => {
  let currentLanguage = languageSetting
  if (currentLanguage) {
    if (currentLanguage.indexOf('-')) {
      currentLanguage = currentLanguage.split('-')[0]
    }
    language.current = currentLanguage
    document.documentElement.lang = currentLanguage
  }
}

/**
 * Loads all app translations for one given language.
 * This should be called each time the language is being changed.
 */
export const loadAppTranslations = ({
  apps,
  gettext,
  lang
}: {
  apps: Record<string, ApplicationInformation>
  gettext: Language
  lang: string
}) => {
  const appTranslations: Translations = {}
  Object.values(apps).forEach((app) => {
    const { translations } = app
    if (gettext.translations[lang] && translations?.[lang]) {
      Object.assign(appTranslations, translations[lang])
    }
  })

  gettext.translations = merge(gettext.translations, {
    [lang]: appTranslations
  })
}
