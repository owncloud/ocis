import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import mermaid from 'mermaid'
import html2canvas from 'html2canvas'
import { useMermaid } from '../../../../../src/composables/webWorkers/exportAsPdfWorker/useMermaid'
import { mock } from 'vitest-mock-extended'

vi.mock('mermaid')
vi.mock('html2canvas')

const mockCanvas = {
  width: 300,
  height: 200,
  toDataURL: vi.fn()
}

const mockSvgElement = {
  getBoundingClientRect: vi.fn(() => ({
    width: 300,
    height: 200
  }))
}

const mockContainer = {
  parentNode: null,
  parentElement: null,
  insertBefore: vi.fn(),
  appendChild: vi.fn(),
  removeChild: vi.fn(),
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
  setAttribute: vi.fn(),
  getAttribute: vi.fn(),
  style: {},
  querySelector: vi.fn(),
  innerHTML: ''
}

const mockDocument = {
  createElement: vi.fn(),
  body: {
    appendChild: vi.fn(),
    removeChild: vi.fn()
  }
}

describe('useMermaid', () => {
  beforeEach(() => {
    console.error = vi.fn()

    vi.mocked(mermaid.render).mockResolvedValue({
      svg: '<svg><g></g></svg>',
      diagramType: 'graph'
    })
    vi.mocked(html2canvas).mockResolvedValue(mock<HTMLCanvasElement>(mockCanvas))

    mockCanvas.toDataURL.mockReturnValue('data:image/png;base64,mockMermaidDataURL')
    mockDocument.createElement.mockReturnValue(mockContainer)
    mockContainer.querySelector.mockReturnValue(mockSvgElement)

    Object.defineProperty(global, 'document', {
      value: mockDocument,
      configurable: true
    })
  })

  describe('method "preprocessMermaidCharts"', () => {
    it('should return unchanged content when no Mermaid charts are found', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = 'This is just text with no Mermaid charts.'

            try {
              const result = await preprocessMermaidCharts(content)
              expect(result).toBe(content)
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should convert Mermaid chart blocks to image data URLs successfully', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Here is a Mermaid chart:
\`\`\`mermaid
graph TD
A --> B
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)

              expect(mermaid.render).toHaveBeenCalledWith('mermaid-temp-div', 'graph TD\nA --> B')
              expect(result).toBe(
                'Here is a Mermaid chart:\n![w=300;h=200](data:image/png;base64,mockMermaidDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle multiple Mermaid charts in the same content', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `First chart:
\`\`\`mermaid
graph LR
A --> B
\`\`\`

Second chart:
\`\`\`mermaid
flowchart TD
C --> D
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)

              expect(mermaid.render).toHaveBeenCalledTimes(2)
              expect(mermaid.render).toHaveBeenNthCalledWith(
                1,
                'mermaid-temp-div',
                'graph LR\nA --> B'
              )
              expect(mermaid.render).toHaveBeenNthCalledWith(
                2,
                'mermaid-temp-div',
                'flowchart TD\nC --> D'
              )

              expect(result).toBe(
                `First chart:
![w=300;h=200](data:image/png;base64,mockMermaidDataURL)

Second chart:
![w=300;h=200](data:image/png;base64,mockMermaidDataURL)`
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should replace failed chart renderings with error message', async () => {
      vi.mocked(mermaid.render).mockRejectedValue(new Error('Mermaid rendering failed'))

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Failed chart:
\`\`\`mermaid
invalid syntax
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)
              expect(result).toBe('Failed chart:\n*Failed to render Mermaid chart.*')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle Mermaid not producing a valid SVG element', async () => {
      mockContainer.querySelector.mockReturnValue(null)

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Chart:
\`\`\`mermaid
graph TD
A --> B
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)
              expect(result).toBe('Chart:\n*Failed to render Mermaid chart.*')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle html2canvas failure', async () => {
      vi.mocked(html2canvas).mockRejectedValue(new Error('html2canvas failed'))

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Chart:
\`\`\`mermaid
graph TD
A --> B
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)
              expect(result).toBe('Chart:\n*Failed to render Mermaid chart.*')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle mixed scenarios with some successful and some failed conversions', async () => {
      let callCount = 0
      vi.mocked(mermaid.render).mockImplementation(async () => {
        callCount++
        if (callCount === 1) {
          mockContainer.querySelector.mockReturnValue(mockSvgElement)
          return { svg: '<svg><g></g></svg>', diagramType: 'graph' }
        } else {
          throw new Error('Mermaid rendering failed')
        }
      })

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Success:
\`\`\`mermaid
graph TD
A --> B
\`\`\`

Fail:
\`\`\`mermaid
invalid syntax
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)
              expect(result).toBe(
                `Success:
![w=300;h=200](data:image/png;base64,mockMermaidDataURL)

Fail:
*Failed to render Mermaid chart.*`
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle charts with complex syntax and whitespace', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Complex chart:
\`\`\`mermaid
sequenceDiagram
    participant A as Alice
    participant B as Bob
    A->>B: Hello Bob!
    B-->>A: Hello Alice!
\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)

              expect(mermaid.render).toHaveBeenCalledWith(
                'mermaid-temp-div',
                'sequenceDiagram\n    participant A as Alice\n    participant B as Bob\n    A->>B: Hello Bob!\n    B-->>A: Hello Alice!'
              )

              expect(result).toBe(
                'Complex chart:\n![w=300;h=200](data:image/png;base64,mockMermaidDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle empty Mermaid blocks', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessMermaidCharts }) => {
            const content = `Empty chart:
\`\`\`mermaid

\`\`\``

            try {
              const result = await preprocessMermaidCharts(content)

              expect(mermaid.render).toHaveBeenCalledWith('mermaid-temp-div', '')
              expect(result).toBe(
                'Empty chart:\n![w=300;h=200](data:image/png;base64,mockMermaidDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })
  })
})

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useMermaid>) => void }) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useMermaid()
      setup(instance)
    })
  }
}
