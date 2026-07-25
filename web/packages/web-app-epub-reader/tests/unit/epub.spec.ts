import * as zip from '@zip.js/zip.js'
import { describe, expect, it } from 'vitest'
import {
  extractEpubMetadata,
  metadataFromPackageDocument,
  titleFromFileName
} from '../../src/utils/epub'

describe('EPUB metadata utilities', () => {
  it('extracts Dublin Core metadata and cover information', () => {
    const metadata = metadataFromPackageDocument(
      `<?xml version="1.0"?>
      <package xmlns="http://www.idpf.org/2007/opf" xmlns:dc="http://purl.org/dc/elements/1.1/">
        <metadata>
          <dc:title>A Wizard of Earthsea</dc:title>
          <dc:creator>Ursula K. Le Guin</dc:creator>
          <dc:language>en</dc:language>
          <dc:publisher>Parnassus Press</dc:publisher>
          <dc:subject>Fantasy</dc:subject>
          <meta name="cover" content="cover-image" />
        </metadata>
        <manifest><item id="cover-image" href="images/cover.jpg" media-type="image/jpeg" /></manifest>
        <spine><itemref idref="chapter-1" /><itemref idref="chapter-2" /></spine>
      </package>`,
      'Fallback title'
    )

    expect(metadata).toMatchObject({
      title: 'A Wizard of Earthsea',
      authors: ['Ursula K. Le Guin'],
      language: 'en',
      publisher: 'Parnassus Press',
      subjects: ['Fantasy'],
      spineItemCount: 2,
      coverHref: 'images/cover.jpg'
    })
  })

  it('creates a readable title from a filename', () => {
    expect(titleFromFileName('the_left-hand.of-earthsea.epub')).toBe('the left hand of earthsea')
  })

  it('reads metadata from an EPUB archive without a web worker', async () => {
    zip.configure({ useWebWorkers: false })
    const writer = new zip.ZipWriter(new zip.BlobWriter('application/epub+zip'))
    await writer.add(
      'META-INF/container.xml',
      new zip.TextReader(
        '<container><rootfiles><rootfile full-path="OEBPS/content.opf" /></rootfiles></container>'
      )
    )
    await writer.add(
      'OEBPS/content.opf',
      new zip.TextReader(
        '<package><metadata><title>Frankenstein</title><creator>Mary Shelley</creator></metadata><manifest /></package>'
      )
    )
    const epub = await writer.close()

    await expect(extractEpubMetadata(epub, 'Fallback')).resolves.toMatchObject({
      title: 'Frankenstein',
      authors: ['Mary Shelley']
    })
  })
})
