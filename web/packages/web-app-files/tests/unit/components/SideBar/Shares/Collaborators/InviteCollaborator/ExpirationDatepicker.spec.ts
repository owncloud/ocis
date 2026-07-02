import ExpirationDatepicker from '../../../../../../../src/components/SideBar/Shares/Collaborators/InviteCollaborator/ExpirationDatepicker.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

describe('InviteCollaborator ExpirationDatepicker', () => {
  it('renders a button to open the datepicker and set an expiration date', () => {
    const { wrapper } = createWrapper()
    expect(wrapper.find('[data-testid="recipient-datepicker-btn"]').exists()).toBe(true)
  })
})

const createWrapper = () => {
  return {
    wrapper: mount(ExpirationDatepicker, {
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
