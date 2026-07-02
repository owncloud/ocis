import katex from 'katex'
import html2canvas from 'html2canvas'
import { useGettext } from 'vue3-gettext'
import { PDF_THEME } from './pdfConfig'
import { pdfColorToCssRgb } from './helpers'

/**
 * Renders a KaTeX mathematical formula to a data URL image.
 *
 * @param formula - The KaTeX mathematical formula to render
 * @param isBlockMode - Whether to render in display mode (block) or inline mode (default: false)
 * @returns Promise resolving to an object containing the data URL and dimensions
 */
async function renderKaTeXToDataURL(
  formula: string,
  isBlockMode: boolean = false
): Promise<{ dataURL: string; width: number; height: number }> {
  const container = document.createElement('div')
  container.style.position = 'absolute'
  container.style.left = '-9999px'
  container.style.background = '#fff'
  container.style.padding = isBlockMode ? '2px' : '0px'
  container.style.fontSize = `${PDF_THEME.font.baseSize}px`
  container.style.color = pdfColorToCssRgb(PDF_THEME.color.text)

  document.body.appendChild(container)

  try {
    katex.render(formula, container, {
      displayMode: isBlockMode,
      throwOnError: false,
      errorColor: pdfColorToCssRgb(PDF_THEME.color.error)
    })

    const katexElement = container.querySelector<HTMLElement>('.katex')
    if (!katexElement) {
      throw new Error('KaTeX did not produce a valid element.')
    }

    const rect = katexElement.getBoundingClientRect()
    const width = Math.ceil(rect.width)
    const height = Math.ceil(rect.height)

    const scaleFactor = PDF_THEME.svg.scaleFactor
    const canvas = await html2canvas(katexElement, {
      backgroundColor: null,
      scale: scaleFactor,
      logging: false
    })

    const dataURL = canvas.toDataURL('image/png')

    return { dataURL, width, height }
  } finally {
    document.body.removeChild(container)
  }
}

/**
 * Composable providing KaTeX formula preprocessing for PDF generation.
 *
 * This composable handles the conversion of KaTeX mathematical formulas in markdown
 * content to images that can be embedded in PDF documents. It supports both inline
 * formulas ($...$) and display block formulas ($$...$$).
 */
export function useKaTeX() {
  const { $pgettext } = useGettext()

  /**
   * Replaces mathematical formulas in content with rendered image representations.
   *
   * @param content - The markdown content containing formulas to replace
   * @param regex - The regular expression pattern to match formulas
   * @param isBlockMode - Whether the formulas should be rendered in block mode
   * @returns Promise resolving to the content with formulas replaced by image tags
   */
  const replaceFormulas = async (
    content: string,
    regex: RegExp,
    isBlockMode: boolean
  ): Promise<string> => {
    const matches = Array.from(content.matchAll(regex))
    if (matches.length === 0) {
      return content
    }

    const renderingPromises = matches.map((match) => {
      const formula = match[1].trim()
      return renderKaTeXToDataURL(formula, isBlockMode).catch((error) => {
        console.error('Failed to render KaTeX formula:', error)
        return null
      })
    })

    const results = await Promise.all(renderingPromises)

    let i = 0
    return content.replace(regex, () => {
      const result = results[i++]
      if (result) {
        const attributes = isBlockMode
          ? `w=${result.width};h=${result.height}`
          : `d=inline;w=${result.width};h=${result.height}`
        return `![${attributes}](${result.dataURL})`
      } else {
        return (
          '*' +
          $pgettext(
            'Error message rendered in a PDF file when there is any error during the rendering of a KaTeX formula.',
            'Failed to render math formula.'
          ) +
          '*'
        )
      }
    })
  }

  /**
   * Preprocesses markdown content to convert KaTeX formulas into image data URLs.
   *
   * This function processes markdown content in two passes:
   * 1. Converts block formulas ($$...$$) to display mode images
   * 2. Converts inline formulas ($...$) to inline mode images
   *
   * Each formula is rendered using KaTeX and converted to a PNG image with appropriate
   * dimensions and styling for PDF integration.
   *
   * @param markdownContent - The markdown content to preprocess
   * @returns Promise resolving to the content with formulas replaced by image tags
   */
  async function preprocessKaTeXFormulas(markdownContent: string): Promise<string> {
    // Extract code blocks, tables, and inline code, then replace with placeholders
    const extractedContent: string[] = []
    let content = markdownContent

    content = content.replace(/```[\s\S]*?```/g, (match) => {
      const placeholder = `__CODE_BLOCK_${extractedContent.length}__`
      extractedContent.push(match)
      return placeholder
    })

    content = content.replace(/^\|.*\|[\s\S]*?(?=\n(?!\|)|$)/gm, (match) => {
      const placeholder = `__TABLE_${extractedContent.length}__`
      extractedContent.push(match)
      return placeholder
    })

    content = content.replace(/`[^`]+`/g, (match) => {
      const placeholder = `__INLINE_CODE_${extractedContent.length}__`
      extractedContent.push(match)
      return placeholder
    })

    const blockRegex = /\$\$([\s\S]*?)\$\$/g
    const inlineRegex = /\$([\s\S]*?)\$/g

    let processedContent = await replaceFormulas(content, blockRegex, true)
    processedContent = await replaceFormulas(processedContent, inlineRegex, false)

    extractedContent.forEach((extracted, index) => {
      const codeBlockPlaceholder = `__CODE_BLOCK_${index}__`
      const tablePlaceholder = `__TABLE_${index}__`
      const inlinePlaceholder = `__INLINE_CODE_${index}__`

      processedContent = processedContent.replace(codeBlockPlaceholder, extracted)
      processedContent = processedContent.replace(tablePlaceholder, extracted)
      processedContent = processedContent.replace(inlinePlaceholder, extracted)
    })

    return processedContent
  }

  return { preprocessKaTeXFormulas }
}
