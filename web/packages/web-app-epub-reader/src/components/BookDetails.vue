<template>
  <div class="book-details-backdrop">
    <button
      type="button"
      class="backdrop-dismiss"
      tabindex="-1"
      :aria-label="$gettext('Close book details')"
      @click="emit('close')"
    />
    <!-- eslint-disable vuejs-accessibility/no-static-element-interactions -->
    <focus-trap :active="true">
      <aside
        class="book-details oc-background-default oc-box-shadow-large"
        role="dialog"
        aria-modal="true"
        aria-labelledby="book-details-title"
        tabindex="-1"
        @keydown.esc="emit('close')"
      >
        <header class="details-header oc-flex oc-flex-between oc-flex-middle oc-p-m">
          <h2 id="book-details-title" class="oc-m-rm">{{ $gettext('Book details') }}</h2>
          <oc-button appearance="raw" :aria-label="$gettext('Close')" @click="emit('close')">
            <oc-icon name="close" />
          </oc-button>
        </header>

        <div class="details-content oc-p-m">
          <div class="details-hero">
            <div class="details-cover">
              <img v-if="book.coverUrl" :src="book.coverUrl" :alt="''" />
              <div v-else class="details-cover-fallback">
                <oc-icon name="book" size="xlarge" />
                <span>{{ book.title }}</span>
              </div>
            </div>
            <div class="details-hero-copy">
              <button
                type="button"
                class="favorite-button"
                :class="{ 'favorite-button-active': book.favorite }"
                :aria-label="book.favorite ? $gettext('Remove favorite') : $gettext('Add favorite')"
                :aria-pressed="book.favorite"
                @click="emit('toggle-favorite')"
              >
                <oc-icon name="star" size="small" variation="inherit" />
                <span>{{ book.favorite ? $gettext('Favorited') : $gettext('Favorite') }}</span>
              </button>
              <h3 class="oc-mb-xs">{{ book.title }}</h3>
              <p class="oc-text-muted oc-m-rm">
                {{ book.authors.join(', ') || $gettext('Unknown author') }}
              </p>
              <p v-if="book.published" class="oc-text-muted oc-mt-xs">
                {{ book.published }}
              </p>
              <span class="status-badge oc-mt-s">{{ currentStatusLabel }}</span>
            </div>
          </div>

          <div class="primary-action oc-mt-l">
            <oc-button class="reading-action" @click="emit('open')">
              <oc-icon name="book" size="small" />
              {{ openActionLabel }}
            </oc-button>
          </div>
          <div class="secondary-actions oc-flex oc-mt-s oc-mb-l">
            <oc-button appearance="outline" @click="emit('show-in-files')">
              <oc-icon name="folder" size="small" />
              {{ $gettext('Show in folder') }}
            </oc-button>
            <oc-button appearance="outline" @click="emit('download')">
              <oc-icon name="download" size="small" />
              {{ $gettext('Download') }}
            </oc-button>
            <oc-button appearance="outline" @click="emit('copy-link')">
              <oc-icon name="link" size="small" />
              <span aria-live="polite">{{
                copied ? $gettext('Copied') : $gettext('Copy link')
              }}</span>
            </oc-button>
          </div>

          <section v-if="book.description" class="details-section description-section">
            <h3>{{ $gettext('About this book') }}</h3>
            <p class="details-description">{{ book.description }}</p>
          </section>

          <section class="organization-card oc-mt-l">
            <div class="organization-block">
              <fieldset class="status-fieldset">
                <legend>{{ $gettext('Reading status') }}</legend>
                <div class="status-options">
                  <label
                    v-for="option in statusOptions"
                    :key="option.value"
                    class="status-option"
                    :class="{ 'status-option-selected': book.readingStatus === option.value }"
                  >
                    <input
                      type="radio"
                      name="epub-reading-status"
                      :value="option.value"
                      :checked="book.readingStatus === option.value"
                      @change="emit('set-status', option.value)"
                    />
                    <span>{{ option.label }}</span>
                  </label>
                </div>
              </fieldset>
            </div>
            <div class="organization-block shelves-block">
              <h3>{{ $gettext('Shelves') }}</h3>
              <p v-if="!shelves.length" class="oc-text-muted">
                {{ $gettext('Create a shelf to organize this book.') }}
              </p>
              <div v-else class="shelf-options oc-flex oc-flex-wrap">
                <label v-for="shelf in shelves" :key="shelf.id" class="shelf-option">
                  <input
                    type="checkbox"
                    :checked="book.shelfIds.includes(shelf.id)"
                    @change="emit('toggle-shelf', shelf.id)"
                  />
                  <span>{{ shelf.name }}</span>
                </label>
              </div>
              <div class="new-shelf oc-flex oc-mt-m">
                <oc-text-input
                  v-model="newShelfName"
                  :label="$gettext('New shelf')"
                  :placeholder="$gettext('Shelf name')"
                  @keydown.enter="createShelf"
                />
                <oc-button
                  class="add-shelf-button"
                  size="small"
                  variant="primary"
                  :disabled="!newShelfName.trim()"
                  @click="createShelf"
                >
                  {{ $gettext('Add shelf') }}
                </oc-button>
              </div>
            </div>
          </section>

          <section v-if="hasPublicationMetadata" class="details-section oc-mt-l">
            <h3>{{ $gettext('Publication') }}</h3>
            <dl class="book-metadata">
              <template v-if="book.publisher">
                <dt>{{ $gettext('Publisher') }}</dt>
                <dd>{{ book.publisher }}</dd>
              </template>
              <template v-if="book.language">
                <dt>{{ $gettext('Language') }}</dt>
                <dd>{{ book.language }}</dd>
              </template>
              <template v-if="book.subjects.length">
                <dt>{{ $gettext('Subjects') }}</dt>
                <dd>{{ book.subjects.join(', ') }}</dd>
              </template>
            </dl>
          </section>

          <details class="file-details oc-mt-l">
            <summary>{{ $gettext('File details') }}</summary>
            <dl class="file-metadata">
              <dt>{{ $gettext('File name') }}</dt>
              <dd>{{ book.resource.name }}</dd>
              <dt>{{ $gettext('Location') }}</dt>
              <dd>{{ locationLabel }}</dd>
              <dt v-if="book.resource.size">{{ $gettext('Size') }}</dt>
              <dd v-if="book.resource.size">{{ formattedSize }}</dd>
            </dl>
          </details>
        </div>
      </aside>
    </focus-trap>
    <!-- eslint-enable vuejs-accessibility/no-static-element-interactions -->
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { FocusTrap } from 'focus-trap-vue'
import type { LibraryBook, LibraryShelf, ReadingStatus } from '../types'

const props = defineProps<{ book: LibraryBook; copied?: boolean; shelves: LibraryShelf[] }>()
const emit = defineEmits<{
  (event: 'close' | 'open' | 'show-in-files' | 'download' | 'copy-link'): void
  (event: 'toggle-favorite'): void
  (event: 'set-status', status: ReadingStatus): void
  (event: 'toggle-shelf', shelfId: string): void
  (event: 'create-shelf', name: string): void
}>()
const { $gettext } = useGettext()
const newShelfName = ref('')
const statusOptions: { label: string; value: ReadingStatus }[] = [
  { label: $gettext('Unread'), value: 'unread' },
  { label: $gettext('Reading'), value: 'reading' },
  { label: $gettext('Finished'), value: 'finished' }
]
const currentStatusLabel = computed(
  () => statusOptions.find(({ value }) => value === props.book.readingStatus)?.label ?? ''
)
const hasPublicationMetadata = computed(() =>
  Boolean(props.book.publisher || props.book.language || props.book.subjects.length)
)
const openActionLabel = computed(() => {
  if (props.book.readingProgress === 100) return $gettext('Open finished book')
  if (!props.book.hasReadingPosition) return $gettext('Start reading')
  if (props.book.readingProgress === undefined) return $gettext('Continue reading')
  return $gettext('Continue reading · %{progress}%', {
    progress: String(props.book.readingProgress)
  })
})
const locationLabel = computed(() => {
  const path = props.book.resource.path.replace(/^\/+/, '')
  return path ? `${props.book.space.name} / ${path}` : props.book.space.name
})

const formattedSize = computed(() => {
  const bytes = Number(props.book.resource.size || 0)
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
})

function createShelf(): void {
  const name = newShelfName.value.trim()
  if (!name) return
  emit('create-shelf', name)
  newShelfName.value = ''
}
</script>

<style scoped>
.book-details-backdrop {
  backdrop-filter: blur(2px);
  background: rgb(0 0 0 / 32%);
  display: flex;
  inset: 0;
  justify-content: flex-end;
  position: fixed;
  z-index: 1000;
}

.book-details {
  border-radius: 14px 0 0 14px;
  height: 100%;
  max-width: 100%;
  overflow-y: auto;
  position: relative;
  width: 34rem;
  z-index: 1;
}

.backdrop-dismiss {
  background: transparent;
  border: 0;
  cursor: default;
  inset: 0;
  padding: 0;
  position: absolute;
}

.details-header {
  background: var(--oc-color-background-default);
  border-bottom: 1px solid var(--oc-color-background-muted);
  position: sticky;
  top: 0;
  z-index: 1;
}

.details-hero {
  align-items: start;
  display: grid;
  gap: var(--oc-space-medium);
  grid-template-columns: 8rem minmax(0, 1fr);
  padding: var(--oc-space-medium);
  background: var(--oc-color-background-muted);
  border-radius: 12px;
  position: relative;
}

.details-hero-copy h3 {
  padding-right: 7rem;
}

.favorite-button {
  align-items: center;
  background: var(--oc-color-background-default);
  border: 1px solid var(--oc-color-swatch-passive-default);
  border-radius: 999px;
  color: inherit;
  cursor: pointer;
  display: inline-flex;
  gap: var(--oc-space-xsmall);
  min-height: 2rem;
  padding: 0 var(--oc-space-small);
  position: absolute;
  right: var(--oc-space-small);
  top: var(--oc-space-small);
}

.favorite-button-active {
  background: var(--oc-color-swatch-primary-default);
  border-color: var(--oc-color-swatch-primary-default);
  color: var(--oc-color-swatch-primary-contrast);
}

.favorite-button:focus-visible {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 2px;
}

.favorite-button:hover {
  border-color: var(--oc-color-swatch-primary-default);
  box-shadow: 0 0 0 2px var(--oc-color-background-highlight);
}

.details-cover {
  aspect-ratio: 2 / 3;
  background: var(--oc-color-background-muted);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 3px 8px 1px rgb(0 0 0 / 14%);
}

.details-cover img {
  height: 100%;
  object-fit: cover;
  width: 100%;
}

.details-cover-fallback {
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: var(--oc-space-small);
  height: 100%;
  justify-content: center;
  padding: var(--oc-space-small);
  text-align: center;
}

.secondary-actions {
  gap: var(--oc-space-small);
}

.reading-action {
  justify-content: center;
  width: 100%;
}

.secondary-actions > * {
  flex: 1;
}

.status-badge {
  background: var(--oc-color-swatch-primary-default);
  border-radius: 999px;
  color: var(--oc-color-swatch-primary-contrast);
  display: inline-flex;
  font-size: var(--oc-font-size-small);
  font-weight: 600;
  padding: var(--oc-space-xsmall) var(--oc-space-small);
}

.new-shelf {
  gap: var(--oc-space-small);
}

.organization-card {
  background: var(--oc-color-background-muted);
  border-radius: 10px;
  padding: var(--oc-space-medium);
}

.status-fieldset {
  border: 0;
  margin: 0;
  padding: 0;
}

.status-fieldset legend {
  font-weight: 600;
  margin-bottom: var(--oc-space-small);
}

.status-options {
  display: grid;
  gap: var(--oc-space-small);
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.status-option {
  align-items: center;
  background: var(--oc-color-background-default);
  border: 2px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  display: grid;
  gap: var(--oc-space-small);
  grid-template-columns: auto minmax(0, 1fr) auto;
  min-height: 2.75rem;
  padding: 0 var(--oc-space-small);
}

.status-option:hover {
  border-color: var(--oc-color-swatch-primary-default);
  box-shadow: 0 0 0 2px var(--oc-color-background-highlight);
}

.status-option:focus-within {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 2px;
}

.status-option-selected {
  background: var(--oc-color-background-highlight);
  border-color: var(--oc-color-swatch-primary-default);
  font-weight: 600;
}

.organization-block + .organization-block {
  border-top: 1px solid var(--oc-color-background-default);
  margin-top: var(--oc-space-medium);
  padding-top: var(--oc-space-medium);
}

.new-shelf {
  align-items: end;
}

.new-shelf > :first-child {
  flex: 1;
}

.add-shelf-button {
  min-height: 2.5rem;
}

.add-shelf-button:hover:not(:disabled) {
  box-shadow: 0 0 0 3px var(--oc-color-background-highlight);
  filter: brightness(0.88);
}

.add-shelf-button:focus-visible {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 2px;
}

.shelf-options {
  gap: var(--oc-space-small);
}

.shelf-option {
  align-items: center;
  display: flex;
  gap: var(--oc-space-small);
  background: var(--oc-color-background-default);
  border-radius: 999px;
  cursor: pointer;
  padding: var(--oc-space-xsmall) var(--oc-space-small);
}

.shelf-option:hover {
  box-shadow: inset 0 0 0 2px var(--oc-color-swatch-primary-default);
}

.shelf-option:focus-within {
  outline: 3px solid var(--oc-color-swatch-primary-default);
  outline-offset: 2px;
}

.details-description {
  line-height: 1.6;
  white-space: pre-line;
}

.book-metadata,
.file-metadata {
  display: grid;
  gap: var(--oc-space-small) var(--oc-space-medium);
  grid-template-columns: max-content minmax(0, 1fr);
}

.book-metadata {
  margin: 0;
}

.file-details {
  background: var(--oc-color-background-muted);
  border-radius: 10px;
  padding: var(--oc-space-medium);
}

.file-details summary {
  cursor: pointer;
  font-weight: 600;
}

.file-metadata {
  margin: var(--oc-space-medium) 0 0;
}

.book-metadata dt,
.file-metadata dt {
  font-weight: 600;
}

.book-metadata dd,
.file-metadata dd {
  margin: 0;
  overflow-wrap: anywhere;
}

@media (max-width: 520px) {
  .book-details {
    width: 100%;
  }

  .status-options {
    grid-template-columns: 1fr;
  }

  .secondary-actions {
    flex-wrap: wrap;
  }

  .secondary-actions > * {
    flex: 1 0 45%;
  }
}
</style>
