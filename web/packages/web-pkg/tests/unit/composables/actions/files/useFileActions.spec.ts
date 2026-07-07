import { mock } from 'vitest-mock-extended'
import { Action, FileActionOptions, useFileActions } from '../../../../../src/composables/actions'
import {
  defaultComponentMocks,
  RouteLocation,
  getComposableWrapper
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
