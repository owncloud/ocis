import { Token } from 'marked'
import { PDFDocument, PDFFont, PDFImage, RGB } from 'pdf-lib'
import {
  splitTextToFit,
  extractTextFromTokens,
  sanitizeText,
  parseInlineTokens,
  embedImage,
  pdfColorToCssRgb,
  parseImageAttributes
} from '../../../../../src/composables/webWorkers/exportAsPdfWorker/helpers'

vi.mock('pdf-lib', () => ({
  PDFDocument: vi.fn(),
  PDFFont: vi.fn(),
  PDFImage: vi.fn(),
  RGB: vi.fn(),
  rgb: vi.fn((r, g, b) => ({ red: r, green: g, blue: b }))
}))

describe('export as PDF worker helpers', () => {
  describe('splitTextToFit', () => {
    let mockFont: PDFFont

    beforeEach(() => {
      mockFont = {
        widthOfTextAtSize: vi.fn((text: string, size: number) => text.length * size * 0.6)
      } as any
    })

    it('should split text into lines that fit within max width', () => {
      const text = 'This is a long text that needs to be split'
      const result = splitTextToFit(text, mockFont, 12, 100)

      expect(result).toEqual(['This is a', 'long text', 'that needs to', 'be split'])
    })

    it('should handle single word that fits', () => {
      const text = 'short'
      const result = splitTextToFit(text, mockFont, 12, 100)

      expect(result).toEqual(['short'])
    })

    it('should handle empty text', () => {
      const text = ''
      const result = splitTextToFit(text, mockFont, 12, 100)

      expect(result).toEqual([])
    })

    it('should handle single word that exceeds max width', () => {
      const text = 'verylongwordthatexceedsmaxwidth'
      const result = splitTextToFit(text, mockFont, 12, 50)

      expect(result).toEqual(['verylo', 'ngword', 'thatex', 'ceedsm', 'axwidt', 'h'])
    })
  })

  describe('extractTextFromTokens', () => {
    it('should extract text from link tokens', () => {
      const tokens: Token[] = [
        {
          type: 'link',
          text: 'Click here',
          href: 'https://example.com',
          raw: '[Click here](https://example.com)'
        } as any
      ]

      const result = extractTextFromTokens(tokens)
      expect(result).toBe('https://example.com')
    })

    it('should extract text from text tokens', () => {
      const tokens: Token[] = [
        {
          type: 'text',
          text: 'Hello world',
          raw: 'Hello world'
        } as any
      ]

      const result = extractTextFromTokens(tokens)
      expect(result).toBe('Hello world')
    })

    it('should fall back to raw content for unknown tokens', () => {
      const tokens: Token[] = [
        {
          type: 'unknown',
          raw: 'raw content'
        } as any
      ]

      const result = extractTextFromTokens(tokens)
      expect(result).toBe('raw content')
    })

    it('should handle mixed token types', () => {
      const tokens: Token[] = [
        {
          type: 'text',
          text: 'Check out ',
          raw: 'Check out '
        } as any,
        {
          type: 'link',
          text: 'this link',
          href: 'https://example.com',
          raw: '[this link](https://example.com)'
        } as any
      ]

      const result = extractTextFromTokens(tokens)
      expect(result).toBe('Check out https://example.com')
    })
  })

  describe('sanitizeText', () => {
    it('should replace typographic characters with ASCII equivalents', () => {
      expect(
        sanitizeText('Here\u2019s a \u201cquote\u201d with an em\u2014dash and ellipsis\u2026')
      ).toBe('Here\'s a "quote" with an em--dash and ellipsis...')
    })

    it('should replace all types of quotes', () => {
      expect(sanitizeText('\u2018\u2019\u201c\u201d')).toBe('\'\'"\"')
    })

    it('should replace all types of dashes', () => {
      expect(sanitizeText('\u2014\u2013\u2011')).toBe('----')
    })

    it('should replace non-breaking spaces', () => {
      expect(sanitizeText('word word')).toBe('word word')
    })

    it('should remove emojis', () => {
      const result = sanitizeText('Hello 😀 world 🌍')

      expect(result).not.toContain('😀')
      expect(result).not.toContain('🌍')
    })

    it('should handle empty string', () => {
      expect(sanitizeText('')).toBe('')
    })

    it('should handle string with no special characters', () => {
      const input = 'Regular text with no special characters'
      expect(sanitizeText(input)).toBe(input)
    })
  })

  describe('parseInlineTokens', () => {
    it('should parse text tokens', () => {
      const tokens: Token[] = [
        {
          type: 'text',
          text: 'Hello world',
          raw: 'Hello world'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'Hello world',
        bold: false,
        italic: false,
        code: false,
        subscript: false,
        superscript: false,
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse strong tokens', () => {
      const tokens: Token[] = [
        {
          type: 'strong',
          text: 'Bold text',
          raw: '**Bold text**'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'Bold text',
        bold: true,
        italic: false,
        code: false,
        subscript: false,
        superscript: false,
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse em tokens', () => {
      const tokens: Token[] = [
        {
          type: 'em',
          text: 'Italic text',
          raw: '*Italic text*'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'Italic text',
        bold: false,
        italic: true,
        code: false,
        subscript: false,
        superscript: false,
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse codespan tokens with color', () => {
      const tokens: Token[] = [
        {
          type: 'codespan',
          text: 'code',
          raw: '`code`'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'code',
        bold: false,
        italic: false,
        code: true,
        subscript: false,
        superscript: false,
        color: { red: 0.7, green: 0.1, blue: 0.1 },
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse sub tokens', () => {
      const tokens: Token[] = [
        {
          type: 'sub',
          text: 'subscript',
          raw: '<sub>subscript</sub>'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'subscript',
        bold: false,
        italic: false,
        code: false,
        subscript: true,
        superscript: false,
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse sup tokens', () => {
      const tokens: Token[] = [
        {
          type: 'sup',
          text: 'superscript',
          raw: '<sup>superscript</sup>'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'superscript',
        bold: false,
        italic: false,
        code: false,
        subscript: false,
        superscript: true,
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse link tokens with color', () => {
      const tokens: Token[] = [
        {
          type: 'link',
          text: 'Link text',
          href: 'https://example.com',
          raw: '[Link text](https://example.com)'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'https://example.com',
        bold: false,
        italic: false,
        code: false,
        subscript: false,
        superscript: false,
        color: { red: 0, green: 0, blue: 0.8 },
        underline: false,
        strikeThrough: false
      })
    })

    it('should parse u tokens', () => {
      const tokens: Token[] = [
        {
          type: 'u',
          text: 'underlined',
          raw: '<u>underlined</u>'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'underlined',
        bold: false,
        italic: false,
        code: false,
        subscript: false,
        superscript: false,
        underline: true,
        strikeThrough: false
      })
    })

    it('should parse del tokens', () => {
      const tokens: Token[] = [
        {
          type: 'del',
          text: 'strikethrough',
          raw: '~~strikethrough~~'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'strikethrough',
        bold: false,
        italic: false,
        code: false,
        subscript: false,
        superscript: false,
        underline: false,
        strikeThrough: true
      })
    })

    it('should handle unknown tokens with text property', () => {
      const tokens: Token[] = [
        {
          type: 'unknown',
          text: 'unknown text'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(1)
      expect(result[0]).toEqual({
        text: 'unknown text',
        bold: false,
        italic: false,
        code: false,
        subscript: false,
        superscript: false,
        underline: false,
        strikeThrough: false
      })
    })

    it('should handle multiple tokens', () => {
      const tokens: Token[] = [
        {
          type: 'text',
          text: 'Normal ',
          raw: 'Normal '
        } as any,
        {
          type: 'strong',
          text: 'bold',
          raw: '**bold**'
        } as any,
        {
          type: 'text',
          text: ' and ',
          raw: ' and '
        } as any,
        {
          type: 'em',
          text: 'italic',
          raw: '*italic*'
        } as any
      ]

      const result = parseInlineTokens(tokens)

      expect(result).toHaveLength(4)
      expect(result[0].text).toBe('Normal ')
      expect(result[0].bold).toBe(false)
      expect(result[1].text).toBe('bold')
      expect(result[1].bold).toBe(true)
      expect(result[2].text).toBe(' and ')
      expect(result[3].text).toBe('italic')
      expect(result[3].italic).toBe(true)
    })
  })

  describe('embedImage', () => {
    let mockPdfDoc: PDFDocument
    let mockEmbeddedImage: PDFImage

    beforeEach(() => {
      mockEmbeddedImage = {
        width: 100,
        height: 80
      } as any

      mockPdfDoc = {
        embedPng: vi.fn().mockResolvedValue(mockEmbeddedImage),
        embedJpg: vi.fn().mockResolvedValue(mockEmbeddedImage)
      } as any

      global.atob = vi.fn((str) => {
        return str
          .split('')
          .map((c) => String.fromCharCode(c.charCodeAt(0)))
          .join('')
      })
    })

    it('should embed PNG image', async () => {
      const dataUri =
        'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg=='

      const result = await embedImage(mockPdfDoc, dataUri)

      expect(mockPdfDoc.embedPng).toHaveBeenCalled()
      expect(result).toEqual({
        image: mockEmbeddedImage,
        width: 100,
        height: 80
      })
    })

    it('should embed JPEG image', async () => {
      const dataUri = 'data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAA=='

      const result = await embedImage(mockPdfDoc, dataUri)

      expect(mockPdfDoc.embedJpg).toHaveBeenCalled()
      expect(result).toEqual({
        image: mockEmbeddedImage,
        width: 100,
        height: 80
      })
    })

    it('should fallback to PNG for unknown image types', async () => {
      const dataUri = 'data:image/webp;base64,UklGRh4AAABXRUJQVlA4TBEAAAAvAAAAAAfQ//73v/+BiOh/AAA='

      const result = await embedImage(mockPdfDoc, dataUri)

      expect(mockPdfDoc.embedPng).toHaveBeenCalled()
      expect(result).toEqual({
        image: mockEmbeddedImage,
        width: 100,
        height: 80
      })
    })
  })

  describe('pdfColorToCssRgb', () => {
    it('should convert PDF color to CSS rgb string', () => {
      const color = { red: 0.5, green: 0.8, blue: 0.2 } as RGB
      const result = pdfColorToCssRgb(color)

      expect(result).toBe('rgb(128, 204, 51)')
    })

    it('should handle edge cases (0 and 1)', () => {
      const color = { red: 0, green: 1, blue: 0.5 } as RGB
      const result = pdfColorToCssRgb(color)

      expect(result).toBe('rgb(0, 255, 128)')
    })

    it('should round values correctly', () => {
      const color = { red: 0.999, green: 0.001, blue: 0.5001 } as RGB
      const result = pdfColorToCssRgb(color)

      expect(result).toBe('rgb(255, 0, 128)')
    })
  })

  describe('parseImageAttributes', () => {
    it('should parse display attribute', () => {
      const title = 'd=inline'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'inline',
        width: 0,
        height: 0,
        text: null
      })
    })

    it('should parse width and height attributes', () => {
      const title = 'w=100;h=50'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'block',
        width: 100,
        height: 50,
        text: null
      })
    })

    it('should parse all attributes together', () => {
      const title = 'd=inline;w=200;h=150'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'inline',
        width: 200,
        height: 150,
        text: null
      })
    })

    it('should ignore invalid key-value pairs', () => {
      const title = 'd=inline;invalid;w=100'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'inline',
        width: 100,
        height: 0,
        text: null
      })
    })

    it('should treat non-attribute string as regular title', () => {
      const title = 'This is a regular image title'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'block',
        width: 0,
        height: 0,
        text: 'This is a regular image title'
      })
    })

    it('should handle null title', () => {
      const result = parseImageAttributes(null)

      expect(result).toEqual({
        display: 'block',
        width: 0,
        height: 0,
        text: null
      })
    })

    it('should handle undefined title', () => {
      const result = parseImageAttributes(undefined)

      expect(result).toEqual({
        display: 'block',
        width: 0,
        height: 0,
        text: null
      })
    })

    it('should handle empty string title', () => {
      const result = parseImageAttributes('')

      expect(result).toEqual({
        display: 'block',
        width: 0,
        height: 0,
        text: null
      })
    })

    it('should ignore invalid width and height values', () => {
      const title = 'w=invalid;h=notanumber;d=inline'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'inline',
        width: 0,
        height: 0,
        text: null
      })
    })

    it('should ignore unknown attributes', () => {
      const title = 'd=inline;unknown=value;w=100'
      const result = parseImageAttributes(title)

      expect(result).toEqual({
        display: 'inline',
        width: 100,
        height: 0,
        text: null
      })
    })
  })
})
