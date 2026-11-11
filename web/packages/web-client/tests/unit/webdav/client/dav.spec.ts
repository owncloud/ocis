import { DAV } from '../../../../src/webdav/client/dav'
import { DavHttpError } from '../../../../src/errors'

const { customRequest } = vi.hoisted(() => ({ customRequest: vi.fn() }))

vi.mock('webdav', () => ({
  createClient: () => ({ customRequest })
}))

const createDav = () => {
  const onSetMaintenance = vi.fn()
  const dav = new DAV({ baseUrl: 'https://example.com/', onSetMaintenance })
  return { dav, onSetMaintenance }
}

describe('DAV', () => {
  beforeEach(() => {
    customRequest.mockReset()
  })

  describe('request errors', () => {
    it('passes on an error without an HTTP response unchanged', async () => {
      const networkError = new TypeError('Failed to fetch')
      customRequest.mockRejectedValue(networkError)
      const { dav, onSetMaintenance } = createDav()

      await expect(dav.copy('/source', '/target')).rejects.toBe(networkError)
      expect(onSetMaintenance).not.toHaveBeenCalled()
    })

    it('turns an HTTP error response into a DavHttpError with its status code', async () => {
      customRequest.mockRejectedValue(
        Object.assign(new Error('Bad Gateway'), { response: new Response('', { status: 502 }) })
      )
      const { dav } = createDav()

      const error = await dav.copy('/source', '/target').catch((e) => e)

      expect(error).toBeInstanceOf(DavHttpError)
      expect(error.statusCode).toBe(502)
    })

    it('switches to maintenance mode on a 503 response', async () => {
      customRequest.mockRejectedValue(
        Object.assign(new Error('Service Unavailable'), {
          response: new Response('', { status: 503 })
        })
      )
      const { dav, onSetMaintenance } = createDav()

      await expect(dav.copy('/source', '/target')).rejects.toBeInstanceOf(DavHttpError)
      expect(onSetMaintenance).toHaveBeenCalledWith(true)
    })
  })
})
