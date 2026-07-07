import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import CopyPrivateLink from '../../../../src/components/Shares/CopyPrivateLink.vue'
import { useMessages } from '@ownclouders/web-pkg'
import { useClipboard } from '@vueuse/core'
import { ref } from 'vue'

const resource = mock<Resource>({
  type: 'folder',
  owner: {
    id: 'marie',
    displayName: 'Marie'
  },
  mdate: 'Wed, 21 Oct 2015 07:28:00 GMT',
  size: '740',
  name: 'lorem.txt',
  privateLink: 'https://example.com/fake-private-link'
})

vi.mock('@vueuse/core', () => ({
  useClipboard: vi.fn(() => ({
    copy: vi.fn(),
    copied: false,
    isSupported: true
  }))
}))

describe('CopyPrivateLink', () => {
  it('should render a button', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.html()).toMatchSnapshot()
  })
  it('upon clicking it should copy the private link to the clipboard button, render a success message and change icon for half a second', async () => {
    const copyMock = vi.fn()
    vi.mocked(useClipboard).mockReturnValue(
      mock<ReturnType<typeof useClipboard>>({ copy: copyMock, copied: ref(true) })
    )

    const { wrapper } = getWrapper()
    const { showMessage } = useMessages()
    expect(showMessage).not.toHaveBeenCalled()

    await wrapper.find('button').trigger('click')
    expect(copyMock).toHaveBeenCalledTimes(1)
    expect(showMessage).toHaveBeenCalledTimes(1)
  })
})

function getWrapper() {
  return {
    wrapper: mount(CopyPrivateLink, {
      props: {
        resource
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
