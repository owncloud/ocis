import { defineComponent, h } from 'vue'
import AppWrapper from './AppWrapper.vue'
import { AppWrapperSlotArgs } from './types'
import { FileContentOptions, UrlForResourceOptions } from '../../composables'
import { Resource } from '@ownclouders/web-client'

export function AppWrapperRoute(
  fileEditor: ReturnType<typeof defineComponent>,
  options: {
    applicationId: string
    urlForResourceOptions?: UrlForResourceOptions
    fileContentOptions?: FileContentOptions
    importResourceWithExtension?: (resource: Resource) => string
    disableAutoSave?: boolean
  }
) {
  return defineComponent({
    render() {
      return h(
        AppWrapper,
        {
          wrappedComponent: fileEditor,
          ...options
        },
        {
          default: (slotArgs: AppWrapperSlotArgs) => {
            return h(fileEditor, slotArgs)
          }
        }
      )
    }
  })
}
