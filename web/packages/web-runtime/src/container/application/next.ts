import { App } from 'vue'
import { RuntimeApi } from '../types'
export abstract class NextApplication {
  protected readonly runtimeApi: RuntimeApi

  protected constructor(runtimeApi: RuntimeApi) {
    this.runtimeApi = runtimeApi
  }

  abstract initialize(): Promise<void>

  abstract ready(): Promise<void>

  abstract mounted(instance: App): Promise<void>
}
