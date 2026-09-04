import { useLocalStorage, usePreferredDark } from '@vueuse/core'
import {
  useThemeStore,
  WebThemeConfig,
  WebThemeConfigType
} from '../../../../src/composables/piniaStores'
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

  describe('tagColorsList', () => {
    const rootStyle = () => (document.querySelector(':root') as HTMLElement).style
    const themeDefaults = {
      designTokens: {},
      loginPage: { backgroundImg: '' },
      logo: { topbar: '', favicon: '', login: '' },
      icons: {}
    }

    beforeEach(() => {
      vi.mocked(useVault).mockReturnValue({ isInVault: false })
      vi.mocked(useLocalStorage).mockReturnValue(ref(null))
      vi.mocked(usePreferredDark).mockReturnValue(computed(() => false))
      rootStyle().cssText = ''
    })

    it('emits one zero-based custom property per entry', () => {
      const themeConfig = mockDeep<WebThemeConfigType>()
      themeConfig.defaults = themeDefaults
      themeConfig.themes = [
        {
          name: 'light',
          isDark: false,
          designTokens: { tagColorsList: ['#aaaaaa', '#bbbbbb', '#cccccc'] }
        }
      ]

      useThemeStore().initializeThemes(themeConfig)

      expect(rootStyle().getPropertyValue('--oc-color-tag-0')).toBe('#aaaaaa')
      expect(rootStyle().getPropertyValue('--oc-color-tag-1')).toBe('#bbbbbb')
      expect(rootStyle().getPropertyValue('--oc-color-tag-2')).toBe('#cccccc')
      expect(rootStyle().getPropertyValue('--oc-color-tag-3')).toBe('')
    })

    it('unsets the properties of the previous theme on theme switch', () => {
      const themeConfig = mockDeep<WebThemeConfigType>()
      themeConfig.defaults = themeDefaults
      themeConfig.themes = [
        { name: 'light', isDark: false, designTokens: { tagColorsList: ['#aaaaaa', '#bbbbbb'] } },
        { name: 'dark', isDark: true, designTokens: { tagColorsList: ['#111111'] } }
      ]

      const store = useThemeStore()
      store.initializeThemes(themeConfig)
      store.setAndApplyTheme(store.availableThemes.find((t) => t.isDark))

      // The dark theme is shorter, so slot 1 must not keep the light theme's value.
      expect(rootStyle().getPropertyValue('--oc-color-tag-0')).toBe('#111111')
      expect(rootStyle().getPropertyValue('--oc-color-tag-1')).toBe('')
    })

    it('lets a theme replace the list in defaults rather than appending to it', () => {
      const themeConfig = mockDeep<WebThemeConfigType>()
      themeConfig.defaults = {
        ...themeDefaults,
        designTokens: { tagColorsList: ['#defa17', '#defa18'] }
      }
      themeConfig.themes = [
        { name: 'light', isDark: false, designTokens: { tagColorsList: ['#aaaaaa', '#bbbbbb'] } }
      ]

      const store = useThemeStore()
      store.initializeThemes(themeConfig)

      // deepmerge concatenates arrays by default, which would yield 4 entries and change
      // `hash % length` for every tag.
      expect(store.currentTheme.designTokens.tagColorsList).toEqual(['#aaaaaa', '#bbbbbb'])
      expect(rootStyle().getPropertyValue('--oc-color-tag-2')).toBe('')
    })

    it('inherits the list from defaults when a theme does not define one', () => {
      const themeConfig = mockDeep<WebThemeConfigType>()
      themeConfig.defaults = { ...themeDefaults, designTokens: { tagColorsList: ['#defa17'] } }
      themeConfig.themes = [{ name: 'light', isDark: false, designTokens: {} }]

      useThemeStore().initializeThemes(themeConfig)

      expect(rootStyle().getPropertyValue('--oc-color-tag-0')).toBe('#defa17')
    })

    it.each([
      ['accepts', 512, true],
      ['rejects', 513, false]
    ])('%s a list of %i entries', (_verb, length, valid) => {
      const config = {
        defaults: {
          ...themeDefaults,
          designTokens: { tagColorsList: Array(length).fill('#000000') }
        },
        themes: [{ name: 'light', isDark: false, designTokens: {} }]
      }

      expect(WebThemeConfig.safeParse(config).success).toBe(valid)
    })
  })
})
