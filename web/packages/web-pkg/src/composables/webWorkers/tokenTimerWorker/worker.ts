import { TokenTimerWorkerTopic } from './useTokenTimerWorker'

type Message = {
  topic: TokenTimerWorkerTopic
  expiry: number
  expiryThreshold: number
}

let timerId: ReturnType<typeof setTimeout>

const resetTimer = () => {
  clearTimeout(timerId)
  timerId = undefined
}

self.onmessage = (e: MessageEvent) => {
  const { topic, expiry, expiryThreshold } = JSON.parse(e.data) as Message

  if (topic === 'reset') {
    resetTimer()
    return
  }

  let timerInSeconds = expiry - expiryThreshold
  if (timerInSeconds <= 0) {
    // timer can't be smaller or equal 0
    timerInSeconds = 1
  }

  resetTimer()

  timerId = setTimeout(() => {
    postMessage(true)
  }, timerInSeconds * 1000)
}
