import { MfaExpiryWorkerTopic } from './useMfaExpiryWorker'

type Message = {
  topic: MfaExpiryWorkerTopic
  expiresAt: number
  warningThreshold: number
}

let timerId: ReturnType<typeof setTimeout>

const resetTimer = () => {
  clearTimeout(timerId)
  timerId = undefined
}

globalThis.onmessage = (e: MessageEvent) => {
  const { topic, expiresAt, warningThreshold } = JSON.parse(e.data) as Message

  if (topic === 'reset') {
    resetTimer()
    return
  }

  const now = Math.floor(Date.now() / 1000)
  let timerInSeconds = expiresAt - warningThreshold - now
  if (timerInSeconds <= 0) {
    timerInSeconds = 1
  }

  resetTimer()

  timerId = setTimeout(() => {
    postMessage(true)
  }, timerInSeconds * 1000)
}
