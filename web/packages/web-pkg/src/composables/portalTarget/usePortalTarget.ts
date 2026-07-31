import { useWormhole } from 'portal-vue'
import { TransportInput } from 'portal-vue/types'
import { useEventBus } from '../eventBus'
import { PortalTargetEventTopics } from './eventTopics'

export const usePortalTarget = () => {
  const eventBus = useEventBus()
  const wormhole = useWormhole()

  const registerPortal = (transportInput: TransportInput) => {
    eventBus.subscribe(PortalTargetEventTopics.mounted, () => {
      wormhole.open(transportInput)
    })
  }

  return {
    registerPortal
  }
}
