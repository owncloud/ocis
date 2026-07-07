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
