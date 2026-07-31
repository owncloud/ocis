import { storeToRefs } from 'pinia'
import { watch, Ref, unref } from 'vue'
import { useEventBus } from '../eventBus'
import { useThemeStore } from '../piniaStores'
import { EventBus } from '../../services'

interface DocumentTitleOptions {
  titleSegments: Ref<string[]>
  eventBus?: EventBus
}

export function useDocumentTitle({ titleSegments, eventBus }: DocumentTitleOptions): void {
  const themeStore = useThemeStore()
  const { currentTheme } = storeToRefs(themeStore)

  eventBus = eventBus || useEventBus()

  watch(
    titleSegments,
    (newTitleSegments) => {
      const titleSegments = unref(newTitleSegments)

      const glue = ' - '
      const payload = {
        shortDocumentTitle: titleSegments.join(glue),
        fullDocumentTitle: [...titleSegments, currentTheme.value.common.name].join(glue)
      }

      eventBus.publish('runtime.documentTitle.changed', payload)
    },
    { immediate: true, deep: true }
  )
}
