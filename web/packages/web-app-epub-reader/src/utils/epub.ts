import * as zip from '@zip.js/zip.js'
import type { EpubMetadata } from '../types'

const textDecoder = new TextDecoder()

function normalizedPath(path: string): string {
  const parts: string[] = []
  for (const segment of path.split('/')) {
    if (!segment || segment === '.') continue
    if (segment === '..') parts.pop()
    else parts.push(segment)
  }
  return parts.join('/')
}

function resolvePath(baseFile: string, relativePath: string): string {
  const base = baseFile.includes('/') ? baseFile.slice(0, baseFile.lastIndexOf('/') + 1) : ''
  return normalizedPath(`${base}${relativePath}`)
}

function elementsByLocalName(root: ParentNode, localName: string): Element[] {
  return [...root.querySelectorAll('*')].filter(
    (element) => element.localName === localName || element.tagName.split(':').pop() === localName
  )
}

function firstText(document: Document, localName: string): string {
  return elementsByLocalName(document, localName)[0]?.textContent?.trim() ?? ''
}

function allText(document: Document, localName: string): string[] {
  return elementsByLocalName(document, localName)
    .map((element) => element.textContent?.trim() ?? '')
    .filter(Boolean)
}

export function metadataFromPackageDocument(
  packageXml: string,
  fallbackTitle: string
): EpubMetadata & { coverHref?: string; coverMediaType?: string } {
  const document = new DOMParser().parseFromString(packageXml, 'application/xml')
  if (document.querySelector('parsererror')) {
    throw new Error('Invalid EPUB package document')
  }

  const metadataElements = elementsByLocalName(document, 'meta')
  const manifestItems = elementsByLocalName(document, 'item')
  const coverImageItem = manifestItems.find((element) =>
    element.getAttribute('properties')?.split(/\s+/).includes('cover-image')
  )
  const coverId =
    metadataElements
      .find((element) => element.getAttribute('name') === 'cover')
      ?.getAttribute('content') ??
    coverImageItem?.getAttribute('id') ??
    ''
  const coverItem =
    manifestItems.find((element) => element.getAttribute('id') === coverId) ?? coverImageItem

  return {
    title: firstText(document, 'title') || fallbackTitle,
    authors: allText(document, 'creator'),
    description: firstText(document, 'description'),
    language: firstText(document, 'language'),
    publisher: firstText(document, 'publisher'),
    published: firstText(document, 'date'),
    subjects: allText(document, 'subject'),
    spineItemCount: elementsByLocalName(document, 'itemref').length,
    coverHref: coverItem?.getAttribute('href') ?? undefined,
    coverMediaType: coverItem?.getAttribute('media-type') ?? undefined
  }
}

async function entryText(entry: zip.FileEntry): Promise<string> {
  const data = await entry.getData(new zip.Uint8ArrayWriter())
  return textDecoder.decode(data)
}

export async function extractEpubMetadata(
  data: Blob,
  fallbackTitle: string
): Promise<EpubMetadata> {
  // Extensions are loaded as standalone bundles from oCIS. ZIP.js cannot
  // reliably resolve its worker script from that runtime location, which can
  // leave metadata extraction pending indefinitely. EPUB metadata is small
  // enough to process on the main thread.
  const reader = new zip.ZipReader(new zip.BlobReader(data), {
    useWebWorkers: false
  })
  try {
    const entries = (await reader.getEntries()).filter(
      (entry): entry is zip.FileEntry => !entry.directory
    )
    const byName = new Map(entries.map((entry) => [normalizedPath(entry.filename), entry]))
    const containerEntry = byName.get('META-INF/container.xml')
    if (!containerEntry) throw new Error('EPUB container.xml is missing')

    const container = new DOMParser().parseFromString(
      await entryText(containerEntry),
      'application/xml'
    )
    const packagePath = container.querySelector('rootfile')?.getAttribute('full-path')
    if (!packagePath) throw new Error('EPUB package path is missing')

    const normalizedPackagePath = normalizedPath(packagePath)
    const packageEntry = byName.get(normalizedPackagePath)
    if (!packageEntry) throw new Error('EPUB package document is missing')

    const parsed = metadataFromPackageDocument(await entryText(packageEntry), fallbackTitle)
    let coverUrl: string | undefined
    let coverBlob: Blob | undefined

    if (parsed.coverHref) {
      const coverEntry = byName.get(resolvePath(normalizedPackagePath, parsed.coverHref))
      if (coverEntry) {
        coverBlob = await coverEntry.getData(new zip.BlobWriter(parsed.coverMediaType))
        coverUrl = URL.createObjectURL(coverBlob)
      }
    }

    return {
      title: parsed.title,
      authors: parsed.authors,
      description: parsed.description,
      language: parsed.language,
      publisher: parsed.publisher,
      published: parsed.published,
      subjects: parsed.subjects,
      spineItemCount: parsed.spineItemCount,
      coverUrl,
      coverBlob
    }
  } finally {
    await reader.close()
  }
}

export function titleFromFileName(name: string): string {
  return name
    .replace(/\.epub$/i, '')
    .replace(/[._-]+/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}
