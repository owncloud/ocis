import { EventBus } from '../../../src/services/eventBus'

describe('EventBus', () => {
  it('can handle load', () => {
    const bus = new EventBus()
    let val

    for (let i = 0; i < 1000; i++) {
      bus.subscribe(`evt.${i}`, (v) => (val = v))
    }

    for (let i = 0; i < 1000; i++) {
      const msg = `msg - ${i}`
      bus.publish(`evt.${i}`, msg)
      expect(val).toBe(msg)
    }

    val = 'untouched'

    for (let i = 0; i < 1000; i++) {
      bus.publish(`evt.unknown.${i}`, `unknown - ${i}`)
    }

    expect(val).toBe('untouched')
  })

  it('calls multiple subscriptions for the same topic', () => {
    const bus = new EventBus()
    const fn = vi.fn()
    const topic = 'evt.1'

    bus.subscribe(topic, fn)
    bus.subscribe(topic, fn)
    bus.publish(topic)

    expect(fn).toHaveBeenCalledTimes(2)
  })

  it('can unsubscribe a subscriber from a topic', () => {
    const bus = new EventBus()
    const fn = vi.fn()
    const topic = 'evt.1'

    const evt1Token = bus.subscribe(topic, fn)

    bus.publish(topic)
    expect(fn).toHaveBeenCalledTimes(1)

    bus.unsubscribe(topic, '')
    bus.publish(topic)
    expect(fn).toHaveBeenCalledTimes(2)

    bus.unsubscribe(topic, evt1Token)
    bus.publish(topic)
    expect(fn).toHaveBeenCalledTimes(2)
  })
})
