import { unref } from 'vue'
import { useFileActionsCopyPermanentLink } from '../../../../../src/composables/actions/files'
import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useClipboard } from '../../../../../src/composables/clipboard'

vi.mock('../../../../../src/composables/clipboard', () => ({
  useClipboard: vi.fn()
}))

describe('useFileActionsCopyPermanentLink', () => {
  describe('isVisible property', () => {
    it('should return false if no resource selected', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(unref(actions)[0].isVisible({ space: null, resources: [] })).toBeFalsy()
        }
      })
    })
    it('should return false in public spaces', () => {
      getWrapper({
        setup: ({ actions }) => {
          const publicSpace = mock<SpaceResource>({ driveType: 'public' })
          expect(
            unref(actions)[0].isVisible({ space: publicSpace, resources: [mock<Resource>()] })
          ).toBeFalsy()
        }
      })
    })
    it('should return true if one resource selected', () => {
      getWrapper({
        setup: ({ actions }) => {
          expect(
            unref(actions)[0].isVisible({ resources: [mock<Resource>()], space: undefined })
          ).toBeTruthy()
        }
      })
    })
  })
  describe('handler', () => {
    it('calls the copyToClipboard method with the private link of the resource', () => {
      getWrapper({
        setup: async ({ actions }, { mocks }) => {
          const privateLink = 'https://example.com'
          await unref(actions)[0].handler({
            resources: [mock<Resource>({ privateLink })],
            space: mock<SpaceResource>()
          })
          expect(mocks.copyToClipboardMock).toHaveBeenCalledWith(privateLink)
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (
    instance: ReturnType<typeof useFileActionsCopyPermanentLink>,
    mocks: Record<string, any>
  ) => void
}) {
  const copyToClipboardMock = vi.fn()
  vi.mocked(useClipboard).mockReturnValue({ copyToClipboard: copyToClipboardMock })
  const mocks = { ...defaultComponentMocks(), copyToClipboardMock }

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActionsCopyPermanentLink()
        setup(instance, { mocks })
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
