import { computed, reactive, type Ref } from 'vue'
import { useGettext } from 'vue3-gettext'
import type { LibraryBook, LibraryShelf } from '../types'

export type LibraryFilterKey = 'collection' | 'author' | 'subject' | 'language' | 'space' | 'status'

interface FilterOption {
  label: string
  value: string
}

export interface LibraryFilter {
  key: LibraryFilterKey
  label: string
  options: FilterOption[]
}

export interface ActiveLibraryFilter {
  key: LibraryFilterKey
  label: string
  valueLabel: string
}

interface UseLibraryFiltersOptions {
  books: Ref<LibraryBook[]>
  visibleBooks: Ref<LibraryBook[]>
  shelves: Ref<LibraryShelf[]>
  query: Ref<string>
}

function uniqueOptions(values: string[]): FilterOption[] {
  return [...new Set(values.filter(Boolean))]
    .sort((left, right) => left.localeCompare(right))
    .map((value) => ({ label: value, value }))
}

export function useLibraryFilters({
  books,
  visibleBooks,
  shelves,
  query
}: UseLibraryFiltersOptions) {
  const { $gettext } = useGettext()
  const filterValues = reactive<Record<LibraryFilterKey, string>>({
    collection: 'all',
    author: 'all',
    subject: 'all',
    language: 'all',
    space: 'all',
    status: 'all'
  })

  const filters = computed<LibraryFilter[]>(() => [
    {
      key: 'collection',
      label: $gettext('Collection'),
      options: [
        { label: $gettext('All books'), value: 'all' },
        { label: $gettext('Favorites'), value: 'favorites' },
        ...shelves.value.map(({ id, name }) => ({ label: name, value: `shelf:${id}` }))
      ]
    },
    {
      key: 'author',
      label: $gettext('Author'),
      options: [
        { label: $gettext('All authors'), value: 'all' },
        ...uniqueOptions(books.value.flatMap(({ authors }) => authors))
      ]
    },
    {
      key: 'subject',
      label: $gettext('Subject'),
      options: [
        { label: $gettext('All subjects'), value: 'all' },
        ...uniqueOptions(books.value.flatMap(({ subjects }) => subjects))
      ]
    },
    {
      key: 'language',
      label: $gettext('Language'),
      options: [
        { label: $gettext('All languages'), value: 'all' },
        ...uniqueOptions(books.value.map(({ language }) => language))
      ]
    },
    {
      key: 'space',
      label: $gettext('Space'),
      options: [
        { label: $gettext('All spaces'), value: 'all' },
        ...uniqueOptions(books.value.map(({ space }) => space.name))
      ]
    },
    {
      key: 'status',
      label: $gettext('Status'),
      options: [
        { label: $gettext('All statuses'), value: 'all' },
        { label: $gettext('Unread'), value: 'unread' },
        { label: $gettext('Reading'), value: 'reading' },
        { label: $gettext('Finished'), value: 'finished' }
      ]
    }
  ])

  const filteredBooks = computed(() =>
    visibleBooks.value.filter((book) => {
      const collection = filterValues.collection
      if (collection === 'favorites' && !book.favorite) return false
      if (collection.startsWith('shelf:') && !book.shelfIds.includes(collection.slice(6))) {
        return false
      }
      if (filterValues.author !== 'all' && !book.authors.includes(filterValues.author)) return false
      if (filterValues.subject !== 'all' && !book.subjects.includes(filterValues.subject)) {
        return false
      }
      if (filterValues.language !== 'all' && book.language !== filterValues.language) return false
      if (filterValues.space !== 'all' && book.space.name !== filterValues.space) return false
      return filterValues.status === 'all' || book.readingStatus === filterValues.status
    })
  )

  const hasActiveFilters = computed(
    () =>
      Boolean(query.value.trim()) || Object.values(filterValues).some((value) => value !== 'all')
  )

  const activeFilters = computed<ActiveLibraryFilter[]>(() =>
    filters.value.flatMap((filter) => {
      const selectedValue = filterValues[filter.key]
      if (selectedValue === 'all') return []
      const selectedOption = filter.options.find(({ value }) => value === selectedValue)
      return selectedOption
        ? [{ key: filter.key, label: filter.label, valueLabel: selectedOption.label }]
        : []
    })
  )

  function clearFilter(key: LibraryFilterKey): void {
    filterValues[key] = 'all'
  }

  function resetFilters(): void {
    Object.keys(filterValues).forEach((key) => (filterValues[key as LibraryFilterKey] = 'all'))
  }

  return {
    activeFilters,
    clearFilter,
    filteredBooks,
    filters,
    filterValues,
    hasActiveFilters,
    resetFilters
  }
}
