import { Resource } from '@ownclouders/web-client'
import { useAppsStore } from '../piniaStores'

export const useCanBeOpenedWithSecureView = () => {
  const appsStore = useAppsStore()

  const canBeOpenedWithSecureView = (resource: Resource) => {
    const secureViewExtensions = appsStore.fileExtensions.filter(({ secureView }) => secureView)
    return secureViewExtensions.some(({ mimeType }) => mimeType === resource.mimeType)
  }

  return { canBeOpenedWithSecureView }
}
