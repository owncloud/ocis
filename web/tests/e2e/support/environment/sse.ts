import { fetchEventSource, EventStreamContentType } from '@ai-zen/node-fetch-event-source'
import { getAuthHeader } from '../api/http'
import { User } from '../types'

const sseEventStore: Record<string, string[]> = {}
const sseConnections = []

export const listenSSE = (baseUrl: string, user: User): Promise<void> => {
  const sseUrl = new URL('ocs/v2.php/apps/notifications/api/v1/notifications/sse', baseUrl).href

  const ctrl = new AbortController()
  sseConnections.push(ctrl)

  return fetchEventSource(sseUrl, {
    headers: {
      ...getAuthHeader(user)
    },
    signal: ctrl.signal,
    onopen(response) {
      if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
        console.log('Listening to SSE events...')
        return Promise.resolve()
      }
    },
    onmessage(message) {
      if (message.event === 'FatalError') {
        throw new Error(message.data)
      }
      const userKey = user.id.toLowerCase()
      if (!Object.hasOwn(sseEventStore, userKey)) {
        sseEventStore[userKey] = []
      }
      // push event to the array
      // TODO: also store message.data if necessary
      sseEventStore[userKey].push(message.event)
    },
    onclose() {
      console.error('Closing SSE...')
    },
    onerror(err) {
      console.error(err)
      ctrl.abort()
    }
  })
}

export const getSSEEvents = (user: User): Array<string> => {
  // recent events should be evaluated first
  return sseEventStore[user.id.toLowerCase()].reverse()
}

export const closeSSEConnections = () => {
  Object.keys(sseEventStore).forEach((key) => delete sseEventStore[key])
  sseConnections.forEach((connection) => connection.abort())
  sseConnections.length = 0
}
