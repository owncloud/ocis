import merge from 'lodash-es/merge'
import App from '../App.vue'
import TokenRenewal from '../pages/tokenRenewal.vue'
import missingOrInvalidConfigPage from '../pages/missingOrInvalidConfig.vue'

export * from './languages'

export const pages = {
  success: App,
  failure: missingOrInvalidConfigPage,
  tokenRenewal: TokenRenewal
}

export const loadTranslations = async () => {
  const { coreTranslations, clientTranslations, pkgTranslations, odsTranslations } =
    await import('./json')

  return merge({}, coreTranslations, clientTranslations, pkgTranslations, odsTranslations)
}

export const loadDesignSystem = async () => {
  return (await import('@ownclouders/design-system')).default
}
