import {
  hexToRgb,
  rgbToHex,
  calculateShadeColor,
  getLuminanace,
  getContrastRatio,
  generateHashedColorForString,
  setDesiredContrastRatio,
  cssRgbToHex,
  getHexFromCssVar
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
