import { v4 as uuidV4 } from 'uuid'

export class EventBus {
  private topics: Map<string, { callback: (data?: unknown) => void; token: string }[]>

  constructor() {
    this.topics = new Map()
  }

  public subscribe<T = unknown>(topic: string, callback: (data?: T) => void): string {
    const subscription = {
      token: uuidV4(),
      callback: callback as (data?: unknown) => void
    }
    const subscriptions = [subscription, ...(this.topics.get(topic) || [])]

    this.topics.set(topic, subscriptions)

    return subscription.token
  }

  public publish(topic: string, data?: unknown): void {
    const subscriptions = this.topics.get(topic) || []

    subscriptions.forEach((subscription) => subscription.callback(data))
  }

  public unsubscribe(topic: string, token: string): void {
    if (!this.topics.has(topic)) {
      return
    }

    this.topics.set(
      topic,
      this.topics.get(topic).filter((subscription) => subscription.token !== token)
    )
  }
}

export const eventBus = new EventBus()
