import {
  hexToRgb,
  rgbToHex,
  calculateShadeColor,
  getLuminanace,
  getContrastRatio,
  generateHashedColorForString,
  setDesiredContrastRatio,
  cssRgbToHex,
  getHexFromCssVar,
  hashString,
  hashToIndex,
  pickReadableTextColor,
  tagColorVarsFor
} from './colors'

describe('hexToRgb', () => {
  it('converts hex to rgb', () => {
    expect(hexToRgb('#FF0000')).toEqual([255, 0, 0])
    expect(hexToRgb('#00FF00')).toEqual([0, 255, 0])
    expect(hexToRgb('#0000FF')).toEqual([0, 0, 255])
    expect(hexToRgb('#FFFFFF')).toEqual([255, 255, 255])
    expect(hexToRgb('#000000')).toEqual([0, 0, 0])
    expect(hexToRgb('FF0000')).toEqual([255, 0, 0]) // Test without #
    expect(hexToRgb('invalid')).toBeNull() // Test with invalid input
  })
})

describe('rgbToHex', () => {
  it('converts rgb to hex', () => {
    expect(rgbToHex([255, 0, 0])).toBe('#ff0000')
    expect(rgbToHex([0, 255, 0])).toBe('#00ff00')
    expect(rgbToHex([0, 0, 255])).toBe('#0000ff')
    expect(rgbToHex([255, 255, 255])).toBe('#ffffff')
    expect(rgbToHex([0, 0, 0])).toBe('#000000')
  })
})

describe('calculateShadeColor', () => {
  it('shades a color', () => {
    const initialColor = [100, 50, 50]
    expect(calculateShadeColor(initialColor, 50)).toBe('#964b4b')
    expect(calculateShadeColor(initialColor, -10)).toBe('#5a2d2d')
  })
})

describe('getLuminanace', () => {
  it('calculates luminance', () => {
    expect(getLuminanace([255, 255, 255])).toBeCloseTo(1)
    expect(getLuminanace([0, 0, 0])).toBeCloseTo(0)
  })
})

describe('getContrastRatio', () => {
  it('calculates contrast ratio', () => {
    const colorA = [255, 255, 255]
    const colorB = [0, 0, 0]
    expect(getContrastRatio(colorA, colorB)).toBeGreaterThan(20)
  })
})

describe('generateHashedColorForString', () => {
  it('generates a hashed color', () => {
    expect(generateHashedColorForString('owncloud')).toBe('#2F26F')
    expect(generateHashedColorForString('example')).toMatch('#25116A')
  })
})

describe('hashString', () => {
  it('is deterministic', () => {
    expect(hashString('physics')).toBe(hashString('physics'))
  })
  it('does not case-fold, because tags do not', () => {
    // reva's tags.normalize() trims and splits but never lowercases, so `Invoice` and `invoice`
    // are two different tags and must get two different colours.
    expect(hashString('Invoice')).not.toBe(hashString('invoice'))
  })
  it('returns 0 for an empty string', () => {
    expect(hashString('')).toBe(0)
  })
  it('normalises, so callers can use the result as-is', () => {
    // `hash << 5` overflows into the sign bit, and the addition that follows leaves int32
    // altogether, so without the normalisation `hashString` returns numbers that are unusable
    // as an index — negative, or beyond 32 bits.
    for (const name of ['physics', 'Invoice', 'a', '', '🎉', 'x'.repeat(200)]) {
      const hash = hashString(name)
      expect(hash).toBeGreaterThanOrEqual(0)
      expect(hash).toBeLessThanOrEqual(0xffffffff)
      expect(Number.isInteger(hash)).toBe(true)
    }
  })
})

describe('hashToIndex', () => {
  it('stays within the list', () => {
    for (const name of ['physics', 'Invoice', 'invoice', 'a', '', 'tag', 'ünïcödé', '🎉']) {
      const index = hashToIndex(name, 30)
      expect(index).toBeGreaterThanOrEqual(0)
      expect(index).toBeLessThan(30)
      expect(Number.isInteger(index)).toBe(true)
    }
  })
  it.each([0, -1, 1.5, NaN])('returns -1 for an unusable length of %s', (length) => {
    expect(hashToIndex('physics', length)).toBe(-1)
  })
  it('is stable across calls', () => {
    expect(hashToIndex('physics', 30)).toBe(hashToIndex('physics', 30))
  })
  it('pins known names to known slots', () => {
    // Pinned deliberately: these values are the cross-device, cross-user guarantee the whole
    // feature rests on. If a change here is intentional it re-colours every existing tag in
    // every deployment, so it must be a conscious edit, not a silent regression.
    expect(hashToIndex('physics', 30)).toBe(13)
    expect(hashToIndex('invoice', 30)).toBe(7)
    expect(hashToIndex('Invoice', 30)).toBe(27)
    expect(hashToIndex('important', 30)).toBe(28)
  })
  it('re-colours names when the list grows, which is why growing it is a one-way act', () => {
    expect(hashToIndex('physics', 30)).not.toBe(hashToIndex('physics', 31))
  })
})

describe('tagColorVarsFor', () => {
  it('names the theme css vars for a slot', () => {
    expect(tagColorVarsFor(13)).toEqual({
      fillColor: 'var(--oc-color-tag-13)',
      textColor: 'var(--oc-color-tag-13-text)'
    })
  })
  it('pairs the fill of a name with the text colour belonging to it', () => {
    const index = hashToIndex('invoice', 30)
    expect(tagColorVarsFor(index)).toEqual({
      fillColor: `var(--oc-color-tag-${index})`,
      textColor: `var(--oc-color-tag-${index}-text)`
    })
  })
  it.each([-1, 1.5, NaN, undefined, null])(
    'returns null for %s rather than naming a var that cannot exist',
    (index) => {
      // -1 is what `hashToIndex` returns when the theme ships no tag colours, which is the state
      // of any theme that has not opted into coloured tags.
      expect(tagColorVarsFor(index)).toBeNull()
    }
  )
})

describe('pickReadableTextColor', () => {
  it('puts light text on a dark fill and dark text on a light fill', () => {
    expect(pickReadableTextColor('#041e42', '#041e42', '#ffffff')).toBe('#ffffff')
    expect(pickReadableTextColor('#e3baba', '#041e42', '#ffffff')).toBe('#041e42')
  })
  it('uses the candidates it is given rather than black and white', () => {
    // A theme decides what its text looks like; the pairing only decides which of the two.
    expect(pickReadableTextColor('#000000', '#333333', '#eeeeee')).toBe('#eeeeee')
    expect(pickReadableTextColor('#ffffff', '#333333', '#eeeeee')).toBe('#333333')
  })
  it('falls back to the dark candidate when a colour cannot be parsed', () => {
    expect(pickReadableTextColor('not-a-colour', '#041e42', '#ffffff')).toBe('#041e42')
  })
})

describe('setDesiredContrastRatio', () => {
  it('adjusts color for desired contrast ratio', () => {
    const targetColor = [100, 100, 100]
    const associatedColor = [255, 255, 255]
    const desiredRatio = 3
    expect(setDesiredContrastRatio(targetColor, associatedColor, desiredRatio)).toEqual([
      142, 142, 142
    ])
  })
})

describe('cssRgbToHex', () => {
  it('converts css rgb value to hex', () => {
    expect(cssRgbToHex('rgb(255, 0, 0)')).toBe('#ff0000')
    expect(cssRgbToHex('rgba(0, 255, 0, 0.5)')).toBe('#00ff0080')
  })
})

describe('getHexFromCssVar', () => {
  it('retrieves hex value from CSS var', () => {
    document.documentElement.style.setProperty('--color-primary', '#ff0000')
    expect(getHexFromCssVar('--color-primary')).toBe('#ff0000')
    document.documentElement.style.setProperty('--color-primary', 'rgb(0, 255, 0)')
    expect(getHexFromCssVar('--color-primary')).toBe('#00ff00')
  })
})
