import { fetchEventSource, FetchEventSourceInit } from '@microsoft/fetch-event-source'

export enum MESSAGE_TYPE {
  NOTIFICATION = 'userlog-notification',
  POSTPROCESSING_FINISHED = 'postprocessing-finished',
  FILE_LOCKED = 'file-locked',
  FILE_UNLOCKED = 'file-unlocked',
  FILE_TOUCHED = 'file-touched',
  ITEM_RENAMED = 'item-renamed',
  ITEM_TRASHED = 'item-trashed',
  ITEM_RESTORED = 'item-restored',
  ITEM_MOVED = 'item-moved',
  FOLDER_CREATED = 'folder-created',
  SPACE_MEMBER_ADDED = 'space-member-added',
  SPACE_MEMBER_REMOVED = 'space-member-removed',
  SPACE_SHARE_UPDATED = 'space-share-updated',
  SHARE_CREATED = 'share-created',
  SHARE_REMOVED = 'share-removed',
  SHARE_UPDATED = 'share-updated',
  LINK_CREATED = 'link-created',
  LINK_REMOVED = 'link-removed',
  LINK_UPDATED = 'link-updated',
  BACKCHANNEL_LOGOUT = 'backchannel-logout'
}

export class RetriableError extends Error {
  name = 'RetriableError'
}

const RECONNECT_RANDOM_OFFSET = 15000

export class SSEAdapter implements EventSource {
  url: string
  fetchOptions: FetchEventSourceInit
  private abortController: AbortController
  private eventListenerMap: Record<string, ((event: MessageEvent) => any)[]>

  readyState: number
  readonly withCredentials: boolean

  readonly CONNECTING: 0
  readonly OPEN: 1
  readonly CLOSED: 2

  onerror: ((this: EventSource, ev: Event) => any) | null
  onmessage: ((this: EventSource, ev: MessageEvent) => any) | null
  onopen: ((this: EventSource, ev: Event) => any) | null

  constructor(url: string, fetchOptions: FetchEventSourceInit) {
    this.url = url
    this.fetchOptions = fetchOptions
    this.abortController = new AbortController()
    this.eventListenerMap = {}
    this.readyState = this.CONNECTING
    this.connect()
  }

  private connect() {
    return fetchEventSource(this.url, {
      openWhenHidden: true,
      signal: this.abortController.signal,
      fetch: this.fetchProvider.bind(this),
      // eslint-disable-next-line require-await
      onopen: async () => {
        const event = new Event('open')
        this.onopen?.bind(this)(event)
        this.readyState = this.OPEN
      },
      onmessage: (msg) => {
        const event = new MessageEvent('message', { data: msg.data })
        this.onmessage?.bind(this)(event)

        const type = msg.event
        const eventListeners = this.eventListenerMap[type]
        eventListeners?.forEach((l) => l(event))
      },
      onclose: () => {
        this.readyState = this.CLOSED
        throw new RetriableError()
      },
      onerror: (err) => {
        console.error(err)
        const event = new CustomEvent('error', { detail: err })
        this.onerror?.bind(this)(event)

        /*
         * Try to reconnect after 30 seconds plus random time in seconds.
         * This prevents all clients try to reconnect concurrent on server error, to reduce load.
         */
        return 30000 + Math.floor(Math.random() * RECONNECT_RANDOM_OFFSET)
      }
    })
  }

  fetchProvider(...args: [RequestInfo | URL, RequestInit]) {
    const [resource, config] = args
    const fetchConfig = { ...config, ...this.fetchOptions }
    return window.fetch(resource, fetchConfig)
  }

  close() {
    this.abortController.abort('closed')
  }

  addEventListener<K extends keyof EventSourceEventMap>(
    type: string,
    listener: (this: EventSource, event: EventSourceEventMap[K]) => any
  ): void {
    this.eventListenerMap[type] = this.eventListenerMap[type] || []
    this.eventListenerMap[type].push(listener)
  }

  removeEventListener<K extends keyof EventSourceEventMap>(
    type: string,
    listener: (this: EventSource, event: EventSourceEventMap[K]) => any
  ): void {
    this.eventListenerMap[type] = this.eventListenerMap[type]?.filter((func) => func !== listener)
  }

  dispatchEvent(): boolean {
    throw new Error('Method not implemented.')
  }

  updateAccessToken(token: string) {
    this.fetchOptions.headers['Authorization'] = `Bearer ${token}`
  }

  updateLanguage(language: string) {
    this.fetchOptions.headers['Accept-Language'] = language

    // Force reconnect, to make the language change effect instantly
    this.close()
    this.connect()
  }
}

let eventSource: SSEAdapter = null

export const sse = (baseURI: string, fetchOptions: FetchEventSourceInit): SSEAdapter => {
  if (!eventSource) {
    eventSource = new SSEAdapter(
      new URL('ocs/v2.php/apps/notifications/api/v1/notifications/sse', baseURI).href,
      fetchOptions
    )
  }

  return eventSource
}
