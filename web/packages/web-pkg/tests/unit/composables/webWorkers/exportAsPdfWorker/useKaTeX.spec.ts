import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import katex from 'katex'
import { useKaTeX } from '../../../../../src/composables/webWorkers/exportAsPdfWorker/useKaTeX'
import html2canvas from 'html2canvas'
import { mock } from 'vitest-mock-extended'

vi.mock('katex')
vi.mock('html2canvas')

const mockCanvas = {
  width: 200,
  height: 100,
  toDataURL: vi.fn()
}

const mockKatexElement = {
  getBoundingClientRect: vi.fn(() => ({
    width: 200,
    height: 100
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
  querySelector: vi.fn()
}

const mockDocument = {
  createElement: vi.fn(),
  body: {
    appendChild: vi.fn(),
    removeChild: vi.fn()
  }
}

describe('useKaTeX', () => {
  beforeEach(() => {
    console.error = vi.fn()

    vi.mocked(katex.render).mockImplementation(vi.fn())
    vi.mocked(html2canvas).mockResolvedValue(mock<HTMLCanvasElement>(mockCanvas))

    mockCanvas.toDataURL.mockReturnValue('data:image/png;base64,mockKaTeXDataURL')
    mockDocument.createElement.mockReturnValue(mockContainer)
    mockContainer.querySelector.mockReturnValue(mockKatexElement)

    Object.defineProperty(global, 'document', {
      value: mockDocument,
      configurable: true
    })
  })

  describe('method "preprocessKaTeXFormulas"', () => {
    it('should return unchanged content when no formulas are found', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'This is just text with no math formulas.'

            try {
              const result = await preprocessKaTeXFormulas(content)
              expect(result).toBe(content)
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should convert inline formulas to image data URLs successfully', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Here is an inline formula: $x = y + z$'

            try {
              const result = await preprocessKaTeXFormulas(content)

              expect(katex.render).toHaveBeenCalledWith('x = y + z', mockContainer, {
                displayMode: false,
                throwOnError: false,
                errorColor: expect.any(String)
              })

              expect(result).toBe(
                'Here is an inline formula: ![d=inline;w=200;h=100](data:image/png;base64,mockKaTeXDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should convert block formulas to image data URLs successfully', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Here is a block formula:\n$$\\sum_{i=1}^{n} x_i$$'

            try {
              const result = await preprocessKaTeXFormulas(content)

              expect(katex.render).toHaveBeenCalledWith('\\sum_{i=1}^{n} x_i', mockContainer, {
                displayMode: true,
                throwOnError: false,
                errorColor: expect.any(String)
              })

              expect(result).toBe(
                'Here is a block formula:\n![w=200;h=100](data:image/png;base64,mockKaTeXDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle multiple formulas in the same content', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Inline: $a = b$ and block: $$c = d$$'

            try {
              const result = await preprocessKaTeXFormulas(content)

              expect(katex.render).toHaveBeenCalledTimes(2)
              expect(result).toBe(
                'Inline: ![d=inline;w=200;h=100](data:image/png;base64,mockKaTeXDataURL) and block: ![w=200;h=100](data:image/png;base64,mockKaTeXDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should replace failed formula renderings with error message', async () => {
      vi.mocked(katex.render).mockImplementation(() => {
        throw new Error('KaTeX rendering failed')
      })

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Failed formula: $\\invalid$'

            try {
              const result = await preprocessKaTeXFormulas(content)
              expect(result).toBe('Failed formula: *Failed to render math formula.*')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle KaTeX not producing a valid element', async () => {
      mockContainer.querySelector.mockReturnValue(null)

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Formula: $x = y$'

            try {
              const result = await preprocessKaTeXFormulas(content)
              expect(result).toBe('Formula: *Failed to render math formula.*')
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
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Formula: $x = y$'

            try {
              const result = await preprocessKaTeXFormulas(content)
              expect(result).toBe('Formula: *Failed to render math formula.*')
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
      vi.mocked(katex.render).mockImplementation(() => {
        callCount++
        if (callCount === 1) {
          mockContainer.querySelector.mockReturnValue(mockKatexElement)
        } else {
          throw new Error('KaTeX rendering failed')
        }
      })

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Success: $a = b$ and fail: $\\invalid$'

            try {
              const result = await preprocessKaTeXFormulas(content)
              expect(result).toBe(
                'Success: ![d=inline;w=200;h=100](data:image/png;base64,mockKaTeXDataURL) and fail: *Failed to render math formula.*'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle formulas with whitespace and special characters', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Formula: $$  \\frac{a}{b} + \\sqrt{c}  $$'

            try {
              const result = await preprocessKaTeXFormulas(content)

              expect(katex.render).toHaveBeenCalledWith('\\frac{a}{b} + \\sqrt{c}', mockContainer, {
                displayMode: true,
                throwOnError: false,
                errorColor: expect.any(String)
              })

              expect(result).toBe('Formula: ![w=200;h=100](data:image/png;base64,mockKaTeXDataURL)')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle empty formulas', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessKaTeXFormulas }) => {
            const content = 'Empty formulas: $$ $$ and $ $'

            try {
              const result = await preprocessKaTeXFormulas(content)

              expect(katex.render).toHaveBeenCalledWith('', mockContainer, expect.any(Object))
              expect(result).toBe(
                'Empty formulas: ![w=200;h=100](data:image/png;base64,mockKaTeXDataURL) and ![d=inline;w=200;h=100](data:image/png;base64,mockKaTeXDataURL)'
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

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useKaTeX>) => void }) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useKaTeX()
      setup(instance)
    })
  }
}
