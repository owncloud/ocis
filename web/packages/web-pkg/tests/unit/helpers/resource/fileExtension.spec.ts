import { Resource } from '@ownclouders/web-client'
import { isResourceMatchingFileExtension, isResourceMatchingFileType } from '../../../../src'

const buildResource = (data: Partial<Resource>) =>
  ({ id: '1', path: '', webDavPath: '', spaceId: '1', ...data }) as Resource

describe('isResourceMatchingFileExtension', () => {
  it('evaluates to false if resource or file extension is nullish', () => {
    expect(isResourceMatchingFileExtension(null, { extension: 'pdf' })).toBe(false)
    expect(isResourceMatchingFileExtension(buildResource({ extension: 'pdf' }), null)).toBe(false)
  })

  it.each([
    { resourceExtension: 'pdf', extension: 'pdf', expected: true },
    { resourceExtension: 'PDF', extension: 'pdf', expected: true },
    { resourceExtension: 'html', extension: 'pdf', expected: false },
    { resourceExtension: 'svg', extension: 'pdf', expected: false }
  ])(
    'matches by extension: $resourceExtension against $extension',
    ({ resourceExtension, extension, expected }) => {
      expect(
        isResourceMatchingFileExtension(buildResource({ extension: resourceExtension }), {
          extension
        })
      ).toBe(expected)
    }
  )

  it('uses the extension and ignores the mime type if both have an extension', () => {
    expect(
      isResourceMatchingFileExtension(
        buildResource({ extension: 'html', mimeType: 'application/pdf' }),
        { extension: 'pdf', mimeType: 'application/pdf' }
      )
    ).toBe(false)
  })

  it.each([
    { resourceMimeType: 'application/pdf', mimeType: 'application/pdf', expected: true },
    { resourceMimeType: 'image/png', mimeType: 'image', expected: true },
    { resourceMimeType: 'text/html', mimeType: 'application/pdf', expected: false }
  ])(
    'matches by mime type: $resourceMimeType against $mimeType',
    ({ resourceMimeType, mimeType, expected }) => {
      expect(
        isResourceMatchingFileExtension(buildResource({ mimeType: resourceMimeType }), {
          mimeType
        })
      ).toBe(expected)
    }
  )

  it('evaluates to false if nothing can be compared', () => {
    expect(isResourceMatchingFileExtension(buildResource({}), { extension: 'pdf' })).toBe(false)
  })
})

describe('isResourceMatchingFileType', () => {
  const pdfFileExtension = { extension: 'pdf', mimeType: 'application/pdf' }

  it('evaluates to false if resource or file extension is nullish', () => {
    expect(isResourceMatchingFileType(null, pdfFileExtension)).toBe(false)
    expect(isResourceMatchingFileType(buildResource({ mimeType: 'application/pdf' }), null)).toBe(
      false
    )
  })

  it.each([
    { mimeType: 'application/pdf', expected: true },
    { mimeType: 'APPLICATION/PDF', expected: true },
    { mimeType: 'application/pdf; charset=UTF-8', expected: true },
    { mimeType: 'text/html', expected: false },
    { mimeType: 'image/svg+xml', expected: false },
    { mimeType: undefined, expected: false }
  ])('matches by mime type: $mimeType', ({ mimeType, expected }) => {
    expect(
      isResourceMatchingFileType(buildResource({ extension: 'pdf', mimeType }), pdfFileExtension)
    ).toBe(expected)
  })

  it('evaluates to false if the name has a matching extension but the mime type differs', () => {
    // e.g. a mounted share named "report.pdf" whose underlying file was renamed to "report.html"
    expect(
      isResourceMatchingFileType(
        buildResource({ name: 'report.pdf', extension: 'pdf', mimeType: 'text/html' }),
        pdfFileExtension
      )
    ).toBe(false)
  })

  it('evaluates to false if the file extension has no mime type registered', () => {
    expect(
      isResourceMatchingFileType(buildResource({ extension: 'pdf', mimeType: 'application/pdf' }), {
        extension: 'pdf'
      })
    ).toBe(false)
  })
})
