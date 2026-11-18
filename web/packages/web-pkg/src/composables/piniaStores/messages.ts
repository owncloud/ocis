import { defineStore } from 'pinia'
import { v4 as uuidV4 } from 'uuid'
import { ref, unref } from 'vue'
import { HttpError } from '@ownclouders/web-client'

type MessageError = Error | HttpError

/**
 * The request ids of the given errors, one "X-Request-Id: …" line each, as shown in the details of an
 * error message. Errors without a request id (e.g. a request that got no answer) are left out.
 */
export const getXRequestIdsLog = (errors: MessageError[]): string =>
  errors
    .map((error) => (error as HttpError).response?.headers?.get('x-request-id'))
    .filter(Boolean)
    .map((xRequestId) => `X-Request-Id: ${xRequestId}`)
    .join('\r\n')

export interface MessageAction {
  label: string
  ariaLabel?: string
  onClick: () => void
}

export interface Message {
  id: string
  title: string
  desc?: string
  errors?: MessageError[]
  errorLogContent?: string
  timeout?: number
  status?: string
  actions?: MessageAction[]
}

export const useMessages = defineStore('messages', () => {
  const messages = ref<Message[]>([])

  const showMessage = (data: Omit<Message, 'id'>) => {
    const message = { ...data, id: uuidV4() as string }
    messages.value.push(message)
    return message
  }

  const showErrorMessage = (data: Omit<Message, 'id'>) => {
    const message = {
      id: uuidV4() as string,
      status: 'danger',
      timeout: 0,
      ...(data.errors && { errorLogContent: getXRequestIdsLog(data.errors) }),
      ...data
    }

    messages.value.push(message)
    return message
  }

  const removeMessage = (message: Message) => {
    messages.value = unref(messages).filter(({ id }) => message.id !== id)
  }

  /**
   * Runs the first action of the most recent message that has one (e.g. "Undo" on a
   * delete notification) and dismisses that message. Used by the global Ctrl/Cmd+Z
   * shortcut. Returns false if there is no message with an action to trigger.
   */
  const triggerLatestAction = (): boolean => {
    const message = [...unref(messages)].reverse().find((m) => m.actions?.length)
    if (!message) {
      return false
    }

    message.actions[0].onClick()
    removeMessage(message)
    return true
  }

  return {
    messages,
    showMessage,
    showErrorMessage,
    removeMessage,
    triggerLatestAction
  }
})

export type MessageStore = ReturnType<typeof useMessages>
