<template>
  <div class="oc-invisible-sr" role="status" aria-live="polite" aria-atomic="true">
    {{ announcements.polite }}
  </div>
  <div class="oc-invisible-sr" role="alert" aria-live="assertive" aria-atomic="true">
    {{ announcements.assertive }}
  </div>
  <oc-notifications>
    <oc-notification-message
      v-for="item in limitedMessages"
      :key="item.id"
      :title="item.title"
      :message="item.desc"
      :status="item.status"
      :timeout="item.timeout"
      :error-log-content="item.errorLogContent"
      :actions="item.actions"
      @close="deleteMessage(item)"
    />
  </oc-notifications>
</template>

<script lang="ts">
import { Message, useMessages } from '@ownclouders/web-pkg'
import { computed, defineComponent, onBeforeUnmount, reactive, watch } from 'vue'

type Politeness = 'polite' | 'assertive'

// lets focus changes triggered by the same action finish speaking first, otherwise
// NVDA/JAWS cancel the pending live region update
const ANNOUNCE_DELAY_MS = 150

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

    const announcements = reactive<Record<Politeness, string>>({ polite: '', assertive: '' })
    const pending: Record<Politeness, string[]> = { polite: [], assertive: [] }
    const timers: Partial<Record<Politeness, ReturnType<typeof setTimeout>>> = {}

    const announce = (politeness: Politeness, text: string) => {
      pending[politeness].push(text)
      if (timers[politeness]) {
        return
      }
      // emptying first makes a repeated identical message a real content change
      announcements[politeness] = ''
      timers[politeness] = setTimeout(() => {
        announcements[politeness] = pending[politeness].join(' ')
        pending[politeness] = []
        timers[politeness] = undefined
      }, ANNOUNCE_DELAY_MS)
    }

    watch(
      () => messageStore.messages.map(({ id }) => id),
      (_ids, previousIds = []) => {
        messageStore.messages
          .filter(({ id }) => !previousIds.includes(id))
          .forEach(({ title, desc, status }) => {
            announce(
              status === 'danger' ? 'assertive' : 'polite',
              desc ? `${title}. ${desc}` : title
            )
          })
      },
      { immediate: true }
    )

    onBeforeUnmount(() => {
      Object.values(timers).forEach(clearTimeout)
    })

    return { limitedMessages, deleteMessage, announcements }
  }
})
</script>
