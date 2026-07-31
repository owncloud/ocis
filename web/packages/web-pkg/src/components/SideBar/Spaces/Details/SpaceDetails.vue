<template>
  <div id="oc-space-details-sidebar">
    <div class="oc-space-details-sidebar-image oc-text-center">
      <oc-spinner v-if="previewsLoading" />
      <div v-else-if="spaceImage" class="oc-position-relative">
        <img :src="spaceImage" alt="" class="oc-mb-s" />
      </div>
      <oc-icon
        v-else
        name="layout-grid"
        size="xxlarge"
        class="space-default-image oc-px-m oc-py-m"
      />
    </div>
    <div
      v-if="showShareIndicators && hasShares && !resource.disabled"
      class="oc-flex oc-flex-middle oc-space-details-sidebar-members oc-mb-s oc-text-small"
      style="gap: 15px"
    >
      <oc-button
        v-if="hasMemberShares"
        appearance="raw"
        :aria-label="openSharesPanelMembersHint"
        @click="expandSharesPanel"
      >
        <oc-icon name="group" />
      </oc-button>
      <oc-button
        v-if="hasLinkShares"
        appearance="raw"
        :aria-label="openSharesPanelLinkHint"
        @click="expandSharesPanel"
      >
        <oc-icon name="link" />
      </oc-button>
      <p v-text="shareLabel" />
      <oc-button
        appearance="raw"
        :aria-label="openSharesPanelHint"
        size="small"
        @click="expandSharesPanel"
      >
        <span class="oc-text-small" v-text="$gettext('Show')" />
      </oc-button>
    </div>
    <table class="details-table oc-width-1-1" :aria-label="detailsTableLabel">
      <colgroup>
        <col class="oc-width-1-3" />
        <col class="oc-width-2-3" />
      </colgroup>
      <tbody>
        <tr>
          <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Last activity')" />
          <td v-text="lastModifiedDate" />
        </tr>
        <tr v-if="resource.description">
          <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Subtitle')" />
          <td v-text="resource.description" />
        </tr>
        <tr>
          <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Manager')" />
          <td>
            <span v-text="ownerUsernames" />
          </td>
        </tr>
        <tr v-if="!resource.disabled">
          <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Quota')" />
          <td>
            <space-quota :space-quota="resource.spaceQuota" />
          </td>
        </tr>
        <tr v-if="showSize" data-testid="sizeInfo">
          <th scope="col" class="oc-pr-s oc-font-semibold" v-text="$gettext('Size')" />
          <td v-text="size" />
        </tr>
        <web-dav-details v-if="showWebDavDetails" :space="resource" />
        <portal-target
          name="app.files.sidebar.space.details.table"
          :slot-props="{ space: resource, resource }"
          :multiple="true"
        />
      </tbody>
    </table>
  </div>
</template>
<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { inject, ref, Ref, computed, unref, watch } from 'vue'
import { buildSpaceImageResource, getSpaceManagers, SpaceResource } from '@ownclouders/web-client'
import {
  useUserStore,
  useSharesStore,
  useResourcesStore,
  useResourceContents,
  useRouter,
  useLoadPreview
} from '../../../../composables'
import SpaceQuota from '../../../SpaceQuota.vue'
import WebDavDetails from '../../WebDavDetails.vue'
import { formatDateFromISO, formatFileSize } from '../../../../helpers'
import { eventBus } from '../../../../services/eventBus'
import { SideBarEventTopics } from '../../../../composables'
import { ImageDimension } from '../../../../constants'
import { ProcessorType } from '../../../../services'
import { isLocationSpacesActive } from '../../../../router'
import { useGettext } from 'vue3-gettext'

interface Props {
  showShareIndicators?: boolean
}
const { showShareIndicators = true } = defineProps<Props>()
const userStore = useUserStore()
const resourcesStore = useResourcesStore()
const { resourceContentsText } = useResourceContents({ showSizeInformation: false })
const router = useRouter()
const { current: currentLanguage, $ngettext, $gettext } = useGettext()
const { loadPreview, previewsLoading } = useLoadPreview()

const sharesStore = useSharesStore()

const resource = inject<Ref<SpaceResource>>('resource')
const spaceImage = ref('')

const { user } = storeToRefs(userStore)

const linkShareCount = computed(() => sharesStore.linkShares.length)
const showWebDavDetails = computed(() => resourcesStore.areWebDavDetailsShown)
const showSize = computed(() => {
  return !isLocationSpacesActive(router, 'files-spaces-projects')
})
const size = computed(() => {
  return `${formatFileSize(unref(resource).size, currentLanguage)}, ${unref(resourceContentsText)}`
})
const hasShares = computed(() => {
  return unref(hasMemberShares) || unref(hasLinkShares)
})
const shareLabel = computed(() => {
  if (unref(hasMemberShares) && !unref(hasLinkShares)) {
    return unref(memberShareLabel)
  }
  if (!unref(hasMemberShares) && unref(hasLinkShares)) {
    return unref(linkShareLabel)
  }

  switch (unref(memberShareCount)) {
    case 1:
      return $ngettext(
        'This space has one member and %{linkShareCount} link.',
        'This space has one member and %{linkShareCount} links.',
        unref(linkShareCount),
        { linkShareCount: unref(linkShareCount).toString() }
      )
    default:
      if (unref(linkShareCount) === 1) {
        return $gettext('This space has %{memberShareCount} members and one link.', {
          memberShareCount: unref(memberShareCount).toString()
        })
      }
      return $gettext('This space has %{memberShareCount} members and %{linkShareCount} links.', {
        memberShareCount: unref(memberShareCount).toString(),
        linkShareCount: unref(linkShareCount).toString()
      })
  }
})
const openSharesPanelHint = computed(() => {
  return $gettext('Open share panel')
})
const openSharesPanelLinkHint = computed(() => {
  return $gettext('Open link list in share panel')
})
const openSharesPanelMembersHint = computed(() => {
  return $gettext('Open member list in share panel')
})
const detailsTableLabel = computed(() => {
  return $gettext('Overview of the information about the selected space')
})
const lastModifiedDate = computed(() => {
  return formatDateFromISO(unref(resource).mdate, currentLanguage)
})
const ownerUsernames = computed(() => {
  const managers = getSpaceManagers(unref(resource))
  return managers
    .map((share) => {
      const member = share.grantedTo.user || share.grantedTo.group
      if (member.id === unref(user)?.id) {
        return $gettext('%{displayName} (me)', { displayName: member.displayName })
      }
      return member.displayName
    })
    .join(', ')
})
const hasMemberShares = computed(() => {
  return unref(memberShareCount) > 0
})
const hasLinkShares = computed(() => {
  return unref(linkShareCount) > 0
})
const memberShareCount = computed(() => {
  return Object.keys(unref(resource).members).length
})
const memberShareLabel = computed(() => {
  return $ngettext(
    'This space has %{memberShareCount} member.',
    'This space has %{memberShareCount} members.',
    unref(memberShareCount),
    { memberShareCount: unref(memberShareCount).toString() }
  )
})
const linkShareLabel = computed(() => {
  return $ngettext(
    '%{linkShareCount} link giving access.',
    '%{linkShareCount} links giving access.',
    unref(linkShareCount),
    { linkShareCount: unref(linkShareCount).toString() }
  )
})

function expandSharesPanel() {
  eventBus.publish(SideBarEventTopics.setActivePanel, 'space-share')
}

watch(
  () => unref(resource).spaceImageData,
  async () => {
    if (!unref(resource).spaceImageData) {
      return
    }

    const imageResource = buildSpaceImageResource(unref(resource))
    spaceImage.value = await loadPreview({
      space: unref(resource),
      resource: imageResource,
      dimensions: ImageDimension.Tile,
      processor: ProcessorType.enum.fit,
      cancelRunning: true,
      updateStore: false
    })
  },
  { immediate: true }
)
</script>
<style lang="scss" scoped>
.oc-space-details-sidebar {
  &-image img {
    max-height: 150px;
    object-fit: cover;
    width: 100%;
  }
}

.details-table {
  text-align: left;
  table-layout: fixed;

  tr {
    height: 1.5rem;
  }
}
</style>
