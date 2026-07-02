import { Space } from '../types'
import { createdSpaceStore } from '../store'
import { getWorld } from '../../environment/world'

export class SpacesEnvironment {
  getSpace({ key }: { key: string }): Space {
    const world = getWorld()
    const storeKey = world ? world.getSpaceId(key) : key

    if (!createdSpaceStore.has(storeKey)) {
      throw new Error(`space with key '${storeKey}' not found`)
    }

    return createdSpaceStore.get(storeKey)
  }

  createSpace({ key, space }: { key: string; space: Space }): Space {
    const world = getWorld()
    const storeKey = world ? world.getSpaceId(key) : key

    if (createdSpaceStore.has(storeKey)) {
      throw new Error(`space with key '${storeKey}' already exists`)
    }

    createdSpaceStore.set(storeKey, space)

    return space
  }
}
