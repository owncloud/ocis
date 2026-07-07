import { useResourceContents } from '../../../src/composables/resources/useResourceContents'
import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { unref } from 'vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { describe } from 'vitest'

describe('resourceContents', () => {
  describe('resourceContentsText', () => {
    it('should contain space count when route equals "files-shares-via-link"', () => {
      const resources = [
        mock<Resource>({ isFolder: true, type: 'folder', name: 'folder1' }),
        mock<SpaceResource>({ driveType: 'project' })
      ]
      getWrapper({
        currentRouteName: 'files-shares-via-link',
        resources,
        setup: ({ resourceContentsText }) => {
          expect(unref(resourceContentsText)).toBe('2 items in total (0 files, 1 folder, 1 space)')
        }
      })
    })
    it('should contain hidden count when areHiddenFilesShown equals false', () => {
      const resources = [
        mock<Resource>({ isFolder: true, type: 'folder', name: 'folder1' }),
        mock<Resource>({ isFolder: true, type: 'folder', name: '.hiddenFolder1' }),
        mock<Resource>({ isFolder: false, type: 'file', name: 'file1' }),
        mock<Resource>({ isFolder: false, type: 'file', name: '.hiddenFile1' })
      ]
      getWrapper({
        resources,
        areHiddenFilesShown: false,
        setup: ({ resourceContentsText }) => {
          expect(unref(resourceContentsText)).toBe(
            '4 items in total (2 files including 1 hidden, 2 folders including 1 hidden)'
          )
        }
      })
    })
    it.each([
      { prop: { resources: [] }, expectedText: '0 items in total (0 files, 0 folders)' },
      {
        prop: {
          resources: [mock<Resource>({ isFolder: true, type: 'folder', name: 'folder1' })]
        },
        expectedText: '1 item in total (0 files, 1 folder)'
      },
      {
        prop: { resources: [mock<Resource>({ isFolder: false, type: 'file', name: 'file1' })] },
        expectedText: '1 item in total (1 file, 0 folders)'
      },
      {
        prop: {
          resources: [
            mock<Resource>({
              isFolder: false,
              type: 'file',
              name: 'file1'
            }),
            mock<Resource>({
              isFolder: true,
              type: 'folder',
              name: 'folder1'
            })
          ]
        },
        expectedText: '2 items in total (1 file, 1 folder)'
      },
      {
        prop: {
          resources: [
            mock<Resource>({
              isFolder: false,
              type: 'file',
              name: 'file1'
            }),
            mock<Resource>({
              isFolder: false,
              type: 'file',
              name: 'file2'
            }),
            mock<Resource>({
              isFolder: true,
              type: 'folder',
              name: 'folder1'
            }),
            mock<Resource>({
              isFolder: true,
              type: 'folder',
              name: 'folder2'
            })
          ]
        },
        expectedText: '4 items in total (2 files, 2 folders)'
      }
    ])('should be singular or plural according to item, files and folders count', (cases) => {
      getWrapper({
        resources: cases.prop.resources,
        setup: ({ resourceContentsText }) => {
          expect(unref(resourceContentsText)).toBe(cases.expectedText)
        }
      })
    })
  })
  it.each`
    size              | expectedSize
    ${1}              | ${'1 B'}
    ${100}            | ${'100 B'}
    ${10000}          | ${'10 kB'}
    ${10000000}       | ${'10 MB'}
    ${10000000000}    | ${'10 GB'}
    ${10000000000000} | ${'10 TB'}
  `('should display correctly size according to items', ({ size, expectedSize }) => {
    const resources = [
      mock<Resource>({
        isFolder: false,
        size: parseInt(size),
        type: 'file',
        name: 'file1'
      }),
      mock<Resource>({
        isFolder: false,
        size: 0,
        type: 'file',
        name: 'file2'
      }),
      mock<Resource>({
        isFolder: true,
        size: 0,
        type: 'folder',
        name: 'folder1'
      }),
      mock<Resource>({
        isFolder: true,
        size: 0,
        type: 'folder',
        name: 'folder2'
      }),
      mock<Resource>({
        isFolder: true,
        size: 0,
        type: 'folder',
        name: 'folder3'
      })
    ]

    getWrapper({
      resources,
      setup: ({ resourceContentsText }) => {
        expect(unref(resourceContentsText)).toBe(
          `5 items with ${expectedSize} in total (2 files, 3 folders)`
        )
      }
    })
  })
})

function getWrapper({
  areHiddenFilesShown = true,
  currentRouteName = 'files-spaces-generic',
  resources = [],
  setup
}: {
  areHiddenFilesShown?: boolean
  currentRouteName?: string
  resources: Resource[]
  setup: (instance: ReturnType<typeof useResourceContents>) => void
}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: currentRouteName })
    })
  }

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useResourceContents()
        setup(instance)
      },
      {
        mocks,
        pluginOptions: {
          piniaOptions: {
            resourcesStore: { resources, areHiddenFilesShown }
          }
        },
        provide: mocks
      }
    )
  }
}
