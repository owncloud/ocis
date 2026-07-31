import { v4 as uuidV4 } from 'uuid'
import { eventBus } from './eventBus'
import { debounce } from 'lodash-es'

export enum LoadingEventTopics {
  add = 'loading-service.add',
  remove = 'loading-service.remove',
  setProgress = 'loading-service.set-progress'
}

export interface LoadingTaskState {
  total: number
  current: number
}

export interface LoadingTask {
  id: string
  active: boolean
  state?: LoadingTaskState
}

export interface LoadingTaskCallbackArguments {
  setProgress: (args: LoadingTaskState) => void
}

// time until a loading task is being set active
const DEFAULT_DEBOUNCE_TIME = 200

export class LoadingService {
  private tasks: LoadingTask[] = []

  public get isLoading(): boolean {
    return this.tasks.some((e) => e.active)
  }

  /**
   * Get the current progress from 0 to 100.
   * Returns null if at least one task is indeterminate.
   */
  public get currentProgress(): number | null {
    if (this.tasks.some((e) => !e.state && e.active)) {
      return null
    }

    const tasks = this.tasks.filter((e) => !!e.state && e.active)
    if (!tasks.length) {
      return null
    }

    const progress = tasks.reduce((acc, task) => {
      acc += task.state.current / task.state.total
      return acc
    }, 0)

    return Math.round((progress / tasks.length) * 100)
  }

  public addTask<T>(
    callback: ({ setProgress }: LoadingTaskCallbackArguments) => Promise<T>,
    {
      debounceTime = DEFAULT_DEBOUNCE_TIME,
      indeterminate = true
    }: { debounceTime?: number; indeterminate?: boolean } = {}
  ): Promise<T> {
    const task = {
      id: uuidV4(),
      active: false,
      ...(!indeterminate && { state: { total: 0, current: 0 } })
    }

    // If no tasks are in progress, attach an event listener for 'beforeunload'.
    if (!this.tasks.length) {
      window.addEventListener('beforeunload', this.onBeforeUnload)
    }

    this.tasks.push(task)

    const debounced = debounce(() => {
      task.active = true
      eventBus.publish(LoadingEventTopics.add)
    }, debounceTime)
    debounced()

    const setProgress = ({ total, current }: LoadingTaskState) => {
      if (!indeterminate) {
        this.setProgress({ task, total, current })
      }
    }

    return callback({ setProgress }).finally(() => {
      this.removeTask(task.id)
    })
  }

  private removeTask(id: string): void {
    this.tasks = this.tasks.filter((e) => e.id !== id)

    if (!this.tasks.length) {
      window.removeEventListener('beforeunload', this.onBeforeUnload)
    }

    eventBus.publish(LoadingEventTopics.remove)
  }

  private setProgress({
    task,
    total,
    current
  }: {
    task: LoadingTask
    total: number
    current: number
  }): void {
    if (!task.state) {
      task.state = { total: 0, current: 0 }
    }
    task.state.total = total
    task.state.current = current
    eventBus.publish(LoadingEventTopics.setProgress)
  }

  private onBeforeUnload(e: BeforeUnloadEvent) {
    e.preventDefault()
  }
}
