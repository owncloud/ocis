<template>
  <oc-notifications>
    <oc-notification-message
      v-for="item in limitedMessages"
      :key="item.id"
      :title="item.title"
      :message="item.desc"
      :status="item.status"
      :timeout="item.timeout"
      :error-log-content="item.errorLogContent"
      @close="deleteMessage(item)"
    />
  </oc-notifications>
</template>

<script lang="ts">
import { Message, useMessages } from '@ownclouders/web-pkg'
import { computed, defineComponent } from 'vue'

export default defineComponent({
  name: 'MessageBar',
  setup() {
    const messageStore = useMessages()

    const limitedMessages = computed(() => {
      return messageStore.messages ? messageStore.messages.slice(0, 5) : []
    })

    const deleteMessage = (message: Message) => {
      messageStore.removeMessage(message)
    }

    return { limitedMessages, deleteMessage }
  }
})
</script>
