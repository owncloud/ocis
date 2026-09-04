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
  tagColorVarFor
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
  it('can return a negative number', () => {
    // Guards the `>>> 0` in hashToIndex: `hash << 5` overflows into the sign bit.
    expect(hashString('physics')).toBeLessThan(0)
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
  it('never returns a negative index for a name whose hash is negative', () => {
    expect(hashString('physics')).toBeLessThan(0)
    expect(hashToIndex('physics', 30)).toBeGreaterThanOrEqual(0)
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

describe('tagColorVarFor', () => {
  it('resolves a name to a zero-based css var', () => {
    expect(tagColorVarFor('physics', 30)).toBe('var(--oc-color-tag-13)')
  })
  it('agrees with hashToIndex', () => {
    expect(tagColorVarFor('invoice', 30)).toBe(`var(--oc-color-tag-${hashToIndex('invoice', 30)})`)
  })
  it('gives different tags on one file different colours', () => {
    expect(tagColorVarFor('Invoice', 30)).not.toBe(tagColorVarFor('invoice', 30))
  })
  it('returns an empty string for an empty list rather than an invalid var', () => {
    // theme.json ships no tagColorsList until design supplies the values, so this is the
    // real state of the app today, not a hypothetical.
    expect(tagColorVarFor('physics', 0)).toBe('')
    expect(tagColorVarFor('physics', undefined)).toBe('')
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
