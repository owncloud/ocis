# EPUB Reader and Library

The ownCloud Web EPUB app combines the built-in EPUB reader with a visual library for books
stored across accessible personal and project spaces.

## Library features

- Discover EPUB files across accessible spaces
- Extract embedded metadata and cover images in the browser
- Search by title, author, publisher, or subject
- Sort and filter by collection, reading status, author, subject, language, and space
- Organize books with favorites and custom shelves
- Switch between grid and list layouts
- Continue from the location saved by the EPUB reader
- Open the containing folder, download a book, or copy its private link
- Cache metadata and covers in IndexedDB

EPUB contents and metadata are not sent to an external metadata service. Favorites, shelves,
reading-status labels, and the selected layout are stored in browser local storage and do not
currently synchronize between browsers or devices.

## Usage

1. Upload EPUB files through the Files app.
2. Open **Library** from the application switcher.
3. Select a cover to open the book, or use its information button for metadata and file actions.
4. Use search, sorting, and filters to organize the library.

The library combines WebDAV search with a recursive scan so newly uploaded books can appear
before search indexing completes. Metadata cache entries are invalidated when file metadata
changes.

## Development

Run commands from the `web` directory:

```bash
pnpm --filter epub-reader test:unit -- --run
pnpm check:types
pnpm eslint 'packages/web-app-epub-reader/**/*.{ts,vue}' --max-warnings=0
```
