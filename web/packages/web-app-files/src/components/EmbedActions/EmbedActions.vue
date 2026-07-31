<template>
  <section class="files-embed-actions oc-width-1-1 oc-flex oc-flex-middle oc-flex-between oc-my-s">
    <oc-text-input
      v-if="chooseFileName"
      v-model="fileName"
      class="files-embed-actions-file-name oc-flex oc-flex-row oc-flex-middle"
      :selection-range="fileNameInputSelectionRange"
      :label="$gettext('File name')"
    />

    <div class="files-embed-actions-buttons oc-flex oc-flex-middle">
      <oc-button
        class="oc-mr-s"
        data-testid="button-cancel"
        appearance="raw-inverse"
        variation="brand"
        @click="emitCancel"
        >{{ $gettext('Cancel') }}
      </oc-button>
      <oc-button
        v-if="!isLocationPicker && !isFilePicker"
        key="btn-share"
        class="oc-mr-s"
        data-testid="button-share"
        variation="inverse"
        appearance="filled"
        :disabled="
          areSelectActionsDisabled ||
          !createLinkAction.isVisible({ resources: selectedFiles, space })
        "
        @click="createLinkAction.handler({ resources: selectedFiles, space })"
        >{{ $gettext('Share link(s)') }}
      </oc-button>
      <oc-button
        v-if="!isFilePicker"
        data-testid="button-select"
        variation="inverse"
        appearance="filled"
        :disabled="areSelectActionsDisabled"
        @click="emitSelect"
        >{{ selectLabel }}
      </oc-button>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { computed, ref, unref } from 'vue'
import {
  embedModeLocationPickMessageData,
  FileAction,
  routeToContextQuery,
  useEmbedMode,
  useFileActionsCreateLink,
  useResourcesStore,
  useRouter,
  useSpacesStore
} from '@ownclouders/web-pkg'
import { Resource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { storeToRefs } from 'pinia'

const { $gettext } = useGettext()
const { isLocationPicker, isFilePicker, postMessage, chooseFileName, chooseFileNameSuggestion } =
  useEmbedMode()
const spacesStore = useSpacesStore()
const router = useRouter()
const { currentSpace: space } = storeToRefs(spacesStore)
const resourcesStore = useResourcesStore()
const { currentFolder, selectedResources } = storeToRefs(resourcesStore)
const fileName = ref(unref(chooseFileNameSuggestion))

const selectedFiles = computed<Resource[]>(() => {
  if (isLocationPicker.value) {
    return [unref(currentFolder)]
  }

  return unref(selectedResources)
})

const { actions: createLinkActions } = useFileActionsCreateLink({ enforceModal: true })
const createLinkAction = computed<FileAction>(() => unref(createLinkActions)[0])

const areSelectActionsDisabled = computed<boolean>(() => selectedFiles.value.length < 1)

const selectLabel = computed<string>(() =>
  isLocationPicker.value ? $gettext('Choose') : $gettext('Attach as copy')
)

const fileNameInputSelectionRange = computed(() => {
  return [0, unref(fileName).split('.')[0].length] as [number, number]
})

const emitSelect = (): void => {
  if (unref(chooseFileName)) {
    postMessage<embedModeLocationPickMessageData>('owncloud-embed:select', {
      resources: JSON.parse(JSON.stringify(selectedFiles.value)),
      fileName: unref(fileName),
      locationQuery: JSON.parse(JSON.stringify(routeToContextQuery(unref(router.currentRoute))))
    })
  }

  // TODO: adjust type to embedModeLocationPickMessageData later (breaking)
  postMessage<Resource[]>('owncloud-embed:select', JSON.parse(JSON.stringify(selectedFiles.value)))
}

const emitCancel = (): void => {
  postMessage<null>('owncloud-embed:cancel', null)
}
</script>

<style lang="scss">
.files-embed-actions {
  // Prevent .snackbar from overlapping the actions
  z-index: calc(var(--oc-z-index-modal) + 2);
  color: var(--oc-color-text-inverse);
  flex-wrap: wrap;
  gap: var(--oc-space-small);

  &-file-name {
    margin-left: 230px;
    gap: var(--oc-space-small);

    input {
      width: 400px;
    }

    @media (max-width: $oc-breakpoint-medium-default) {
      margin-left: 0;

      input {
        width: auto;
      }
    }
  }

  &-buttons {
    margin-left: auto;
  }
}
</style>
