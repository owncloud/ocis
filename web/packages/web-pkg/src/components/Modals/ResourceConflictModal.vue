<template>
  <span class="oc-display-inline-block oc-mb-m" v-text="message" />
  <div class="oc-my-m">
    <oc-checkbox
      v-if="conflictCount > 1"
      v-model="checkboxValue"
      size="medium"
      :label="checkboxLabel"
      :aria-label="checkboxLabel"
    />
  </div>
  <div class="oc-flex oc-flex-right oc-flex-middle oc-mt-m">
    <div class="oc-modal-body-actions-grid">
      <oc-button
        class="oc-modal-body-actions-cancel oc-ml-s"
        appearance="outline"
        variation="passive"
        @click="onCancel"
        >{{ $gettext('Skip') }}
      </oc-button>
      <oc-button
        class="oc-modal-body-actions-secondary oc-ml-s"
        appearance="outline"
        variation="passive"
        @click="onConfirmSecondary"
        >{{ confirmSecondaryText }}
      </oc-button>
      <oc-button
        class="oc-modal-body-actions-confirm oc-ml-s"
        appearance="filled"
        variation="primary"
        @click="onConfirm"
        >{{ $gettext('Keep both') }}
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { Modal, useModals } from '../../composables'
import { Resource } from '@ownclouders/web-client'
import { ResolveConflict, ResolveStrategy } from '../../helpers/resource'

interface Props {
  modal: Modal
  resource: Resource
  conflictCount: number
  callbackFn: (resolveConflict: ResolveConflict) => void
  suggestMerge?: boolean
  separateSkipHandling?: boolean
  confirmSecondaryTextOverwrite?: string
}
const {
  modal,
  resource,
  conflictCount,
  callbackFn,
  suggestMerge = true,
  separateSkipHandling = false,
  confirmSecondaryTextOverwrite = null
} = defineProps<Props>()

const { removeModal } = useModals()
const { $gettext } = useGettext()

const checkboxValue = ref(false)
const checkboxLabel = computed(() => {
  if (conflictCount < 2) {
    return ''
  }
  if (!separateSkipHandling) {
    return $gettext('Apply to all %{count} conflicts', { count: conflictCount.toString() }, true)
  } else if (resource.isFolder) {
    return $gettext('Apply to all %{count} folders', { count: conflictCount.toString() }, true)
  } else {
    return $gettext('Apply to all %{count} files', { count: conflictCount.toString() }, true)
  }
})

const message = computed(() =>
  resource.isFolder
    ? $gettext('Folder with name "%{name}" already exists.', { name: resource.name }, true)
    : $gettext('File with name "%{name}" already exists.', { name: resource.name }, true)
)

const confirmSecondaryText = computed(() => {
  return confirmSecondaryTextOverwrite || $gettext('Replace')
})

const onConfirm = () => {
  removeModal(modal.id)
  callbackFn({
    strategy: ResolveStrategy.KEEP_BOTH,
    doForAllConflicts: unref(checkboxValue)
  })
}

const onConfirmSecondary = () => {
  removeModal(modal.id)
  const strategy = suggestMerge ? ResolveStrategy.MERGE : ResolveStrategy.REPLACE
  callbackFn({
    strategy,
    doForAllConflicts: unref(checkboxValue)
  })
}

const onCancel = () => {
  removeModal(modal.id)
  callbackFn({
    strategy: ResolveStrategy.SKIP,
    doForAllConflicts: unref(checkboxValue)
  })
}
</script>
