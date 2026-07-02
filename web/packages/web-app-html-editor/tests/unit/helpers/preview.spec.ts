import {
  PREVIEW_CSP,
  PREVIEW_SIZE_LIMIT,
  isPreviewTooLarge,
  wrapWithPreviewCsp
} from '../../../src/helpers/preview'

const META = `<meta http-equiv="Content-Security-Policy" content="${PREVIEW_CSP}">`

describe('preview CSP', () => {
  it('denies network egress and form submission by default', () => {
    expect(PREVIEW_CSP).toContain("default-src 'none'")
    expect(PREVIEW_CSP).toContain("form-action 'none'")
    expect(PREVIEW_CSP).toContain("base-uri 'none'")
    // no wildcard / external network sources
    expect(PREVIEW_CSP).not.toContain('http:')
    expect(PREVIEW_CSP).not.toContain('https:')
    expect(PREVIEW_CSP).not.toContain('*')
  })

  it('prepends the CSP meta as the very first bytes of a full document', () => {
    const out = wrapWithPreviewCsp(
      '<!doctype html><html><head><title>t</title></head><body>x</body></html>'
    )
    expect(out.startsWith(META)).toBe(true)
    // the original document follows the meta unchanged
    expect(out).toBe(
      `${META}<!doctype html><html><head><title>t</title></head><body>x</body></html>`
    )
  })

  it('prepends the CSP meta for a bare fragment', () => {
    expect(wrapWithPreviewCsp('<p>x</p>')).toBe(`${META}<p>x</p>`)
  })

  it('returns just the meta for empty or nullish content', () => {
    expect(wrapWithPreviewCsp('')).toBe(META)
    expect(wrapWithPreviewCsp(undefined as unknown as string)).toBe(META)
  })

  it('governs a resource placed before any <head> (regression: first meta wins)', () => {
    // An attacker can position an exfiltration element before the document's head
    // so that an "inject into the real head" strategy would leave it uncovered.
    // Because the meta is prepended, the CSP still precedes — and governs — it.
    const out = wrapWithPreviewCsp('<img src="https://evil.example/beacon"><head></head>')
    expect(out.startsWith(META)).toBe(true)
    expect(out.indexOf(META)).toBeLessThan(out.indexOf('<img'))
  })

  it('is not diverted by a commented-out <head> decoy', () => {
    const out = wrapWithPreviewCsp(
      '<!-- <head> --><html><head><title>t</title></head><body>x</body></html>'
    )
    // exactly one CSP meta, and it leads the whole document
    expect(out.startsWith(META)).toBe(true)
    expect(out.split(META).length - 1).toBe(1)
  })
})

describe('preview size guard', () => {
  it('flags content above the limit', () => {
    expect(isPreviewTooLarge('a'.repeat(PREVIEW_SIZE_LIMIT + 1))).toBe(true)
  })

  it('allows content at or below the limit', () => {
    expect(isPreviewTooLarge('a'.repeat(PREVIEW_SIZE_LIMIT))).toBe(false)
    expect(isPreviewTooLarge('')).toBe(false)
  })
})
