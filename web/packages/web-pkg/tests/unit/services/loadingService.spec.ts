import { LoadingService } from '../../../src/services/loadingService'

describe('LoadingService', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('adds a task and sets it inactive initially', () => {
    const service = new LoadingService()
    const action = new Promise((resolve) => {
      resolve(true)
    })

    service.addTask(() => action)
    expect(service.isLoading).toBeFalsy()
  })
  it('adds a task and sets it active after the debounce', () => {
    const service = new LoadingService()
    const action = new Promise((resolve) => {
      resolve(true)
    })

    service.addTask(() => action)
    vi.runAllTimers()
    expect(service.isLoading).toBeTruthy()
  })
  it('removes a task after being finished', async () => {
    const service = new LoadingService()
    const action = new Promise((resolve) => {
      resolve(true)
    })

    await service.addTask(() => action)
    vi.runAllTimers()
    expect(service.isLoading).toBeFalsy()
  })
  it('returns the current progress of a running task', () => {
    const service = new LoadingService()
    const expectedResult = {
      1: 25,
      2: 50,
      3: 75,
      4: 100
    }

    service.addTask(
      ({ setProgress }) => {
        const promises: Promise<void>[] = []
        const actions = [1, 2, 3, 4] as const
        for (const action of actions) {
          promises.push(
            new Promise((resolve) => {
              resolve(true)
            }).then(() => {
              setProgress({ total: actions.length, current: action })
              expect(service.currentProgress).toBe(expectedResult[action])
            })
          )
        }
        return Promise.all(promises)
      },
      { indeterminate: false }
    )
    vi.runAllTimers()
  })
})
