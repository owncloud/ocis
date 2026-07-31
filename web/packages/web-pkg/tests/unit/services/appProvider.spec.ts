import { mockDeep } from 'vitest-mock-extended'
import { AppProviderService, ClientService } from '../../../src/services'
import { MimeType } from '../../../src/services/appProvider/schemas'

const docx = 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'

const buildService = (mimeTypes: MimeType[]) => {
  const service = new AppProviderService('https://example.test', mockDeep<ClientService>())
  service.mimeTypes = mimeTypes
  return service
}

describe('AppProviderService', () => {
  describe('getDefaultAppNameForMimeType', () => {
    it('returns the configured default_application when it is offered for the mime type', () => {
      const service = buildService([
        {
          mime_type: docx,
          default_application: 'ByCS-Office',
          app_providers: [
            { name: 'Webeditor', icon: '', secure_view: false },
            { name: 'ByCS-Office', icon: '', secure_view: false }
          ]
        }
      ])
      expect(service.getDefaultAppNameForMimeType(docx)).toEqual('ByCS-Office')
    })

    it('falls back to the first provider when no default_application is configured', () => {
      const service = buildService([
        {
          mime_type: docx,
          app_providers: [
            { name: 'ByCS-Office', icon: '', secure_view: false },
            { name: 'Webeditor', icon: '', secure_view: false }
          ]
        }
      ])
      expect(service.getDefaultAppNameForMimeType(docx)).toEqual('ByCS-Office')
    })

    it('ignores a default_application that is not among the registered providers', () => {
      const service = buildService([
        {
          mime_type: docx,
          default_application: 'GoneApp',
          app_providers: [{ name: 'ByCS-Office', icon: '', secure_view: false }]
        }
      ])
      expect(service.getDefaultAppNameForMimeType(docx)).toEqual('ByCS-Office')
    })

    it('returns undefined for a mime type that no provider handles', () => {
      const service = buildService([
        { mime_type: docx, app_providers: [{ name: 'ByCS-Office', icon: '', secure_view: false }] }
      ])
      expect(service.getDefaultAppNameForMimeType('application/x-unknown')).toBeUndefined()
    })

    it('returns undefined when the mime type entry has no providers', () => {
      const service = buildService([{ mime_type: docx, app_providers: [] }])
      expect(service.getDefaultAppNameForMimeType(docx)).toBeUndefined()
    })
  })
})
