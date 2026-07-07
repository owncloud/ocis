import { Link } from '../types'
import { createdLinkStore } from '../store'

export class LinksEnvironment {
  // Store copied password for parallel test safety
  copiedPassword: string = ''

  getLink({ name }: { name: string }): Link {
    if (!createdLinkStore.has(name)) {
      throw new Error(`link with name '${name}' not found`)
    }
    return createdLinkStore.get(name)
  }

  updateLinkName({ key, link }: { key: string; link: Link }): any {
    if (!createdLinkStore.has(key)) {
      throw new Error(`link with name '${key}' not found`)
    }
    createdLinkStore.set(link.name, link)
    createdLinkStore.delete(key)
  }

  createLink({ key, link }: { key: string; link: Link }): Link {
    if (createdLinkStore.has(key)) {
      throw new Error(`link with key '${key}' already exists`)
    }
    createdLinkStore.set(key, link)
    return link
  }

  deleteLink({ key }: { key: string }): boolean {
    if (createdLinkStore.has(key)) {
      return createdLinkStore.delete(key)
    }
    throw new Error(`link with key '${key}' doesn't exist`)
  }
}
