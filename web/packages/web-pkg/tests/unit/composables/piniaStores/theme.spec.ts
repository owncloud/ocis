import { useLocalStorage, usePreferredDark } from '@vueuse/core'
import { useThemeStore, WebThemeConfigType } from '../../../../src/composables/piniaStores'
import { mockDeep } from 'vitest-mock-extended'
import { createPinia, setActivePinia } from 'pinia'
import { ref, computed } from 'vue'
import { useVault } from '../../../../src/composables/vault'

vi.mock('@vueuse/core', () => {
  return { useLocalStorage: vi.fn(() => ref('')), usePreferredDark: vi.fn(() => ref(false)) }
})

vi.mock('../../../../src/composables/vault', () => {
  return { useVault: vi.fn(() => ({ isInVault: false })) }
})

describe('useThemeStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('initializeThemes', () => {
    it('sets availableThemes', () => {
      const themeConfig = mockDeep<WebThemeConfigType>()
      themeConfig.themes = [
        { name: 'light', designTokens: {}, isDark: false },
        { name: 'dark', designTokens: {}, isDark: true }
      ]

      const store = useThemeStore()
      store.initializeThemes(themeConfig)

      expect(store.availableThemes.length).toBe(themeConfig.themes.length)
    })
    describe('currentTheme', () => {
      it.each([true, false])('gets set based on the OS setting', (isDark) => {
        vi.mocked(usePreferredDark).mockReturnValue(computed(() => isDark))
        vi.mocked(useLocalStorage).mockReturnValue(ref(null))

        const themeConfig = mockDeep<WebThemeConfigType>()
        themeConfig.themes = [
          { name: 'light', designTokens: {}, isDark: false },
          { name: 'dark', designTokens: {}, isDark: true }
        ]
        themeConfig.defaults = {
          designTokens: {},
          loginPage: { backgroundImg: '' },
          logo: { topbar: '', favicon: '', login: '' },
          icons: {}
        }

        const store = useThemeStore()
        store.initializeThemes(themeConfig)

        expect(store.currentTheme.name).toEqual(
          themeConfig.themes.find((t) => t.isDark === isDark).name
        )
      })
      it('falls back to the first theme if no match for the OS setting is found', () => {
        vi.mocked(usePreferredDark).mockReturnValue(computed(() => true))
        vi.mocked(useLocalStorage).mockReturnValue(ref(null))

        const themeConfig = mockDeep<WebThemeConfigType>()
        themeConfig.themes = [{ name: 'light', designTokens: {}, isDark: false }]
        themeConfig.defaults = {
          designTokens: {},
          loginPage: { backgroundImg: '' },
          logo: { topbar: '', favicon: '', login: '' },
          icons: {}
        }

        const store = useThemeStore()
        store.initializeThemes(themeConfig)

        expect(store.currentTheme.name).toEqual('light')
      })
    })

    describe('availableThemes', () => {
      it('returns regular themes if not in vault', () => {
        vi.mocked(useVault).mockReturnValue({ isInVault: false })

        const themeConfig = mockDeep<WebThemeConfigType>()
        themeConfig.themes = [
          { name: 'light', designTokens: {}, isDark: false, mode: 'regular' },
          { name: 'dark', designTokens: {}, isDark: true, mode: 'regular' },
          { name: 'light', designTokens: {}, isDark: false, mode: 'vault' },
          { name: 'dark', designTokens: {}, isDark: true, mode: 'vault' }
        ]

        const store = useThemeStore()
        store.initializeThemes(themeConfig)

        for (const theme of store.availableThemes) {
          expect(theme.mode).toBe('regular')
        }
      })
      it('returns vault themes if in vault', () => {
        vi.mocked(useVault).mockReturnValue({ isInVault: true })

        const themeConfig = mockDeep<WebThemeConfigType>()
        themeConfig.themes = [
          { name: 'light', designTokens: {}, isDark: false, mode: 'regular' },
          { name: 'dark', designTokens: {}, isDark: true, mode: 'regular' },
          { name: 'light', designTokens: {}, isDark: false, mode: 'vault' },
          { name: 'dark', designTokens: {}, isDark: true, mode: 'vault' }
        ]

        const store = useThemeStore()
        store.initializeThemes(themeConfig)

        for (const theme of store.availableThemes) {
          expect(theme.mode).toBe('vault')
        }
      })
      it('treats themes without mode as regular themes', () => {
        vi.mocked(useVault).mockReturnValue({ isInVault: false })

        const themeConfig = mockDeep<WebThemeConfigType>()
        themeConfig.themes = [
          { name: 'light', designTokens: {}, isDark: false, mode: 'regular' },
          { name: 'dark', designTokens: {}, isDark: true },
          { name: 'light', designTokens: {}, isDark: false, mode: 'vault' },
          { name: 'dark', designTokens: {}, isDark: true, mode: 'vault' }
        ]

        const store = useThemeStore()
        store.initializeThemes(themeConfig)
        expect(store.availableThemes.length).toBe(2)
        for (const theme of store.availableThemes) {
          expect(theme.mode).not.toBe('vault')
        }
      })
    })
  })
})
