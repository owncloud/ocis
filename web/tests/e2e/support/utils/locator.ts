import { Locator } from '@playwright/test'
import { getSSEEvents } from '../environment/sse'
import { config } from '../../config'
import { User } from '../types'

export const waitForEvent = (locator: Locator, type: keyof SVGElementEventMap): Promise<void> =>
  locator.evaluate(
    (element, arg) =>
      new Promise<void>((resolve) => {
        const finalizer = () => {
          element.removeEventListener(arg.type, finalizer)
          resolve()
        }

        element.addEventListener(arg.type, finalizer)
      }),
    { type }
  )

export const buildXpathLiteral = (value: string) => {
  switch (true) {
    case !value.includes("'"):
      // if we don't have any single quotes, then wrap them with single quotes
      return "'" + value + "'"
    case !value.includes('"'):
      // if we don't have any double quotes, then wrap them with double quotes
      return '"' + value + '"'
    default:
      // use concat to find the literal in the xpath if they contain both quotes
      return "concat('" + value.replace(/'/g, "',\"'\",'") + "')"
  }
}

export const waitForSSEEvent = (user: User, event: string) => {
  return new Promise((resolve, reject) => {
    const startTime = Date.now()
    const interval = setInterval(function () {
      const events = getSSEEvents(user)
      if (Date.now() - startTime > config.minTimeout * 1000) {
        reject(new Error(`SSE event ${event} was not obtained in the events list:[${events}]`))
        clearInterval(interval)
      }
      if (events.indexOf(event) !== -1) {
        // remove the read event
        events.splice(events.indexOf(event), 1)
        resolve(true)
        clearInterval(interval)
      }
    }, 500)
  })
}
