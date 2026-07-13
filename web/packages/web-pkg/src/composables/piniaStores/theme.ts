import merge from 'deepmerge'
import { defineStore } from 'pinia'
import { computed, ref, unref } from 'vue'
import { useLocalStorage, usePreferredDark } from '@vueuse/core'
import { z } from 'zod'
import { applyCustomProp, removeCustomProp } from '@ownclouders/design-system/helpers'
import { ShareRole } from '@ownclouders/web-client'
import { useVault } from '../vault'

const AppBanner = z.object({
  title: z.string().optional(),
  publisher: z.string().optional(),
  additionalInformation: z.string().optional(),
  ctaText: z.string().optional(),
  icon: z.string().optional(),
  appScheme: z.string().optional()
})

const CommonSection = z.object({
  name: z.string(),
  slogan: z.string(),
  logo: z.string(),
  urls: z.object({
    accessDeniedHelp: z.string(),
    imprint: z.string(),
    privacy: z.string(),
    accessibilityStatement: z.string().optional(),
    universalAccessEasyLanguage: z.string().optional(),
    universalAccessSignLanguage: z.string().optional(),
    softwareLicense: z.string().optional(),
    helpPage: z.string().optional()
  }),
  shareRoles: z.record(
    z.string(),
    z.object({
      iconName: z.string()
    })
  )
})

const DesignTokens = z.object({
  breakpoints: z.record(z.string(), z.string()).optional(),
  colorPalette: z.record(z.string(), z.string()).optional(),
  fontFamily: z.string().optional(),
  fontSizes: z.record(z.string(), z.string()).optional(),
  sizes: z.record(z.string(), z.string()).optional(),
  spacing: z.record(z.string(), z.string()).optional()
})

const LoginPage = z.object({
  backgroundImg: z.string()
})

const Logo = z.object({
  topbar: z.string(),
  topbarSm: z.string().optional(),
  favicon: z.string(),
  login: z.string(),
  notFound: z.string().optional(),
  href: z.string().optional()
})

const Icons = z.object({
  universalAccess: z.string().optional(),
  universalAccessEasyLanguage: z.string().optional(),
  universalAccessSignLanguage: z.string().optional()
})

const ThemeDefaults = z.object({
  appBanner: AppBanner.optional(),
  common: CommonSection.optional(),
  designTokens: DesignTokens,
  loginPage: LoginPage,
  logo: Logo,
  icons: Icons.optional()
})

const WebTheme = z.object({
  appBanner: AppBanner.optional(),
  common: CommonSection.optional(),
  designTokens: DesignTokens.optional(),
  isDark: z.boolean(),
  /**
   * Specifies whether the theme is suitable for regular mode or vault mode.
   * If not specified, the theme is suitable for regular mode.
   */
  mode: z.optional(z.enum(['regular', 'vault']).default('regular')),
  name: z.string(),
  loginPage: LoginPage.optional(),
  logo: Logo.optional(),
  icons: Icons.optional()
})

export const WebThemeConfig = z.object({
  defaults: ThemeDefaults,
  themes: z.array(WebTheme)
})

export const ThemingConfig = z.object({
  common: CommonSection.optional(),
  clients: z.object({
    web: WebThemeConfig
  })
})

export type WebThemeType = z.infer<typeof WebTheme>
export type WebThemeConfigType = z.infer<typeof WebThemeConfig>

const themeStorageKey = 'oc_currentThemeName'

export const useThemeStore = defineStore('theme', () => {
  const currentLocalStorageThemeName = useLocalStorage(themeStorageKey, null)

  const isDark = usePreferredDark()
  const { isInVault } = useVault()

  const currentTheme = ref<WebThemeType | undefined>()

  const themes = ref<WebThemeType[]>([])

  const availableThemes = computed(() => {
    return unref(themes).filter((theme) => {
      if (unref(isInVault)) {
        return theme.mode === 'vault'
      }

      return theme.mode === 'regular' || theme.mode === undefined
    })
  })

  const initializeThemes = (themeConfig: WebThemeConfigType) => {
    themes.value = themeConfig.themes.map((theme) =>
      merge<WebThemeType>(themeConfig.defaults, theme)
    )
    setThemeFromStorageOrSystem()
  }

  const setThemeFromStorageOrSystem = () => {
    const firstLightTheme = unref(availableThemes).find((theme) => !theme.isDark)
    const firstDarkTheme = unref(availableThemes).find((theme) => theme.isDark)
    setAndApplyTheme(
      unref(availableThemes).find((t) => t.name === unref(currentLocalStorageThemeName)) ||
        (unref(isDark) ? firstDarkTheme : firstLightTheme) ||
        unref(availableThemes)[0],
      false
    )
  }

  const setAutoSystemTheme = () => {
    currentLocalStorageThemeName.value = null
    setThemeFromStorageOrSystem()
  }

  const isCurrentThemeAutoSystem = computed(() => {
    return currentLocalStorageThemeName.value === null
  })

  const setAndApplyTheme = (theme: WebThemeType, updateStorage = true) => {
    const previousTheme = unref(currentTheme)
    currentTheme.value = theme
    if (updateStorage) {
      currentLocalStorageThemeName.value = unref(currentTheme).name
    }

    const customizableDesignTokens = [
      { name: 'breakpoints', prefix: 'breakpoint' },
      { name: 'colorPalette', prefix: 'color' },
      { name: 'fontSizes', prefix: 'font-size' },
      { name: 'sizes', prefix: 'size' },
      { name: 'spacing', prefix: 'spacing' }
    ] as const

    if (previousTheme) {
      customizableDesignTokens.forEach((token) => {
        for (const param in previousTheme.designTokens[token.name]) {
          removeCustomProp(`${token.prefix}-${param}`)
        }
      })
    }

    applyCustomProp('font-family', unref(currentTheme).designTokens.fontFamily)

    customizableDesignTokens.forEach((token) => {
      for (const param in unref(currentTheme).designTokens[token.name]) {
        applyCustomProp(
          `${token.prefix}-${param}`,
          unref(currentTheme).designTokens[token.name][param]
        )
      }
    })
  }

  const getRoleIcon = (role: ShareRole) => {
    return unref(currentTheme).common?.shareRoles[role.id]?.iconName || 'user'
  }

  return {
    availableThemes,
    currentTheme,
    themes,
    initializeThemes,
    setAndApplyTheme,
    setAutoSystemTheme,
    isCurrentThemeAutoSystem,
    getRoleIcon
  }
})
