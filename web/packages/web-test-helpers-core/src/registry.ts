import { Plugin } from 'vue'

/**
 * `defaultPlugins` needs a DesignSystem Vue plugin and a pinia mock-store factory, but this
 * package must not import either directly: `@ownclouders/design-system` and
 * `@ownclouders/web-pkg` both depend on `@ownclouders/web-test-helpers` for testing, and
 * `web-test-helpers` depends on this package, so a direct import back would recreate that
 * cycle.
 *
 * Instead, `@ownclouders/web-test-helpers` (the published package built on top of this one)
 * registers the real DesignSystem plugin and store factory once, as an import side effect of
 * its own entrypoint. `defaultPlugins` reads them from here at call time. Callers only see
 * this indirection if they import straight from `@ownclouders/web-test-helpers-core` without
 * going through the facade - see the two internal consumers (`design-system`, `web-pkg`),
 * which register providers from the shared vitest setup file instead (not part of the package
 * dependency graph either).
 */

type MockStoreFactory = (options?: unknown) => Plugin

let designSystemPlugin: Plugin | undefined
let mockStoreFactory: MockStoreFactory | undefined

export const registerTestingPlugins = ({
  designSystem,
  mockStore
}: {
  designSystem?: Plugin
  mockStore?: MockStoreFactory
}): void => {
  if (designSystem) {
    designSystemPlugin = designSystem
  }
  if (mockStore) {
    mockStoreFactory = mockStore
  }
}

export const getDesignSystemPlugin = (): Plugin => {
  if (!designSystemPlugin) {
    throw new Error(
      'No design system plugin registered. Make sure registerTestingPlugins() ran before ' +
        'defaultPlugins() - either via @ownclouders/web-test-helpers (which does this ' +
        'automatically) or the shared vitest setup file.'
    )
  }
  return designSystemPlugin
}

export const getMockStoreFactory = (): MockStoreFactory => {
  if (!mockStoreFactory) {
    throw new Error(
      'No mock store factory registered. Make sure registerTestingPlugins() ran before ' +
        'defaultPlugins() - either via @ownclouders/web-test-helpers (which does this ' +
        'automatically) or the shared vitest setup file.'
    )
  }
  return mockStoreFactory
}
