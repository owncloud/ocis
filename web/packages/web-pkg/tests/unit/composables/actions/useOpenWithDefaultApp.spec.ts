import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import {
  useFileActions,
  Action,
  useOpenWithDefaultApp,
  FileAction
} from '../../../../src/composables'

vi.mock('../../../../src/composables/actions/files', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActions: vi.fn()
}))

describe('useOpenWithDefaultApp', () => {
  it('should be valid', () => {
    expect(useOpenWithDefaultApp).toBeDefined()
  })
  describe('method "openWithDefaultApp"', () => {
    it('should call the default action handler for files', () => {
      getWrapper({
        setup: ({ openWithDefaultApp }, { defaultEditorAction }) => {
          openWithDefaultApp({
            space: mock<SpaceResource>(),
            resource: mock<Resource>({ isFolder: false })
          })
          expect(defaultEditorAction.handler).toHaveBeenCalled()
        }
      })
    })
    it('should not call the default action handler for folders', () => {
      getWrapper({
        setup: ({ openWithDefaultApp }, { defaultEditorAction }) => {
          openWithDefaultApp({
            space: mock<SpaceResource>(),
            resource: mock<Resource>({ isFolder: true })
          })
          expect(defaultEditorAction.handler).not.toHaveBeenCalled()
        }
      })
    })
  })
})

function getWrapper({
  setup,
  defaultEditorAction = mock<Action>({ handler: vi.fn() })
}: {
  setup: (
    instance: ReturnType<typeof useOpenWithDefaultApp>,
    mocks: { defaultEditorAction: FileAction }
  ) => void
  defaultEditorAction?: FileAction
}) {
  vi.mocked(useFileActions).mockReturnValue(
    mock<ReturnType<typeof useFileActions>>({
      getDefaultAction: () => defaultEditorAction
    })
  )

  const mocks = { defaultEditorAction }

  return {
    wrapper: getComposableWrapper(() => {
      const instance = useOpenWithDefaultApp()
      setup(instance, mocks)
    })
  }
}
