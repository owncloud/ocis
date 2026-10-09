import ExpirationDatepicker from '../../../../../../../src/components/SideBar/Shares/Collaborators/InviteCollaborator/ExpirationDatepicker.vue'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { useModals } from '@ownclouders/web-pkg'
import { DateTime } from 'luxon'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useModals: vi.fn()
}))

describe('InviteCollaborator ExpirationDatepicker', () => {
  const dispatchModal = vi.fn()

  beforeEach(() => {
    dispatchModal.mockClear()
    vi.mocked(useModals).mockReturnValue({ dispatchModal } as any)
  })

  it('renders a button to open the datepicker and set an expiration date', () => {
    const { wrapper } = createWrapper()
    expect(wrapper.find('[data-testid="recipient-datepicker-btn"]').exists()).toBe(true)
  })

  it('opens the picker with no suggested date by default', async () => {
    const { wrapper } = createWrapper()
    await wrapper.find('[data-testid="recipient-datepicker-btn"]').trigger('click')

    const attrs = dispatchModal.mock.calls[0][0].customComponentAttrs()
    expect(attrs.currentDate).toBeNull()
  })

  it('pre-fills the suggested date when no date has been set yet', async () => {
    const suggestedDate = DateTime.now().plus({ days: 30 })
    const { wrapper } = createWrapper({ suggestedDate })
    await wrapper.find('[data-testid="recipient-datepicker-btn"]').trigger('click')

    const attrs = dispatchModal.mock.calls[0][0].customComponentAttrs()
    expect(attrs.currentDate.toMillis()).toBe(suggestedDate.toMillis())
  })
})

const createWrapper = ({ suggestedDate }: { suggestedDate?: DateTime } = {}) => {
  return {
    wrapper: mount(ExpirationDatepicker, {
      props: { suggestedDate },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
