<template>
  <div id="oc-file-details-sidebar">
    <div v-if="hasContent">
      <div
        v-if="isPreviewLoading || preview"
        key="file-thumbnail"
        :style="{
          'background-image': isPreviewLoading ? 'none' : `url(${preview})`
        }"
        class="details-preview oc-flex oc-flex-middle oc-flex-center oc-mb"
        data-testid="preview"
      >
        <oc-spinner v-if="isPreviewLoading" :aria-label="$gettext('Loading preview')" />
      </div>
      <div
        v-else
        class="details-icon-wrapper oc-width-1-1 oc-flex oc-flex-middle oc-flex-center oc-mb"
      >
        <resource-icon class="details-icon" :resource="resource" size="xxxlarge" />
      </div>
      <div
        v-if="!publicLinkContextReady && shareIndicators.length"
        key="file-shares"
        data-testid="sharingInfo"
        class="oc-flex oc-flex-middle oc-my-m"
      >
        <oc-status-indicators :resource="resource" :indicators="shareIndicators" />
        <p class="oc-my-rm oc-mx-s" v-text="detailSharingInformation" />
      </div>
      <table
        class="details-table oc-width-1-1"
        :aria-label="$gettext('Overview of the information about the selected file')"
      >
        <colgroup>
          <col class="oc-width-1-3" />
          <col class="oc-width-2-3" />
        </colgroup>
        <tbody>
          <tr v-if="hasDeletionDate" data-testid="delete-timestamp">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Deleted at')" />
            <td>
              <span v-text="capitalizedDeletionDate"></span>
            </td>
          </tr>
          <tr v-if="hasTimestamp" data-testid="timestamp">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Last modified')" />
            <td>
              <oc-button
                v-if="showVersions"
                v-oc-tooltip="seeVersionsLabel"
                appearance="raw"
                aria-labelledby="time-stamp"
                @click="expandVersionsPanel"
              >
                <span id="time-stamp" class="oc-invisible-sr">
                  {{
                    $pgettext(
                      'Accessible label for showing last modified date of a resource',
                      'Last modified %{ date }.',
                      {
                        date: capitalizedTimestamp
                      }
                    )
                  }}
                  {{
                    $pgettext(
                      'Accessible label for showing total number of versions a resource has',
                      'See all %{ versions } versions',
                      {
                        versions: unref(versions).length.toString()
                      }
                    )
                  }}
                </span>
                {{ capitalizedTimestamp }}
              </oc-button>
              <span v-else v-text="capitalizedTimestamp" />
            </td>
          </tr>
          <tr v-if="resource.locked" data-testid="locked-by">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Locked via')" />
            <td>
              <span>{{ resource.lockOwner }}</span>
              <span v-if="resource.lockTime">({{ formatDateRelative(resource.lockTime) }})</span>
            </td>
          </tr>
          <tr v-if="showSharedVia" data-testid="shared-via">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Shared via')" />
            <td>
              <router-link :to="sharedAncestorRoute">
                <span v-oc-tooltip="sharedViaTooltip" v-text="sharedAncestor.path" />
              </router-link>
            </td>
          </tr>
          <tr v-if="showSharedBy" data-testid="shared-by">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Shared by')" />
            <td>
              <span v-text="sharedByDisplayNames" />
            </td>
          </tr>
          <tr
            v-if="ownerDisplayName && ownerDisplayName !== sharedByDisplayNames"
            data-testid="ownerDisplayName"
          >
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Owner')" />
            <td>
              <p class="oc-m-rm">
                {{ ownerDisplayName }}
                <span v-if="ownedByCurrentUser" v-translate>(me)</span>
              </p>
            </td>
          </tr>
          <tr v-if="showSize" data-testid="sizeInfo">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Size')" />
            <td v-text="resourceSize" />
          </tr>
          <web-dav-details v-if="showWebDavDetails" :space="space" />
          <tr v-if="showVersions" data-testid="versionsInfo">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Versions')" />
            <td>
              <oc-button
                v-oc-tooltip="seeVersionsLabel"
                appearance="raw"
                aria-labelledby="versions-count"
                @click="expandVersionsPanel"
              >
                <span id="versions-count" class="oc-invisible-sr">
                  {{
                    $pgettext(
                      'Accessible label for showing total number of versions a resource has',
                      'See all %{ versions } versions',
                      {
                        versions: unref(versions).length.toString()
                      }
                    )
                  }}
                </span>
                {{ versions.length }}
              </oc-button>
            </td>
          </tr>
          <portal-target
            name="app.files.sidebar.file.details.table"
            :slot-props="{ space, resource }"
            :multiple="true"
          />
          <tr v-if="hasTags" data-testid="tags">
            <th scope="col" class="oc-pr-s oc-font-semibold">
              {{ $gettext('Tags') }}
              <oc-contextual-helper
                v-if="contextualHelper.isEnabled"
                v-bind="contextualHelper.data"
                class="oc-pl-xs"
              ></oc-contextual-helper>
            </th>
            <td>
              <tags-select :resource="resource"></tags-select>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-else data-testid="noContentText" v-text="$gettext('No information to display')" />
  </div>
</template>
<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { computed, inject, Ref, ref, unref, watch } from 'vue'
import {
  ImageDimension,
  useAuthStore,
  useUserStore,
  useCapabilityStore,
  useConfigStore,
  useResourcesStore,
  formatDateFromJSDate,
  useResourceContents,
  useLoadPreview
} from '@ownclouders/web-pkg'
import upperFirst from 'lodash-es/upperFirst'
import {
  isShareResource,
  isShareSpaceResource,
  isTrashResource,
  ShareTypes
} from '@ownclouders/web-client'
import { useGetMatchingSpace } from '@ownclouders/web-pkg'
import { getIndicators } from '@ownclouders/web-pkg'
import {
  formatDateFromHTTP,
  formatFileSize,
  formatRelativeDateFromJSDate
} from '@ownclouders/web-pkg'
import { eventBus } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { getSharedAncestorRoute } from '@ownclouders/web-pkg'
import { ResourceIcon } from '@ownclouders/web-pkg'
import { useContextualHelpers } from '../../../composables/contextualHelpers/useContextualHelpers'
import { ContextualHelper } from '@ownclouders/design-system/helpers'
import TagsSelect from './TagsSelect.vue'
import { WebDavDetails } from '@ownclouders/web-pkg'

interface Props {
  previewEnabled?: boolean
  tagsEnabled?: boolean
}
const { previewEnabled = true, tagsEnabled = true } = defineProps<Props>()
const configStore = useConfigStore()
const userStore = useUserStore()
const capabilityStore = useCapabilityStore()
const { getMatchingSpace } = useGetMatchingSpace()
const { resourceContentsText } = useResourceContents({ showSizeInformation: false })
const { loadPreview, previewsLoading } = useLoadPreview()
const { tagsHelper } = useContextualHelpers()

const language = useGettext()
const { $gettext, current: currentLanguage } = language

const resourcesStore = useResourcesStore()
const { ancestorMetaData, currentFolder } = storeToRefs(resourcesStore)

const { user } = storeToRefs(userStore)

const resource = inject<Ref<Resource>>('resource')
const versions = inject<Ref<Resource[]>>('versions')
const space = inject<Ref<SpaceResource>>('space')

const preview = ref<string>(undefined)

const authStore = useAuthStore()
const { publicLinkContextReady } = storeToRefs(authStore)

const isPreviewLoading = computed(() => previewEnabled && unref(previewsLoading))

const sharedAncestor = computed(() => {
  return Object.values(unref(ancestorMetaData)).find(
    (a) =>
      a.path !== unref(resource).path &&
      ShareTypes.containsAnyValue(ShareTypes.authenticated, a.shareTypes)
  )
})
const sharedAncestorRoute = computed(() => {
  return getSharedAncestorRoute({
    sharedAncestor: unref(sharedAncestor),
    matchingSpace: unref(space) || getMatchingSpace(unref(resource))
  })
})
const showWebDavDetails = computed(() => {
  /**
   * webDavPath might not be set when user is navigating on public link,
   * even if the user is authenticated and the file owner.
   */
  return resourcesStore.areWebDavDetailsShown && unref(resource).webDavPath
})
const formatDateRelative = (date: string) => {
  return formatRelativeDateFromJSDate(new Date(date), language.current)
}

const contextualHelper = computed<ContextualHelper>(() => ({
  isEnabled: configStore.options.contextHelpers,
  data: unref(tagsHelper)
}))

const hasTags = computed(() => {
  return tagsEnabled && capabilityStore.filesTags
})

const hasDeletionDate = computed(() => {
  return isTrashResource(unref(resource))
})

const capitalizedDeletionDate = computed(() => {
  const item = unref(resource)
  if (!isTrashResource(item)) {
    return ''
  }
  const displayDate = formatDateFromJSDate(new Date(item.ddate), language.current)
  return upperFirst(displayDate)
})

const shareIndicators = computed(() => {
  return getIndicators({
    space: unref(space),
    resource: unref(resource),
    ancestorMetaData: unref(ancestorMetaData),
    user: unref(user)
  }).filter(({ category }) => category === 'sharing')
})

const hasAnyShares = computed(() => {
  return unref(resource).shareTypes?.length > 0 || unref(sharedAncestor)
})
const sharedViaTooltip = computed(() => {
  return $gettext("Navigate to '%{folder}'", { folder: unref(sharedAncestor).path || '' }, true)
})
const showSharedBy = computed(() => {
  return unref(showShares) && !unref(ownedByCurrentUser) && unref(sharedByDisplayNames)
})
const showSharedVia = computed(() => {
  return unref(showShares) && unref(sharedAncestor) && !isShareSpaceResource(unref(space))
})
const showShares = computed(() => {
  if (unref(publicLinkContextReady)) {
    return false
  }
  return unref(hasAnyShares)
})
const ownedByCurrentUser = computed(() => {
  return unref(resource).owner?.id === unref(user)?.id
})
const sharedByDisplayNames = computed(() => {
  const res = unref(resource)
  if (!isShareResource(res)) {
    return ''
  }
  return res.sharedBy?.map(({ displayName }) => displayName).join(', ')
})
const hasContent = computed(() => {
  return (
    unref(hasTimestamp) ||
    unref(ownerDisplayName) ||
    unref(showSize) ||
    unref(showShares) ||
    unref(showVersions) ||
    unref(hasDeletionDate)
  )
})
const detailSharingInformation = computed(() => {
  if (unref(resource).type === 'folder') {
    return $gettext('This folder has been shared.')
  }
  return $gettext('This file has been shared.')
})
const hasTimestamp = computed(() => {
  return unref(resource).mdate?.length > 0
})
const ownerDisplayName = computed(() => {
  return unref(resource).owner?.displayName
})
const resourceSize = computed(() => {
  if (unref(resource).id === unref(currentFolder)?.id) {
    return `${formatFileSize(unref(resource).size, currentLanguage)}, ${unref(
      resourceContentsText
    )}`
  }

  return formatFileSize(unref(resource).size, currentLanguage)
})
const showSize = computed(() => {
  return formatFileSize(unref(resource).size, currentLanguage) !== '?'
})
const showVersions = computed(() => {
  if (unref(resource).type === 'folder' || unref(publicLinkContextReady)) {
    return
  }
  return unref(versions).length > 0
})
const seeVersionsLabel = computed(() => {
  return $gettext('See all versions')
})

const capitalizedTimestamp = computed(() => {
  const displayDate = formatDateFromHTTP(unref(resource).mdate, currentLanguage)
  return upperFirst(displayDate)
})

const expandVersionsPanel = () => {
  eventBus.publish(SideBarEventTopics.setActivePanel, 'versions')
}

watch(
  () => unref(resource).id,
  async () => {
    if (unref(resource)) {
      preview.value = await loadPreview({
        space: unref(space),
        resource: unref(resource),
        dimensions: ImageDimension.Preview,
        cancelRunning: true,
        updateStore: false
      })
    }
  },
  { immediate: true }
)
</script>
<style lang="scss" scoped>
.details-table {
  text-align: left;
  table-layout: fixed;

  tr {
    height: 1.5rem;
  }
}

.details-preview,
.details-icon-wrapper {
  background-color: var(--oc-color-background-highlight);
  border: 10px solid var(--oc-color-background-highlight);
  height: 230px;

  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

.details-icon > svg {
  height: 192px !important;
  max-height: 192px !important;
  max-width: 192px !important;
  width: 192px !important;
}
</style>
