import { Resource } from '@ownclouders/web-client'
import { AppConfigObject } from '../../apps/types'
import { Ref } from 'vue'

export interface AppWrapperSlotArgs {
  applicationConfig: AppConfigObject
  resource: Resource
  currentContent: Ref<string>
  isDirty: boolean
  isReadOnly: boolean
  url: string
}
