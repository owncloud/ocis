<template>
  <main class="library oc-p-m" :aria-label="$gettext('Library')">
    <header class="library-header oc-flex oc-flex-between oc-flex-middle oc-mb-l">
      <div>
        <h1 class="oc-m-rm">{{ $gettext('Library') }}</h1>
        <p class="oc-text-muted oc-m-rm">
          {{ $gettext('Browse EPUB books stored in your ownCloud spaces.') }}
        </p>
      </div>
      <oc-button appearance="outline" :disabled="loading" @click="loadBooks">
        <oc-icon name="refresh" size="small" />
        {{ $gettext('Refresh') }}
      </oc-button>
    </header>

    <div class="library-controls oc-flex oc-flex-between oc-flex-middle oc-mb-l">
      <oc-text-input
        v-model="query"
        class="library-search"
        :label="$gettext('Search books')"
        :placeholder="$gettext('Title, author, publisher or subject')"
      />
      <oc-select
        :model-value="selectedSortOption"
        class="library-sort"
        :label="$gettext('Sort by')"
        :options="sortOptions"
        option-label="label"
        @update:model-value="onSortUpdate"
      />
    </div>

    <section
      v-if="showContinueReading"
      class="continue-section oc-mb-xl"
      aria-labelledby="continue-reading-heading"
    >
      <div class="continue-section-heading oc-flex oc-flex-between oc-flex-middle oc-mb-m">
        <div>
          <h2 id="continue-reading-heading" class="oc-m-rm">
            {{ $gettext('Continue reading') }}
          </h2>
          <p class="oc-text-muted oc-m-rm">{{ $gettext('Pick up where you left off.') }}</p>
        </div>
      </div>
      <div class="continue-list">
        <button
          v-for="book in continueReadingBooks"
          :key="book.id"
          type="button"
          class="continue-card"
          :aria-label="$gettext('Continue reading %{title}', { title: book.title })"
          @click="openBook(book)"
        >
          <span class="continue-cover">
            <img v-if="book.coverUrl" :src="book.coverUrl" :alt="''" />
            <oc-icon v-else name="book" size="large" />
          </span>
          <span class="continue-card-content">
            <strong>{{ book.title }}</strong>
            <span class="oc-text-muted">
              {{ book.authors.join(', ') || $gettext('Unknown author') }}
            </span>
            <span class="continue-card-progress">
              {{
                book.readingProgress === undefined
                  ? $gettext('Resume')
                  : $gettext('%{progress}% complete', { progress: String(book.readingProgress) })
              }}
            </span>
            <span
              v-if="book.readingProgress !== undefined"
              class="reading-progress"
              role="progressbar"
              :aria-label="$gettext('Reading progress for %{title}', { title: book.title })"
              aria-valuemin="0"
              aria-valuemax="100"
              :aria-valuenow="book.readingProgress"
            >
              <span :style="{ width: `${book.readingProgress}%` }" />
            </span>
          </span>
        </button>
      </div>
    </section>

    <div class="library-toolbar oc-flex oc-flex-between oc-flex-middle oc-mb-l">
      <div class="library-filters">
        <label v-for="filter in primaryFilters" :key="filter.key" class="library-filter">
          <span>{{ filter.label }}</span>
          <select v-model="filterValues[filter.key]">
            <option v-for="option in filter.options" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </label>
        <oc-button
          appearance="outline"
          class="advanced-filter-toggle"
          aria-controls="advanced-library-filters"
          :aria-expanded="showAdvancedFilters"
          @click="showAdvancedFilters = !showAdvancedFilters"
        >
          <oc-icon name="filter" size="small" />
          {{ $gettext('More filters') }}
          <span v-if="advancedFilterCount">({{ advancedFilterCount }})</span>
        </oc-button>
      </div>
      <div class="view-toggle oc-flex" role="group" :aria-label="$gettext('View mode')">
        <button
          type="button"
          class="view-toggle-button"
          :class="{ 'view-toggle-button-selected': viewMode === 'grid' }"
          :aria-pressed="viewMode === 'grid'"
          @click="setViewMode('grid')"
        >
          {{ $gettext('Grid') }}
        </button>
        <button
          type="button"
          class="view-toggle-button"
          :class="{ 'view-toggle-button-selected': viewMode === 'list' }"
          :aria-pressed="viewMode === 'list'"
          @click="setViewMode('list')"
        >
          {{ $gettext('List') }}
        </button>
      </div>
    </div>

    <div
      v-if="showAdvancedFilters"
      id="advanced-library-filters"
      class="advanced-filters oc-mb-m oc-p-m"
    >
      <label v-for="filter in advancedFilters" :key="filter.key" class="library-filter">
        <span>{{ filter.label }}</span>
        <select v-model="filterValues[filter.key]">
          <option v-for="option in filter.options" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>
      </label>
    </div>

    <div v-if="activeFilters.length" class="active-filters oc-flex oc-flex-wrap oc-mb-l">
      <oc-button
        v-for="filter in activeFilters"
        :key="filter.key"
        appearance="raw"
        class="filter-chip"
        :aria-label="
          $gettext('Remove %{label} filter: %{value}', {
            label: filter.label,
            value: filter.valueLabel
          })
        "
        @click="clearFilter(filter.key)"
      >
        <span>{{ filter.label }}: {{ filter.valueLabel }}</span>
        <oc-icon name="close" size="small" />
      </oc-button>
      <oc-button appearance="raw" @click="resetFilters">{{ $gettext('Clear all') }}</oc-button>
    </div>

    <div v-if="loading && !books.length" class="library-state" role="status">
      <oc-spinner size="large" />
      <p>{{ $gettext('Finding EPUB books…') }}</p>
    </div>

    <div v-else-if="error" class="library-state" role="alert">
      <oc-icon name="error-warning" size="xlarge" />
      <h2>{{ $gettext('Library unavailable') }}</h2>
      <p>{{ error }}</p>
      <oc-button @click="loadBooks">{{ $gettext('Try again') }}</oc-button>
    </div>

    <div v-else-if="!filteredBooks.length" class="library-state">
      <oc-icon name="book" size="xlarge" />
      <h2>
        {{ hasActiveFilters ? $gettext('No matching books') : $gettext('No EPUB books found') }}
      </h2>
      <p v-if="!hasActiveFilters">
        {{ $gettext('Upload EPUB files to an accessible ownCloud space and refresh the library.') }}
      </p>
    </div>

    <template v-else>
      <p class="oc-text-muted oc-mb-m" role="status">
        {{ $gettext('%{count} books', { count: String(filteredBooks.length) }) }}
        <span v-if="loading"> · {{ $gettext('Reading metadata…') }}</span>
      </p>
      <div class="bookshelf" :class="`bookshelf-${viewMode}`">
        <article v-for="book in filteredBooks" :key="book.id" class="book-card">
          <oc-button
            appearance="raw"
            class="book-details-button"
            :aria-label="$gettext('View details for %{title}', { title: book.title })"
            @click="showDetails(book, $event)"
          >
            <oc-icon name="information" size="small" />
          </oc-button>
          <button
            type="button"
            class="book-open"
            :aria-label="$gettext('Open %{title}', { title: book.title })"
            @click="openBook(book)"
          >
            <span class="book-cover">
              <img v-if="book.coverUrl" :src="book.coverUrl" :alt="''" />
              <span v-else class="book-cover-fallback">
                <oc-icon name="book" size="xlarge" />
                <span>{{ book.title }}</span>
              </span>
            </span>
            <span class="book-title">{{ book.title }}</span>
            <span class="book-author">
              {{ book.authors.join(', ') || $gettext('Unknown author') }}
            </span>
            <span v-if="book.hasReadingPosition" class="continue-reading">
              <template v-if="book.readingProgress === 100">{{ $gettext('Finished') }}</template>
              <template v-else>
                {{ $gettext('Continue reading') }}
                <span v-if="book.readingProgress !== undefined">
                  · {{ book.readingProgress }}%</span
                >
              </template>
            </span>
            <span
              v-if="book.readingProgress !== undefined"
              class="reading-progress"
              role="progressbar"
              :aria-label="$gettext('Reading progress for %{title}', { title: book.title })"
              aria-valuemin="0"
              aria-valuemax="100"
              :aria-valuenow="book.readingProgress"
            >
              <span :style="{ width: `${book.readingProgress}%` }" />
            </span>
            <span v-if="book.loadingMetadata" class="book-status oc-text-muted">
              {{ $gettext('Reading metadata…') }}
            </span>
            <span v-else-if="book.metadataError" class="book-status oc-text-muted">
              {{ book.metadataError }}
            </span>
            <span v-else-if="book.published" class="book-status oc-text-muted">
              {{ book.published.slice(0, 4) }}
            </span>
          </button>
          <div class="book-card-meta">
            <div
              v-if="book.favorite || book.readingStatus !== 'unread'"
              class="book-labels oc-flex oc-flex-wrap oc-mt-xs"
            >
              <span v-if="book.favorite" class="book-label">{{ $gettext('Favorite') }}</span>
              <span v-if="book.readingStatus !== 'unread'" class="book-label">
                {{ readingStatusLabel(book.readingStatus) }}
              </span>
            </div>
          </div>
        </article>
      </div>
    </template>

    <book-details
      v-if="selectedBook"
      :book="selectedBook"
      :copied="copied"
      :shelves="shelves"
      @close="closeDetails"
      @open="openBook(selectedBook)"
      @show-in-files="showBookInFiles(selectedBook)"
      @download="downloadBook(selectedBook)"
      @copy-link="copySelectedBookLink"
      @toggle-favorite="toggleFavorite(selectedBook)"
      @set-status="setReadingStatus(selectedBook, $event)"
      @toggle-shelf="toggleBookShelf(selectedBook, $event)"
      @create-shelf="createAndAssignShelf"
    />
  </main>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useGettext } from 'vue3-gettext'
import BookDetails from '../components/BookDetails.vue'
import { useLibrary } from '../composables/useLibrary'
import { useLibraryFilters } from '../composables/useLibraryFilters'
import type { LibraryBook, LibrarySort, LibraryViewMode, ReadingStatus } from '../types'
import { loadViewMode, saveViewMode } from '../utils/preferences'

const { $gettext } = useGettext()
const {
  books,
  copyBookLink,
  createShelf,
  downloadBook,
  error,
  loadBooks,
  loading,
  openBook,
  query,
  shelves,
  setReadingStatus,
  showBookInFiles,
  sort,
  toggleBookShelf,
  toggleFavorite,
  visibleBooks
} = useLibrary()
const selectedBook = ref<LibraryBook | null>(null)
const copied = ref(false)
const viewMode = ref<LibraryViewMode>(loadViewMode())
const showAdvancedFilters = ref(false)
let copiedTimer: ReturnType<typeof setTimeout> | undefined
let detailsTrigger: HTMLElement | null = null
const {
  activeFilters,
  clearFilter,
  filteredBooks,
  filters,
  filterValues,
  hasActiveFilters,
  resetFilters
} = useLibraryFilters({ books, query, shelves, visibleBooks })

const primaryFilters = computed(() =>
  filters.value.filter(({ key }) => ['collection', 'status'].includes(key))
)
const advancedFilters = computed(() =>
  filters.value.filter(({ key }) => !['collection', 'status'].includes(key))
)
const advancedFilterCount = computed(
  () => advancedFilters.value.filter(({ key }) => filterValues[key] !== 'all').length
)
const continueReadingBooks = computed(() =>
  books.value
    .filter(
      ({ hasReadingPosition, readingProgress, readingStatus }) =>
        hasReadingPosition && readingProgress !== 100 && readingStatus !== 'finished'
    )
    .sort((left, right) => (right.readingProgress ?? 0) - (left.readingProgress ?? 0))
    .slice(0, 6)
)
const showContinueReading = computed(
  () =>
    !loading.value &&
    !error.value &&
    !hasActiveFilters.value &&
    continueReadingBooks.value.length > 0
)

const sortOptions = [
  { label: $gettext('Recently added'), value: 'recent' },
  { label: $gettext('Title'), value: 'title' },
  { label: $gettext('Author'), value: 'author' }
]

const selectedSortOption = computed(
  () => sortOptions.find((option) => option.value === sort.value) ?? sortOptions[0]
)

function onSortUpdate(option: { value: LibrarySort } | null): void {
  if (option) sort.value = option.value
}

async function closeDetails(): Promise<void> {
  selectedBook.value = null
  copied.value = false
  if (copiedTimer) clearTimeout(copiedTimer)
  await nextTick()
  detailsTrigger?.focus()
  detailsTrigger = null
}

function showDetails(book: LibraryBook, event: MouseEvent): void {
  detailsTrigger = event.currentTarget as HTMLElement
  copied.value = false
  selectedBook.value = book
}

function setViewMode(mode: LibraryViewMode): void {
  viewMode.value = mode
  saveViewMode(mode)
}

function readingStatusLabel(status: ReadingStatus): string {
  return {
    unread: $gettext('Unread'),
    reading: $gettext('Reading'),
    finished: $gettext('Finished')
  }[status]
}

function createAndAssignShelf(name: string): void {
  if (!selectedBook.value) return
  const shelf = createShelf(name)
  if (shelf && !selectedBook.value.shelfIds.includes(shelf.id)) {
    toggleBookShelf(selectedBook.value, shelf.id)
  }
}

async function copySelectedBookLink(): Promise<void> {
  if (!selectedBook.value) return
  const success = await copyBookLink(selectedBook.value)
  if (!success) return
  copied.value = true
  if (copiedTimer) clearTimeout(copiedTimer)
  copiedTimer = setTimeout(() => (copied.value = false), 2000)
}

onMounted(loadBooks)
</script>

<style scoped>
.library {
  box-sizing: border-box;
  height: 100%;
  min-height: 0;
  overflow-y: auto;
  overscroll-behavior: contain;
  width: 100%;
}

.library-header,
.library-controls {
  gap: var(--oc-space-medium);
}

.library-toolbar {
  align-items: end;
  gap: var(--oc-space-medium);
}

.continue-section {
  background: var(--oc-color-background-muted);
  border-radius: 12px;
  padding: var(--oc-space-medium);
}

.continue-list {
  display: grid;
  gap: var(--oc-space-medium);
  grid-auto-columns: minmax(17rem, 21rem);
  grid-auto-flow: column;
  overflow-x: auto;
  padding: var(--oc-space-xsmall) var(--oc-space-xsmall) var(--oc-space-small);
  scroll-snap-type: x proximity;
}

.continue-card {
  align-items: center;
  background: var(--oc-color-background-default);
  border: 0;
  border-radius: 10px;
  box-shadow: 0 2px 6px rgb(0 0 0 / 12%);
  color: inherit;
  cursor: pointer;
  display: grid;
  gap: var(--oc-space-medium);
  grid-template-columns: 4.5rem minmax(0, 1fr);
  padding: var(--oc-space-small);
  scroll-snap-align: start;
  text-align: left;
}

.continue-card:focus-visible {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 2px;
}

.continue-cover {
  align-items: center;
  aspect-ratio: 2 / 3;
  background: var(--oc-color-background-muted);
  border-radius: 6px;
  display: flex;
  justify-content: center;
  overflow: hidden;
}

.continue-cover img {
  height: 100%;
  object-fit: cover;
  width: 100%;
}

.continue-card-content,
.continue-card-content > span {
  display: block;
  min-width: 0;
}

.continue-card-content {
  overflow: hidden;
}

.continue-card-content strong,
.continue-card-content > span:not(.reading-progress) {
  display: block;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.continue-card-progress {
  color: var(--oc-color-swatch-primary-default);
  font-size: var(--oc-font-size-small);
  font-weight: 600;
  margin-top: var(--oc-space-small);
}

.library-filters {
  display: flex;
  flex: 1;
  flex-wrap: wrap;
  gap: var(--oc-space-small);
}

.advanced-filter-toggle {
  align-self: end;
}

.advanced-filters {
  background: var(--oc-color-background-muted);
  border-radius: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: var(--oc-space-medium);
}

.active-filters {
  align-items: center;
  gap: var(--oc-space-small);
}

.filter-chip {
  background: var(--oc-color-background-muted);
  border-radius: 999px;
  gap: var(--oc-space-xsmall);
  padding: var(--oc-space-xsmall) var(--oc-space-small);
}

.library-filter {
  display: grid;
  font-size: var(--oc-font-size-small);
  gap: var(--oc-space-xsmall);
}

.library-filter select {
  background: var(--oc-color-background-default);
  border: 1px solid var(--oc-color-background-muted);
  border-radius: 5px;
  color: inherit;
  min-height: 2.5rem;
  padding: 0 var(--oc-space-small);
}

.view-toggle {
  background: var(--oc-color-background-muted);
  border-radius: 8px;
  gap: var(--oc-space-xsmall);
  padding: var(--oc-space-xsmall);
}

.view-toggle-button {
  background: var(--oc-color-background-default);
  border: 2px solid var(--oc-color-swatch-passive-default);
  border-radius: 6px;
  color: inherit;
  cursor: pointer;
  font: inherit;
  font-weight: 600;
  min-height: 2.5rem;
  padding: 0 var(--oc-space-medium);
}

.view-toggle-button:hover {
  border-color: var(--oc-color-swatch-primary-default);
  box-shadow: 0 0 0 2px var(--oc-color-background-highlight);
}

.view-toggle-button:focus-visible {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 2px;
}

.view-toggle-button-selected {
  background: var(--oc-color-swatch-primary-default);
  border-color: var(--oc-color-swatch-primary-default);
  color: var(--oc-color-swatch-primary-contrast);
}

.library-controls {
  align-items: end;
}

.library-search {
  max-width: 32rem;
  width: 100%;
}

.library-sort {
  min-width: 12rem;
}

.library-state {
  align-items: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 20rem;
  text-align: center;
}

.bookshelf {
  display: grid;
  gap: var(--oc-space-large) var(--oc-space-medium);
  grid-template-columns: repeat(auto-fill, minmax(10rem, 1fr));
}

.book-card {
  min-width: 0;
  position: relative;
}

.book-labels {
  gap: var(--oc-space-xsmall);
}

.book-label {
  background: var(--oc-color-background-muted);
  border-radius: 5px;
  font-size: var(--oc-font-size-xsmall);
  padding: 0 var(--oc-space-xsmall);
}

.bookshelf-list {
  display: block;
}

.bookshelf-list .book-card {
  align-items: stretch;
  border-bottom: 1px solid var(--oc-color-background-muted);
  display: grid;
  gap: var(--oc-space-medium);
  grid-template-columns: minmax(0, 1fr) minmax(10rem, 14rem);
  padding: var(--oc-space-small) 0;
}

.bookshelf-list .book-open {
  align-items: center;
  column-gap: var(--oc-space-medium);
  display: grid;
  grid-template-areas:
    'cover title continue'
    'cover author progress'
    'cover status status';
  grid-template-columns: 3.5rem minmax(8rem, 2fr) minmax(7rem, 1fr);
  row-gap: var(--oc-space-xsmall);
}

.bookshelf-list .book-cover {
  grid-area: cover;
}

.bookshelf-list .book-title {
  grid-area: title;
}

.bookshelf-list .book-author {
  grid-area: author;
}

.bookshelf-list .book-status {
  grid-area: status;
  margin-top: 0;
}

.bookshelf-list .continue-reading {
  grid-area: continue;
}

.bookshelf-list .reading-progress {
  align-self: center;
  grid-area: progress;
  margin-top: 0;
  width: 100%;
}

.bookshelf-list .book-card-meta {
  align-items: center;
  align-self: end;
  display: flex;
  gap: var(--oc-space-small);
  justify-content: flex-end;
  min-height: 1.5rem;
  white-space: nowrap;
}

.bookshelf-list .book-labels {
  margin-top: 0;
}

.bookshelf-list .book-labels {
  flex-wrap: nowrap;
}

.bookshelf-list .book-title,
.bookshelf-list .book-author {
  margin-top: 0;
}

.book-open {
  background: transparent;
  border: 0;
  color: inherit;
  cursor: pointer;
  display: block;
  padding: 0;
  text-align: left;
  width: 100%;
}

.book-open:focus-visible .book-cover {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 3px;
}

.book-cover {
  aspect-ratio: 2 / 3;
  background: var(--oc-color-background-muted);
  border-radius: 8px;
  box-shadow: 0 3px 8px 1px rgb(0 0 0 / 14%);
  display: block;
  overflow: hidden;
  width: 100%;
}

.book-cover img {
  height: 100%;
  object-fit: cover;
  width: 100%;
}

.book-cover-fallback {
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: var(--oc-space-medium);
  height: 100%;
  justify-content: center;
  padding: var(--oc-space-medium);
  text-align: center;
}

.book-title,
.book-author,
.book-status {
  display: block;
}

.book-details-button {
  background: var(--oc-color-background-default);
  border-radius: 50%;
  box-shadow: 0 2px 6px rgb(0 0 0 / 18%);
  min-height: 2rem;
  min-width: 2rem;
  padding: 0;
  position: absolute;
  right: var(--oc-space-small);
  top: var(--oc-space-small);
  z-index: 2;
}

.book-title {
  font-weight: 600;
  margin-top: var(--oc-space-small);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.book-author,
.book-status,
.continue-reading {
  font-size: var(--oc-font-size-small);
  margin-top: var(--oc-space-xsmall);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.continue-reading {
  color: var(--oc-color-swatch-primary-default);
  display: block;
  font-size: var(--oc-font-size-small);
  font-weight: 600;
  margin-top: var(--oc-space-small);
}

.reading-progress {
  background: var(--oc-color-background-muted);
  border-radius: 999px;
  display: block;
  height: 0.25rem;
  margin-top: var(--oc-space-xsmall);
  overflow: hidden;
}

.reading-progress > span {
  background: var(--oc-color-swatch-primary-default);
  display: block;
  height: 100%;
}

@media (max-width: 640px) {
  .library-header,
  .library-controls,
  .library-toolbar {
    align-items: stretch;
    flex-direction: column;
  }

  .library-sort {
    min-width: 0;
    width: 100%;
  }

  .library-filters,
  .library-filter,
  .advanced-filter-toggle {
    width: 100%;
  }

  .advanced-filters {
    display: grid;
  }

  .view-toggle {
    align-self: flex-end;
    width: fit-content;
  }

  .continue-list {
    grid-auto-columns: minmax(15rem, 85%);
  }

  .bookshelf-list .book-card {
    grid-template-columns: 1fr;
  }

  .bookshelf-list .book-open {
    grid-template-areas:
      'cover title'
      'cover author'
      'cover status'
      'cover continue'
      'cover progress';
    grid-template-columns: 3.5rem minmax(0, 1fr);
  }

  .bookshelf-list .book-title {
    padding-right: 2.75rem;
  }

  .bookshelf-list .book-card-meta {
    justify-content: flex-start;
  }
}
</style>
