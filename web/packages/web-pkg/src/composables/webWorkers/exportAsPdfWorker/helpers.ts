import { Token, Tokens } from 'marked'
import { PDFDocument, PDFFont, PDFImage, RGB } from 'pdf-lib'
import emojiRegex from 'emoji-regex'
import { PDF_THEME } from './pdfConfig'

/**
 * Represents a text segment with its formatting properties.
 *
 * This type is used to describe text content along with all its formatting
 * attributes such as bold, italic, code, subscript/superscript, color, and
 * text decorations like underline and strikethrough.
 */
export type TextSegment = {
  /** The actual text content */
  text: string
  /** Whether the text should be rendered in bold */
  bold: boolean
  /** Whether the text should be rendered in italic */
  italic: boolean
  /** Whether the text should be rendered as inline code */
  code: boolean
  /** Whether the text should be rendered as subscript */
  subscript: boolean
  /** Whether the text should be rendered as superscript */
  superscript: boolean
  /** Optional color for the text (defaults to theme text color) */
  color?: any
  /** Whether the text should have an underline */
  underline: boolean
  /** Whether the text should have a strikethrough */
  strikeThrough: boolean
}

/**
 * Available font weights for PDF text rendering.
 *
 * Corresponds to the different font variants loaded for the PDF document,
 * including regular text, emphasis, monospace, and combinations.
 */
export type FontWeight = 'regular' | 'bold' | 'italic' | 'boldItalic' | 'mono' | 'monoBold'

/**
 * Splits text into lines that fit within a specified maximum width using the given font and size.
 *
 * This function performs word-wrapping by measuring the width of text segments and breaking
 * lines when they exceed the maximum width. If a single word is too long to fit within the
 * maximum width, it will be split at character boundaries.
 *
 * @param text - The text to split into lines
 * @param font - The PDF font object used for measuring text width
 * @param fontSize - The font size in points
 * @param maxWidth - The maximum width in points that each line should not exceed
 * @returns Array of text lines that fit within the specified width
 */
export function splitTextToFit(text: string, font: PDFFont, fontSize: number, maxWidth: number) {
  const words = text.split(' ')
  const lines: string[] = []
  let currentLine = ''

  for (const word of words) {
    const testLine = currentLine + (currentLine ? ' ' : '') + word
    const width = font.widthOfTextAtSize(testLine, fontSize)

    if (width <= maxWidth) {
      currentLine = testLine
    } else {
      if (currentLine) {
        lines.push(currentLine)
        currentLine = ''
      }

      const wordWidth = font.widthOfTextAtSize(word, fontSize)
      if (wordWidth > maxWidth) {
        const wordLines = splitWordToFit(word, font, fontSize, maxWidth)
        lines.push(...wordLines.slice(0, -1))
        currentLine = wordLines[wordLines.length - 1]
      } else {
        currentLine = word
      }
    }
  }

  if (currentLine) {
    lines.push(currentLine)
  }

  return lines
}

/**
 * Splits a single word into lines that fit within the maximum width by breaking at character boundaries.
 *
 * @param word - The word to split
 * @param font - The PDF font object used for measuring text width
 * @param fontSize - The font size in points
 * @param maxWidth - The maximum width in points that each line should not exceed
 * @returns Array of text lines that fit within the specified width
 */
function splitWordToFit(word: string, font: PDFFont, fontSize: number, maxWidth: number): string[] {
  const lines: string[] = []
  let currentLine = ''

  for (let i = 0; i < word.length; i++) {
    const char = word[i]
    const testLine = currentLine + char
    const width = font.widthOfTextAtSize(testLine, fontSize)

    if (width <= maxWidth) {
      currentLine = testLine
    } else {
      if (currentLine) {
        lines.push(currentLine)
      }

      currentLine = char
    }
  }

  if (currentLine) {
    lines.push(currentLine)
  }

  return lines
}

/**
 * Determines the appropriate font to use for a text segment based on its formatting properties.
 *
 * This function selects the correct font variant based on the text segment's formatting:
 * - Monospace font for code segments
 * - Bold-italic font for segments with both bold and italic formatting
 * - Bold font for bold segments
 * - Italic font for italic segments
 * - Regular font as the default
 *
 * @param segment - The text segment with formatting properties
 * @param fonts - Record of font weights and their corresponding PDF fonts
 * @param preferredFont - The font which should override the default font
 * @returns The appropriate PDF font for the segment
 */
export function getFontForSegment(
  segment: TextSegment,
  fonts: Record<FontWeight, PDFFont>,
  preferredFont?: FontWeight
): PDFFont {
  if (preferredFont && preferredFont === 'italic' && segment.bold) {
    return fonts['boldItalic']
  }

  if (preferredFont && Object.prototype.hasOwnProperty.call(fonts, preferredFont)) {
    return fonts[preferredFont]
  }

  if (segment.code) {
    return fonts['mono']
  }

  if (segment.bold && segment.italic) {
    return fonts['boldItalic']
  }

  if (segment.bold) {
    return fonts['bold']
  }

  if (segment.italic) {
    return fonts['italic']
  }

  return fonts['regular']
}

/**
 * Extracts plain text content from an array of markdown tokens.
 *
 * This function processes different token types and converts them to readable text:
 * - Links are converted to "text (url)" format
 * - Text tokens are extracted directly
 * - Other tokens fall back to their raw content
 *
 * @param tokens - Array of markdown tokens to extract text from
 * @returns Concatenated plain text representation of all tokens
 */
export function extractTextFromTokens(tokens: Token[]): string {
  return tokens
    .map((token) => {
      if (token.type === 'link') {
        return token.href
      }

      if ('text' in token) {
        return token.text
      }

      return token.raw || ''
    })
    .join('')
}

/**
 * Partitions tokens into text and image tokens.
 *
 * This function separates an array of markdown tokens into two categories:
 * image tokens and all other tokens, allowing for different processing
 * of images and text content in layout algorithms.
 *
 * @param tokens - Array of markdown tokens to partition
 * @returns Object containing separated text and image tokens
 */
export function partitionTokens(tokens: Token[]): { textTokens: Token[]; imageTokens: Token[] } {
  return tokens.reduce(
    (acc, token) => {
      if (token.type === 'image') {
        acc.imageTokens.push(token)
      } else {
        acc.textTokens.push(token)
      }

      return acc
    },
    { textTokens: [] as Token[], imageTokens: [] as Token[] }
  )
}

/**
 * Sanitizes text by converting typographic characters to ASCII equivalents for PDF compatibility
 * and fixing common markdown formatting issues.
 *
 * This function is necessary because:
 * 1. The PDF generation uses StandardFonts (Helvetica, Courier) from pdf-lib, which have limited
 *    Unicode character support. Typographic characters like curly quotes and em/en dashes are not
 *    supported by these fonts and would cause rendering issues or PDF generation failures.
 * 2. User-provided markdown may contain formatting issues that break parsers (e.g., trailing
 *    whitespace after code fence closings).
 *
 * The replacements maintain semantic meaning while ensuring compatibility:
 * - Typographic quotes → straight quotes (same meaning)
 * - Em/en dashes → hyphens (similar visual effect)
 * - Ellipsis → three dots (same meaning)
 * - Trailing whitespace after code fences → removed (fixes parsing)
 *
 * @param text - The input text that may contain typographic characters and markdown issues
 * @returns The sanitized text with typographic characters replaced and markdown issues fixed
 *
 * @example
 * ```typescript
 * sanitizeText("Here's a "quote" with an em—dash and ellipsis…")
 * // Returns: "Here's a "quote" with an em--dash and ellipsis..."
 *
 * sanitizeText("```\ncode\n```\t\nMore content")
 * // Returns: "```\ncode\n```\nMore content"
 * ```
 */
export function sanitizeText(text: string): string {
  return text
    .replaceAll('…', '...')
    .replaceAll("'", "'")
    .replaceAll("'", "'")
    .replaceAll('’', "'")
    .replaceAll('‘', "'")
    .replaceAll('"', '"')
    .replaceAll('"', '"')
    .replaceAll('“', '"')
    .replaceAll('”', '"')
    .replaceAll('—', '--')
    .replaceAll('–', '-')
    .replaceAll(' ', ' ')
    .replaceAll('‑', '-')
    .replace(emojiRegex(), '')
    .replace(/^(```+)[\t ]+$/gm, '$1')
    .replace(/^(```+\w*)[\t ]+$/gm, '$1')
    .replaceAll('⋅', '  ')
}

/**
 * Parses inline markdown tokens into TextSegment objects with formatting information.
 *
 * This function converts markdown inline tokens (text, strong, em, codespan, sub, sup, link)
 * into TextSegment objects that contain the text content and formatting flags. Each segment
 * represents a piece of text with consistent formatting that can be rendered with the
 * appropriate font and styling.
 *
 * @param tokens - Array of inline markdown tokens to parse
 * @returns Array of TextSegment objects with formatting information
 */
export function parseInlineTokens(tokens: Token[]): TextSegment[] {
  const segments: TextSegment[] = []

  for (const token of tokens) {
    switch (token.type) {
      case 'text':
        segments.push({
          text: token.text,
          bold: false,
          italic: false,
          code: false,
          subscript: false,
          superscript: false,
          underline: false,
          strikeThrough: false
        })
        break
      case 'strong':
        segments.push({
          text: token.text,
          bold: true,
          italic: false,
          code: false,
          subscript: false,
          superscript: false,
          underline: false,
          strikeThrough: false
        })
        break
      case 'em':
        segments.push({
          text: token.text,
          bold: false,
          italic: true,
          code: false,
          subscript: false,
          superscript: false,
          underline: false,
          strikeThrough: false
        })
        break
      case 'codespan':
        segments.push({
          text: token.text,
          bold: false,
          italic: false,
          code: true,
          subscript: false,
          superscript: false,
          color: PDF_THEME.color.codeSpan,
          underline: false,
          strikeThrough: false
        })
        break
      case 'sub':
        segments.push({
          text: token.text,
          bold: false,
          italic: false,
          code: false,
          subscript: true,
          superscript: false,
          underline: false,
          strikeThrough: false
        })
        break
      case 'sup':
        segments.push({
          text: token.text,
          bold: false,
          italic: false,
          code: false,
          subscript: false,
          superscript: true,
          underline: false,
          strikeThrough: false
        })
        break
      case 'link':
        segments.push({
          text: token.href,
          bold: false,
          italic: false,
          code: false,
          subscript: false,
          superscript: false,
          color: PDF_THEME.color.link,
          underline: false,
          strikeThrough: false
        })
        break
      case 'u':
        segments.push({
          text: token.text,
          bold: false,
          italic: false,
          code: false,
          subscript: false,
          superscript: false,
          underline: true,
          strikeThrough: false
        })
        break
      case 'del':
        segments.push({
          text: token.text,
          bold: false,
          italic: false,
          code: false,
          subscript: false,
          superscript: false,
          underline: false,
          strikeThrough: true
        })
        break
      default:
        if ((token as Tokens.Text).text) {
          segments.push({
            text: (token as Tokens.Text).text,
            bold: false,
            italic: false,
            code: false,
            subscript: false,
            superscript: false,
            underline: false,
            strikeThrough: false
          })
        }
    }
  }

  return segments
}

/**
 * Fetches an image from a data URI and embeds it into the PDF document.
 *
 * This function handles data URIs for both PNG and JPEG formats, decoding the base64
 * content and embedding it into the PDF document. For unsupported formats, it attempts
 * to embed as PNG. The function returns the embedded image along with its dimensions.
 *
 * @param pdfDoc - The PDF document to embed the image into
 * @param imageUrl - The data URI of the image to embed
 * @returns Promise resolving to an object with embedded image, width, and height
 * @throws Error if the image format is unsupported or embedding fails
 */
export async function embedImage(pdfDoc: PDFDocument, imageUrl: string) {
  const parts = imageUrl.split(',')
  const meta = parts[0]
  const base64Data = parts[1]
  const binaryStr = atob(base64Data)
  const len = binaryStr.length
  const bytes = new Uint8Array(len)

  for (let i = 0; i < len; i++) {
    bytes[i] = binaryStr.charCodeAt(i)
  }

  const imageBytes = bytes.buffer
  let embeddedImage: PDFImage

  if (meta.includes('image/png')) {
    embeddedImage = await pdfDoc.embedPng(imageBytes)
  } else if (meta.includes('image/jpeg')) {
    embeddedImage = await pdfDoc.embedJpg(imageBytes)
  } else {
    embeddedImage = await pdfDoc.embedPng(imageBytes)
  }

  return {
    image: embeddedImage,
    width: embeddedImage.width,
    height: embeddedImage.height
  }
}

/**
 * Converts a PDF color to CSS rgb() string.
 *
 * This utility function converts PDF-lib RGB color objects to CSS-compatible
 * rgb() strings for use in HTML elements during intermediate processing.
 *
 * @param color - The PDF RGB color object to convert
 * @returns The CSS rgb() string representation
 */
export function pdfColorToCssRgb(color: RGB): string {
  const r = Math.round(color.red * 255)
  const g = Math.round(color.green * 255)
  const b = Math.round(color.blue * 255)
  return `rgb(${r}, ${g}, ${b})`
}

/**
 * Attributes parsed from image titles for customizing image display.
 *
 * These attributes control how images are rendered in the PDF, including
 * display mode, dimensions, and whether to show title text.
 */
type ImageAttributes = {
  /** Display mode: 'inline' for inline images, 'block' for block images */
  display: 'inline' | 'block'
  /** Desired width in pixels (0 for original width) */
  width: number
  /** Desired height in pixels (0 for original height) */
  height: number
  /** Title text to display below the image (null to hide) */
  text: string | null
}

/**
 * Parses an image title for custom attributes like display mode, width, and height.
 *
 * Attributes are expected in a semicolon-separated key-value format (e.g., "d=inline;w=50;h=20").
 * If the title does not conform to this format, it is treated as a regular image title.
 * Supported attributes:
 * - d: display mode ('inline' or 'block')
 * - w: width in pixels
 * - h: height in pixels
 *
 * @param title - The title string from the markdown image token
 * @returns An ImageAttributes object with parsed values or defaults
 */
export function parseImageAttributes(title: string | null | undefined): ImageAttributes {
  const result: ImageAttributes = {
    display: 'block',
    width: 0,
    height: 0,
    text: title || null
  }

  if (!title || !title.includes('=')) {
    return result
  }

  const pairs = title.split(';')
  let isAttributeString = false

  for (const pair of pairs) {
    const parts = pair.split('=')

    if (parts.length !== 2) {
      continue
    }

    isAttributeString = true
    const key = parts[0].trim()
    const value = parts[1].trim()

    switch (key) {
      case 'd':
        if (value === 'inline') {
          result.display = 'inline'
        }
        break
      case 'w':
        result.width = parseInt(value, 10) || 0
        break
      case 'h':
        result.height = parseInt(value, 10) || 0
        break
    }
  }

  if (isAttributeString) {
    result.text = null
  }

  return result
}
