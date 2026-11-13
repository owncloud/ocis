import { mock } from 'vitest-mock-extended'
import { useUpload } from '../../../../src/composables/upload'
import { getComposableWrapper } from '../../../../src/testing'
import { OcUppyFile, UppyService } from '../../../../src/services/uppy/uppyService'
import { ClientService } from '../../../../src/services'

describe('useUpload', () => {
  it('should be valid', () => {
    expect(useUpload).toBeDefined()
  })

  describe('x-oc-mtime header', () => {
    const getTusOptions = (uppyService: UppyService) => {
      return (uppyService.useTus as ReturnType<typeof vi.fn>).mock.calls[0][0]
    }
    const getXhrOptions = (uppyService: UppyService) => {
      return (uppyService.useXhr as ReturnType<typeof vi.fn>).mock.calls[0][0]
    }

    const mountWithTusSupport = (tusSupported: boolean) => {
      const uppyService = mock<UppyService>()
      const clientService = mock<ClientService>()
      getComposableWrapper(
        () => {
          useUpload({ uppyService })
        },
        {
          provide: { $clientService: clientService },
          pluginOptions: {
            piniaOptions: {
              capabilityState: {
                capabilities: {
                  files: {
                    tus_support: {
                      max_chunk_size: tusSupported ? 1000 : 0,
                      extension: 'creation-with-upload',
                      http_method_override: false
                    }
                  }
                }
              }
            }
          }
        }
      )
      return uppyService
    }

    describe('tus onBeforeRequest (remote/companion uploads)', () => {
      const runOnBeforeRequest = async (file: OcUppyFile) => {
        const uppyService = mountWithTusSupport(true)
        const tusOptions = getTusOptions(uppyService)
        const req = { setHeader: vi.fn() }
        await tusOptions.onBeforeRequest(req, file)
        return req
      }

      it('omits the header for a Blob-backed remote file', async () => {
        const file = mock<OcUppyFile>({ isRemote: true, data: new Blob() })
        const req = await runOnBeforeRequest(file)
        expect(req.setHeader).not.toHaveBeenCalledWith('x-oc-mtime', expect.anything())
      })

      it('omits the header for a remote file with lastModified 0', async () => {
        const file = mock<OcUppyFile>({ isRemote: true, data: { lastModified: 0 } as File })
        const req = await runOnBeforeRequest(file)
        expect(req.setHeader).not.toHaveBeenCalledWith('x-oc-mtime', expect.anything())
      })

      it('sets the header to the rounded integer-seconds string for a valid remote file', async () => {
        const lastModified = 1700000000000
        const file = mock<OcUppyFile>({ isRemote: true, data: { lastModified } as File })
        const req = await runOnBeforeRequest(file)
        expect(req.setHeader).toHaveBeenCalledWith('x-oc-mtime', (lastModified / 1000).toFixed(0))
      })
    })

    describe('tus headers for xhr/remote files', () => {
      it('omits the header for a Blob-backed xhrUpload file', () => {
        const uppyService = mountWithTusSupport(true)
        const tusOptions = getTusOptions(uppyService)
        const file = mock<OcUppyFile>({ xhrUpload: {}, data: new Blob() })
        const headers = tusOptions.headers(file)
        expect(headers).not.toHaveProperty('x-oc-mtime')
      })

      it('omits the header for an xhrUpload file with lastModified 0', () => {
        const uppyService = mountWithTusSupport(true)
        const tusOptions = getTusOptions(uppyService)
        const file = mock<OcUppyFile>({ xhrUpload: {}, data: { lastModified: 0 } as File })
        const headers = tusOptions.headers(file)
        expect(headers).not.toHaveProperty('x-oc-mtime')
      })

      it('sets the header to the rounded integer-seconds string for a valid xhrUpload file', () => {
        const uppyService = mountWithTusSupport(true)
        const tusOptions = getTusOptions(uppyService)
        const lastModified = 1700000000000
        const file = mock<OcUppyFile>({ xhrUpload: {}, data: { lastModified } as File })
        const headers = tusOptions.headers(file)
        expect(headers['x-oc-mtime']).toEqual((lastModified / 1000).toFixed(0))
      })
    })

    describe('xhrOptions.headers', () => {
      it('omits the header for a Blob-backed file', () => {
        const uppyService = mountWithTusSupport(false)
        const xhrOptions = getXhrOptions(uppyService)
        const file = mock<OcUppyFile>({ data: new Blob() })
        const headers = xhrOptions.headers(file)
        expect(headers).not.toHaveProperty('x-oc-mtime')
      })

      it('omits the header for a file with lastModified 0', () => {
        const uppyService = mountWithTusSupport(false)
        const xhrOptions = getXhrOptions(uppyService)
        const file = mock<OcUppyFile>({ data: { lastModified: 0 } as File })
        const headers = xhrOptions.headers(file)
        expect(headers).not.toHaveProperty('x-oc-mtime')
      })

      it('sets the header to the rounded integer-seconds string for a valid file', () => {
        const uppyService = mountWithTusSupport(false)
        const xhrOptions = getXhrOptions(uppyService)
        const lastModified = 1700000000000
        const file = mock<OcUppyFile>({ data: { lastModified } as File })
        const headers = xhrOptions.headers(file)
        expect(headers['x-oc-mtime']).toEqual((lastModified / 1000).toFixed(0))
      })
    })
  })
})
