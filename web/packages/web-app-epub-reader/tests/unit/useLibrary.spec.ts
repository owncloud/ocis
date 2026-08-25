import { mock } from 'vitest-mock-extended'
import { Resource, SearchResource, SpaceResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  nextTicks,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { useLibrary } from '../../src/composables/useLibrary'
import * as cache from '../../src/utils/cache'
import * as epub from '../../src/utils/epub'
import * as preferences from '../../src/utils/preferences'
import * as readerProgress from '../../src/utils/readerProgress'

const { copyToClipboard, downloadFile, showErrorMessage } = vi.hoisted(() => ({
  copyToClipboard: vi.fn(),
  downloadFile: vi.fn(),
  showErrorMessage: vi.fn()
}))

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@ownclouders/web-pkg')>()),
  useClipboard: () => ({ copyToClipboard }),
  useDownloadFile: () => ({ downloadFile }),
  useMessages: () => ({ showErrorMessage, showMessage: vi.fn() })
}))
vi.mock('../../src/utils/cache')
vi.mock('../../src/utils/epub')
vi.mock('../../src/utils/preferences')
vi.mock('../../src/utils/readerProgress')

const personalSpace = mock<SpaceResource>({ id: 'space-1', name: 'Personal', path: '/' })
const projectSpace = mock<SpaceResource>({ id: 'space-2', name: 'Project', path: '/' })

function epubResource(overrides: Partial<Resource> = {}): Resource {
  return {
    id: 'file-1',
    fileId: 'file-1',
    name: 'Book.epub',
    path: '/Books/Book.epub',
    spaceId: 'space-1',
    storageId: 'space-1',
    isFolder: false,
    type: 'file',
    mdate: 'Wed, 15 Jul 2026 10:00:00 GMT',
    privateLink: 'https://cloud.example/f/file-1',
    ...overrides
  } as Resource
}

function searchResult(resources: Resource[]) {
  return {
    resources: resources.map((resource) => ({ ...resource, highlights: '' })) as SearchResource[],
    totalResults: resources.length
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  copyToClipboard.mockResolvedValue(undefined)
  downloadFile.mockResolvedValue(undefined)
  vi.mocked(cache.getCachedMetadata).mockResolvedValue(null)
  vi.mocked(cache.setCachedMetadata).mockResolvedValue()
  vi.mocked(epub.extractEpubMetadata).mockResolvedValue({
    title: 'Book',
    authors: ['Author'],
    description: '',
    language: '',
    publisher: '',
    published: '',
    subjects: [],
    spineItemCount: 1
  })
  vi.mocked(epub.titleFromFileName).mockImplementation((name: string) =>
    name.replace(/\.epub$/, '')
  )
  vi.mocked(preferences.getBookPreferences).mockReturnValue({
    favorite: false,
    readingStatus: 'unread',
    shelfIds: []
  })
  vi.mocked(preferences.loadShelves).mockReturnValue([])
  vi.mocked(readerProgress.getReaderProgress).mockReturnValue({
    finished: false,
    hasPosition: false
  })
})

describe('useLibrary', () => {
  describe('discovery', () => {
    it('uses webdav.search when search-files is available and skips the folder scan', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library, clientService }) => {
        clientService.webdav.search.mockResolvedValue(searchResult([epubResource()]))

        await library.loadBooks()

        expect(clientService.webdav.search).toHaveBeenCalledWith(
          'name:*.epub',
          expect.objectContaining({ searchLimit: expect.any(Number) })
        )
        expect(clientService.webdav.listFiles).not.toHaveBeenCalled()
        expect(library.books.value).toHaveLength(1)
      })
    })

    it('falls back to scanning folders when search-files is not supported', async () => {
      await withLibrary({ searchAvailable: false }, async ({ library, clientService }) => {
        clientService.webdav.listFiles.mockResolvedValue({
          resource: mock<Resource>(),
          children: [epubResource()]
        })

        await library.loadBooks()

        expect(clientService.webdav.search).not.toHaveBeenCalled()
        expect(clientService.webdav.listFiles).toHaveBeenCalled()
        expect(library.books.value).toHaveLength(1)
      })
    })

    it('falls back to scanning folders when the search request fails', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library, clientService }) => {
        clientService.webdav.search.mockRejectedValue(new Error('boom'))
        clientService.webdav.listFiles.mockResolvedValue({
          resource: mock<Resource>(),
          children: [epubResource()]
        })

        await library.loadBooks()

        expect(clientService.webdav.listFiles).toHaveBeenCalled()
        expect(library.books.value).toHaveLength(1)
      })
    })

    it('deduplicates books discovered under the same id', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library, clientService }) => {
        clientService.webdav.search.mockResolvedValue(
          searchResult([epubResource(), epubResource({ path: '/Other/Book.epub' })])
        )

        await library.loadBooks()

        expect(library.books.value).toHaveLength(1)
      })
    })

    it('drops results whose space is not available', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library, clientService }) => {
        clientService.webdav.search.mockResolvedValue(
          searchResult([
            epubResource({ id: 'x', fileId: 'x', spaceId: 'unknown', storageId: 'unknown' })
          ])
        )

        await library.loadBooks()

        expect(library.books.value).toHaveLength(0)
      })
    })

    it('surfaces an error when all discovery fails', async () => {
      await withLibrary({ searchAvailable: false }, async ({ library, clientService }) => {
        clientService.webdav.listFiles.mockRejectedValue(new Error('unreachable'))

        await library.loadBooks()

        expect(library.error.value).toBe('unreachable')
        expect(library.books.value).toHaveLength(0)
      })
    })
  })

  describe('hydration', () => {
    it('reads metadata for every discovered book', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library, clientService }) => {
        clientService.webdav.search.mockResolvedValue(
          searchResult([
            epubResource(),
            epubResource({ id: 'file-2', fileId: 'file-2', path: '/Books/Second.epub' })
          ])
        )
        clientService.webdav.getFileContents.mockResolvedValue({
          response: { data: new Blob() }
        } as never)

        await library.loadBooks()

        expect(epub.extractEpubMetadata).toHaveBeenCalledTimes(2)
        expect(library.books.value.every((book) => !book.loadingMetadata)).toBe(true)
      })
    })

    it('records a metadata error without failing the whole load', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library, clientService }) => {
        clientService.webdav.search.mockResolvedValue(searchResult([epubResource()]))
        clientService.webdav.getFileContents.mockRejectedValue(new Error('no bytes'))

        await library.loadBooks()

        expect(library.books.value[0].metadataError).toBeTruthy()
        expect(library.books.value[0].loadingMetadata).toBe(false)
      })
    })
  })

  describe('actions', () => {
    it('copies the resource private link to the clipboard', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library }) => {
        const book = { resource: epubResource(), space: personalSpace } as never

        const success = await library.copyBookLink(book)

        expect(success).toBe(true)
        expect(copyToClipboard).toHaveBeenCalledWith('https://cloud.example/f/file-1')
      })
    })

    it('reports an error when no private link is available', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library }) => {
        const book = {
          resource: epubResource({ privateLink: undefined }),
          space: personalSpace
        } as never

        const success = await library.copyBookLink(book)

        expect(success).toBe(false)
        expect(copyToClipboard).not.toHaveBeenCalled()
        expect(showErrorMessage).toHaveBeenCalled()
      })
    })

    it('delegates downloads to useDownloadFile', async () => {
      await withLibrary({ searchAvailable: true }, async ({ library }) => {
        const resource = epubResource()
        const book = { resource, space: personalSpace } as never

        await library.downloadBook(book)

        expect(downloadFile).toHaveBeenCalledWith(personalSpace, resource)
      })
    })
  })
})

type LibraryInstance = ReturnType<typeof useLibrary>

async function withLibrary(
  { searchAvailable }: { searchAvailable: boolean },
  scenario: (ctx: {
    library: LibraryInstance
    clientService: ReturnType<typeof defaultComponentMocks>['$clientService']
  }) => Promise<void>
) {
  const componentMocks = defaultComponentMocks({
    currentRoute: mock<RouteLocation>({ name: 'epub-library' })
  })

  let library!: LibraryInstance
  const wrapper = getComposableWrapper(
    () => {
      library = useLibrary()
    },
    {
      mocks: componentMocks,
      provide: componentMocks,
      pluginOptions: {
        piniaOptions: {
          spacesState: { spaces: [personalSpace, projectSpace] },
          capabilityState: {
            capabilities: {
              dav: { reports: searchAvailable ? ['search-files'] : [] }
            } as never
          }
        }
      }
    }
  )

  try {
    await scenario({ library, clientService: componentMocks.$clientService })
    await nextTicks(2)
  } finally {
    wrapper.unmount()
  }
}
