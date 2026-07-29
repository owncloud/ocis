import type { Resource, SpaceResource } from '@ownclouders/web-client'
import { DavProperties } from '@ownclouders/web-client/webdav'
import {
  encodePath,
  useCapabilityStore,
  useClientService,
  useClipboard,
  useDownloadFile,
  useFolderLink,
  useMessages,
  useRouter,
  useSpacesStore
} from '@ownclouders/web-pkg'
import { computed, onUnmounted, ref } from 'vue'
import { useGettext } from 'vue3-gettext'
import type { LibraryBook, LibraryShelf, LibrarySort, ReadingStatus } from '../types'
import { getCachedMetadata, setCachedMetadata } from '../utils/cache'
import { extractEpubMetadata, titleFromFileName } from '../utils/epub'
import {
  getBookPreferences,
  loadShelves,
  saveBookPreferences,
  saveShelves
} from '../utils/preferences'
import { getReaderProgress } from '../utils/readerProgress'

const SEARCH_LIMIT = 500
const METADATA_WORKERS = 4
const MAX_SCANNED_FOLDERS = 2000
const EPUB_SEARCH_PATTERN = 'name:*.epub'

export function useLibrary() {
  const { $gettext } = useGettext()
  const { showErrorMessage } = useMessages()
  const clientService = useClientService()
  const capabilityStore = useCapabilityStore()
  const spacesStore = useSpacesStore()
  const router = useRouter()
  const { downloadFile } = useDownloadFile()
  const { copyToClipboard } = useClipboard()
  const { getParentFolderLink } = useFolderLink()

  const searchAvailable = computed(() => !!capabilityStore.davReports?.includes('search-files'))

  const books = ref<LibraryBook[]>([])
  const loading = ref(false)
  const error = ref('')
  const query = ref('')
  const sort = ref<LibrarySort>('recent')
  const shelves = ref<LibraryShelf[]>(loadShelves())

  const visibleBooks = computed(() => {
    const needle = query.value.trim().toLocaleLowerCase()
    const filtered = needle
      ? books.value.filter((book) =>
          [book.title, ...book.authors, book.publisher, ...book.subjects]
            .join(' ')
            .toLocaleLowerCase()
            .includes(needle)
        )
      : [...books.value]

    return filtered.sort((left, right) => {
      if (sort.value === 'title') return left.title.localeCompare(right.title)
      if (sort.value === 'author') {
        return (left.authors[0] ?? '').localeCompare(right.authors[0] ?? '')
      }
      return (
        new Date(right.resource.mdate ?? 0).getTime() - new Date(left.resource.mdate ?? 0).getTime()
      )
    })
  })

  async function searchBooks(): Promise<Resource[]> {
    const { resources } = await clientService.webdav.search(EPUB_SEARCH_PATTERN, {
      searchLimit: SEARCH_LIMIT,
      davProperties: DavProperties.Default
    })
    return resources.filter((resource) => resource.name.toLowerCase().endsWith('.epub'))
  }

  async function scanSpace(space: SpaceResource): Promise<Resource[]> {
    const pendingPaths = [space.path || '/']
    const visitedPaths = new Set<string>()
    const epubs: Resource[] = []

    while (pendingPaths.length && visitedPaths.size < MAX_SCANNED_FOLDERS) {
      const path = pendingPaths.shift()!
      if (visitedPaths.has(path)) continue
      visitedPaths.add(path)

      const { children } = await clientService.webdav.listFiles(space, { path }, { depth: 1 })

      for (const resource of children ?? []) {
        if (resource.isFolder || resource.type === 'folder') {
          if (resource.path && !visitedPaths.has(resource.path)) pendingPaths.push(resource.path)
          continue
        }
        if (!resource.name.toLowerCase().endsWith('.epub')) continue
        epubs.push({
          ...resource,
          mimeType: resource.mimeType || 'application/epub+zip',
          spaceId: resource.spaceId || space.id,
          storageId: resource.storageId || space.id
        })
      }
    }

    if (visitedPaths.size >= MAX_SCANNED_FOLDERS && pendingPaths.length) {
      console.warn(
        `EPUB library scan reached the ${MAX_SCANNED_FOLDERS} folder limit in "${space.name}"; results may be incomplete.`
      )
    }

    return epubs
  }

  async function scanSpaces(spaces: SpaceResource[]): Promise<Resource[]> {
    const results = await Promise.allSettled(spaces.map(scanSpace))
    const discovered = results.flatMap((result) =>
      result.status === 'fulfilled' ? result.value : []
    )
    if (!discovered.length && results.every((result) => result.status === 'rejected')) {
      throw (results[0] as PromiseRejectedResult).reason
    }
    return discovered
  }

  async function discoverBooks(spaces: SpaceResource[]): Promise<Resource[]> {
    // Prefer the server-side search REPORT (fast, one request). Only walk every
    // folder as a fallback when the server lacks search-files support or the
    // search request fails, rather than always paying for both.
    if (searchAvailable.value) {
      try {
        return await searchBooks()
      } catch (cause) {
        console.warn('EPUB library search failed, falling back to folder scan', cause)
      }
    }

    return scanSpaces(spaces)
  }

  async function hydrateBook(book: LibraryBook): Promise<void> {
    try {
      const cachedMetadata = await getCachedMetadata(book.resource)
      if (cachedMetadata) {
        Object.assign(book, cachedMetadata)
        updateReaderProgress(book)
        return
      }

      const { response } = await clientService.webdav.getFileContents(
        book.space,
        { path: book.resource.path },
        { responseType: 'blob' }
      )
      const metadata = await extractEpubMetadata(response.data as Blob, book.title)
      Object.assign(book, metadata)
      updateReaderProgress(book)
      await setCachedMetadata(book.resource, metadata)
    } catch (cause) {
      console.warn(`Could not read EPUB metadata for ${book.resource.name}`, cause)
      book.metadataError = $gettext('Metadata could not be read')
    } finally {
      book.loadingMetadata = false
    }
  }

  function updateReaderProgress(book: LibraryBook): void {
    const readerProgress = getReaderProgress(book.resource, book.spineItemCount)
    book.hasReadingPosition = readerProgress.hasPosition
    book.readingProgress = readerProgress.percentage
    if (readerProgress.finished && book.readingStatus !== 'finished') {
      book.readingStatus = 'finished'
      persistBook(book)
    }
  }

  async function hydrateBooks(items: LibraryBook[]): Promise<void> {
    let next = 0
    async function worker() {
      while (next < items.length) {
        const book = items[next++]
        await hydrateBook(book)
      }
    }
    await Promise.all(Array.from({ length: Math.min(METADATA_WORKERS, items.length) }, worker))
  }

  async function loadBooks(): Promise<void> {
    if (loading.value) return
    loading.value = true
    error.value = ''
    books.value.forEach((book) => book.coverUrl && URL.revokeObjectURL(book.coverUrl))
    books.value = []

    try {
      const spaces = (spacesStore.spaces as SpaceResource[]).filter((space) => space.id)
      if (!spaces.length) throw new Error($gettext('No file spaces are available'))
      const spaceById = new Map(spaces.map((space) => [space.id, space]))
      const discovered = await discoverBooks(spaces)
      const resources = [
        ...new Map(
          discovered
            .filter((resource) => spaceById.has(resource.spaceId || resource.storageId || ''))
            .map((resource) => [resource.id, resource])
        ).values()
      ]

      books.value = resources.map((resource) => {
        const preferences = getBookPreferences(resource)
        const readerProgress = getReaderProgress(resource)
        return {
          id: resource.fileId || resource.id,
          resource,
          space: spaceById.get(resource.spaceId || resource.storageId || '')!,
          title: titleFromFileName(resource.name),
          authors: [],
          description: '',
          language: '',
          publisher: '',
          published: '',
          subjects: [],
          spineItemCount: 0,
          loadingMetadata: true,
          hasReadingPosition: readerProgress.hasPosition,
          readingProgress: readerProgress.percentage,
          ...preferences
        }
      })

      await hydrateBooks(books.value)
    } catch (cause) {
      error.value =
        cause instanceof Error ? cause.message : $gettext('The library could not be loaded')
    } finally {
      loading.value = false
    }
  }

  function openBook(book: LibraryBook): void {
    const driveAliasAndItem = book.space.getDriveAliasAndItem(book.resource)
    router.push(`/epub-reader/read/${encodePath(driveAliasAndItem)}`)
  }

  function showBookInFiles(book: LibraryBook): void {
    router.push(getParentFolderLink(book.resource))
  }

  function downloadBook(book: LibraryBook): Promise<void> {
    return downloadFile(book.space, book.resource)
  }

  async function copyBookLink(book: LibraryBook): Promise<boolean> {
    const privateLink = book.resource.privateLink
    if (!privateLink) {
      showErrorMessage({ title: $gettext('The link could not be copied') })
      return false
    }
    try {
      await copyToClipboard(privateLink)
      return true
    } catch (cause) {
      console.error(cause)
      showErrorMessage({ title: $gettext('The link could not be copied'), errors: [cause] })
      return false
    }
  }

  function persistBook(book: LibraryBook): void {
    saveBookPreferences(book.resource, {
      favorite: book.favorite,
      readingStatus: book.readingStatus,
      shelfIds: book.shelfIds
    })
  }

  function toggleFavorite(book: LibraryBook): void {
    book.favorite = !book.favorite
    persistBook(book)
  }

  function setReadingStatus(book: LibraryBook, status: ReadingStatus): void {
    book.readingStatus = status
    persistBook(book)
  }

  function toggleBookShelf(book: LibraryBook, shelfId: string): void {
    book.shelfIds = book.shelfIds.includes(shelfId)
      ? book.shelfIds.filter((id) => id !== shelfId)
      : [...book.shelfIds, shelfId]
    persistBook(book)
  }

  function createShelf(name: string): LibraryShelf | null {
    const trimmedName = name.trim()
    if (!trimmedName) return null
    const existing = shelves.value.find(
      (shelf) => shelf.name.toLocaleLowerCase() === trimmedName.toLocaleLowerCase()
    )
    if (existing) return existing
    const shelf = {
      id: globalThis.crypto?.randomUUID?.() || `${Date.now()}-${Math.random()}`,
      name: trimmedName
    }
    shelves.value = [...shelves.value, shelf]
    saveShelves(shelves.value)
    return shelf
  }

  onUnmounted(() => {
    books.value.forEach((book) => book.coverUrl && URL.revokeObjectURL(book.coverUrl))
  })

  return {
    books,
    copyBookLink,
    downloadBook,
    error,
    loadBooks,
    loading,
    openBook,
    query,
    createShelf,
    shelves,
    showBookInFiles,
    sort,
    setReadingStatus,
    toggleBookShelf,
    toggleFavorite,
    visibleBooks
  }
}
