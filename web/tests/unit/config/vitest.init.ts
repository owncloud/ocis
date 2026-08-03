const IntersectionObserverMock = vi.fn(function () {
  return {
    disconnect: vi.fn(),
    observe: vi.fn(),
    takeRecords: vi.fn(),
    unobserve: vi.fn()
  }
})

vi.stubGlobal('IntersectionObserver', IntersectionObserverMock)

const ResizeObserverMock = vi.fn(function () {
  return { observe: vi.fn(), unobserve: vi.fn() }
})

vi.stubGlobal('ResizeObserver', ResizeObserverMock)

// jsdom (used via `@vitest-environment jsdom` overrides) doesn't implement matchMedia,
// unlike the default happy-dom environment
if (typeof window !== 'undefined' && !window.matchMedia) {
  vi.stubGlobal(
    'matchMedia',
    vi.fn().mockImplementation((query) => ({
      matches: false,
      media: query,
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn()
    }))
  )
}

vi.stubGlobal('define', vi.fn())

// This is needed for KaTeX to work in the tests
Object.defineProperty(document, 'compatMode', {
  value: 'CSS1Compat'
})

// Mock Math.random to return predictable values for tests
let mathRandomCounter = 0
const originalMathRandom = Math.random
Math.random = () => {
  mathRandomCounter++
  return mathRandomCounter / 10000
}

// Reset counter before each test for consistent IDs
beforeEach(() => {
  mathRandomCounter = 0
})
