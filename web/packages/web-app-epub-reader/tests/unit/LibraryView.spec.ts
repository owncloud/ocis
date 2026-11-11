import { ref } from 'vue'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import LibraryView from '../../src/views/LibraryView.vue'
import type { LibraryBook } from '../../src/types'

const { libraryState, loadBooks, openBook } = vi.hoisted(() => ({
  libraryState: {
    books: undefined as never,
    visibleBooks: undefined as never,
    shelves: undefined as never,
    query: undefined as never,
    sort: undefined as never,
    loading: undefined as never,
    error: undefined as never
  },
  loadBooks: vi.fn(),
  openBook: vi.fn()
}))

vi.mock('../../src/composables/useLibrary', () => ({
  useLibrary: () => ({
    books: libraryState.books,
    visibleBooks: libraryState.visibleBooks,
    shelves: libraryState.shelves,
    query: libraryState.query,
    sort: libraryState.sort,
    loading: libraryState.loading,
    error: libraryState.error,
    loadBooks,
    openBook,
    copyBookLink: vi.fn(),
    createShelf: vi.fn(),
    downloadBook: vi.fn(),
    setReadingStatus: vi.fn(),
    showBookInFiles: vi.fn(),
    toggleBookShelf: vi.fn(),
    toggleFavorite: vi.fn()
  })
}))

function libraryBook(overrides: Partial<LibraryBook> = {}): LibraryBook {
  return {
    id: 'file-1',
    resource: mock<Resource>({ name: 'Book.epub', path: '/Books/Book.epub' }),
    space: mock<SpaceResource>({ name: 'Personal' }),
    title: 'A Book',
    authors: ['Ada Lovelace'],
    description: '',
    language: '',
    publisher: '',
    published: '',
    subjects: [],
    spineItemCount: 1,
    loadingMetadata: false,
    favorite: false,
    readingStatus: 'unread',
    shelfIds: [],
    hasReadingPosition: false,
    ...overrides
  }
}

beforeEach(() => {
  vi.clearAllMocks()
})

describe('LibraryView', () => {
  it('loads books on mount', () => {
    getWrapper()
    expect(loadBooks).toHaveBeenCalled()
  })

  it('shows a loading state while loading with no books', () => {
    const { wrapper } = getWrapper({ loading: true, books: [] })
    expect(wrapper.find('.library-state[role="status"]').exists()).toBe(true)
  })

  it('shows an error state when loading fails', () => {
    const { wrapper } = getWrapper({ error: 'Boom' })
    const alert = wrapper.find('.library-state[role="alert"]')
    expect(alert.exists()).toBe(true)
    expect(alert.text()).toContain('Boom')
  })

  it('shows the empty state when there are no books', () => {
    const { wrapper } = getWrapper({ books: [], visibleBooks: [] })
    expect(wrapper.html()).toContain('No EPUB books found')
  })

  it('renders a card per visible book', () => {
    const books = [libraryBook(), libraryBook({ id: 'file-2', title: 'Second' })]
    const { wrapper } = getWrapper({ books, visibleBooks: books })
    expect(wrapper.findAll('.book-card')).toHaveLength(2)
  })

  it('opens a book when its open button is clicked', async () => {
    const books = [libraryBook()]
    const { wrapper } = getWrapper({ books, visibleBooks: books })
    await wrapper.find('.book-open').trigger('click')
    expect(openBook).toHaveBeenCalledWith(books[0])
  })
})

function getWrapper(
  state: {
    books?: LibraryBook[]
    visibleBooks?: LibraryBook[]
    loading?: boolean
    error?: string
  } = {}
) {
  libraryState.books = ref(state.books ?? []) as never
  libraryState.visibleBooks = ref(state.visibleBooks ?? state.books ?? []) as never
  libraryState.shelves = ref([]) as never
  libraryState.query = ref('') as never
  libraryState.sort = ref('recent') as never
  libraryState.loading = ref(state.loading ?? false) as never
  libraryState.error = ref(state.error ?? '') as never

  return {
    wrapper: mount(LibraryView, {
      global: {
        plugins: [...defaultPlugins()],
        stubs: { 'oc-select': true, 'oc-spinner': true, BookDetails: true }
      }
    })
  }
}
