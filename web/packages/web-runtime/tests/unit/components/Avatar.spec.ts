import Avatar from '../../../src/components/Avatar.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { CapabilityStore, ClientService } from '@ownclouders/web-pkg'
import { AxiosResponse } from 'axios'
import { nextTick } from 'vue'
import { OcAvatar } from '@ownclouders/design-system/components'

const propsData = {
  userName: 'admin',
  userid: 'admin',
  width: 24
}

const ocSpinner = 'oc-spinner-stub'
const ocAvatar = 'oc-avatar-stub'

describe('Avatar component', () => {
  window.URL.createObjectURL = vi.fn()

  it('should set user when the component is mounted', () => {
    const spySetUser = vi.spyOn(Avatar.methods, 'setUser')
    getShallowWrapper()
    expect(spySetUser).toHaveBeenCalledTimes(1)
    expect(spySetUser).toHaveBeenCalledWith(propsData.userid)
  })

  describe('when the component is still loading', () => {
    it('should render oc-spinner but not oc-avatar', () => {
      const { wrapper } = getShallowWrapper(true)
      const spinner = wrapper.find(ocSpinner)
      const avatar = wrapper.find(ocAvatar)

      expect(avatar.exists()).toBeFalsy()
      expect(spinner.exists()).toBeTruthy()
      expect(spinner.attributes().style).toEqual(
        `width: ${propsData.width}px; height: ${propsData.width}px;`
      )
    })
  })

  describe('when the component is not loading anymore', () => {
    it('should render oc-avatar but not oc-spinner', () => {
      const { wrapper } = getShallowWrapper()
      const spinner = wrapper.find(ocSpinner)
      const avatar = wrapper.find(ocAvatar)

      expect(spinner.exists()).toBeFalsy()
      expect(avatar.exists()).toBeTruthy()
    })
    it('should set props on oc-avatar component', () => {
      const { wrapper } = getShallowWrapper()
      const avatar = wrapper.findComponent<typeof OcAvatar>(ocAvatar)

      expect(avatar.props().width).toEqual(propsData.width)
      expect(avatar.props().userName).toEqual(propsData.userName)
    })

    describe('when an avatar is not found', () => {
      it('should set empty string to src prop on oc-avatar component', () => {
        const { wrapper } = getShallowWrapper()
        const avatar = wrapper.findComponent<typeof OcAvatar>(ocAvatar)
        expect(avatar.props().src).toEqual('')
      })
    })

    describe('when an avatar is found', () => {
      const blob = 'blob:https://web.org/6fe8f675-6727'
      it('should set blob as src prop on oc-avatar component', async () => {
        global.URL.createObjectURL = vi.fn(() => blob)
        const clientService = mockDeep<ClientService>()
        clientService.httpAuthenticated.get.mockResolvedValue(
          mock<AxiosResponse>({
            status: 200,
            data: blob
          })
        )
        const { wrapper } = getShallowWrapper(false, clientService)
        await nextTick()
        await nextTick()
        await nextTick()
        await nextTick()
        await nextTick()
        const avatar = wrapper.findComponent<typeof OcAvatar>(ocAvatar)
        expect(avatar.props().src).toEqual(blob)
      })
    })
  })
})

function getShallowWrapper(
  loading = false,
  clientService: ReturnType<typeof mockDeep<ClientService>> = undefined
) {
  if (!clientService) {
    clientService = mockDeep<ClientService>()
    clientService.httpAuthenticated.get.mockResolvedValue(mock<AxiosResponse>({ status: 200 }))
  }
  const mocks = { ...defaultComponentMocks(), $clientService: clientService }

  const capabilities = {
    files_sharing: { user: { profile_picture: true } }
  } satisfies Partial<CapabilityStore['capabilities']>

  return {
    wrapper: shallowMount(Avatar, {
      props: propsData,
      data() {
        return {
          loading
        }
      },
      global: {
        mocks,
        plugins: [...defaultPlugins({ piniaOptions: { capabilityState: { capabilities } } })]
      }
    })
  }
}
