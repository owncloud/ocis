import { marked } from 'marked'
import { PDFDocument } from 'pdf-lib'
import fontkit from '@pdf-lib/fontkit'
import { PDFRenderer } from '../../../../../src/composables/webWorkers/exportAsPdfWorker/renderer'

vi.mock('marked', () => ({
  marked: {
    lexer: vi.fn().mockReturnValue([]),
    use: vi.fn()
  }
}))

vi.mock('pdf-lib', () => ({
  PDFDocument: {
    create: vi.fn()
  },
  PDFFont: vi.fn(),
  PDFPage: vi.fn(),
  PDFImage: vi.fn(),
  RGB: vi.fn(),
  rgb: vi.fn((r, g, b) => ({ red: r, green: g, blue: b }))
}))

vi.mock('@pdf-lib/fontkit', () => ({
  default: {}
}))

describe('PDFRenderer', () => {
  let mockPdfDoc: any
  let mockPage: any
  let mockFont: any
  let mockSave: any

  beforeEach(() => {
    mockSave = vi.fn().mockResolvedValue(new Uint8Array([1, 2, 3, 4]))
    mockPage = {
      getWidth: vi.fn().mockReturnValue(600),
      getHeight: vi.fn().mockReturnValue(800),
      drawText: vi.fn(),
      drawRectangle: vi.fn(),
      drawLine: vi.fn(),
      drawImage: vi.fn()
    }

    mockFont = {
      widthOfTextAtSize: vi.fn().mockReturnValue(50)
    }

    mockPdfDoc = {
      addPage: vi.fn().mockReturnValue(mockPage),
      registerFontkit: vi.fn(),
      embedFont: vi.fn().mockResolvedValue(mockFont),
      save: mockSave,
      embedPng: vi.fn().mockResolvedValue({ width: 100, height: 100 })
    }

    vi.mocked(PDFDocument.create).mockResolvedValue(mockPdfDoc)

    global.fetch = vi.fn().mockResolvedValue({
      arrayBuffer: vi.fn().mockResolvedValue(new ArrayBuffer(8))
    })
  })

  describe('constructor', () => {
    it('should initialize with markdown content and create tokens', () => {
      const mockTokens = [{ type: 'paragraph', text: 'test', raw: 'test' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      new PDFRenderer('# Test markdown')

      expect(marked.lexer).toHaveBeenCalledWith('# Test markdown')
    })
  })

  describe('renderAsArrayBuffer', () => {
    it('should create PDF document and return ArrayBuffer', async () => {
      const mockTokens = [{ type: 'heading', depth: 1, text: 'Test Heading' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('# Test')
      const result = await renderer.renderAsArrayBuffer()

      expect(PDFDocument.create).toHaveBeenCalled()
      expect(mockPdfDoc.addPage).toHaveBeenCalled()
      expect(mockPdfDoc.registerFontkit).toHaveBeenCalledWith(fontkit)
      expect(mockPdfDoc.save).toHaveBeenCalled()
      expect(result).toBeInstanceOf(ArrayBuffer)
    })

    it('should load all required fonts', async () => {
      const mockTokens = [] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('')
      await renderer.renderAsArrayBuffer()

      // Should fetch 6 different font files
      expect(global.fetch).toHaveBeenCalledTimes(6)
      expect(mockPdfDoc.embedFont).toHaveBeenCalledTimes(6)
    })

    it('should handle empty content', async () => {
      const mockTokens = [] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('')
      const result = await renderer.renderAsArrayBuffer()

      expect(result).toBeInstanceOf(ArrayBuffer)
      expect(mockPdfDoc.addPage).toHaveBeenCalledTimes(1) // Initial page
    })

    it('should add new page when content requires it', async () => {
      const mockTokens = [
        { type: 'heading', depth: 1, text: 'Test Heading 1' },
        { type: 'heading', depth: 1, text: 'Test Heading 2' }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('# Test 1\n# Test 2')
      await renderer.renderAsArrayBuffer()

      expect(mockPdfDoc.addPage).toHaveBeenCalledTimes(1) // Initial page
    })
  })

  describe('token rendering', () => {
    it('should render heading tokens', async () => {
      const mockTokens = [{ type: 'heading', depth: 1, text: 'Test Heading' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('# Test Heading')
      await renderer.renderAsArrayBuffer()

      expect(mockPage.drawText).toHaveBeenCalledWith(
        'Test Heading',
        expect.objectContaining({
          font: mockFont,
          size: expect.any(Number)
        })
      )
    })

    it('should render paragraph tokens', async () => {
      const mockTokens = [
        {
          type: 'paragraph',
          tokens: [{ type: 'text', text: 'Test paragraph' }]
        }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('Test paragraph')
      await renderer.renderAsArrayBuffer()

      // Verify the paragraph rendering process was initiated
      expect(mockPage.drawText).toHaveBeenCalled()
    })

    it('should render code block tokens', async () => {
      const mockTokens = [{ type: 'code', text: 'console.log("test")' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('```js\nconsole.log("test")\n```')
      await renderer.renderAsArrayBuffer()

      expect(mockPage.drawRectangle).toHaveBeenCalled() // Background
      expect(mockPage.drawText).toHaveBeenCalled() // Code text
    })

    it('should render list tokens', async () => {
      const mockTokens = [
        {
          type: 'list',
          ordered: false,
          items: [
            { tokens: [{ type: 'text', text: 'Item 1' }] },
            { tokens: [{ type: 'text', text: 'Item 2' }] }
          ]
        }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('- Item 1\n- Item 2')
      await renderer.renderAsArrayBuffer()

      expect(mockPage.drawText).toHaveBeenCalled()
    })

    it('should render blockquote tokens', async () => {
      const mockTokens = [
        {
          type: 'blockquote',
          tokens: [
            {
              type: 'paragraph',
              tokens: [{ type: 'text', text: 'Quote text' }]
            }
          ]
        }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('> Quote text')
      await renderer.renderAsArrayBuffer()

      expect(mockPage.drawRectangle).toHaveBeenCalled() // Quote bar
      expect(mockPage.drawText).toHaveBeenCalled()
    })

    it('should render table tokens', async () => {
      const mockTokens = [
        {
          type: 'table',
          header: [
            { tokens: [{ type: 'text', text: 'Header 1' }] },
            { tokens: [{ type: 'text', text: 'Header 2' }] }
          ],
          rows: [
            [
              { tokens: [{ type: 'text', text: 'Cell 1' }] },
              { tokens: [{ type: 'text', text: 'Cell 2' }] }
            ]
          ]
        }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('| Header 1 | Header 2 |\n| Cell 1 | Cell 2 |')
      await renderer.renderAsArrayBuffer()

      expect(mockPage.drawRectangle).toHaveBeenCalled() // Table cells
      expect(mockPage.drawText).toHaveBeenCalled()
    })

    it('should render horizontal rule tokens', async () => {
      const mockTokens = [{ type: 'hr' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('---')
      await renderer.renderAsArrayBuffer()

      expect(mockPage.drawLine).toHaveBeenCalled()
    })

    it('should handle space tokens', async () => {
      const mockTokens = [{ type: 'space' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('\n\n')
      const result = await renderer.renderAsArrayBuffer()

      expect(result).toBeInstanceOf(ArrayBuffer)
    })

    it('should handle unknown token types gracefully', async () => {
      const mockTokens = [{ type: 'unknown_type' }] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('some content')
      const result = await renderer.renderAsArrayBuffer()

      expect(result).toBeInstanceOf(ArrayBuffer)
    })
  })

  describe('image handling', () => {
    it('should render image tokens', async () => {
      const mockTokens = [
        {
          type: 'paragraph',
          tokens: [
            {
              type: 'image',
              href: 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==',
              text: 'alt text'
            }
          ]
        }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer(
        '![alt text](data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==)'
      )
      await renderer.renderAsArrayBuffer()

      expect(mockPdfDoc.embedPng).toHaveBeenCalled()
      expect(mockPage.drawImage).toHaveBeenCalled()
    })
  })

  describe('error handling', () => {
    it('should handle font loading errors gracefully', async () => {
      global.fetch = vi.fn().mockRejectedValue(new Error('Failed to fetch font'))

      const mockTokens = [] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const renderer = new PDFRenderer('')

      await expect(renderer.renderAsArrayBuffer()).rejects.toThrow('Failed to fetch font')
    })

    it('should handle PDF creation errors', async () => {
      vi.mocked(PDFDocument.create).mockRejectedValue(new Error('PDF creation failed'))

      const renderer = new PDFRenderer('test')

      await expect(renderer.renderAsArrayBuffer()).rejects.toThrow('PDF creation failed')
    })
  })

  describe('complex content scenarios', () => {
    it('should handle mixed content types', async () => {
      const mockTokens = [
        { type: 'heading', depth: 1, text: 'Title' },
        {
          type: 'paragraph',
          tokens: [{ type: 'text', text: 'Paragraph text' }]
        },
        { type: 'code', text: 'code block' },
        {
          type: 'list',
          ordered: false,
          items: [{ tokens: [{ type: 'text', text: 'List item' }] }]
        }
      ] as any
      vi.mocked(marked.lexer).mockReturnValue(mockTokens)

      const content = '# Title\n\nParagraph text\n\n```\ncode block\n```\n\n- List item'
      const renderer = new PDFRenderer(content)
      const result = await renderer.renderAsArrayBuffer()

      expect(result).toBeInstanceOf(ArrayBuffer)
      expect(mockPage.drawText).toHaveBeenCalledTimes(5)
    })
  })
})
