import { Message, useMessages } from '@ownclouders/web-pkg'
import { OcNotificationMessage } from '@ownclouders/design-system/components'
import MessageBar from '../../../src/components/MessageBar.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { nextTick } from 'vue'

const messages = [
  {
    id: '101',
    title: 'Error while moving',
    desc: '',
    status: 'danger'
  },
  {
    id: '102',
    title: 'Error while deleting',
    desc: '',
    status: 'danger'
  },
  {
    id: '103',
    title: 'Error while renaming',
    desc: '',
    status: 'danger'
  },
  {
    id: '104',
    title: 'Error while copying',
    desc: '',
    status: 'danger'
  },
  {
    id: '105',
    title: 'Error while restoring',
    desc: '',
    status: 'danger'
  },
  {
    id: '106',
    title: 'Error while uploading',
    desc: '',
    status: 'danger'
  }
]

const selectors = {
  notificationMessage: 'oc-notification-message-stub',
  politeRegion: '[aria-live="polite"]',
  assertiveRegion: '[aria-live="assertive"]'
}

describe('MessageBar component', () => {
  describe('when there is an active message', () => {
    it('should set props in oc-notification-message component', () => {
      const { wrapper } = getShallowWrapper([messages[0]])
      const notificationMessage = wrapper.findComponent<typeof OcNotificationMessage>(
        selectors.notificationMessage
      )

      expect(notificationMessage.attributes().title).toEqual(messages[0].title)
      expect(notificationMessage.attributes().status).toEqual(messages[0].status)
      expect(notificationMessage.attributes().message).toEqual(messages[0].desc)
    })
    it('should call "removeMessage" method on close event', () => {
      const { wrapper } = getShallowWrapper([messages[0]])
      const messageStore = useMessages()
      const notificationMessage = wrapper.findComponent<typeof OcNotificationMessage>(
        selectors.notificationMessage
      )
      notificationMessage.vm.$emit('close')

      expect(messageStore.removeMessage).toHaveBeenCalledTimes(1)
    })
  })

  describe('when there are more than five active messages', () => {
    it('should return only the first five messages', () => {
      const { wrapper } = getShallowWrapper(messages)

      expect(wrapper.findAll(selectors.notificationMessage).length).toBe(5)
    })
  })

  describe('screen reader announcements', () => {
    beforeEach(() => {
      vi.useFakeTimers()
    })
    afterEach(() => {
      vi.useRealTimers()
    })

    const addMessage = async (message: Message) => {
      useMessages().messages.push(message)
      await nextTick()
      vi.runOnlyPendingTimers()
      await nextTick()
    }

    it('renders empty live regions before any message is shown', () => {
      const { wrapper } = getShallowWrapper()

      expect(wrapper.find(selectors.politeRegion).text()).toBe('')
      expect(wrapper.find(selectors.assertiveRegion).text()).toBe('')
    })

    it('announces a new message in the polite region', async () => {
      const { wrapper } = getShallowWrapper()
      await addMessage({ id: '1', title: 'Copied to clipboard!', status: 'success' } as Message)

      expect(wrapper.find(selectors.politeRegion).text()).toBe('Copied to clipboard!')
      expect(wrapper.find(selectors.assertiveRegion).text()).toBe('')
    })

    it('announces title and description', async () => {
      const { wrapper } = getShallowWrapper()
      await addMessage({ id: '1', title: 'Link created', desc: 'Link copied' } as Message)

      expect(wrapper.find(selectors.politeRegion).text()).toBe('Link created. Link copied')
    })

    it('announces danger messages in the assertive region', async () => {
      const { wrapper } = getShallowWrapper()
      await addMessage(messages[0] as Message)

      expect(wrapper.find(selectors.assertiveRegion).text()).toBe(messages[0].title)
      expect(wrapper.find(selectors.politeRegion).text()).toBe('')
    })

    it('delays the announcement so it is not swallowed by focus changes', async () => {
      const { wrapper } = getShallowWrapper()
      useMessages().messages.push({ id: '1', title: 'Restored' } as Message)
      await nextTick()

      expect(wrapper.find(selectors.politeRegion).text()).toBe('')
    })

    it('clears the region before re-announcing an identical message', async () => {
      const { wrapper } = getShallowWrapper()
      await addMessage({ id: '1', title: 'Copied to clipboard!' } as Message)
      useMessages().messages.push({ id: '2', title: 'Copied to clipboard!' } as Message)
      await nextTick()

      expect(wrapper.find(selectors.politeRegion).text()).toBe('')

      vi.runOnlyPendingTimers()
      await nextTick()
      expect(wrapper.find(selectors.politeRegion).text()).toBe('Copied to clipboard!')
    })

    it('does not re-announce messages when one is removed', async () => {
      const { wrapper } = getShallowWrapper()
      await addMessage({ id: '1', title: 'First' } as Message)
      await addMessage({ id: '2', title: 'Second' } as Message)
      const store = useMessages()
      store.messages = store.messages.filter(({ id }) => id !== '1')
      await nextTick()
      vi.runOnlyPendingTimers()
      await nextTick()

      expect(wrapper.find(selectors.politeRegion).text()).toBe('Second')
    })
  })
})

function getShallowWrapper(messages: Message[] = []) {
  return {
    wrapper: shallowMount(MessageBar, {
      global: {
        renderStubDefaultSlot: true,
        plugins: [...defaultPlugins({ piniaOptions: { messagesState: { messages } } })]
      }
    })
  }
}
