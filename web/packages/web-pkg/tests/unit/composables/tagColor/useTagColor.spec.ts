// eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-imports
import { useLocalStorage, usePreferredDark } from '@vueuse/core'
import { useThemeStore, type WebThemeConfigType } from '../../../../src/composables/piniaStores'
import { useTagColor } from '../../../../src/composables'
import { mockDeep } from 'vitest-mock-extended'
import { createPinia, setActivePinia } from 'pinia'
import { ref } from 'vue'
// eslint-disable-next-line @typescript-eslint/no-unused-vars, unused-imports/no-unused-imports
import { useVault } from '../../../../src/composables/vault'

vi.mock('@vueuse/core', () => {
  return { useLocalStorage: vi.fn(() => ref('')), usePreferredDark: vi.fn(() => ref(false)) }
})

vi.mock('../../../../src/composables/vault', () => {
  return { useVault: vi.fn(() => ({ isInVault: false })) }
})

describe('useTagColor', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  // A realistic list length, so hashed slots spread the way they will in production.
  const thirtyColors = Array.from({ length: 30 }, (_, i) => `#${i.toString().padStart(6, '0')}`)

  const initThemeWithTagColors = (tagColorsList?: string[]) => {
    const themeConfig = mockDeep<WebThemeConfigType>()
    themeConfig.defaults = {
      designTokens: {},
      loginPage: { backgroundImg: '' },
      logo: { topbar: '', favicon: '', login: '' },
      icons: {}
    }
    themeConfig.themes = [
      {
        name: 'light',
        designTokens: tagColorsList ? { tagColorsList } : {},
        isDark: false
      }
    ]
    useThemeStore().initializeThemes(themeConfig)
  }

  describe('Tag chip with colors', () => {
    it('resolves a tag name to a var(--oc-color-tag-N) string when the theme has a tag color list', () => {
      initThemeWithTagColors(['#111111', '#222222', '#333333'])

      const { tagColor } = useTagColor()
      const result = tagColor('test-tag')

      expect(result).toMatch(/^var\(--oc-color-tag-\d+\)$/)
    })

    it('returns the same value when called twice with the same tag name', () => {
      initThemeWithTagColors(['#111111', '#222222', '#333333'])

      const { tagColor } = useTagColor()
      const result1 = tagColor('physics')
      const result2 = tagColor('physics')

      expect(result1).toBe(result2)
    })

    it('returns different values for two different tag names', () => {
      initThemeWithTagColors(thirtyColors)

      const { tagColor } = useTagColor()
      const physicsColor = tagColor('physics')
      const invoiceColor = tagColor('invoice')

      expect(physicsColor).not.toBe(invoiceColor)
    })

    it('returns different values for "Invoice" and "invoice" (case sensitive)', () => {
      initThemeWithTagColors(thirtyColors)

      const { tagColor } = useTagColor()
      const uppercaseColor = tagColor('Invoice')
      const lowercaseColor = tagColor('invoice')

      expect(uppercaseColor).not.toBe(lowercaseColor)
    })

    it('returns an empty string when the theme has no tagColorsList', () => {
      initThemeWithTagColors()

      const { tagColor } = useTagColor()
      const result = tagColor('test-tag')

      expect(result).toBe('')
    })

    it('returns an empty string when the theme has an empty tagColorsList array', () => {
      initThemeWithTagColors([])

      const { tagColor } = useTagColor()
      const result = tagColor('test-tag')

      expect(result).toBe('')
    })

    it('returns an empty string when no theme has been initialized yet', () => {
      // currentTheme is `ref<WebThemeType | undefined>()`, so it is undefined until
      // initializeThemes runs. Any component that renders a tag before that — or under a
      // testing pinia that never initializes themes — must not crash.
      const { tagColor } = useTagColor()

      expect(tagColor('test-tag')).toBe('')
    })
  })

  describe('tagLabelColor', () => {
    // The fill comes from the theme, so the label has to be chosen per fill at runtime: a pale
    // fill needs a dark label and a dark fill needs a pale one. A single hardcoded token
    // cannot satisfy both, and `--oc-color-text-inverse` gets it backwards in light mode.
    it('returns a dark label for a pale fill', () => {
      // One-colour list, so every tag name lands on --oc-color-tag-0.
      initThemeWithTagColors(['#fdf5c9'])

      expect(useTagColor().tagLabelColor('physics')).toBe('#000000')
    })

    it('returns a pale label for a dark fill', () => {
      initThemeWithTagColors(['#1b365d'])

      expect(useTagColor().tagLabelColor('physics')).toBe('#ffffff')
    })

    it('reads the fill from the theme, so an unapplied css var cannot flip the label', () => {
      // Going through `getHexFromCssVar` would be unsafe here: it funnels an unresolvable var
      // into `cssRgbToHex`, which returns '#000000' for an empty string. A pale fill would then
      // be read as black and get a white — illegible — label.
      initThemeWithTagColors(['#fdf5c9'])
      document.documentElement.style.removeProperty('--oc-color-tag-0')

      expect(useTagColor().tagLabelColor('physics')).toBe('#000000')
    })

    it('returns an empty string when the theme colour is not a hex value', () => {
      // Rather than guess a label for a colour it cannot measure.
      initThemeWithTagColors(['not-a-colour'])

      expect(useTagColor().tagLabelColor('physics')).toBe('')
    })

    it('returns an empty string when the theme has no tag color list', () => {
      initThemeWithTagColors([])

      expect(useTagColor().tagLabelColor('physics')).toBe('')
    })

    it('returns an empty string when no theme has been initialized yet', () => {
      expect(useTagColor().tagLabelColor('physics')).toBe('')
    })
  })
})
