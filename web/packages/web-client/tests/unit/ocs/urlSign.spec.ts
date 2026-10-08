import { ocs } from '../../../src/ocs'
import { FetchClient } from '../../../src/http'

const signingKeyResponse =
  '<?xml version="1.0" encoding="UTF-8"?><ocs><data><signing-key>key</signing-key></data></ocs>'

describe('ocs signUrl', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi
      .fn()
      .mockImplementation(() => Promise.resolve(new Response(signingKeyResponse, { status: 200 })))
    vi.stubGlobal('fetch', fetchMock)
  })

  it('requests the signing key from the ocs endpoint', async () => {
    await ocs('https://host/', new FetchClient()).signUrl({
      url: 'https://host/dav/file',
      username: 'admin'
    })

    expect(fetchMock.mock.calls[0][0]).toBe('https://host/ocs/v2.php/cloud/user/signing-key')
  })

  it('keeps the path out of the query when the base URI has one', async () => {
    await ocs('https://host/?vault=true', new FetchClient()).signUrl({
      url: 'https://host/dav/file',
      username: 'admin'
    })

    expect(fetchMock.mock.calls[0][0]).toBe(
      'https://host/ocs/v2.php/cloud/user/signing-key?vault=true'
    )
  })

  it('signs the url', async () => {
    const signedUrl = new URL(
      await ocs('https://host/?vault=true', new FetchClient()).signUrl({
        url: 'https://host/dav/file',
        username: 'admin'
      })
    )

    expect(signedUrl.searchParams.get('OC-Credential')).toBe('admin')
    expect(signedUrl.searchParams.get('OC-Signature')).toBeTruthy()
  })
})
