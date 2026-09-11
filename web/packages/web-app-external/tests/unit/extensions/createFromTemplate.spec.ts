import { mock } from 'vitest-mock-extended'
import { flushPromises, getComposableWrapper } from '@ownclouders/web-test-helpers'
import {
  ActionExtension,
  ApplicationInformation,
  AppProviderService,
  ClientService,
  contextRouteNameKey,
  contextRouteParamsKey,
  contextRouteQueryKey,
  FileActionOptions,
  resolveFileNameDuplicate,
  useAppProviderService,
  useClientService,
  useFileActions,
  useMessages,
  useRouter,
  useSpacesStore
} from '@ownclouders/web-pkg'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { Router } from 'vue-router'
import { useActionExtensionCreateFromTemplate } from '../../../src/extensions/createFromTemplate'
import { useCreateFileHandler } from '../../../src/composables'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useAppProviderService: vi.fn(),
  useSpacesStore: vi.fn(),
  useClientService: vi.fn(),
  useFileActions: vi.fn(),
  useMessages: vi.fn(),
  useRouter: vi.fn(),
  resolveFileNameDuplicate: vi.fn()
}))

vi.mock('../../../src/composables', () => ({
  useCreateFileHandler: vi.fn()
}))

const appName = 'Collabora'
const templateMimeType = 'application/vnd.oasis.opendocument.text'

const getAppInfo = () => mock<ApplicationInformation>({ name: appName })

const getTemplateResource = (canDownload = true) => {
  const template = mock<Resource>({
    name: 'template.odt',
    extension: 'odt',
    mimeType: templateMimeType,
    fileId: 'template-file-id'
  })
  template.canDownload.mockReturnValue(canDownload)
  return template
}

const getPersonalSpace = () =>
  mock<SpaceResource>({
    fileId: 'personal-space-id',
    driveAlias: 'personal/admin'
  })

const matchingTemplateMimeTypes = [
  {
    mime_type: templateMimeType,
    app_providers: [{ name: appName, target_ext: 'odt' }]
  }
]

type GetWrapperOptions = {
  hasPersonalSpace?: boolean
  templateMimeTypes?: unknown[]
  existingResources?: Resource[]
  listFilesRejects?: boolean
  createFileHandlerRejects?: boolean
}

const getWrapper = ({
  hasPersonalSpace = true,
  templateMimeTypes = matchingTemplateMimeTypes,
  existingResources = [],
  listFilesRejects = false,
  createFileHandlerRejects = false
}: GetWrapperOptions = {}) => {
  const personalSpace = hasPersonalSpace ? getPersonalSpace() : undefined
  const showErrorMessage = vi.fn()
  const push = vi.fn().mockResolvedValue(undefined)
  const getEditorRouteOpts = vi.fn().mockReturnValue({ query: {} })
  const listFiles = vi.fn()
  const createFileHandler = vi.fn()
  const personalSpaceRoot = mock<Resource>({ path: '/' })

  if (listFilesRejects) {
    listFiles.mockRejectedValue(new Error('listFiles failed'))
  } else {
    listFiles.mockResolvedValue({ resource: personalSpaceRoot, children: existingResources })
  }

  if (createFileHandlerRejects) {
    createFileHandler.mockRejectedValue(new Error('createFileHandler failed'))
  } else {
    createFileHandler.mockResolvedValue(mock<Resource>({ name: 'created.odt' }))
  }

  vi.mocked(useAppProviderService).mockReturnValue(
    mock<AppProviderService>({ templateMimeTypes: templateMimeTypes as any })
  )
  vi.mocked(useSpacesStore).mockReturnValue({ personalSpace } as any)
  vi.mocked(useClientService).mockReturnValue(mock<ClientService>({ webdav: { listFiles } } as any))
  vi.mocked(useFileActions).mockReturnValue({ getEditorRouteOpts } as any)
  vi.mocked(useMessages).mockReturnValue({ showErrorMessage } as any)
  vi.mocked(useRouter).mockReturnValue(mock<Router>({ push }))
  vi.mocked(useCreateFileHandler).mockReturnValue({ createFileHandler })

  let extension: ActionExtension
  getComposableWrapper(() => {
    extension = useActionExtensionCreateFromTemplate(getAppInfo())
  })

  return {
    action: extension.action,
    mocks: { showErrorMessage, push, getEditorRouteOpts, listFiles, createFileHandler }
  }
}

describe('createFromTemplate action extension', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    vi.mocked(resolveFileNameDuplicate).mockReset()
  })

  describe('isVisible', () => {
    it('is false when more than one resource is selected', () => {
      const { action } = getWrapper()
      expect(
        action.isVisible({ resources: [getTemplateResource(), getTemplateResource()] } as any)
      ).toBe(false)
    })

    it('is false when no personal space is present', () => {
      const { action } = getWrapper({ hasPersonalSpace: false })
      expect(action.isVisible({ resources: [getTemplateResource()] } as any)).toBe(false)
    })

    it('is false when the resource is not downloadable', () => {
      const { action } = getWrapper()
      expect(action.isVisible({ resources: [getTemplateResource(false)] } as any)).toBe(false)
    })

    it('is false when no template mime type matches the app', () => {
      const { action } = getWrapper({ templateMimeTypes: [] })
      expect(action.isVisible({ resources: [getTemplateResource()] } as any)).toBe(false)
    })

    it('is true for a single matching template resource', () => {
      const { action } = getWrapper()
      expect(action.isVisible({ resources: [getTemplateResource()] } as any)).toBe(true)
    })
  })

  describe('handler', () => {
    const callHandler = async (
      action: ActionExtension['action'],
      template = getTemplateResource()
    ) => {
      await action.handler({ resources: [template] } as unknown as FileActionOptions)
      await flushPromises()
    }

    it('creates the file and navigates to the editor route with context params', async () => {
      const { action, mocks } = getWrapper()
      await callHandler(action)

      expect(mocks.createFileHandler).toHaveBeenCalledWith(
        expect.objectContaining({ fileName: 'template.odt' })
      )
      expect(mocks.getEditorRouteOpts).toHaveBeenCalledWith(
        `external-${appName.toLowerCase()}-apps`,
        expect.anything(),
        expect.anything(),
        expect.anything(),
        undefined,
        'template-file-id'
      )
      expect(mocks.push).toHaveBeenCalledTimes(1)
      const pushedRoute = mocks.push.mock.calls[0][0]
      expect(pushedRoute.query).toEqual(
        expect.objectContaining({
          [contextRouteNameKey]: expect.anything(),
          [contextRouteParamsKey]: expect.anything(),
          [contextRouteQueryKey]: expect.anything()
        })
      )
      expect(mocks.showErrorMessage).not.toHaveBeenCalled()
    })

    it('resolves a duplicate file name when a file of the same name already exists', async () => {
      vi.mocked(resolveFileNameDuplicate).mockReturnValue('template (1).odt')
      const { action, mocks } = getWrapper({
        existingResources: [mock<Resource>({ name: 'template.odt' })]
      })
      await callHandler(action)

      expect(resolveFileNameDuplicate).toHaveBeenCalledWith(
        'template.odt',
        'odt',
        expect.anything()
      )
      expect(mocks.createFileHandler).toHaveBeenCalledWith(
        expect.objectContaining({ fileName: 'template (1).odt' })
      )
    })

    it('shows an error message when listing existing files fails', async () => {
      const { action, mocks } = getWrapper({ listFilesRejects: true })
      await callHandler(action)

      expect(mocks.createFileHandler).not.toHaveBeenCalled()
      expect(mocks.showErrorMessage).toHaveBeenCalledWith(
        expect.objectContaining({ title: expect.any(String) })
      )
      expect(mocks.push).not.toHaveBeenCalled()
    })

    it('shows an error message when creating the file fails', async () => {
      const { action, mocks } = getWrapper({ createFileHandlerRejects: true })
      await callHandler(action)

      expect(mocks.showErrorMessage).toHaveBeenCalledWith(
        expect.objectContaining({ title: expect.any(String) })
      )
      expect(mocks.push).not.toHaveBeenCalled()
    })
  })
})
