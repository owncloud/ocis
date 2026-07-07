import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useImages } from '../../../../../src/composables/webWorkers/exportAsPdfWorker/useImages'

const mockCanvas = {
  width: 0,
  height: 0,
  getContext: vi.fn(),
  toDataURL: vi.fn(),
  parentNode: null,
  parentElement: null,
  insertBefore: vi.fn(),
  appendChild: vi.fn(),
  removeChild: vi.fn(),
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
  setAttribute: vi.fn(),
  getAttribute: vi.fn(),
  style: {}
}

const mockContext = {
  drawImage: vi.fn()
}

const mockImage = {
  crossOrigin: '',
  onload: null as any,
  onerror: null as any,
  src: '',
  width: 100,
  height: 100
}

describe('useImages', () => {
  beforeEach(() => {
    console.error = vi.fn()

    Object.defineProperty(global, 'document', {
      value: {
        createElement: vi.fn((tagName: string) => {
          if (tagName === 'canvas') {
            return mockCanvas
          }

          return {
            parentNode: null,
            parentElement: null,
            insertBefore: vi.fn(),
            appendChild: vi.fn(),
            removeChild: vi.fn(),
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            setAttribute: vi.fn(),
            getAttribute: vi.fn(),
            style: {}
          }
        })
      }
    })

    Object.defineProperty(global, 'Image', {
      value: vi.fn(function () {
        const image = { ...mockImage }

        setTimeout(() => {
          if (image.onload) {
            image.onload()
          }
        }, 0)

        return image
      })
    })

    mockCanvas.getContext.mockReturnValue(mockContext)
    mockCanvas.toDataURL.mockReturnValue('data:image/png;base64,mockDataURL')
  })

  describe('method "preprocessImages"', () => {
    it('should return unchanged content when no images are found', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = 'This is just text with no images.'

            try {
              const result = await preprocessImages(content)
              expect(result).toBe(content)
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should ignore data URLs and leave them unchanged', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = '![test](data:image/png;base64,existingDataURL)'
            try {
              const result = await preprocessImages(content)
              expect(result).toBe(content)
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should convert external image URLs to data URLs successfully', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = '![alt text](https://example.com/image.png)'

            try {
              const result = await preprocessImages(content)

              expect(result).toBe('![alt text](data:image/png;base64,mockDataURL)')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle multiple images in the same content', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content =
              '![first](https://example.com/image1.png) and ![second](https://example.com/image2.png)'

            try {
              const result = await preprocessImages(content)
              expect(result).toBe(
                '![first](data:image/png;base64,mockDataURL) and ![second](data:image/png;base64,mockDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should replace failed image conversions with error message', async () => {
      Object.defineProperty(global, 'Image', {
        value: vi.fn(function () {
          const image = { ...mockImage }

          setTimeout(() => {
            if (image.onerror) {
              image.onerror(new Error('Failed to load image'))
            }
          }, 0)

          return image
        })
      })

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = '![failed image](https://example.com/broken-image.png)'

            try {
              const result = await preprocessImages(content)
              expect(result).toBe('*Failed to render image.*')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle canvas context creation failure', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = '![test](https://example.com/canvas-failure.png)'
            mockCanvas.getContext.mockReturnValue(null)

            try {
              const result = await preprocessImages(content)
              expect(result).toBe('*Failed to render image.*')
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle mixed scenarios with some successful and some failed conversions', async () => {
      Object.defineProperty(global, 'Image', {
        value: vi.fn(function () {
          const image = { ...mockImage }

          setTimeout(() => {
            if (image.onload && image.src.includes('good.png')) {
              image.onload()
            }

            if (image.onerror && image.src.includes('bad.png')) {
              image.onerror(new Error('Failed to load image'))
            }
          }, 0)

          return image
        })
      })

      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content =
              '![success](https://example.com/good.png) ![fail](https://example.com/bad.png)'

            try {
              const result = await preprocessImages(content)
              expect(result).toBe(
                '![success](data:image/png;base64,mockDataURL) *Failed to render image.*'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should preserve alt text in successful conversions', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = '![This is alt text with spaces](https://example.com/image.png)'

            try {
              const result = await preprocessImages(content)
              expect(result).toBe(
                '![This is alt text with spaces](data:image/png;base64,mockDataURL)'
              )
              resolve()
            } catch (error) {
              reject(error)
            }
          }
        })
      })
    })

    it('should handle empty alt text', async () => {
      await new Promise<void>((resolve, reject) => {
        getWrapper({
          setup: async ({ preprocessImages }) => {
            const content = '![](https://example.com/image.png)'

            try {
              const result = await preprocessImages(content)
              expect(result).toBe('![](data:image/png;base64,mockDataURL)')
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

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useImages>) => void }) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useImages()
      setup(instance)
    })
  }
}
