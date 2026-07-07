import { MockInstance } from 'vitest'
import Cache from '../../../../src/helpers/cache/cache'

const newCache = <T>(vs: T[], ttl?: number, capacity?: number): Cache<number, T> => {
  const cache = new Cache<number, T>({ ttl, capacity })
  vs.forEach((v, i) => cache.set(i, v))
  return cache
}

describe('Cache', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('can set and get entries', () => {
    const cacheValues: number[] = [1, 2, 3, 4]
    const cache = newCache(cacheValues)

    cacheValues.forEach((v, i) => {
      expect(cache.get(i)).toBe(v)
    })

    cache.set(4, 5)
    expect(cache.get(4)).toBe(5)
    expect(cache.set(5, 6)).toBe(6)
  })

  it('return all keys', () => {
    const cache = newCache([1, 2, 3, 4])
    expect(cache.keys()).toMatchObject([0, 1, 2, 3])
  })

  it('return all values', () => {
    const cacheValues: number[] = [1, 2, 3, 4]
    const cache = newCache(cacheValues)
    expect(cache.values()).toMatchObject(cacheValues)
  })

  it('return all entries', () => {
    const cache = newCache([1, 2, 3, 4])
    expect(cache.entries()).toMatchObject([
      [0, 1],
      [1, 2],
      [2, 3],
      [3, 4]
    ])
  })

  it('can handle ttl', () => {
    const cacheValues: number[] = []
    const cache = newCache(cacheValues, 50)

    cache.set(1, 1)
    vi.setSystemTime(new Date().getTime() + 10)
    cache.set(2, 2)

    expect(cache.get(1)).toBe(1)
    expect(cache.get(2)).toBe(2)
    expect(cache.values().length).toBe(2)
    expect(cache.keys().length).toBe(2)
    expect(cache.entries().length).toBe(2)

    vi.setSystemTime(new Date().getTime() + 41)

    expect(cache.get(1)).toBeFalsy()
    expect(cache.get(2)).toBe(2)
    expect(cache.values().length).toBe(1)
    expect(cache.keys().length).toBe(1)
    expect(cache.entries().length).toBe(1)

    vi.setSystemTime(new Date().getTime() + 10)

    expect(cache.get(2)).toBeFalsy()
    expect(cache.values().length).toBe(0)
    expect(cache.keys().length).toBe(0)
    expect(cache.entries().length).toBe(0)

    cache.set(3, 3, 10)
    cache.set(4, 4, 20)
    cache.set(5, 5, 0)
    cache.set(6, 6, 0)

    expect(cache.get(3)).toBe(3)
    expect(cache.get(4)).toBe(4)
    expect(cache.get(5)).toBe(5)
    expect(cache.get(6)).toBe(6)
    expect(cache.values().length).toBe(4)
    expect(cache.keys().length).toBe(4)
    expect(cache.entries().length).toBe(4)

    vi.setSystemTime(new Date().getTime() + 11)

    expect(cache.get(3)).toBeFalsy()
    expect(cache.get(4)).toBe(4)
    expect(cache.get(5)).toBe(5)
    expect(cache.get(6)).toBe(6)
    expect(cache.values().length).toBe(3)
    expect(cache.keys().length).toBe(3)
    expect(cache.entries().length).toBe(3)

    vi.setSystemTime(new Date().getTime() + 10)

    expect(cache.get(4)).toBeFalsy()
    expect(cache.get(5)).toBe(5)
    expect(cache.get(6)).toBe(6)
    expect(cache.values().length).toBe(2)
    expect(cache.keys().length).toBe(2)
    expect(cache.entries().length).toBe(2)
  })

  it('can handle capacity', () => {
    const initialValues: number[] = [1, 2, 3, 4, 5]
    const newValues: number[] = [6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
    const capacity = 5
    const cache = newCache(initialValues, 0, capacity)

    newValues.forEach((v) => {
      cache.set(v, v)
      expect(cache.get(v)).toBe(v)
      expect(cache.entries().length).toBe(capacity)
    })
  })

  it('can clear the cache', () => {
    const cache = newCache([1, 2, 3, 4, 5])
    expect(cache.entries().length).toBe(5)
    cache.clear()
    expect(cache.entries().length).toBe(0)
  })

  it('can check if a cache contains a entry for given key', () => {
    const values = [1, 2, 3, 4, 5]
    const cache = newCache(values)

    values.forEach((value) => expect(cache.has(value - 1)).toBeTruthy())
    expect(cache.has(5)).toBeFalsy()
  })
})

describe('cache', () => {
  describe('CacheElement', () => {
    let cache: Cache<string, unknown>
    let key: string, value: string, key2: string, value2: string
    let evictSpy: MockInstance
    beforeEach(() => {
      const options = { ttl: 0, opacity: 0 }
      evictSpy = vi.spyOn(Cache.prototype, 'evict')
      cache = new Cache<string, string>(options)
      key = 'key'
      value = 'value'
      key2 = 'key2'
      value2 = 'value2'
    })
    it('should set value and be receivable with get', () => {
      expect(cache.set(key, value)).toBe(value)
      expect(cache.get(key)).toBe(value)
    })
    it('should evict before any access', () => {
      cache.set(key, value)
      cache.get(key)
      cache.entries()
      cache.keys()
      cache.values()
      cache.has(key)
      expect(evictSpy).toHaveBeenCalledTimes(6)
    })
    it('should delete key', () => {
      cache.set(key, value)
      cache.set(key2, value2)
      expect(cache.delete(key)).toBeTruthy()
      expect(cache.get(key)).toBe(undefined)
      expect(cache.get(key2)).toBe(value2)
    })
    it('should clear cache', () => {
      cache.set(key, value)
      cache.set(key2, value2)
      cache.clear()
      expect(cache.keys().length).toBe(0)
    })
    it('should return cache entries', () => {
      cache.set(key, value)
      cache.set(key2, value2)
      expect(cache.entries()).toStrictEqual([
        [key, value],
        [key2, value2]
      ])
    })
    it('should return keys', () => {
      cache.set(key, value)
      cache.set(key2, value2)
      expect(cache.keys()).toStrictEqual([key, key2])
    })
    it('should return values', () => {
      cache.set(key, value)
      cache.set(key2, value2)
      expect(cache.values()).toStrictEqual([value, value2])
    })
    it('should return if has key', () => {
      cache.set(key, value)
      expect(cache.has(key)).toBe(true)
      expect(cache.has(key2)).toBe(false)
    })
    it('should evict expired item', () => {
      const oldGetTime = Date.prototype.getTime
      Date.prototype.getTime = vi.fn(() => 1487076708000)
      cache.set(key, value, 1)
      Date.prototype.getTime = oldGetTime
      expect(cache.values().length).toBe(0)
    })
    it('should not evict item that has not expired', () => {
      cache.set(key, value, 1000)
      expect(cache.values().length).toBe(1)
    })
    it('should not evict item without ttl', () => {
      const oldGetTime = Date.prototype.getTime
      Date.prototype.getTime = vi.fn(() => 1487076708000)
      cache.set(key, value, 0)
      Date.prototype.getTime = oldGetTime
      expect(cache.values().length).toBe(1)
    })
  })
})
