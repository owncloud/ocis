/**
 * Converts hex to rgb
 * @param {string} hex: The hex color value (with or without #)
 * @return {Array<number>|null} The converted rgb array
 **/
export function hexToRgb(hex: string): Array<number> | null {
  hex = hex.startsWith('#') ? hex : `#${hex}`
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result ? [parseInt(result[1], 16), parseInt(result[2], 16), parseInt(result[3], 16)] : null
}

/**
 * Converts rgb to hex
 * @param {Array<number>} rgbArray: The rgb color
 * @return {string} The converted hex value
 **/
export function rgbToHex(rgbArray: Array<number>): string {
  const [r, g, b] = rgbArray
  const rHex = r.toString(16).padStart(2, '0')
  const gHex = g.toString(16).padStart(2, '0')
  const bHex = b.toString(16).padStart(2, '0')
  return `#${rHex}${gHex}${bHex}`
}

function shadeValue(value: number, percent: number): string {
  let a = (value * (100 + percent)) / 100
  a = a < 255 ? a : 255
  a = Math.round(a)
  return a.toString(16).length == 1 ? '0' + a.toString(16) : a.toString(16)
}

/**
 * Dim or brighten a hex color
 * @param {string} rgb: The rgb color value
 * @return {string} The brightened or dimmed hex color
 **/
export function calculateShadeColor(rgb: Array<number>, percent: number): string {
  const r = rgb[0]
  const g = rgb[1]
  const b = rgb[2]

  return `#${shadeValue(r, percent)}${shadeValue(g, percent)}${shadeValue(b, percent)}`
}

/**
 * Get the luminance of an rgb color
 * @param {Array<number>} rgb: The rgb value as an array
 * @return {Number} Returns value between 0 and 1, where 1 is white
 **/
export function getLuminanace(rgb: Array<number>): number {
  const mappedRgb = rgb.map((v) => {
    const val = v / 255
    return val <= 0.03928 ? val / 12.92 : ((val + 0.055) / 1.055) ** 2.4
  })
  return Number((0.2126 * mappedRgb[0] + 0.7152 * mappedRgb[1] + 0.0722 * mappedRgb[2]).toFixed(3))
}

/**
 * Get the contrast ratio between two rgb colors
 * @param {Array<number>} rgbColorA: The first rgb value as an array
 * @param {Array<number>} rgbColorB: The second rgb value as an array
 * @return {Number} Returns value between 1 and 21, where 1 is no contrast and 21 is max contrast
 **/
export function getContrastRatio(rgbColorA: Array<number>, rgbColorB: Array<number>): number {
  const lumA = getLuminanace(rgbColorA)
  const lumB = getLuminanace(rgbColorB)

  return (Math.max(lumA, lumB) + 0.05) / (Math.min(lumA, lumB) + 0.05)
}

/**
 * Hashes a string into an integer, deterministically and with no state.
 * The same string always yields the same number, in any browser, on any device, for any user.
 * The final `>>> 0` normalises the result so callers never have to: the loop's `<< 5` truncates
 * to int32 and can land in the sign bit, while the addition that follows can exceed int32
 * entirely, so without it the return value is neither reliably positive nor reliably 32 bit.
 * Normalising keeps the low 32 bits the loop produced, so hashes are unchanged.
 * @param {string} value: Can be any string
 * @return {Number} Returns an unsigned 32 bit integer
 **/
export function hashString(value: string): number {
  let hash = 0
  for (let i = 0; i < value.length; i++) {
    hash = value.charCodeAt(i) + ((hash << 5) - hash)
  }
  return hash >>> 0
}

/**
 * Gives you a random hashed color for a string, e.g. if you give it 'owncloud' it will always return the same color
 * @param {string} name: Can be any string
 * @return {string} Returns a hex color
 **/
export function generateHashedColorForString(name: string): string {
  return `#${(hashString(name) & 0x00ffffff).toString(16).toUpperCase()}`
}

/**
 * Picks a stable index into a list of `length` entries for the given name.
 * @param {string} name: Can be any string
 * @param {Number} length: Number of entries to spread names across
 * @return {Number} Returns an index in [0, length), or -1 if there is nothing to index
 **/
export function hashToIndex(name: string, length: number): number {
  if (!Number.isInteger(length) || length <= 0) {
    return -1
  }
  return hashString(name) % length
}

/**
 * The css vars a themed tag colour is emitted as: the fill, and the text colour paired with it.
 *
 * Both come from the theme (see the theme store, which derives the `-text` half from the fill and
 * the theme's own text colours), so a component only needs the index to colour a chip legibly —
 * no contrast maths, and no reading colours back out of the stylesheet.
 *
 * @param {Number} index: Index into the theme's tag colour list, as returned by `hashToIndex`
 * @return {Object} Returns the two var() references, or null if the index cannot name a colour
 **/
export function tagColorVarsFor(
  index: number | null | undefined
): { fillColor: string; textColor: string } | null {
  if (!Number.isInteger(index) || index < 0) {
    return null
  }
  return {
    fillColor: `var(--oc-color-tag-${index})`,
    textColor: `var(--oc-color-tag-${index}-text,currentColor)`
  }
}

/**
 * Picks whichever of two text colours is more legible on the given background.
 *
 * Both candidates are passed in rather than assumed to be black and white, so that a theme decides
 * what its text looks like. Which of them is the lighter one is not assumed either — only the
 * contrast decides.
 *
 * @param {string} background: The hex background colour the text sits on
 * @param {string} candidateA: One of the two text colours to choose between
 * @param {string} candidateB: The other one
 * @return {string} Returns whichever candidate contrasts more with the background, or null if any
 * of the three cannot be measured
 **/
export function pickReadableTextColor(
  background: string,
  candidateA: string,
  candidateB: string
): string | null {
  const backgroundRgb = hexToRgb(background)
  const rgbA = hexToRgb(candidateA)
  const rgbB = hexToRgb(candidateB)
  if (!backgroundRgb || !rgbA || !rgbB) {
    return null
  }
  return getContrastRatio(backgroundRgb, rgbA) >= getContrastRatio(backgroundRgb, rgbB)
    ? candidateA
    : candidateB
}

/**
 * Adjusts a given color to match the contrast ratio of another color
 * @param {Array<number>} targetColorRgb: color to adjust
 * @param {Array<number>} associatedColorRgb: brightest reference color
 * @param {Array<number>} desiredRatio: desired contrast ratio
 * @return {string} Returns a rgb color array
 **/
export function setDesiredContrastRatio(
  targetColorRgb: Array<number>,
  associatedColorRgb: Array<number>,
  desiredRatio: number
): Array<number> {
  const ratio = getContrastRatio(targetColorRgb, associatedColorRgb)
  if (Math.abs(desiredRatio - ratio) <= 0.3) {
    return targetColorRgb
  }
  if (ratio < desiredRatio) {
    return setDesiredContrastRatio(
      hexToRgb(calculateShadeColor(targetColorRgb, -1)),
      associatedColorRgb,
      desiredRatio
    )
  }
  return setDesiredContrastRatio(
    hexToRgb(calculateShadeColor(targetColorRgb, 1)),
    associatedColorRgb,
    desiredRatio
  )
}

/**
 * Convert a css rgb value like rgb(255, 255, 255) to a hex value like #FFFFFF
 * works also with rgba()
 * @param {string} cssRgb: color to adjust
 * @return {string} Returns a hex color
 **/
export function cssRgbToHex(cssRgb: string): string {
  cssRgb = cssRgb.toLowerCase().replace(/\s/g, '')
  const isRgba = cssRgb.includes('rgba(')
  const values = cssRgb.match(isRgba ? /rgba?\(([^)]+)\)/ : /rgb?\(([^)]+)\)/)
  if (!values) {
    return '#000000'
  }
  const colorArray = values[1].split(',')
  let alpha = ''
  if (isRgba) {
    alpha = Math.round(parseFloat(colorArray.pop()) * 255).toString(16)
    if (alpha.length === 1) {
      alpha = '0' + alpha
    }
  }
  const hex = colorArray
    .map((value) => {
      const intValue = parseInt(value, 10)
      const hexValue = intValue.toString(16)
      return hexValue.length === 1 ? '0' + hexValue : hexValue
    })
    .join('')
  if (isRgba) {
    return `#${hex}${alpha}`
  } else {
    return `#${hex}`
  }
}

/**
 * Get the hex value of a css var()
 * @param {string} variable: The css var name e.g. var(--color-primary) or --color-primary
 * @return {string} Returns a hex color
 **/
export function getHexFromCssVar(variable: string): string {
  if (!variable) {
    return ''
  }
  // if color is a hex value, return it
  if (variable.startsWith('#')) {
    return variable
  }
  const varName = variable.match(/var\(([^)]+)\)/)?.[1] || variable
  const result = getComputedStyle(document.documentElement).getPropertyValue(varName)
  // if css var is hex value, return it
  if (result.startsWith('#')) {
    return result
  }
  return cssRgbToHex(result)
}
