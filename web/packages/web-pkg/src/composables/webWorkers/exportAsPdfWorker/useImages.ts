import { useGettext } from 'vue3-gettext'

const inMemoryCache = new Map<string, string>()

/**
 * Converts an external image URL to a data URL by loading it into a canvas.
 *
 * This function creates an HTML image element, loads the external image with CORS enabled,
 * draws it onto a canvas, and then converts the canvas content to a PNG data URL.
 * This is necessary for embedding external images into PDF documents.
 *
 * @param imageUrl - The external URL of the image to convert
 * @returns Promise resolving to the image as a PNG data URL
 * @throws Error if the image fails to load or canvas context cannot be obtained
 */
function convertImageToDataURL(imageUrl: string): Promise<string> {
  if (inMemoryCache.has(imageUrl)) {
    return Promise.resolve(inMemoryCache.get(imageUrl)!)
  }

  return new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'

    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = img.width
      canvas.height = img.height

      const ctx = canvas.getContext('2d')
      if (!ctx) {
        return reject(new Error('Could not get canvas context.'))
      }

      ctx.drawImage(img, 0, 0, img.width, img.height)

      const dataURL = canvas.toDataURL('image/png')

      inMemoryCache.set(imageUrl, dataURL)
      resolve(dataURL)
    }

    img.onerror = (err) => {
      reject(err)
    }

    img.src = imageUrl
  })
}

/**
 * Helper function to identify code regions in markdown content.
 * Returns an array of ranges [start, end] that represent code blocks and inline code.
 */
function getCodeRegions(markdownContent: string): Array<[number, number]> {
  const codeRegions: Array<[number, number]> = []

  const fencedCodeBlockRegex = /^```[^\n]*\n[\s\S]*?^```$/gm
  let match: RegExpExecArray | null
  while ((match = fencedCodeBlockRegex.exec(markdownContent)) !== null) {
    codeRegions.push([match.index, match.index + match[0].length])
  }

  const lines = markdownContent.split('\n')
  let inCodeBlock = false
  let codeBlockStart = 0
  let currentPos = 0

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const isIndentedLine = /^(?:    |\t)/.test(line)
    const lineLength = line.length + 1

    if (isIndentedLine && !inCodeBlock) {
      inCodeBlock = true
      codeBlockStart = currentPos
    } else if (!isIndentedLine && inCodeBlock) {
      codeRegions.push([codeBlockStart, currentPos])
      inCodeBlock = false
    }

    currentPos += lineLength
  }

  if (inCodeBlock && codeBlockStart >= 0) {
    codeRegions.push([codeBlockStart, markdownContent.length])
  }

  const inlineCodeRegex = /`[^`\n]+`/g
  while ((match = inlineCodeRegex.exec(markdownContent)) !== null) {
    codeRegions.push([match.index, match.index + match[0].length])
  }

  codeRegions.sort((a, b) => a[0] - b[0])

  const mergedRegions: Array<[number, number]> = []
  for (const [start, end] of codeRegions) {
    if (mergedRegions.length === 0 || mergedRegions[mergedRegions.length - 1][1] < start) {
      mergedRegions.push([start, end])
    } else {
      mergedRegions[mergedRegions.length - 1][1] = Math.max(
        mergedRegions[mergedRegions.length - 1][1],
        end
      )
    }
  }

  return mergedRegions
}

/**
 * Helper function to check if a position is within any code region.
 */
function isInCodeRegion(position: number, codeRegions: Array<[number, number]>): boolean {
  return codeRegions.some(([start, end]) => position >= start && position < end)
}

/**
 * Composable providing image preprocessing for PDF generation.
 *
 * This composable handles the conversion of external image URLs in markdown content
 * to data URLs that can be embedded directly in PDF documents. It processes all
 * external images (non-data URLs) and converts them to base64-encoded PNG data.
 * Images within code blocks and inline code are excluded from processing.
 */
export function useImages() {
  const { $pgettext } = useGettext()

  /**
   * Preprocesses markdown content to convert external image URLs into data URLs.
   *
   * This function scans markdown content for both inline and reference-style image syntax
   * with external URLs (excluding data URLs), converts each external image to a data URL using
   * canvas rendering, and replaces the original URL with the data URL.
   * Images within code blocks and inline code are excluded from processing.
   * Failed conversions are replaced with error messages.
   *
   * @param markdownContent - The markdown content to preprocess
   * @returns Promise resolving to the content with image sources replaced by data URLs
   */
  async function preprocessImages(markdownContent: string): Promise<string> {
    const codeRegions = getCodeRegions(markdownContent)

    const referenceRegex = /^\[([^\]]+)\]:\s*(.+?)(?:\s+"[^"]*")?\s*$/gm
    const referenceMap = new Map<string, string>()

    let match: RegExpExecArray | null
    while ((match = referenceRegex.exec(markdownContent)) !== null) {
      if (!isInCodeRegion(match.index, codeRegions)) {
        const referenceId = match[1]
        const imageUrl = match[2].trim()
        referenceMap.set(referenceId, imageUrl)
      }
    }

    const inlineImageRegex = /!\[([^\]]*)\]\((?!data:)([^"\s)]+)(?:\s+"[^"]*")?\)/g
    const referenceImageRegex = /!\[([^\]]*)\]\[([^\]]*)\]/g

    const inlineMatches = Array.from(markdownContent.matchAll(inlineImageRegex))
    const referenceMatches = Array.from(markdownContent.matchAll(referenceImageRegex))

    const filteredInlineMatches = inlineMatches.filter(
      (match) => !isInCodeRegion(match.index, codeRegions)
    )
    const filteredReferenceMatches = referenceMatches.filter(
      (match) => !isInCodeRegion(match.index, codeRegions)
    )

    if (filteredInlineMatches.length === 0 && filteredReferenceMatches.length === 0) {
      return markdownContent
    }

    const imageUrls = new Set<string>()

    filteredInlineMatches.forEach((match) => {
      const imageUrl = match[2]
      if (!imageUrl.startsWith('data:')) {
        imageUrls.add(imageUrl)
      }
    })

    filteredReferenceMatches.forEach((match) => {
      const referenceId = match[2]
      const imageUrl = referenceMap.get(referenceId)
      if (imageUrl && !imageUrl.startsWith('data:')) {
        imageUrls.add(imageUrl)
      }
    })

    const conversionPromises = Array.from(imageUrls).map(async (imageUrl) => {
      try {
        const dataURL = await convertImageToDataURL(imageUrl)
        return { imageUrl, dataURL }
      } catch (error) {
        console.error('Failed to convert image to data URL:', error)
        return { imageUrl, dataURL: null }
      }
    })

    const conversionResults = await Promise.all(conversionPromises)
    const urlToDataUrlMap = new Map<string, string | null>()
    conversionResults.forEach(({ imageUrl, dataURL }) => {
      urlToDataUrlMap.set(imageUrl, dataURL)
    })

    let processedContent = markdownContent

    let inlineMatchIndex = 0
    processedContent = processedContent.replace(inlineImageRegex, (match, altText, imageUrl) => {
      const currentMatch = inlineMatches[inlineMatchIndex++]
      if (!currentMatch || isInCodeRegion(currentMatch.index, codeRegions)) {
        return match
      }

      const dataURL = urlToDataUrlMap.get(imageUrl)
      if (dataURL) {
        return `![${altText}](${dataURL})`
      }

      return (
        '*' +
        $pgettext(
          'Error message rendered in a PDF file when there is any error during the rendering of an image.',
          'Failed to render image.'
        ) +
        '*'
      )
    })

    let referenceMatchIndex = 0
    processedContent = processedContent.replace(
      referenceImageRegex,
      (match, altText, referenceId) => {
        const currentMatch = referenceMatches[referenceMatchIndex++]
        if (!currentMatch || isInCodeRegion(currentMatch.index, codeRegions)) {
          return match
        }

        const imageUrl = referenceMap.get(referenceId)
        if (!imageUrl) {
          return (
            '*' +
            $pgettext(
              'Error message rendered in a PDF file when there is any error during the rendering of an image.',
              'Failed to render image.'
            ) +
            '*'
          )
        }

        const dataURL = urlToDataUrlMap.get(imageUrl)
        if (dataURL) {
          return `![${altText}](${dataURL})`
        }

        return (
          '*' +
          $pgettext(
            'Error message rendered in a PDF file when there is any error during the rendering of an image.',
            'Failed to render image.'
          ) +
          '*'
        )
      }
    )

    return processedContent
  }

  return {
    preprocessImages
  }
}
