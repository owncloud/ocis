import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useCanShare } from '../../../../src/composables/shares'
import { useCapabilityStore } from '../../../../src/composables/piniaStores'

describe('useCanShare', () => {
  describe('canShare', () => {
    describe('server managed spaces', () => {
      it('should disable sharing of spaces when server managed spaces capability is enabled', () => {
        getWrapper({
          setup: ({ canShare }) => {
            const space = mock<SpaceResource>()
            const resource = mock<SpaceResource>({ type: 'space', canShare: vi.fn(() => true) })

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).capabilities.spaces.server_managed = true

            expect(canShare({ space, resource })).toBeFalsy()
          }
        })
      })

      it('should not disable sharing of spaces when server managed spaces capability is disabled', () => {
        getWrapper({
          setup: ({ canShare }) => {
            const space = mock<SpaceResource>()
            const resource = mock<SpaceResource>({ type: 'space', canShare: vi.fn(() => true) })

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).capabilities.spaces.server_managed = false

            expect(canShare({ space, resource })).toBeTruthy()
          }
        })
      })

      it('should not disable sharing of resources when server managed spaces capability is enabled', () => {
        getWrapper({
          setup: ({ canShare }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>({ canShare: vi.fn(() => true) })

            const capabilityStore = useCapabilityStore()
            vi.mocked(capabilityStore).capabilities.spaces.server_managed = true

            expect(canShare({ space, resource })).toBeTruthy()
          }
        })
      })
    })
  })
})

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useCanShare>) => void }) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useCanShare()
      setup(instance)
    })
  }
}
