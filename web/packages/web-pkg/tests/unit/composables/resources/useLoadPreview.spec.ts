import { ref } from 'vue'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { buildSpaceImageResource, Resource, SpaceResource } from '@ownclouders/web-client'
import { useLoadPreview } from '../../../../src/composables/resources'
import { usePreviewService } from '../../../../src/composables/previewService'
import { PreviewService, ProcessorType } from '../../../../src/services'
import { FolderViewModeConstants, ImageDimension } from '../../../../src'

vi.mock('../../../../src/composables/previewService/usePreviewService')
vi.mock('@ownclouders/web-client', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  buildSpaceImageResource: vi.fn()
}))

describe('useLoadPreview', () => {
  describe('loadPreview', () => {
    it('returns a loaded preview for a given file', () => {
      const loadedPreview = 'blob:image'
      getWrapper({
        setup: async ({ loadPreview }) => {
          const space = mock<SpaceResource>()
          const resource = mock<Resource>()
          const preview = await loadPreview({ space, resource })
          expect(preview).toEqual(loadedPreview)
        },
        loadedPreview
      })
    })
    describe('project space resources', () => {
      it('does not return a preview for a project space without spaceImageData', () => {
        getWrapper({
          setup: async ({ loadPreview }) => {
            const space = mock<SpaceResource>({ driveType: 'project', spaceImageData: undefined })
            const resource = space
            const preview = await loadPreview({ space, resource })
            expect(preview).toBe(null)
          }
        })
      })
      it('does not return a preview for a disabled project space', () => {
        getWrapper({
          setup: async ({ loadPreview }) => {
            const space = mock<SpaceResource>({ driveType: 'project', disabled: true })
            const resource = space
            const preview = await loadPreview({ space, resource })
            expect(preview).toBe(null)
          }
        })
      })
      it('calls buildSpaceImageResource to build a space image resource', () => {
        getWrapper({
          setup: async ({ loadPreview }) => {
            const buildSpaceImageResourceMock = vi.mocked(buildSpaceImageResource)
            const space = mock<SpaceResource>({ driveType: 'project', disabled: false })
            const resource = space
            const preview = await loadPreview({ space, resource })
            expect(preview).toBeDefined()
            expect(buildSpaceImageResourceMock).toHaveBeenCalledTimes(1)
          }
        })
      })
    })
    describe('dimensions', () => {
      it('uses the thumbnail default dimensions', () => {
        getWrapper({
          setup: async ({ loadPreview }, { previewService }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>()
            await loadPreview({ space, resource })
            expect(previewService.loadPreview).toHaveBeenCalledWith(
              expect.objectContaining({ dimensions: ImageDimension.Thumbnail }),
              expect.anything(),
              expect.anything(),
              expect.anything()
            )
          }
        })
      })
      it('uses tile default dimensions in tiles view', () => {
        getWrapper({
          setup: async ({ loadPreview }, { previewService }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>()
            await loadPreview({ space, resource })
            expect(previewService.loadPreview).toHaveBeenCalledWith(
              expect.objectContaining({ dimensions: ImageDimension.Tile }),
              expect.anything(),
              expect.anything(),
              expect.anything()
            )
          },
          viewMode: FolderViewModeConstants.name.tiles
        })
      })
      it('can overwrite the default dimensions', () => {
        getWrapper({
          setup: async ({ loadPreview }, { previewService }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>()
            await loadPreview({ space, resource, dimensions: ImageDimension.Preview })
            expect(previewService.loadPreview).toHaveBeenCalledWith(
              expect.objectContaining({ dimensions: ImageDimension.Preview }),
              expect.anything(),
              expect.anything(),
              expect.anything()
            )
          }
        })
      })
    })
    describe('processor', () => {
      it('uses the thumbnail default processor', () => {
        getWrapper({
          setup: async ({ loadPreview }, { previewService }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>()
            await loadPreview({ space, resource })
            expect(previewService.loadPreview).toHaveBeenCalledWith(
              expect.objectContaining({ processor: ProcessorType.enum.thumbnail }),
              expect.anything(),
              expect.anything(),
              expect.anything()
            )
          }
        })
      })
      it('uses the fit default processor in tiles view', () => {
        getWrapper({
          setup: async ({ loadPreview }, { previewService }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>()
            await loadPreview({ space, resource })
            expect(previewService.loadPreview).toHaveBeenCalledWith(
              expect.objectContaining({ processor: ProcessorType.enum.fit }),
              expect.anything(),
              expect.anything(),
              expect.anything()
            )
          },
          viewMode: FolderViewModeConstants.name.tiles
        })
      })
      it('can overwrite the default processor', () => {
        getWrapper({
          setup: async ({ loadPreview }, { previewService }) => {
            const space = mock<SpaceResource>()
            const resource = mock<Resource>()
            await loadPreview({ space, resource, processor: ProcessorType.enum.resize })
            expect(previewService.loadPreview).toHaveBeenCalledWith(
              expect.objectContaining({ processor: ProcessorType.enum.resize }),
              expect.anything(),
              expect.anything(),
              expect.anything()
            )
          }
        })
      })
    })
  })
})

function getWrapper({
  setup,
  loadedPreview = 'blob:image',
  viewMode
}: {
  setup: (
    instance: ReturnType<typeof useLoadPreview>,
    mocks: { previewService: PreviewService }
  ) => void
  loadedPreview?: string
  viewMode?: string
}) {
  const previewService = mock<PreviewService>()
  previewService.loadPreview.mockResolvedValue(loadedPreview)
  vi.mocked(usePreviewService).mockReturnValue(previewService)

  return {
    wrapper: getComposableWrapper(() => {
      const instance = useLoadPreview(viewMode ? ref(viewMode) : undefined)
      setup(instance, { previewService })
    })
  }
}
