import { kebabCase } from 'lodash-es'
import { DateTime } from 'luxon'
import EventEmitter from 'events'
import { Actor } from '../../types'
import { ActorsOptions } from './shared'
import { ActorEnvironment } from './actor'
import { actorStore } from '../../store'

export class ActorsEnvironment extends EventEmitter {
  private readonly options: ActorsOptions

  constructor(options: ActorsOptions) {
    super()
    this.options = options
  }

  public getActor({ key }: { key: string }): Actor {
    if (!actorStore.has(key)) {
      throw new Error(`actor with key '${key}' not found`)
    }

    return actorStore.get(key)
  }

  public async createActor({ key, namespace }: { key: string; namespace: string }): Promise<Actor> {
    if (actorStore.has(key)) {
      return this.getActor({ key })
    }

    const actor = new ActorEnvironment({ id: key, namespace, ...this.options })
    await actor.setup()
    actor.on('closed', () => actorStore.delete(key))
    actor.page.on('console', (message) => {
      this.emit('console', key, message)
    })

    actorStore.set(key, actor)

    return actor
  }

  public async close(): Promise<void> {
    await Promise.all([...actorStore.values()].map((actor) => actor.close()))
  }

  public generateNamespace(scenarioTitle: string, user: string): string {
    return kebabCase([scenarioTitle, user, DateTime.now().toFormat('yyyy-M-d-hh-mm-ss')].join('-'))
  }
}
