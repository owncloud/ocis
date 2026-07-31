import { useCapabilityStore } from '../../../src/composables/index'
import { getBackendVersion, getWebVersion } from '../../../src/helpers/versions'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { Capabilities } from '@ownclouders/web-client/ocs'

describe('collect version information', () => {
  describe('web version', () => {
    beforeEach(() => {
      process.env.PACKAGE_VERSION = '4.7.0'
    })
    it('provides the web version with a static string without exceptions', () => {
      expect(getWebVersion()).toBe('ownCloud Web UI 4.7.0')
    })
  })
  describe('backend version', () => {
    it('returns undefined when the backend version object is not available', () => {
      const capabilityStore = versionStore(undefined)
      expect(getBackendVersion({ capabilityStore })).toBeUndefined()
    })
    it('returns undefined when the backend version object has no "string" field', () => {
      const capabilityStore = versionStore({
        product: 'ownCloud',
        versionstring: undefined
      })
      expect(getBackendVersion({ capabilityStore })).toBeUndefined()
    })
    it('falls back to "ownCloud" as a product when none is defined', () => {
      const capabilityStore = versionStore({
        versionstring: '10.8.0',
        edition: 'Community'
      })
      expect(getBackendVersion({ capabilityStore })).toBe('ownCloud 10.8.0 Community')
    })
    it('provides the backend version as concatenation of product, version and edition', () => {
      const capabilityStore = versionStore({
        product: 'oCIS',
        versionstring: '1.16.0',
        edition: 'Reva'
      })
      expect(getBackendVersion({ capabilityStore })).toBe('oCIS 1.16.0 Reva')
    })
    it('prefers the productversion over versionstring field if both are provided', () => {
      const capabilityStore = versionStore({
        product: 'oCIS',
        versionstring: '10.8.0',
        productversion: '2.0.0',
        edition: 'Community'
      })
      expect(getBackendVersion({ capabilityStore })).toBe('oCIS 2.0.0 Community')
    })
  })
})

const versionStore = (version: Capabilities['capabilities']['core']['status']) => {
  createTestingPinia()
  const capabilityStore = useCapabilityStore()
  capabilityStore.capabilities.core.status = version
  return capabilityStore
}
