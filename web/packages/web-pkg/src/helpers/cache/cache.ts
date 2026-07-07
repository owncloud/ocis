class CacheElement<T> {
  public value: T
  public expires: number

  constructor(value: T, ttl: number) {
    this.value = value
    this.expires = ttl ? new Date().getTime() + ttl : 0
  }

  get expired(): boolean {
    return this.expires > 0 && this.expires < new Date().getTime()
  }
}

interface CacheOptions {
  ttl?: number
  capacity?: number
}

export default class Cache<K, V> {
  private map: Map<K, CacheElement<V>>
  private readonly ttl: number
  private readonly capacity: number

  constructor(options: CacheOptions) {
    this.ttl = options.ttl || 0
    this.capacity = options.capacity || 0

    this.map = new Map<K, CacheElement<V>>()
  }

  public set(key: K, value: V, ttl?: number): V {
    this.evict()
    this.map.set(key, new CacheElement<V>(value, isNaN(ttl) ? this.ttl : ttl))

    return value
  }

  public get(key: K): V {
    this.evict()
    const entry = this.map.get(key)

    if (entry) {
      return entry.value
    }
  }

  public delete(key: K): boolean {
    return this.map.delete(key)
  }

  public clear(): void {
    return this.map.clear()
  }

  public entries(): [K, V][] {
    this.evict()
    return [...this.map.entries()].map((kv) => [kv[0], kv[1].value])
  }

  public keys(): K[] {
    this.evict()
    return [...this.map.keys()]
  }

  public has(key: K): boolean {
    this.evict()
    return this.map.has(key)
  }

  public values(): V[] {
    this.evict()
    return [...this.map.values()].map((e) => e.value)
  }

  public evict(): void {
    this.map.forEach((mv, mk) => {
      if (mv.expired) {
        this.delete(mk)
      }
    })

    if (!this.capacity) {
      return
    }

    for (const [k] of [...this.map.entries()]) {
      if (this.map.size <= this.capacity) {
        break
      }

      this.delete(k)
    }
  }
}
