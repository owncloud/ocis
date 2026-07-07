import { SSEEventOptions } from './types'

export const onSSEBackchannelLogoutEvent = ({ router }: SSEEventOptions) => {
  return router.push({ name: 'logout' })
}
