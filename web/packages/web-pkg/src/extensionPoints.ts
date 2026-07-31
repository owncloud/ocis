import { Extension, ExtensionPoint, ResourceIndicatorExtension } from './'
import { computed } from 'vue'

export const resourceIndicatorExtensionPoint: ExtensionPoint<ResourceIndicatorExtension> = {
  id: 'global.files.resource-indicator',
  extensionType: 'resourceIndicator',
  multiple: true
}

export const extensionPoints = () => {
  return computed<ExtensionPoint<Extension>[]>(() => {
    return [resourceIndicatorExtensionPoint]
  })
}
