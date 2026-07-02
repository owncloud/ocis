import mermaid from 'mermaid'
import html2canvas from 'html2canvas'
import { useGettext } from 'vue3-gettext'

mermaid.initialize({
  startOnLoad: false,
  theme: 'default'
})

/**
 * Renders a Mermaid diagram to a data URL image.
 *
 * This function takes a Mermaid diagram syntax, renders it to SVG using the Mermaid library,
 * then converts the SVG to a high-quality PNG image using html2canvas. The rendered image
 * is scaled up for better quality and includes proper dimensions for PDF integration.
 *
 * @param diagram - The Mermaid diagram syntax to render
 * @returns Promise resolving to an object containing the data URL and dimensions
 * @throws Error if Mermaid fails to produce a valid SVG element
 */
async function renderMermaidToDataURL(
  diagram: string
): Promise<{ dataURL: string; width: number; height: number }> {
  const { svg } = await mermaid.render('mermaid-temp-div', diagram)

  const container = document.createElement('div')
  container.style.position = 'absolute'
  container.style.left = '-9999px'
  container.style.background = '#fff'
  container.innerHTML = svg
  document.body.appendChild(container)

  const svgElement = container.querySelector('svg')
  if (!svgElement) {
    document.body.removeChild(container)
    throw new Error('Mermaid did not produce a valid SVG element.')
  }

  const rect = svgElement.getBoundingClientRect()
  const width = Math.ceil(rect.width)
  const height = Math.ceil(rect.height)

  try {
    const scaleFactor = 3
    const canvas = await html2canvas(container, {
      backgroundColor: '#fff',
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
 * Composable providing Mermaid chart preprocessing for PDF generation.
 *
 * This composable handles the conversion of Mermaid diagram blocks in markdown content
 * to images that can be embedded in PDF documents. It processes all Mermaid code blocks
 * and converts them to high-quality PNG images with proper dimensions.
 */
export function useMermaid() {
  const { $pgettext } = useGettext()

  /**
   * Preprocesses markdown content to convert Mermaid chart blocks into image data URLs.
   *
   * This function scans markdown content for Mermaid code blocks (```mermaid...```),
   * renders each diagram using the Mermaid library, converts them to PNG images,
   * and replaces the original code blocks with markdown image syntax.
   * Failed renderings are replaced with error messages.
   *
   * @param markdownContent - The markdown content to preprocess
   * @returns Promise resolving to the content with charts replaced by image tags
   */
  async function preprocessMermaidCharts(markdownContent: string): Promise<string> {
    const mermaidRegex = /```mermaid\n([\s\S]*?)\n```/g
    const matches = Array.from(markdownContent.matchAll(mermaidRegex))
    if (matches.length === 0) {
      return markdownContent
    }

    const renderingPromises = matches.map((match) => {
      const chartSyntax = match[1]
      return renderMermaidToDataURL(chartSyntax).catch((error) => {
        console.error('Failed to render Mermaid chart:', error)
        return null
      })
    })

    const results = await Promise.all(renderingPromises)

    let i = 0
    return markdownContent.replace(mermaidRegex, () => {
      const result = results[i++]
      if (result) {
        const imageMarkdown = `![w=${result.width};h=${result.height}](${result.dataURL})`
        return imageMarkdown
      } else {
        return (
          '*' +
          $pgettext(
            'Error message rendered in a PDF file when there is any error during the rendering of a Mermaid chart.',
            'Failed to render Mermaid chart.'
          ) +
          '*'
        )
      }
    })
  }

  return {
    preprocessMermaidCharts
  }
}
