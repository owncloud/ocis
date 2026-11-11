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

  describe('tagColorIndex', () => {
    it('resolves a tag name to a slot of the theme list', () => {
      initThemeWithTagColors(['#111111', '#222222', '#333333'])

      const index = useTagColor().tagColorIndex('test-tag')

      expect(index).toBeGreaterThanOrEqual(0)
      expect(index).toBeLessThan(3)
    })

    it('returns the same slot when called twice with the same tag name', () => {
      initThemeWithTagColors(['#111111', '#222222', '#333333'])

      const { tagColorIndex } = useTagColor()

      expect(tagColorIndex('physics')).toBe(tagColorIndex('physics'))
    })

    it('returns different slots for two different tag names', () => {
      initThemeWithTagColors(thirtyColors)

      const { tagColorIndex } = useTagColor()

      expect(tagColorIndex('physics')).not.toBe(tagColorIndex('invoice'))
    })

    it('returns different slots for "Invoice" and "invoice" (case sensitive)', () => {
      initThemeWithTagColors(thirtyColors)

      const { tagColorIndex } = useTagColor()

      expect(tagColorIndex('Invoice')).not.toBe(tagColorIndex('invoice'))
    })

    it('spreads names across the whole list', () => {
      initThemeWithTagColors(thirtyColors)

      const { tagColorIndex } = useTagColor()
      const slots = new Set(Array.from({ length: 200 }, (_, i) => tagColorIndex(`tag-${i}`)))

      expect(slots.size).toBeGreaterThan(20)
    })

    it('returns -1 when the theme has no tagColorsList', () => {
      // -1 is what tells `OcTag` to stay a plain badge, which is how a theme that has not opted
      // into coloured tags renders.
      initThemeWithTagColors()

      expect(useTagColor().tagColorIndex('test-tag')).toBe(-1)
    })

    it('returns -1 when the theme has an empty tagColorsList array', () => {
      initThemeWithTagColors([])

      expect(useTagColor().tagColorIndex('test-tag')).toBe(-1)
    })

    it('returns -1 when no theme has been initialized yet', () => {
      // currentTheme is `ref<WebThemeType | undefined>()`, so it is undefined until
      // initializeThemes runs. Any component that renders a tag before that — or under a
      // testing pinia that never initializes themes — must not crash.
      expect(useTagColor().tagColorIndex('test-tag')).toBe(-1)
    })
  })
})
