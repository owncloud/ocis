import { PartialComponentProps, mount } from '@ownclouders/web-test-helpers'
import App from '../../src/App.vue'

vi.mock('@ownclouders/web-pkg')

describe('Text editor app', () => {
  it('shows the editor', () => {
    const { wrapper } = getWrapper({
      applicationConfig: {}
    })
    expect(wrapper.html()).toMatchSnapshot()
  })
})

function getWrapper(props: PartialComponentProps<typeof App>) {
  return {
    wrapper: mount(App, {
      props: {
        applicationConfig: {},
        currentContent: '',
        isReadOnly: false,
        resource: undefined,
        ...props
      }
    })
  }
}
