<template>
  <span
    class="oc-display-inline-block oc-mb-m"
    v-text="$gettext('Your changes were not saved. Do you want to save them?')"
  />
  <div class="oc-my-m"></div>
  <div class="oc-flex oc-flex-right oc-flex-middle oc-mt-m">
    <div class="oc-modal-body-actions-grid">
      <oc-button
        class="oc-modal-body-actions-cancel oc-ml-s"
        appearance="outline"
        variation="passive"
        @click="$emit('cancel')"
        >{{ $gettext('Cancel') }}
      </oc-button>
      <oc-button
        class="oc-modal-body-actions-secondary oc-ml-s"
        appearance="outline"
        variation="passive"
        @click="onClose"
      >
        {{ $gettext("Don't Save") }}
      </oc-button>
      <oc-button
        class="oc-modal-body-actions-confirm oc-ml-s"
        appearance="filled"
        variation="primary"
        @click="$emit('confirm')"
        >{{ $gettext('Save') }}
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Modal, useModals } from '../../composables'
import { useGettext } from 'vue3-gettext'

interface Props {
  modal: Modal
  closeCallback: () => void
}

interface Emits {
  (e: 'cancel'): void
  (e: 'confirm'): void
}

defineEmits<Emits>()
const { modal, closeCallback } = defineProps<Props>()
const { removeModal } = useModals()
const { $gettext } = useGettext()
const onClose = () => {
  removeModal(modal.id)
  closeCallback()
}
</script>
