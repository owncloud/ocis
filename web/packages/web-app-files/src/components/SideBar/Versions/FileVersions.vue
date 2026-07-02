<template>
  <div id="oc-file-versions-sidebar" class="-oc-mt-s">
    <ul v-if="versions.length" class="oc-m-rm oc-position-relative">
      <li class="spacer oc-pb-l" aria-hidden="true"></li>
      <li
        v-for="(item, index) in versions"
        :key="index"
        class="version-item oc-pb-m oc-position-relative"
      >
        <div class="version-details">
          <span
            v-oc-tooltip="formatVersionDate(item)"
            class="version-date oc-font-semibold"
            data-testid="file-versions-file-last-modified-date"
          >
            {{ formatVersionDateRelative(item) }}
          </span>

          -
          <span class="version-filesize" data-testid="file-versions-file-size">
            {{ formatVersionFileSize(item) }}
          </span>
          <span tabindex="0" class="oc-invisible-sr">
            {{ formatVersionDateRelative(item) }} {{ formatVersionDate(item) }}
          </span>
        </div>
        <oc-list id="oc-file-versions-sidebar-actions" class="oc-pt-xs">
          <li v-if="isRevertible">
            <oc-button
              data-testid="file-versions-revert-button"
              appearance="raw"
              :aria-label="$gettext('Restore')"
              class="version-action-item oc-width-1-1 oc-rounded oc-button-justify-content-left oc-button-gap-m oc-py-s oc-px-m oc-display-block"
              @click="revertToVersion(item)"
            >
              <oc-icon name="history" class="oc-icon-m oc-mr-s -oc-mt-xs" fill-type="line" />
              {{ $gettext('Restore') }}
            </oc-button>
          </li>
          <li>
            <oc-button
              data-testid="file-versions-download-button"
              appearance="raw"
              :aria-label="$gettext('Download')"
              class="version-action-item oc-width-1-1 oc-rounded oc-button-justify-content-left oc-button-gap-m oc-py-s oc-px-m oc-display-block"
              @click="downloadVersion(item)"
            >
              <oc-icon name="file-download" class="oc-icon-m oc-mr-s" fill-type="line" />
              {{ $gettext('Download') }}
            </oc-button>
          </li>
        </oc-list>
      </li>
    </ul>
    <div v-else>
      <p v-translate data-testid="file-versions-no-versions">No versions available for this file</p>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { DavPermission } from '@ownclouders/web-client/webdav'
import {
  formatRelativeDateFromHTTP,
  formatDateFromJSDate,
  formatFileSize,
  useClientService,
  useDownloadFile,
  useResourcesStore
} from '@ownclouders/web-pkg'
import { computed, inject, Ref, unref } from 'vue'
import { isShareSpaceResource, Resource, SpaceResource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'

interface Props {
  isReadOnly?: boolean
}
const { isReadOnly } = defineProps<Props>()
const clientService = useClientService()
const language = useGettext()
const { downloadFile } = useDownloadFile({ clientService })
const { updateResourceField } = useResourcesStore()

const space = inject<Ref<SpaceResource>>('space')
const resource = inject<Ref<Resource>>('resource')
const versions = inject<Ref<Resource[]>>('versions')

const isRevertible = computed(() => {
  if (isReadOnly) {
    return false
  }

  if (isShareSpaceResource(unref(space)) || unref(resource).isReceivedShare()) {
    if (unref(resource).permissions !== undefined) {
      return unref(resource).permissions.includes(DavPermission.Updateable)
    }
  }

  return true
})

const revertToVersion = async (version: Resource) => {
  await clientService.webdav.restoreFileVersion(unref(space), unref(resource), version.name)
  const restoredResource = await clientService.webdav.getFileInfo(unref(space), unref(resource))

  const fieldsToUpdate = ['size', 'mdate'] as const
  for (const field of fieldsToUpdate) {
    if (Object.prototype.hasOwnProperty.call(unref(resource), field)) {
      updateResourceField({
        id: unref(resource).id,
        field: field,
        value: restoredResource[field]
      })
    }
  }
}
const downloadVersion = (version: Resource) => {
  return downloadFile(unref(space), unref(resource), version.name)
}
const formatVersionDateRelative = (version: Resource) => {
  return formatRelativeDateFromHTTP(version.mdate, language.current)
}
const formatVersionDate = (version: Resource) => {
  return formatDateFromJSDate(new Date(version.mdate), language.current)
}
const formatVersionFileSize = (version: Resource) => {
  return formatFileSize(version.size, language.current)
}
</script>

<style lang="scss" scoped>
#oc-file-versions-sidebar {
  > ul {
    list-style: none;

    .spacer {
      border-left: 1px solid var(--oc-color-border);
      margin-left: calc(-1 * var(--oc-space-large)) !important;
    }

    > li.version-item {
      border-left: 1px solid var(--oc-color-border);
      margin-left: calc(-1 * var(--oc-space-large)) !important;
      padding-left: var(--oc-space-medium);
      padding-bottom: var(--oc-space-medium);
      margin-top: calc(-1 * var(--oc-space-small));

      &::before {
        content: '';
        display: block;
        width: 11px;
        height: 11px;
        position: absolute;
        left: -6px;
        top: 4px;
        background-color: var(--oc-color-border);
        border-radius: 50%;
      }

      button.version-action-item {
        &:hover {
          color: var(--oc-color-primary-contrast);
          background-color: var(--oc-color-background-hover);
        }

        .oc-icon {
          vertical-align: middle;
        }
      }

      &:last-child {
        border-left: 1px solid transparent;
      }
    }
  }
}
</style>
