import { useThemeStore } from '../piniaStores'
import { hashToIndex } from '@ownclouders/design-system/helpers'

/**
 * Resolves a tag name to the slot of the theme's tag colour list its chip should use.
 *
 * Colours are derived, never stored: the same tag string always lands on the same slot, so a tag
 * looks identical for every user on every device without the server knowing anything about
 * colours. What that slot looks like is entirely the theme's business — pass the index to `OcTag`
 * and it takes both the fill and the label colour from the theme's css vars.
 */
export const useTagColor = () => {
  const themeStore = useThemeStore()

  /**
   * The chip's colour slot, or -1 when the theme ships no tag colours, which leaves the chip with
   * its default appearance.
   */
  const tagColorIndex = (name: string): number =>
    hashToIndex(name, themeStore.currentTheme?.designTokens?.tagColorsList?.length ?? 0)

  return { tagColorIndex }
}
