<template>
  <resource-list-item
    ref="resourceListItemRef"
    :resource="resource"
    :path-prefix="pathPrefix"
    :is-path-displayed="true"
    :link="resourceLink"
    :is-extension-displayed="areFileExtensionsShown"
    :parent-folder-link-icon-additional-attributes="parentFolderLinkIconAdditionalAttributes"
    :parent-folder-name="parentFolderName"
    :is-thumbnail-displayed="!!previewData"
    :is-search-result="isSearchResult"
    v-bind="additionalAttrs"
  />
</template>

<script lang="ts" setup>
import { ImageDimension } from '../../constants'
import { VisibilityObserver } from '../../observer'
import { debounce } from 'lodash-es'
import { computed, ref, unref, onMounted, onBeforeMount, useTemplateRef } from 'vue'
import {
  useGetMatchingSpace,
  useFileActions,
  useFolderLink,
  useResourcesStore,
  useLoadPreview
} from '../../composables'
import { isSpaceResource, Resource } from '@ownclouders/web-client'
import ResourceListItem from '../FilesList/ResourceListItem.vue'
import { SearchResultValue } from './types'
import { RouteLocationPathRaw } from 'vue-router'

const visibilityObserver = new VisibilityObserver()

interface Props {
  searchResult?: SearchResultValue
  isClickable?: boolean
  term?: string
  isSearchResult?: boolean
}
const {
  searchResult = { data: {} },
  isClickable = true,
  term = '',
  isSearchResult
} = defineProps<Props>()
const { triggerDefaultAction } = useFileActions()
const { getMatchingSpace } = useGetMatchingSpace()
const { getDefaultAction } = useFileActions()
const { loadPreview } = useLoadPreview()

const {
  getPathPrefix,
  getParentFolderName,
  getParentFolderLink,
  getParentFolderLinkIconAdditionalAttributes,
  getFolderLink
} = useFolderLink()
const resourcesStore = useResourcesStore()
const resourceListItemRef = useTemplateRef('resourceListItemRef')

const previewData = ref<string>()

const areFileExtensionsShown = computed(() => resourcesStore.areFileExtensionsShown)

const resource = computed((): Resource => {
  return {
    ...(searchResult.data as Resource),
    ...(unref(previewData) &&
      ({
        thumbnail: unref(previewData)
      } as Resource))
  }
})

const pathPrefix = getPathPrefix(unref(resource))
const parentFolderName = getParentFolderName(unref(resource))
const parentFolderLinkIconAdditionalAttributes = getParentFolderLinkIconAdditionalAttributes(
  unref(resource.value)
)

const space = computed(() => getMatchingSpace(unref(resource)))

const resourceDisabled = computed(() => {
  const res = unref(resource)
  return isSpaceResource(res) && res.disabled === true
})

const resourceClicked = () => {
  triggerDefaultAction({
    space: unref(space),
    resources: [unref(resource)]
  })
}

const additionalAttrs = computed(() => {
  if (!isClickable) {
    return {
      isResourceClickable: false
    }
  }

  return {
    parentFolderLink: getParentFolderLink(unref(resource)),
    onClick: resourceClicked
  }
})

const resourceLink = computed(() => {
  if (unref(resource).isFolder) {
    return getFolderLink(unref(resource))
  }

  const action = getDefaultAction({ resources: [unref(resource)], space: unref(space) })

  if (!action?.route) {
    return null
  }

  const route = action.route({
    space: unref(space),
    resources: [unref(resource)]
  }) as RouteLocationPathRaw

  // add search term to query param
  route.query = {
    ...route.query,
    contextRouteQuery: {
      ...((route.query?.contextRouteQuery as any) || {}),
      term
    }
  }

  return route
})
onMounted(() => {
  /*
   * Accessing the parent element via defineExpose in <ResourceListItem />
   * */
  if (unref(resourceDisabled)) {
    resourceListItemRef.value.resourceListItem.parentElement.classList.add('disabled')
  }

  const loadPreviewHandler = async () => {
    const preview = await loadPreview({
      space: unref(space),
      resource: unref(resource),
      dimensions: ImageDimension.Thumbnail,
      cancelRunning: true
    })

    preview && (previewData.value = preview)
  }

  const debounced = debounce(({ unobserve }) => {
    unobserve()
    loadPreviewHandler()
  }, 250)

  visibilityObserver.observe(resourceListItemRef.value.resourceListItem, {
    onEnter: debounced,
    onExit: debounced.cancel
  })
})
onBeforeMount(() => {
  visibilityObserver.disconnect()
})
</script>
