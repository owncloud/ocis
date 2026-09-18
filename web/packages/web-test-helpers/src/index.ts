import { Plugin } from 'vue'
import DesignSystem from '@ownclouders/design-system'
import { registerTestingPlugins } from '@ownclouders/web-test-helpers-core'
import { createMockStore } from '@ownclouders/web-pkg/src/testing'

// `@ownclouders/web-test-helpers-core` provides an injectable `defaultPlugins` so it never
// has to import `@ownclouders/design-system` or `@ownclouders/web-pkg` directly (both of
// those depend on this package for testing, so a direct import back would recreate that
// cycle - see core's registry.ts). This package sits above core and web-pkg in the
// dependency graph and can safely depend on both, so it performs the registration once, as
// an import side effect, giving consumers of `defaultPlugins` from here the exact same
// behavior as before this split.
registerTestingPlugins({
  designSystem: DesignSystem as unknown as Plugin,
  mockStore: createMockStore
})

// Named (not `export *`) so `defaultPlugins`/`DefaultPluginsOptions` below - re-exported
// from `@ownclouders/web-pkg/testing`, typed with `PiniaMockOptions` - are the only ones
// visible here, matching this package's pre-split public API exactly.
export {
  mount,
  shallowMount,
  flushPromises,
  VueWrapper,
  getComposableWrapper,
  getOcSelectOptions,
  RouterLinkStub,
  createRouter,
  writable,
  sleep,
  nextTicks,
  registerTestingPlugins,
  getDesignSystemPlugin,
  getMockStoreFactory,
  defaultStubs,
  mockHttpResponse,
  mockHttpError,
  mockAxiosResolve,
  mockAxiosReject
} from '@ownclouders/web-test-helpers-core'
export type {
  RouteLocation,
  ComponentProps,
  PartialComponentProps
} from '@ownclouders/web-test-helpers-core'

export * from '@ownclouders/web-pkg/src/testing'
