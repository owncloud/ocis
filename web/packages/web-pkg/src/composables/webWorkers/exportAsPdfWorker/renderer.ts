import { marked, Token, Tokens } from 'marked'
import { PDFDocument, PDFFont, PDFPage, RGB } from 'pdf-lib'
import fontkit from '@pdf-lib/fontkit'
import {
  extractTextFromTokens,
  embedImage,
  getFontForSegment,
  parseImageAttributes,
  parseInlineTokens,
  partitionTokens,
  splitTextToFit,
  FontWeight
} from './helpers'
import { PDF_THEME } from './pdfConfig'
import { markedExtensions } from './markedExtensions'
import { captureException } from '@sentry/vue'

marked.use({ extensions: markedExtensions })

/**
 * Result object returned by token rendering methods.
 */
type RenderResult = {
  /** Whether a new page is needed to render the content */
  needsNewPage: boolean
}

/**
 * Font URLs for different font weights used in PDF generation.
 *
 * All fonts are loaded from the Ubuntu font family hosted on pdf-lib.js.org
 * to ensure consistent typography and proper Unicode support.
 */
const FONT_URLS = Object.freeze({
  regular: '/fonts/ubuntu/Ubuntu-R.ttf',
  bold: '/fonts/ubuntu/Ubuntu-B.ttf',
  italic: '/fonts/ubuntu/Ubuntu-RI.ttf',
  boldItalic: '/fonts/ubuntu/Ubuntu-BI.ttf',
  mono: '/fonts/ubuntu/UbuntuMono-R.ttf',
  monoBold: '/fonts/ubuntu/UbuntuMono-B.ttf'
})

/**
 * PDF renderer class that converts markdown content to PDF documents.
 *
 * This class handles the complete process of parsing markdown content, rendering
 * it to PDF pages with proper typography, layout, and formatting. It supports
 * all common markdown elements including headings, paragraphs, code blocks,
 * lists, tables, blockquotes, images, and horizontal rules.
 */
export class PDFRenderer {
  #pdfDoc: PDFDocument
  #page: PDFPage
  #maxWidth: number
  #tokens: Token[]
  #pageHeight: number
  #fonts = {} as Record<FontWeight, PDFFont>
  #yPosition = 0

  /**
   * Creates a new PDF renderer instance.
   *
   * @param content - The markdown content to render as PDF
   */
  constructor(content: string) {
    this.#tokens = marked.lexer(content)
  }

  /**
   * Converts markdown text to a PDF document.
   *
   * This method processes markdown content through the following steps:
   * 1. Sanitizes text for PDF compatibility (typographic characters)
   * 2. Creates a new PDF document and loads fonts
   * 3. Parses markdown into tokens
   * 4. Renders each token type to the PDF
   * 5. Handles page breaks automatically
   *
   * @returns Promise resolving to the PDF as an ArrayBuffer
   */
  async renderAsArrayBuffer() {
    this.#pdfDoc = await PDFDocument.create()

    this.#addNewPage()
    await this.#loadFonts()

    for (const token of this.#tokens) {
      let result = await this.#renderToken(token)

      if (result.needsNewPage) {
        this.#addNewPage()
        result = await this.#renderToken(token)
      }

      if (this.#yPosition < PDF_THEME.layout.margin) {
        this.#page = this.#pdfDoc.addPage()
        this.#yPosition = this.#pageHeight - PDF_THEME.layout.margin
      }
    }

    const pdfBytes = await this.#pdfDoc.save()
    return pdfBytes.buffer as ArrayBuffer
  }

  /**
   * Renders a markdown token to the PDF page by delegating to the appropriate renderer method.
   *
   * This method acts as a dispatcher that routes different token types to their specific
   * rendering methods. It handles all supported markdown token types including headings,
   * paragraphs, code blocks, lists, blockquotes, tables, horizontal rules, and spaces.
   *
   * @param token - The markdown token to render
   * @returns Promise or RenderResult indicating if a new page is needed
   */
  #renderToken(token: Token): Promise<RenderResult> | RenderResult {
    switch (token.type) {
      case 'heading':
        return this.#renderHeading(token as Tokens.Heading)
      case 'paragraph':
        return this.#renderParagraph({ token: token as Tokens.Paragraph })
      case 'code':
        return this.#renderCodeBlock(token as Tokens.Code)
      case 'list':
        return this.#renderList(token as Tokens.List)
      case 'blockquote':
        return this.#renderBlockquote(token as Tokens.Blockquote)
      case 'table':
        return this.#renderTable(token as Tokens.Table)
      case 'hr':
        return this.#renderHorizontalRule()
      case 'space':
        this.#yPosition -= PDF_THEME.spacing.md
      default:
        return { needsNewPage: false }
    }
  }

  /**
   * Renders a heading token to the PDF page with appropriate font size and spacing.
   *
   * Headings are rendered in bold font with size determined by the heading level.
   * Text is automatically wrapped to fit within page margins and proper spacing
   * is added before and after the heading.
   *
   * @param token - The markdown heading token containing text and depth level
   * @returns RenderResult indicating if a new page is needed
   */
  #renderHeading(token: Tokens.Heading): RenderResult {
    const fontSize = PDF_THEME.font[`h${token.depth}`]
    const lineHeight = fontSize * 1.4

    const lines = splitTextToFit(
      token.text,
      this.#fonts['bold'],
      fontSize,
      this.#page.getWidth() - PDF_THEME.layout.margin * 2
    )

    if (this.#yPosition - lines.length * lineHeight < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    this.#yPosition -= fontSize

    for (const line of lines) {
      this.#page.drawText(line, {
        x: PDF_THEME.layout.margin,
        y: this.#yPosition,
        size: fontSize,
        font: this.#fonts['bold'],
        color: PDF_THEME.color.text
      })
      this.#yPosition -= lineHeight
    }

    this.#yPosition -= PDF_THEME.spacing.md
    return { needsNewPage: false }
  }

  /**
   * Renders a paragraph token to the PDF page, handling inline formatting and images.
   *
   * This method processes all inline elements within the paragraph including text,
   * images, formatting (bold, italic, code), links, and special characters.
   * It handles line wrapping, inline images, and maintains proper spacing.
   *
   * @param options.token - The markdown paragraph token containing inline tokens
   * @param options.margin - The left margin for the paragraph
   * @param options.preferredFont - The font which should override the default font
   * @returns Promise resolving to RenderResult indicating if a new page is needed
   */
  async #renderParagraph({
    token,
    margin = PDF_THEME.layout.margin,
    preferredFont,
    includeBottomSpacing = true
  }: {
    token: Tokens.Paragraph
    margin?: number
    preferredFont?: FontWeight
    includeBottomSpacing?: boolean
  }): Promise<RenderResult> {
    const lineHeight = PDF_THEME.font.lineHeight

    let currentX = margin

    const wrapLine = () => {
      this.#yPosition -= lineHeight
      currentX = margin

      return this.#yPosition < PDF_THEME.layout.margin
    }

    if (token.tokens.length === 1 && token.tokens[0].type === 'image') {
      const imageToken = token.tokens[0] as Tokens.Image
      const attrs = parseImageAttributes(imageToken.text)

      if (attrs.display === 'block') {
        return this.#renderImage(imageToken, margin, this.#maxWidth)
      }
    }

    if (this.#yPosition - lineHeight < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    for (const inlineToken of token.tokens) {
      if (inlineToken.type === 'image') {
        const attrs = parseImageAttributes(inlineToken.text)

        if (attrs.display === 'inline') {
          const imageWidth = attrs.width + PDF_THEME.math.inlineModePadding

          if (currentX + imageWidth > margin + this.#maxWidth && currentX > margin && wrapLine()) {
            return { needsNewPage: true }
          }

          const imageResult = await embedImage(this.#pdfDoc, inlineToken.href)
          const yOffset = (attrs.height - PDF_THEME.font.baseSize) / 2

          this.#page.drawImage(imageResult.image, {
            x: currentX,
            y: this.#yPosition - yOffset,
            width: attrs.width,
            height: attrs.height
          })

          currentX += imageWidth
        } else {
          if (currentX > margin && wrapLine()) {
            return { needsNewPage: true }
          }

          const imageRenderResult = await this.#renderImage(
            inlineToken as Tokens.Image,
            margin,
            this.#maxWidth
          )

          if (imageRenderResult.needsNewPage) {
            return { needsNewPage: true }
          }

          currentX = margin
        }

        continue
      }

      const segments = parseInlineTokens([inlineToken])

      for (const segment of segments) {
        const font = getFontForSegment(segment, this.#fonts, preferredFont)
        let fontSize: number = PDF_THEME.font.baseSize
        let yOffset = 0

        if (segment.subscript || segment.superscript) {
          fontSize = PDF_THEME.font.subSupSize
          yOffset = segment.subscript ? PDF_THEME.offset.subscript : PDF_THEME.offset.superscript
        }

        const textLines = segment.text.split('\n')

        for (let lineIndex = 0; lineIndex < textLines.length; lineIndex++) {
          const lineText = textLines[lineIndex]
          const words = lineText.split(' ')

          for (const word of words) {
            if (!word) {
              continue
            }

            const wordWithSpace = word + ' '
            const wordWidth = font.widthOfTextAtSize(wordWithSpace, fontSize)

            if (currentX + wordWidth > margin + this.#maxWidth && currentX > margin && wrapLine()) {
              return { needsNewPage: true }
            }

            this.#page.drawText(wordWithSpace, {
              x: currentX,
              y: this.#yPosition + yOffset,
              font,
              size: fontSize,
              color: segment.color || PDF_THEME.color.text
            })

            if (segment.underline) {
              const underlineY = this.#yPosition + yOffset + PDF_THEME.underline.offset

              this.#page.drawLine({
                start: { x: currentX, y: underlineY },
                end: { x: currentX + wordWidth, y: underlineY },
                thickness: PDF_THEME.underline.thickness,
                color: segment.color || PDF_THEME.color.text
              })
            }

            if (segment.strikeThrough) {
              const strikeThroughY =
                this.#yPosition + yOffset + (fontSize / 2 - PDF_THEME.strikeThrough.thickness)

              this.#page.drawLine({
                start: { x: currentX, y: strikeThroughY },
                end: { x: currentX + wordWidth, y: strikeThroughY },
                thickness: PDF_THEME.strikeThrough.thickness
              })
            }

            currentX += wordWidth
          }

          if (lineIndex < textLines.length - 1 && wrapLine()) {
            return { needsNewPage: true }
          }
        }
      }
    }

    if (includeBottomSpacing) {
      this.#yPosition -= PDF_THEME.spacing.md
    }

    return { needsNewPage: false }
  }

  /**
   * Renders a code block token to the PDF page with background and monospace font.
   *
   * Code blocks are rendered with a colored background, monospace font, and proper
   * padding. Each line is rendered separately to maintain formatting and indentation.
   *
   * @param token - The markdown code token containing the code text
   * @returns RenderResult indicating if a new page is needed
   */
  #renderCodeBlock(token: Tokens.Code): RenderResult {
    const fontSize = PDF_THEME.font.codeSize
    const lineHeight = PDF_THEME.font.codeLineHeight
    const padding = PDF_THEME.codeBlock.padding
    const margin = PDF_THEME.layout.margin

    const lines = token.text.split('\n')
    const blockHeight = lines.length * lineHeight + padding * 2

    if (this.#yPosition - blockHeight < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    this.#page.drawRectangle({
      x: margin,
      y: this.#yPosition - blockHeight,
      width: this.#maxWidth,
      height: blockHeight,
      color: PDF_THEME.color.codeBlockBg
    })

    this.#yPosition -= padding

    for (const line of lines) {
      this.#page.drawText(line, {
        x: margin + padding,
        y: this.#yPosition - fontSize,
        size: fontSize,
        font: this.#fonts['mono'],
        color: PDF_THEME.color.codeBlockText
      })
      this.#yPosition -= lineHeight
    }

    this.#yPosition -= PDF_THEME.spacing.md + padding
    return { needsNewPage: false }
  }

  /**
   * Renders a list token to the PDF page, supporting both ordered and unordered lists with nested items.
   *
   * This method handles both bullet points and numbered lists, with proper indentation
   * for nested lists. It processes both text content and images within list items.
   *
   * @param token - The markdown list token containing list items
   * @param level - Nesting level for indentation (default: 0)
   * @returns Promise resolving to RenderResult indicating if a new page is needed
   */
  async #renderList(token: Tokens.List, level = 0): Promise<RenderResult> {
    const margin = PDF_THEME.layout.margin
    const indent = margin + level * PDF_THEME.spacing.listIndent
    const bulletChar = token.ordered ? token.start + '.' : '•'
    const lineHeight = PDF_THEME.font.listItemLineHeight

    for (let i = 0; i < token.items.length; i++) {
      const item = token.items[i]

      if (this.#yPosition - PDF_THEME.spacing.listGap < PDF_THEME.layout.margin) {
        return { needsNewPage: true }
      }

      this.#yPosition -= PDF_THEME.spacing.listGap

      const bullet = token.ordered ? `${i + 1}.` : bulletChar
      this.#page.drawText(bullet, {
        x: indent,
        y: this.#yPosition,
        size: PDF_THEME.font.listBulletSize,
        font: this.#fonts['regular'],
        color: PDF_THEME.color.text
      })

      const [content, ...subitems] = item.tokens

      if (!content) {
        continue
      }

      if ('tokens' in content) {
        const { textTokens, imageTokens } = partitionTokens(content.tokens)

        for (const imageToken of imageTokens) {
          const imageResult = await this.#renderImage(
            imageToken as Tokens.Image,
            indent + 20,
            this.#maxWidth - (indent + 20 - margin)
          )

          if (imageResult.needsNewPage) {
            return { needsNewPage: true }
          }
        }

        if (textTokens.length > 0) {
          const itemText = extractTextFromTokens(textTokens)
          const lines = splitTextToFit(
            itemText,
            this.#fonts['regular'],
            PDF_THEME.font.baseSize,
            this.#maxWidth - (indent + 20 - margin)
          )

          for (const line of lines) {
            if (this.#yPosition < PDF_THEME.layout.margin) {
              return { needsNewPage: true }
            }

            this.#page.drawText(line, {
              x: indent + 20,
              y: this.#yPosition,
              size: PDF_THEME.font.baseSize,
              font: this.#fonts['regular'],
              color: PDF_THEME.color.text
            })
            this.#yPosition -= lineHeight
          }
        }
      }

      if (subitems.length < 1) {
        continue
      }

      for (const token of subitems) {
        if (token.type !== 'list') {
          continue
        }

        const result = await this.#renderList(token as Tokens.List, level + 1)

        if (result.needsNewPage) {
          return { needsNewPage: true }
        }
      }
    }

    this.#yPosition -= PDF_THEME.spacing.md
    return { needsNewPage: false }
  }

  /**
   * Renders a blockquote token to the PDF page with a colored bar and italic text.
   *
   * Blockquotes are rendered with a colored left border bar and italic text.
   * The content is indented and can contain both text and images.
   *
   * @param token - The markdown blockquote token containing quoted content
   * @returns Promise resolving to RenderResult indicating if a new page is needed
   */
  async #renderBlockquote(token: Tokens.Blockquote, level = 1): Promise<RenderResult> {
    const margin = PDF_THEME.layout.margin * level
    const quoteMargin = margin + PDF_THEME.blockquote.barXOffset
    const lineHeight = PDF_THEME.font.blockquoteLineHeight

    if (this.#yPosition - lineHeight < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    if (level > 1) {
      this.#yPosition -= PDF_THEME.blockquote.nestedQuoteYOffset
    }

    const originalYPosition = this.#yPosition
    this.#yPosition -= lineHeight

    for (const item of token.tokens) {
      switch (item.type) {
        case 'paragraph': {
          const { needsNewPage } = await this.#renderParagraph({
            token: item as Tokens.Paragraph,
            margin: quoteMargin,
            preferredFont: 'italic',
            includeBottomSpacing: false
          })
          if (needsNewPage) {
            return { needsNewPage: true }
          }
          break
        }
        case 'blockquote': {
          const { needsNewPage } = await this.#renderBlockquote(
            item as Tokens.Blockquote,
            level + 1
          )
          if (needsNewPage) {
            return { needsNewPage: true }
          }
          break
        }
        default: {
          const error = new Error(`Unsupported blockquote item type: ${item.type}`)
          console.error(error)
          captureException(error)
          break
        }
      }
    }

    this.#yPosition -= lineHeight
    this.#page.drawRectangle({
      x: margin,
      y: originalYPosition,
      width: PDF_THEME.blockquote.barWidth,
      height: this.#yPosition - originalYPosition,
      color: PDF_THEME.color.blockquoteBar
    })

    this.#yPosition -= PDF_THEME.spacing.md
    return { needsNewPage: false }
  }

  /**
   * Renders an image token to the PDF page, including the image and optional title.
   *
   * Images are scaled to fit within the page width while maintaining aspect ratio.
   * Optional titles are rendered below the image in italic font.
   *
   * @param imageToken - The markdown image token containing href and optional title
   * @param margin - Left margin for positioning
   * @param maxWidth - Maximum width available for rendering
   * @returns Promise resolving to RenderResult indicating if a new page is needed
   */
  async #renderImage(
    imageToken: Tokens.Image,
    margin: number,
    maxWidth: number
  ): Promise<RenderResult> {
    const attrs = parseImageAttributes(imageToken.text)
    const imageResult = await embedImage(this.#pdfDoc, imageToken.href)

    let finalWidth = attrs.width > 0 ? attrs.width : imageResult.width
    let finalHeight = attrs.height > 0 ? attrs.height : imageResult.height

    const pageContentWidth = maxWidth - PDF_THEME.image.contentPadding
    if (finalWidth > pageContentWidth) {
      const scale = pageContentWidth / finalWidth
      finalWidth = pageContentWidth
      finalHeight = finalHeight * scale
    }

    const availableHeight = this.#pageHeight - 2 * PDF_THEME.layout.margin
    if (finalHeight > availableHeight) {
      const scale = availableHeight / finalHeight
      finalWidth = finalWidth * scale
      finalHeight = availableHeight
    }

    if (this.#yPosition - finalHeight < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    this.#yPosition -= finalHeight + PDF_THEME.spacing.md

    this.#page.drawImage(imageResult.image, {
      x: margin + (maxWidth - finalWidth) / 2,
      y: this.#yPosition,
      width: finalWidth,
      height: finalHeight
    })

    if (attrs.text) {
      this.#yPosition -= PDF_THEME.spacing.sm
      this.#page.drawText(`${attrs.text}`, {
        x: margin,
        y: this.#yPosition,
        size: PDF_THEME.font.imageTitleSize,
        font: this.#fonts['italic'],
        color: PDF_THEME.color.imagePlaceholder
      })
    }

    this.#yPosition -= PDF_THEME.spacing.md
    return { needsNewPage: false }
  }

  /**
   * Renders a table token to the PDF page with headers, borders, and cell content.
   *
   * Tables are rendered with borders, proper column spacing, and distinct header styling.
   * Text within cells is automatically wrapped to fit within column widths.
   *
   * @param token - The markdown table token containing header and row data
   * @returns RenderResult indicating if a new page is needed
   */
  #renderTable(token: Tokens.Table): RenderResult {
    const headerFont = this.#fonts['bold']
    const headerFontSize = PDF_THEME.font.tableHeaderTextSize
    const headerLineHeight = PDF_THEME.font.tableHeaderLineHeight
    const colWidth = this.#maxWidth / token.header.length

    const headerResult = this.#renderTableRow(
      token.header,
      headerFont,
      headerFontSize,
      headerLineHeight,
      PDF_THEME.color.tableHeaderBg,
      colWidth
    )

    if (headerResult.needsNewPage) {
      return { needsNewPage: true }
    }

    if (this.#yPosition < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    const rowFont = this.#fonts['regular']
    const rowFontSize = PDF_THEME.font.tableCellTextSize
    const rowLineHeight = PDF_THEME.font.tableCellLineHeight

    for (let i = 0; i < token.rows.length; i++) {
      const row = token.rows[i]
      const rowResult = this.#renderTableRow(
        row,
        rowFont,
        rowFontSize,
        rowLineHeight,
        PDF_THEME.color.tableCellBg,
        colWidth
      )

      if (rowResult.needsNewPage) {
        token.rows = token.rows.slice(i)
        return { needsNewPage: true }
      }
    }

    this.#yPosition -= PDF_THEME.spacing.md
    return { needsNewPage: false }
  }

  /**
   * Renders a single table cell with borders and text content.
   *
   * @param col - Column index for positioning
   * @param text - Text content to render in the cell
   * @param font - Font to use for the text
   * @param fontSize - Font size for the text
   * @param lineHeight - Line height for text layout
   * @param colWidth - Width of the column
   * @param rowHeight - Height of the row
   */
  #renderTableCell(
    col: number,
    text: string,
    font: PDFFont,
    fontSize: number,
    lineHeight: number,
    colWidth: number,
    rowHeight: number
  ) {
    const margin = PDF_THEME.layout.margin
    const cellPadding = PDF_THEME.table.cellPadding
    const cellX = margin + col * colWidth

    this.#page.drawRectangle({
      x: cellX,
      y: this.#yPosition - rowHeight,
      width: colWidth,
      height: rowHeight,
      borderColor: PDF_THEME.color.tableBorder,
      borderWidth: 1
    })

    this.#page.drawText(text, {
      x: cellX + cellPadding,
      y: this.#yPosition - cellPadding - fontSize,
      size: fontSize,
      font,
      color: PDF_THEME.color.text,
      maxWidth: colWidth - cellPadding * 2,
      lineHeight
    })
  }

  /**
   * Renders a complete table row with background color and cell borders.
   *
   * @param cells - Array of table cells to render
   * @param font - Font to use for the row text
   * @param fontSize - Font size for the row text
   * @param lineHeight - Line height for text layout
   * @param backgroundColor - Background color for the row
   * @param colWidth - Width of each column
   * @returns RenderResult indicating if a new page is needed
   */
  #renderTableRow(
    cells: Tokens.TableCell[],
    font: PDFFont,
    fontSize: number,
    lineHeight: number,
    backgroundColor: RGB,
    colWidth: number
  ): RenderResult {
    const cellPadding = PDF_THEME.table.cellPadding
    const margin = PDF_THEME.layout.margin

    const cellTexts = cells.map((cell) => extractTextFromTokens(cell.tokens).replace(/\n/g, ' '))
    const cellLines = cellTexts.map((text) =>
      splitTextToFit(text, font, fontSize, colWidth - cellPadding * 2)
    )

    const rowHeight =
      cellPadding * 2 +
      cellLines.reduce((max, lines) => Math.max(max, lines.length), 0) *
        PDF_THEME.font.tableHeaderLineHeight

    if (this.#yPosition - rowHeight < PDF_THEME.layout.margin) {
      return { needsNewPage: true }
    }

    this.#page.drawRectangle({
      x: margin,
      y: this.#yPosition - rowHeight,
      width: colWidth * cells.length,
      height: rowHeight,
      color: backgroundColor
    })

    for (let col = 0; col < cells.length; col++) {
      this.#renderTableCell(col, cellTexts.at(col), font, fontSize, lineHeight, colWidth, rowHeight)
    }

    this.#yPosition -= rowHeight
    return { needsNewPage: false }
  }

  /**
   * Renders a horizontal rule to the PDF page as a simple line.
   *
   * @returns RenderResult indicating if a new page is needed
   */
  #renderHorizontalRule(): RenderResult {
    const margin = PDF_THEME.layout.margin

    this.#page.drawLine({
      start: { x: margin, y: this.#yPosition },
      end: { x: margin + this.#maxWidth, y: this.#yPosition },
      thickness: PDF_THEME.hr.thickness,
      color: PDF_THEME.color.hr
    })

    this.#yPosition -= PDF_THEME.spacing.md
    return { needsNewPage: false }
  }

  /**
   * Loads and embeds all required fonts into the PDF document.
   *
   * This method fetches Ubuntu fonts from external URLs and embeds them
   * into the PDF document for consistent typography across different systems.
   *
   * @returns Promise that resolves when all fonts are loaded and embedded
   */
  async #loadFonts() {
    this.#pdfDoc.registerFontkit(fontkit)

    const fontWeights = Object.keys(FONT_URLS) as FontWeight[]

    const promises = fontWeights.map(async (fontWeight) => {
      const fontUrl = FONT_URLS[fontWeight]
      const fontBytes = await fetch(fontUrl).then((res) => res.arrayBuffer())
      const font = await this.#pdfDoc.embedFont(fontBytes)

      return { fontWeight, font }
    })

    const results = await Promise.all(promises)

    for (const { fontWeight, font } of results) {
      this.#fonts[fontWeight] = font
    }
  }

  /**
   * Creates a new page and initializes layout properties.
   *
   * This method adds a new page to the PDF document and sets up the
   * initial positioning and layout constraints for content rendering.
   */
  #addNewPage() {
    this.#page = this.#pdfDoc.addPage()
    this.#pageHeight = this.#page.getHeight()
    this.#maxWidth = this.#page.getWidth() - PDF_THEME.layout.margin * 2
    this.#yPosition = this.#pageHeight - PDF_THEME.layout.margin
  }
}
