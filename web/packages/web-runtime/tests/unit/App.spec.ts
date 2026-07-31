import App from '../../src/App.vue'
import { ref } from 'vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { CapabilityStore, ClientService } from '@ownclouders/web-pkg'
import { AxiosResponse } from 'axios'
import * as LanguageHelpderModule from '../../src/helpers/language'

vi.spyOn(LanguageHelpderModule, 'setCurrentLanguage')

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRouter: vi.fn(() => ({
    resolve: vi.fn().mockImplementation(() => ref({ name: 'lorem name', path: '/lorem' })),
    currentRoute: ref('lorem current route')
  }))
}))
vi.mock('vue-router', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useRoute: vi.fn(() => ({ meta: { title: 'lorem title' }, to: '/lorem' }))
}))

vi.mock('@vueuse/core', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  MaybeElement: vi.fn(),
  useElementSize: vi.fn().mockImplementation(() => ({ height: 0 }))
}))

describe('App', () => {
  test('to set html language to de on mounted', () => {
    getShallowWrapper({ getTextDefaultLanguage: 'de' })
    expect(LanguageHelpderModule.setCurrentLanguage).toHaveBeenCalledWith(
      expect.objectContaining({ language: expect.anything(), languageSetting: 'de' })
    )
    expect(document.documentElement.lang).toBe('de')
  })
  test('to set html language to en on mounted', () => {
    getShallowWrapper({})
    expect(LanguageHelpderModule.setCurrentLanguage).toHaveBeenCalledWith(
      expect.objectContaining({ language: expect.anything(), languageSetting: 'en' })
    )
    expect(document.documentElement.lang).toBe('en')
  })
})

function getShallowWrapper({
  clientService = undefined,
  getTextDefaultLanguage = 'en'
}: {
  clientService?: ReturnType<typeof mockDeep<ClientService>>
  getTextDefaultLanguage?: string
}) {
  if (!clientService) {
    clientService = mockDeep<ClientService>()
    clientService.httpAuthenticated.get.mockResolvedValue(mock<AxiosResponse>({ status: 200 }))
  }
  const mocks = { ...defaultComponentMocks(), $clientService: clientService }

  const capabilities = {
    files_sharing: { user: { profile_picture: true } }
  } satisfies Partial<CapabilityStore['capabilities']>

  return {
    wrapper: shallowMount(App, {
      global: {
        mocks,
        plugins: [
          ...defaultPlugins({
            piniaOptions: { capabilityState: { capabilities } },
            getTextDefaultLanguage
          })
        ]
      }
    })
  }
}
