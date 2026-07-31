import { Message, useMessages } from '@ownclouders/web-pkg'
import { OcNotificationMessage } from '@ownclouders/design-system/components'
import MessageBar from '../../../src/components/MessageBar.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

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
  notificationMessage: 'oc-notification-message-stub'
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
