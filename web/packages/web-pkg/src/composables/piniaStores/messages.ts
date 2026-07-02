import { defineStore } from 'pinia'
import { v4 as uuidV4 } from 'uuid'
import { ref, unref } from 'vue'
import { HttpError } from '@ownclouders/web-client'

type MessageError = Error | HttpError

export interface Message {
  id: string
  title: string
  desc?: string
  errors?: MessageError[]
  errorLogContent?: string
  timeout?: number
  status?: string
}

export const useMessages = defineStore('messages', () => {
  const messages = ref<Message[]>([])

  const showMessage = (data: Omit<Message, 'id'>) => {
    const message = { ...data, id: uuidV4() as string }
    messages.value.push(message)
    return message
  }

  const getXRequestIdsFromErrors = (errors: Message['errors']) => {
    const getXRequestID = (error: MessageError) => {
      return (error as HttpError).response?.headers?.get('x-request-id')
    }

    return errors
      .map((error) => getXRequestID(error))
      .filter((xRequestId) => xRequestId !== null)
      .map((item) => `X-Request-Id: ${item}`)
      .join('\r\n')
  }

  const showErrorMessage = (data: Omit<Message, 'id'>) => {
    const message = {
      id: uuidV4() as string,
      status: 'danger',
      timeout: 0,
      ...(data.errors && { errorLogContent: getXRequestIdsFromErrors(data.errors) }),
      ...data
    }

    messages.value.push(message)
    return message
  }

  const removeMessage = (message: Message) => {
    messages.value = unref(messages).filter(({ id }) => message.id !== id)
  }

  return {
    messages,
    showMessage,
    showErrorMessage,
    removeMessage
  }
})

export type MessageStore = ReturnType<typeof useMessages>
