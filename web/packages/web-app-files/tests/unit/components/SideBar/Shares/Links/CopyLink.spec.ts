import { ref } from 'vue'
import { LinkShare } from '@ownclouders/web-client'
import { useMessages } from '@ownclouders/web-pkg'
import CopyLink from '../../../../../../src/components/SideBar/Shares/Links/CopyLink.vue'
import { useClipboard } from '@vueuse/core'
import { mock } from 'vitest-mock-extended'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

const linkShare = {
  displayName: 'Example link',
  webUrl: 'https://some-url.com/abc'
} as LinkShare

vi.mock('@vueuse/core', () => ({
  useClipboard: vi.fn(() => ({
    copy: vi.fn(),
    copied: false,
    isSupported: true
  }))
}))

describe('CopyLink', () => {
  // ignore tippy warning
  vi.spyOn(console, 'warn').mockImplementation(undefined)
  it('upon clicking it should copy the link to the clipboard, render a success message and change icon for half a second', async () => {
    const copyMock = vi.fn()
    const copiedRef = ref(true)
    vi.mocked(useClipboard).mockReturnValue(
      mock<ReturnType<typeof useClipboard>>({ copy: copyMock, copied: copiedRef })
    )

    const { wrapper } = getWrapper()
    const { showMessage } = useMessages()
    expect(showMessage).not.toHaveBeenCalled()

    await wrapper.find('.oc-files-public-link-copy-url').trigger('click')
    expect(copyMock).toHaveBeenCalledTimes(1)
    expect(wrapper.html()).toMatchSnapshot()
    expect(showMessage).toHaveBeenCalledTimes(1)

    copiedRef.value = false

    await wrapper.vm.$nextTick()
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function getWrapper() {
  return {
    wrapper: mount(CopyLink, {
      props: {
        linkShare
      },
      global: {
        plugins: [...defaultPlugins()],
        stubs: {
          'oc-icon': {
            template: '<span class="oc-icon" />'
          }
        }
      }
    })
  }
}
