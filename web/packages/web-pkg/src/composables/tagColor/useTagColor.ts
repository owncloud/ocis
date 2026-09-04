import { useThemeStore } from '../piniaStores'
import {
  getContrastRatio,
  hashToIndex,
  hexToRgb,
  tagColorVarFor
} from '@ownclouders/design-system/helpers'

const BLACK = '#000000'
const WHITE = '#ffffff'

/**
 * Resolves a tag name to the theme colours its chip should use.
 *
 * Colours are derived, never stored: the same tag string always lands on the same slot of the
 * theme's `tagColorsList`, so a tag looks identical for every user on every device without the
 * server knowing anything about colours. See ADR-0029.
 */
export const useTagColor = () => {
  const themeStore = useThemeStore()

  const colorList = (): string[] | undefined => themeStore.currentTheme?.designTokens?.tagColorsList

  /**
   * The chip's fill, as a `var(--oc-color-tag-N)` reference. Empty when the theme ships no tag
   * colours, which leaves the chip with its default appearance.
   */
  const tagColor = (name: string): string => {
    const tagColorsList = colorList()
    if (!tagColorsList?.length) {
      return ''
    }
    return tagColorVarFor(name, tagColorsList.length)
  }

  /**
   * The chip's label colour, picked per fill for legibility. This has to be computed rather than
   * pinned to a token: the fills come from the theme, so a single token would be wrong for half
   * of them. It lives here rather than in OcTag because the theme store is a reactive dependency
   * and `getComputedStyle` is not — reading the fill inside the component would leave the label
   * stale after a theme switch.
   *
   * The fill is read straight out of the theme rather than resolved from the emitted css var:
   * `getHexFromCssVar` reports an unresolvable var as `#000000`, which would read a pale fill as
   * black and hand it an illegible white label.
   */
  const tagLabelColor = (name: string): string => {
    const tagColorsList = colorList()
    if (!tagColorsList?.length) {
      return ''
    }
    const rgb = hexToRgb(tagColorsList[hashToIndex(name, tagColorsList.length)])
    if (!rgb) {
      return ''
    }
    return getContrastRatio(rgb, hexToRgb(BLACK)) >= getContrastRatio(rgb, hexToRgb(WHITE))
      ? BLACK
      : WHITE
  }

  return { tagColor, tagLabelColor }
}
