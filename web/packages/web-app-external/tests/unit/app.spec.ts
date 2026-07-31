import { mock } from 'vitest-mock-extended'
import {
  defaultComponentMocks,
  defaultPlugins,
  flushPromises,
  shallowMount
} from '@ownclouders/web-test-helpers'
import {
  AppProviderService,
  useMessages,
  useModals,
  useRequest,
  useRoute
} from '@ownclouders/web-pkg'
import { computed } from 'vue'

import { Resource } from '@ownclouders/web-client'
import App from '../../src/App.vue'
import { RouteLocation } from 'vue-router'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRequest: vi.fn(),
  useRoute: vi.fn(),
  useModals: vi.fn(),
  useMessages: vi.fn()
}))

const appUrl = 'https://example.test/d12ab86/loe009157-MzBw'
const appOrigin = 'https://example.test'

const providerSuccessResponsePost = {
  app_url: appUrl,
  method: 'POST',
  form_parameters: {
    access_token: 'asdfsadfsadf',
    access_token_ttl: '123456'
  }
}

const providerSuccessResponseGet = {
  app_url: appUrl,
  method: 'GET'
}

describe('The app provider extension', () => {
  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
  })

  it('should fail for unauthenticated users', async () => {
    const makeRequest = vi.fn().mockResolvedValue({
      ok: true,
      status: 401,
      message: 'Login Required'
    })
    const { wrapper } = createShallowMountWrapper(makeRequest)
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('should be able to load an iFrame via get', async () => {
    const makeRequest = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      data: providerSuccessResponseGet
    })

    const { wrapper } = createShallowMountWrapper(makeRequest)
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('should be able to load an iFrame via post', async () => {
    const makeRequest = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      data: providerSuccessResponsePost
    })
    const { wrapper } = createShallowMountWrapper(makeRequest)
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })
})

describe('Collabora Save As postMessage handling', () => {
  const successGet = { ok: true, status: 200, data: providerSuccessResponseGet }

  const dispatchEditorMessage = (origin: string, messageId: string) => {
    window.dispatchEvent(
      new MessageEvent('message', { origin, data: JSON.stringify({ MessageId: messageId }) })
    )
  }

  const stubIframeContentWindow = (wrapper: any) => {
    const postMessage = vi.fn()
    const iframe = wrapper.find('iframe')
    Object.defineProperty(iframe.element, 'contentWindow', {
      value: { postMessage },
      configurable: true
    })
    return { iframe, postMessage }
  }

  it('ignores editor messages until the editor origin is known (fail closed)', async () => {
    // makeRequest never resolves -> appUrl stays undefined -> origin unknown
    const makeRequest = vi.fn().mockReturnValue(new Promise(() => undefined))
    const { wrapper, dispatchModal } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora'
    })
    await wrapper.vm.$nextTick()

    dispatchEditorMessage(appOrigin, 'UI_SaveAs')

    expect(dispatchModal).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('drops editor messages coming from a foreign origin', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper, dispatchModal } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora'
    })
    await flushPromises()

    dispatchEditorMessage('https://evil.test', 'UI_SaveAs')

    expect(dispatchModal).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('opens the modal on UI_SaveAs from the editor origin and posts Action_SaveAs on confirm', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper, dispatchModal } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora'
    })
    await flushPromises()
    const { postMessage } = stubIframeContentWindow(wrapper)

    dispatchEditorMessage(appOrigin, 'UI_SaveAs')

    expect(dispatchModal).toHaveBeenCalledTimes(1)
    const modal = dispatchModal.mock.calls[0][0] as any
    expect(modal.hasInput).toBe(true)

    modal.onConfirm('export.pdf')

    expect(postMessage).toHaveBeenCalledTimes(1)
    const [payload, targetOrigin] = postMessage.mock.calls[0]
    expect(targetOrigin).toBe(appOrigin)
    expect(JSON.parse(payload as string)).toMatchObject({
      MessageId: 'Action_SaveAs',
      Values: { Filename: 'export.pdf', Notify: true }
    })
    wrapper.unmount()
  })

  it('ignores a second UI_SaveAs while the modal is already open', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper, dispatchModal } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora'
    })
    await flushPromises()
    stubIframeContentWindow(wrapper)

    dispatchEditorMessage(appOrigin, 'UI_SaveAs')
    dispatchEditorMessage(appOrigin, 'UI_SaveAs')

    expect(dispatchModal).toHaveBeenCalledTimes(1)

    // once cancelled/confirmed, a subsequent UI_SaveAs may open the modal again
    const modal = dispatchModal.mock.calls[0][0] as any
    modal.onCancel()
    dispatchEditorMessage(appOrigin, 'UI_SaveAs')
    expect(dispatchModal).toHaveBeenCalledTimes(2)
    wrapper.unmount()
  })

  it('shows an error instead of silently closing when the iframe is unavailable on confirm', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper, dispatchModal, showErrorMessage } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora'
    })
    await flushPromises()
    Object.defineProperty(wrapper.find('iframe').element, 'contentWindow', {
      value: null,
      configurable: true
    })

    dispatchEditorMessage(appOrigin, 'UI_SaveAs')
    const modal = dispatchModal.mock.calls[0][0] as any

    modal.onConfirm('export.pdf')

    expect(showErrorMessage).toHaveBeenCalledWith(
      expect.objectContaining({ title: expect.stringContaining('editor is not available') })
    )
    wrapper.unmount()
  })

  it('does not open the modal on UI_SaveAs when the file is read-only', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper, dispatchModal } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora',
      isReadOnly: true
    })
    await flushPromises()

    dispatchEditorMessage(appOrigin, 'UI_SaveAs')

    expect(dispatchModal).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('rejects empty, path-separator and reserved file names in the modal', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper, dispatchModal } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora'
    })
    await flushPromises()
    const { postMessage } = stubIframeContentWindow(wrapper)

    dispatchEditorMessage(appOrigin, 'UI_SaveAs')
    const modal = dispatchModal.mock.calls[0][0] as any

    const setError = vi.fn()
    modal.onInput('', setError)
    expect(setError).toHaveBeenLastCalledWith(expect.stringContaining('empty'))
    modal.onInput('foo/bar.pdf', setError)
    expect(setError).toHaveBeenLastCalledWith(expect.stringContaining('/'))
    modal.onInput('..', setError)
    expect(setError).toHaveBeenLastCalledWith(expect.stringContaining('..'))
    modal.onInput('valid.pdf', setError)
    expect(setError).toHaveBeenLastCalledWith('')

    // a confirm with an invalid name must not post anything
    modal.onConfirm('foo/bar.pdf')
    expect(postMessage).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('posts Host_PostmessageReady when the editor iframe loads', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper } = createShallowMountWrapper(makeRequest, { appName: 'Collabora' })
    await flushPromises()
    const { iframe, postMessage } = stubIframeContentWindow(wrapper)

    await iframe.trigger('load')

    expect(postMessage).toHaveBeenCalledTimes(1)
    const [payload, targetOrigin] = postMessage.mock.calls[0]
    expect(targetOrigin).toBe(appOrigin)
    expect(JSON.parse(payload as string)).toMatchObject({ MessageId: 'Host_PostmessageReady' })
    wrapper.unmount()
  })

  it('does not post Host_PostmessageReady or request ui_defaults for non-Collabora providers', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper } = createShallowMountWrapper(makeRequest, { appName: 'example-app' })
    await flushPromises()
    const { iframe, postMessage } = stubIframeContentWindow(wrapper)

    expect(wrapper.find('iframe').attributes('src')).toBe(appUrl)

    await iframe.trigger('load')

    expect(postMessage).not.toHaveBeenCalled()
    wrapper.unmount()
  })

  it('does not switch to write mode on UI_Edit when openAsPreview is false', async () => {
    const makeRequest = vi.fn().mockResolvedValue(successGet)
    const { wrapper } = createShallowMountWrapper(makeRequest, {
      appName: 'Collabora',
      openAsPreview: false
    })
    await flushPromises()
    makeRequest.mockClear()

    dispatchEditorMessage(appOrigin, 'UI_Edit')
    await wrapper.vm.$nextTick()

    expect(makeRequest).not.toHaveBeenCalled()
    wrapper.unmount()
  })
})

function createShallowMountWrapper(
  makeRequest = vi.fn().mockResolvedValue({ status: 200 }),
  {
    appName = 'example-app',
    isReadOnly = false,
    openAsPreview = true
  }: { appName?: string; isReadOnly?: boolean; openAsPreview?: boolean | string[] } = {}
) {
  vi.mocked(useRequest).mockImplementation(() => ({
    makeRequest
  }))

  vi.mocked(useRoute).mockReturnValue(
    computed(() => mock<RouteLocation>({ name: `external-${appName.toLowerCase()}-apps` }))
  )

  const dispatchModal = vi.fn()
  vi.mocked(useModals).mockReturnValue(mock<ReturnType<typeof useModals>>({ dispatchModal }))

  const showErrorMessage = vi.fn()
  vi.mocked(useMessages).mockReturnValue(mock<ReturnType<typeof useMessages>>({ showErrorMessage }))

  const mocks = {
    ...defaultComponentMocks(),
    $appProviderService: mock<AppProviderService>({ appNames: [appName] })
  }

  const capabilities = {
    files: {
      app_providers: [{ apps_url: '/app/list', enabled: true, open_url: '/app/open' }]
    }
  }

  return {
    dispatchModal,
    showErrorMessage,
    wrapper: shallowMount(App, {
      props: {
        space: null,
        resource: mock<Resource>({ name: 'document.odt' }),
        isReadOnly
      },
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              capabilityState: { capabilities },
              configState: { options: { editor: { openAsPreview } } }
            }
          })
        ],
        provide: mocks,
        mocks
      }
    })
  }
}
