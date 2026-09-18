import { Plugin } from 'vue'
import DesignSystem from '@ownclouders/design-system'
import { registerTestingPlugins } from '@ownclouders/web-test-helpers-core'
import { createMockStore } from '@ownclouders/web-pkg/src/testing'

// `@ownclouders/design-system` and `@ownclouders/web-pkg` both depend on
// `@ownclouders/web-test-helpers-core` for their own tests, so `defaultPlugins` there can't
// import either directly without recreating that cycle - see core's registry.ts. Consumers
// going through the published `@ownclouders/web-test-helpers` facade get this registration
// automatically (as an import side effect of the facade itself); design-system's and
// web-pkg's own specs import straight from core/`../testing` instead, so they need it done
// here, in the one shared setup file that isn't part of the package dependency graph.
registerTestingPlugins({
  designSystem: DesignSystem as unknown as Plugin,
  mockStore: createMockStore
})

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

// This is needed for KaTeX to work in the tests. Guarded because specs that assert against
// platform-native behaviour opt out of the DOM with `@vitest-environment node`.
if (typeof document !== 'undefined') {
  Object.defineProperty(document, 'compatMode', {
    value: 'CSS1Compat'
  })
}

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
