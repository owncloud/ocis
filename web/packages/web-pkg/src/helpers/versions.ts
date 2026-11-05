import { CapabilityStore } from '../composables/index'

export const getBackendVersion = ({
  capabilityStore
}: {
  capabilityStore: CapabilityStore
}): string => {
  const backendStatus = capabilityStore.status
  if (!backendStatus || !backendStatus.versionstring) {
    return undefined
  }
  const product = backendStatus.product || 'ownCloud'
  const version = backendStatus.productversion || backendStatus.versionstring
  const edition = backendStatus.edition
  return `${product} ${version} ${edition}`
}
