import { config } from '../config'
import { environment } from '../support'
import { state } from './shared'

export class World {
  workerIndex: number
  testId: string

  actorsEnvironment: environment.ActorsEnvironment
  filesEnvironment: environment.FilesEnvironment
  linksEnvironment: environment.LinksEnvironment
  spacesEnvironment: environment.SpacesEnvironment
  usersEnvironment: environment.UsersEnvironment

  constructor(workerIndex: number = 0, testId: string = '') {
    this.workerIndex = workerIndex
    this.testId = testId

    this.usersEnvironment = new environment.UsersEnvironment()
    this.spacesEnvironment = new environment.SpacesEnvironment()
    this.filesEnvironment = new environment.FilesEnvironment()
    this.linksEnvironment = new environment.LinksEnvironment()
    this.actorsEnvironment = new environment.ActorsEnvironment({
      context: {
        acceptDownloads: config.acceptDownloads,
        reportDir: config.reportDir,
        tracingReportDir: config.tracingReportDir,
        reportHar: config.reportHar,
        reportTracing: config.reportTracing,
        reportVideo: config.reportVideo,
        failOnUncaughtConsoleError: config.failOnUncaughtConsoleError
      },
      browser: state.browser
    })
  }

  getGroupId(key: string): string {
    return `${key}-w${this.workerIndex}-${this.testId}`
  }

  getUserId(key: string): string {
    return `${key}-w${this.workerIndex}-${this.testId}`
  }

  getSpaceId(key: string): string {
    return `${key}-w${this.workerIndex}-${this.testId}`
  }

  /**
   * Transform resource name for parallel test safety.
   * Transforms: testfile.txt -> testfile-w1.txt (only when workerIndex > 0)
   */
  getResourceId(key: string): string {
    if (this.workerIndex === 0) {
      return key
    }

    const parts = key.split('/')
    const fileName = parts.at(-1)
    const dir = parts.slice(0, -1).join('/')
    const newFileName = fileName.includes('.')
      ? fileName.replace(/(\.[^.]+)$/, `-w${this.workerIndex}$1`)
      : `${fileName}-w${this.workerIndex}`

    return dir ? `${dir}/${newFileName}` : newFileName
  }
}

// --- Module-level context store ---
// Each Playwright worker is a separate Node.js process,
// so module state is isolated. No race conditions.
let _world: World | null = null

export const getWorld = (): World => {
  if (!_world) throw new Error('World not initialized — must run in Playwright test')
  return _world
}

export const setWorld = (world: World | null): void => {
  _world = world
}
