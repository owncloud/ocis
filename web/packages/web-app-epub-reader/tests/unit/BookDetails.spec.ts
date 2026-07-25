import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { defaultPlugins, mount, PartialComponentProps } from '@ownclouders/web-test-helpers'
import BookDetails from '../../src/components/BookDetails.vue'
import type { LibraryBook } from '../../src/types'

function libraryBook(overrides: Partial<LibraryBook> = {}): LibraryBook {
  return {
    id: 'file-1',
    resource: mock<Resource>({ name: 'Book.epub', path: '/Books/Book.epub', size: '2048' }),
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

describe('BookDetails', () => {
  it('renders the book title and author', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.find('#book-details-title').exists()).toBe(true)
    expect(wrapper.html()).toContain('A Book')
    expect(wrapper.html()).toContain('Ada Lovelace')
  })

  it('emits "close" when Escape is pressed on the dialog', async () => {
    const { wrapper } = getWrapper()
    await wrapper.get('[role="dialog"]').trigger('keydown.esc')
    expect(wrapper.emitted('close')).toBeTruthy()
  })

  it('emits action events from the secondary actions', async () => {
    const { wrapper } = getWrapper()
    const buttons = wrapper.findAll('.secondary-actions button')
    await buttons[0].trigger('click')
    await buttons[1].trigger('click')
    await buttons[2].trigger('click')
    expect(wrapper.emitted('show-in-files')).toBeTruthy()
    expect(wrapper.emitted('download')).toBeTruthy()
    expect(wrapper.emitted('copy-link')).toBeTruthy()
  })

  it('shows the copied state when the copied prop is set', () => {
    const { wrapper } = getWrapper({ copied: true })
    expect(wrapper.html()).toContain('Copied')
  })
})

function getWrapper(props: PartialComponentProps<typeof BookDetails> = {}) {
  return {
    wrapper: mount(BookDetails, {
      props: {
        book: libraryBook(),
        shelves: [],
        ...props
      },
      global: {
        plugins: [...defaultPlugins()],
        stubs: {
          'focus-trap': {
            template: '<slot />'
          }
        }
      }
    })
  }
}
