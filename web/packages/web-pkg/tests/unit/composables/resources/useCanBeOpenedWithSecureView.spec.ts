import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import { useCanBeOpenedWithSecureView } from '../../../../src/composables/resources'
import { ApplicationFileExtension } from '../../../../src/apps/types'

describe('canBeOpenedWithSecureView', () => {
  describe('resource', () => {
    it('can be opened if a matching file extension with secure view exists', () => {
      getWrapper({
        setup: ({ canBeOpenedWithSecureView }) => {
          const canBeOpened = canBeOpenedWithSecureView(mock<Resource>({ mimeType: 'text/plain' }))
          expect(canBeOpened).toBeTruthy()
        },
        fileExtensions: [
          mock<ApplicationFileExtension>({ secureView: true, mimeType: 'text/plain' })
        ]
      })
    })
    it('can not be opened if no file extension with secure view exists', () => {
      getWrapper({
        setup: ({ canBeOpenedWithSecureView }) => {
          const canBeOpened = canBeOpenedWithSecureView(mock<Resource>({ mimeType: 'text/plain' }))
          expect(canBeOpened).toBeFalsy()
        },
        fileExtensions: [
          mock<ApplicationFileExtension>({ secureView: false, mimeType: 'text/plain' })
        ]
      })
    })
    it('can not be opened if no file extension exists', () => {
      getWrapper({
        setup: ({ canBeOpenedWithSecureView }) => {
          const canBeOpened = canBeOpenedWithSecureView(mock<Resource>({ mimeType: 'text/plain' }))
          expect(canBeOpened).toBeFalsy()
        },
        fileExtensions: []
      })
    })
    it("can not be opened if the file extension's mime type doesn't match the one of the resource", () => {
      getWrapper({
        setup: ({ canBeOpenedWithSecureView }) => {
          const canBeOpened = canBeOpenedWithSecureView(mock<Resource>({ mimeType: 'text/plain' }))
          expect(canBeOpened).toBeFalsy()
        },
        fileExtensions: [
          mock<ApplicationFileExtension>({ secureView: true, mimeType: 'image/jpg' })
        ]
      })
    })
  })
})

function getWrapper({
  setup,
  fileExtensions = [mock<ApplicationFileExtension>()]
}: {
  setup: (instance: ReturnType<typeof useCanBeOpenedWithSecureView>) => void
  fileExtensions?: ApplicationFileExtension[]
}) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useCanBeOpenedWithSecureView()
        setup(instance)
      },
      { pluginOptions: { piniaOptions: { appsState: { fileExtensions } } } }
    )
  }
}
