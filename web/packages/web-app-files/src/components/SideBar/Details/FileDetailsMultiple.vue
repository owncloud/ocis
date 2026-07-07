<template>
  <div id="oc-file-details-multiple-sidebar">
    <div class="files-preview oc-mb">
      <div class="files-preview-body">
        <oc-icon class="preview-icon" size="xxlarge" variation="passive" name="file-copy" />
        <p class="preview-text" data-testid="selectedFilesText" v-text="selectedFilesString" />
      </div>
    </div>
    <div>
      <table class="details-table" :aria-label="detailsTableLabel" role="presentation">
        <tbody>
          <tr data-testid="filesCount">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="filesText" />
            <td v-text="filesCount" />
          </tr>
          <tr data-testid="foldersCount">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="foldersText" />
            <td v-text="foldersCount" />
          </tr>
          <tr v-if="showSpaceCount" data-testid="spacesCount">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="spacesText" />
            <td v-text="spacesCount" />
          </tr>
          <tr v-if="hasSize" data-testid="size">
            <th scope="col" class="oc-pr-s oc-font-semibold" v-text="sizeText" />
            <td v-text="sizeValue" />
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, unref } from 'vue'
import { storeToRefs } from 'pinia'
import { useGettext } from 'vue3-gettext'
import { formatFileSize, useResourcesStore } from '@ownclouders/web-pkg'

interface Props {
  showSpaceCount?: boolean
}
const { showSpaceCount = false } = defineProps<Props>()
const resourcesStore = useResourcesStore()
const language = useGettext()
const { $ngettext, $gettext } = language
const { selectedResources } = storeToRefs(resourcesStore)

const hasSize = computed(() => {
  return unref(selectedResources).some((resource) => Object.hasOwn(resource, 'size'))
})
const selectedFilesCount = computed(() => {
  return unref(selectedResources).length
})
const selectedFilesString = computed(() => {
  return $ngettext(
    '%{ itemCount } item selected',
    '%{ itemCount } items selected',
    unref(selectedFilesCount),
    {
      itemCount: unref(selectedFilesCount).toString()
    }
  )
})
const sizeValue = computed(() => {
  let size = 0
  unref(selectedResources).forEach((i) => (size += parseInt(i.size?.toString() || '0')))
  return formatFileSize(size, language.current)
})
const sizeText = computed(() => {
  return $gettext('Size')
})
const filesCount = computed(() => {
  return unref(selectedResources).filter((i) => i.type === 'file').length
})
const filesText = computed(() => {
  return $gettext('Files')
})
const foldersCount = computed(() => {
  return unref(selectedResources).filter((i) => i.type === 'folder').length
})
const foldersText = computed(() => {
  return $gettext('Folders')
})
const spacesCount = computed(() => {
  return unref(selectedResources).filter((i) => i.type === 'space').length
})
const spacesText = computed(() => {
  return $gettext('Spaces')
})
const detailsTableLabel = computed(() => {
  return $gettext('Overview of the information about the selected files')
})
</script>
<style lang="scss" scoped>
.files-preview {
  position: relative;
  background-color: var(--oc-color-background-muted);
  border: 10px solid var(--oc-color-background-muted);
  height: 230px;
  text-align: center;
  color: var(--oc-color-swatch-passive-muted);

  &-body {
    margin: 0;
    position: absolute;
    top: 50%;
    left: 50%;
    -ms-transform: translate(-50%, -50%);
    transform: translate(-50%, -50%);

    .preview-icon {
      display: inline-block;
    }

    .preview-text {
      display: block;
    }
  }
}

.details-table {
  text-align: left;

  tr {
    height: 1.5rem;
  }
}
</style>
