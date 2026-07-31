import FileActions from '../../../../../src/components/SideBar/Actions/FileActions.vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import {
  defaultPlugins,
  defaultStubs,
  mount,
  defaultComponentMocks,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { useFileActions } from '@ownclouders/web-pkg'
import { Action } from '@ownclouders/web-pkg'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useFileActions: vi.fn()
}))

type ActionWithSelector = Action & { selector: string }
const fileActions: Record<string, ActionWithSelector> = {
  copy: mock<ActionWithSelector>({
    handler: vi.fn(),
    label: () => 'Copy',
    class: 'oc-files-actions-copy-trigger',
    selector: '.oc-files-actions-copy-trigger'
  }),
  move: mock<ActionWithSelector>({
    handler: vi.fn(),
    label: () => 'Move',
    class: 'oc-files-actions-move-trigger',
    selector: '.oc-files-actions-move-trigger'
  }),
  download: mock<ActionWithSelector>({
    handler: vi.fn(),
    label: () => 'Download',
    class: 'oc-files-actions-download-file-trigger',
    selector: '.oc-files-actions-download-file-trigger'
  }),
  'text-editor': mock<ActionWithSelector>({
    handler: vi.fn(),
    label: () => 'Open in Text Editor',
    class: 'oc-files-actions-text-editor-trigger',
    selector: '.oc-files-actions-text-editor-trigger'
  })
}

describe('FileActions', () => {
  describe('when user is on personal route', () => {
    describe('action handlers', () => {
      it('renders action handlers as clickable elements', async () => {
        vi.mocked(useFileActions).mockImplementation(() =>
          mock<ReturnType<typeof useFileActions>>({
            getAllAvailableActions: () => Object.values(fileActions)
          })
        )

        const actions = ['copy', 'move', 'download', 'text-editor']
        const { wrapper } = getWrapper()

        for (const button of actions) {
          const action = fileActions[button]

          const buttonElement = wrapper.find(action.selector)
          expect(buttonElement.exists()).toBeTruthy()

          await buttonElement.trigger('click.stop')
          expect(action.handler).toHaveBeenCalledTimes(1)
        }
      })
    })

    describe('menu items', () => {
      it('renders a list of actions', () => {
        const { wrapper } = getWrapper()
        for (const action of ['copy', 'text-editor']) {
          expect(wrapper.find(fileActions[action].selector).exists()).toBeTruthy()
        }
      })
    })
  })
})

function getWrapper() {
  return {
    wrapper: mount(FileActions, {
      global: {
        plugins: [...defaultPlugins()],
        mocks: defaultComponentMocks({
          currentRoute: mock<RouteLocation>({
            name: 'files-spaces-generic',
            path: '/files/spaces/personal/admin'
          })
        }),
        stubs: { ...defaultStubs, OcButton: false },
        provide: {
          space: mock<SpaceResource>(),
          resource: mock<Resource>({
            extension: 'md'
          })
        }
      }
    })
  }
}
