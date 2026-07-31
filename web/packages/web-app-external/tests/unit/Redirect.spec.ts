import { mock } from 'vitest-mock-extended'
import { ref } from 'vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  flushPromises,
  shallowMount
} from '@ownclouders/web-test-helpers'
import {
  AppProviderService,
  queryItemAsString,
  useRouteMeta,
  useRouteQuery
} from '@ownclouders/web-pkg'
import { Resource } from '@ownclouders/web-client'
import Redirect from '../../src/Redirect.vue'
import { useApplicationReadyStore } from '../../src/piniaStores'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRouteQuery: vi.fn(),
  useRouteMeta: vi.fn(),
  queryItemAsString: vi.fn()
}))

const { mockRouterReplace } = vi.hoisted(() => ({ mockRouterReplace: vi.fn() }))

vi.mock('vue-router', async (importOriginal) => {
  const actual = await importOriginal<any>()
  const { ref } = await import('vue')
  return {
    ...actual,
    useRouter: () => ({ replace: mockRouterReplace, currentRoute: ref({ query: {} }) })
  }
})

const docxMimeType = 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'

describe('Redirect.vue', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    mockRouterReplace.mockClear()
  })

  it('does not redirect while the application is not ready', async () => {
    const { wrapper } = getWrapper({ query: { app: 'Collabora' }, ready: false })
    await flushPromises()
    expect(mockRouterReplace).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('You are being redirected.')
  })

  it('redirects to the app given via the "app" query', async () => {
    getWrapper({ query: { app: 'Collabora' } })
    await flushPromises()
    expect(mockRouterReplace).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'external-collabora-apps' })
    )
  })

  it('redirects to the app given via the "appName" query', async () => {
    getWrapper({ query: { appName: 'OnlyOffice' } })
    await flushPromises()
    expect(mockRouterReplace).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'external-onlyoffice-apps' })
    )
  })

  it('resolves the app from the file mime type when no app query is given', async () => {
    const { appProviderService } = getWrapper({
      query: { fileId: 'file-id' },
      mimeType: docxMimeType,
      resolvedApp: 'ByCS-Office'
    })
    await flushPromises()
    expect(appProviderService.getDefaultAppNameForMimeType).toHaveBeenCalledWith(docxMimeType)
    expect(mockRouterReplace).toHaveBeenCalledWith(
      expect.objectContaining({ name: 'external-bycs-office-apps' })
    )
  })

  it('shows an error and does not redirect when no app is configured for the mime type', async () => {
    const { wrapper } = getWrapper({
      query: { fileId: 'file-id' },
      mimeType: 'application/x-unknown',
      resolvedApp: undefined
    })
    await flushPromises()
    expect(mockRouterReplace).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('We could not open this file')
  })

  it('shows an error and does not redirect when the file cannot be statted', async () => {
    const { wrapper } = getWrapper({ query: { fileId: 'file-id' }, statThrows: true })
    await flushPromises()
    expect(mockRouterReplace).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('We could not open this file')
  })

  it('shows an error and does not redirect when neither app nor fileId is given', async () => {
    const { wrapper } = getWrapper({ query: {} })
    await flushPromises()
    expect(mockRouterReplace).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('We could not open this file')
  })
})

function getWrapper({
  query = {},
  mimeType = '',
  resolvedApp = undefined,
  statThrows = false,
  ready = true
}: {
  query?: { app?: string; appName?: string; fileId?: string }
  mimeType?: string
  resolvedApp?: string
  statThrows?: boolean
  ready?: boolean
} = {}) {
  vi.mocked(useRouteQuery).mockImplementation(
    (name: string) => ref((query as Record<string, string>)[name] ?? '') as never
  )
  vi.mocked(useRouteMeta).mockReturnValue(ref('Redirecting to external app') as never)
  vi.mocked(queryItemAsString).mockImplementation((value) => (value ?? '').toString())

  const appProviderService = mock<AppProviderService>()
  appProviderService.getDefaultAppNameForMimeType.mockReturnValue(resolvedApp)

  const mocks = {
    ...defaultComponentMocks(),
    $appProviderService: appProviderService
  }

  if (statThrows) {
    mocks.$clientService.webdav.getFileInfo.mockRejectedValue(new Error('stat failed'))
  } else {
    mocks.$clientService.webdav.getFileInfo.mockResolvedValue(mock<Resource>({ mimeType }))
  }

  const wrapper = shallowMount(Redirect, {
    global: {
      plugins: [...defaultPlugins()],
      provide: mocks,
      mocks
    }
  })

  useApplicationReadyStore().isReady = ready

  return { wrapper, mocks, appProviderService }
}
