import { mock } from 'vitest-mock-extended'
import { Action, FileActionOptions, useFileActions } from '../../../../../src/composables/actions'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper,
  createRouter
} from '@ownclouders/web-test-helpers'
import { computed, unref } from 'vue'
import { describe } from 'vitest'
import { Resource, SpaceResource } from '@ownclouders/web-client'

const mockUseEmbedMode = vi.fn().mockReturnValue({ isEnabled: computed(() => false) })
vi.mock('../../../../../src/composables/embedMode', () => ({
  useEmbedMode: vi.fn().mockImplementation(() => mockUseEmbedMode())
}))

describe('fileActions', () => {
  describe('computed property "editorActions"', () => {
    it('should provide a list of editors', () => {
      getWrapper({
        setup: ({ editorActions }) => {
          expect(unref(editorActions).length).toEqual(2)
        }
      })
    })
    it('should provide an empty list if embed mode is enabled', () => {
      mockUseEmbedMode.mockReturnValueOnce({
        isEnabled: computed(() => true)
      })
      getWrapper({
        setup: ({ editorActions }) => {
          expect(unref(editorActions).length).toBeFalsy()
        }
      })
    })

    it('should hide action when editor with matching routeName is opened', () => {
      getWrapper({
        currentRoute: mock<RouteLocation>({ name: 'text-editor' }),
        setup: ({ editorActions }) => {
          const [textEditor] = unref(editorActions)

          expect(
            (textEditor as Action<FileActionOptions>).isVisible({
              space: mock<SpaceResource>(),
              resources: [
                mock<Resource>({
                  id: '2',
                  extension: 'txt',
                  mimeType: 'text/txt',
                  canDownload: () => true
                })
              ]
            })
          ).toStrictEqual(false)
        }
      })
    })
  })
  describe('editor action "route"', () => {
    it.each([
      ['#.txt', '/personal%2Fadmin%2F%23.txt'],
      ['a?b.txt', '/personal%2Fadmin%2Fa%3Fb.txt'],
      ['a#b?c.txt', '/personal%2Fadmin%2Fa%23b%3Fc.txt']
    ])('percent-encodes "%s" when the router resolves the route', (fileName, expectedPath) => {
      getWrapper({
        setup: ({ editorActions }) => {
          const [textEditor] = unref(editorActions)
          const route = (textEditor as Action<FileActionOptions>).route({
            space: mock<SpaceResource>({
              getDriveAliasAndItem: () => `personal/admin/${fileName}`
            }),
            resources: [mock<Resource>({ path: `/${fileName}`, extension: 'txt' })]
          })

          expect(createTextEditorRouter().resolve(route).path).toEqual(expectedPath)
        }
      })
    })
  })
  describe('secure view context', () => {
    describe('computed property "editorActions"', () => {
      it('only displays editors that support secure view', () => {
        getWrapper({
          setup: ({ editorActions }) => {
            const secureViewResource = mock<Resource>({
              id: '1',
              canDownload: () => false,
              mimeType: 'text/txt',
              extension: 'txt'
            })
            const actions = unref(editorActions)
            expect(actions.length).toEqual(2)
            expect(
              actions[0].isVisible({ resources: [secureViewResource], space: null })
            ).toBeFalsy()
            expect(
              actions[1].isVisible({ resources: [secureViewResource], space: null })
            ).toBeTruthy()
          }
        })
      })
    })
  })
})

function createTextEditorRouter() {
  return createRouter({
    routes: [
      { path: '/:driveAliasAndItem(.*)?', name: 'text-editor', component: { template: '<div />' } }
    ]
  })
}

function getWrapper({
  setup,
  currentRoute = mock<RouteLocation>({ name: 'files-spaces-generic' })
}: {
  setup: (instance: ReturnType<typeof useFileActions>) => void
  currentRoute?: RouteLocation
}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute
    })
  }
  return {
    mocks,
    wrapper: getComposableWrapper(
      () => {
        const instance = useFileActions()
        setup(instance)
      },
      {
        mocks,
        provide: mocks,
        pluginOptions: {
          piniaOptions: {
            appsState: {
              apps: {
                'text-editor': {
                  defaultExtension: 'txt',
                  icon: 'file-text',
                  name: 'Text Editor',
                  id: 'text-editor',
                  color: '#0D856F',
                  extensions: [
                    {
                      extension: 'txt'
                    }
                  ],
                  hasEditor: true
                },
                external: {
                  defaultExtension: '',
                  icon: 'check_box_outline_blank',
                  name: 'External',
                  id: 'external',
                  hasEditor: true
                },
                'editor-less': {
                  defaultExtension: '',
                  icon: 'check_box_outline_blank',
                  name: 'Editor Less',
                  id: 'editor-less',
                  hasEditor: false
                }
              },
              fileExtensions: [
                {
                  app: 'text-editor',
                  extension: 'txt',
                  hasPriority: false,
                  routeName: 'text-editor'
                },
                {
                  app: 'external',
                  label: 'Open in Collabora',
                  mimeType: 'text/txt',
                  routeName: 'external-apps',
                  icon: 'https://host.docker.internal:9980/favicon.ico',
                  name: 'Collabora',
                  hasPriority: false,
                  secureView: true
                }
              ]
            }
          }
        }
      }
    )
  }
}
