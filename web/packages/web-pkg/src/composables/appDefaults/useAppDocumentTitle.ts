import { Ref, computed, unref } from 'vue'
import { basename } from 'path'
import { FileContext } from './types'
import { useAppMeta } from './useAppMeta'
import { useDocumentTitle } from './useDocumentTitle'
import { RouteLocationNormalizedLoaded } from 'vue-router'
import { MaybeRef } from 'vue'
import { useGettext } from 'vue3-gettext'
import { AppsStore } from '../piniaStores'

interface AppDocumentTitleOptions {
  appsStore: AppsStore
  applicationId: string
  applicationName?: MaybeRef<string>
  currentFileContext: Ref<FileContext>
  currentRoute?: Ref<RouteLocationNormalizedLoaded>
}

export function useAppDocumentTitle({
  appsStore,
  applicationId,
  applicationName,
  currentFileContext,
  currentRoute
}: AppDocumentTitleOptions): void {
  const appMeta = useAppMeta({ applicationId, appsStore })
  const { $gettext } = useGettext()

  const titleSegments = computed(() => {
    const baseTitle =
      basename(unref(unref(currentFileContext)?.fileName)) ||
      $gettext((unref(currentRoute)?.meta?.title as string) || '')
    const meta = unref(unref(appMeta).applicationMeta)

    return [baseTitle, unref(applicationName) || meta.name || meta.id].filter(Boolean)
  })

  useDocumentTitle({
    titleSegments
  })
}
